package auth

import (
	middleware "authentication/pkg/middleware"
	"authentication/pkg/models"
	repository "authentication/pkg/repository"

	"github.com/golang-jwt/jwt/v5"

	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// RegisterUserCommand for user registration
type RegisterUserCommand struct {
	Username string
	Email    string
	Password string
}

// LoginCommand for user authentication
type LoginCommand struct {
	Username string
	Password string
}

// CreateRoleCommand for role management
type CreateRoleCommand struct {
	Name        string
	Description string
	Permissions []uuid.UUID // IDs van permissies
}

// AssignRoleCommand for assigning roles to users
type AssignRoleCommand struct {
	UserID  uuid.UUID
	RoleID  uuid.UUID
	AdminID uuid.UUID // Who is assigning this
}

// RegisterAPIEndpointCommand for endpoint registration
type RegisterAPIEndpointCommand struct {
	ServiceName string
	Path        string
	Method      string
	Description string
	Version     string
	RoleNames   []string // Which roles can access this
}

// RequestPasswordResetCommand for password reset
type RequestPasswordResetCommand struct {
	Email string
}

// ResetPasswordCommand for resetting password
type ResetPasswordCommand struct {
	Token    string
	Password string
}

// AuthorizationQuery for checking permissions
type AuthorizationQuery struct {
	UserID   uuid.UUID
	Resource string
	Action   string
}

// UserRegisteredEvent triggered after successful registration
type UserRegisteredEvent struct {
	UserID    uuid.UUID
	Username  string
	Email     string
	Timestamp time.Time
}

// UserLoggedInEvent for tracking logins
type UserLoggedInEvent struct {
	UserID    uuid.UUID
	IPAddress string
	UserAgent string
	Timestamp time.Time
}

// RoleAssignedEvent when a user gets a new role
type RoleAssignedEvent struct {
	UserID     uuid.UUID
	RoleID     uuid.UUID
	AssignedBy uuid.UUID
	Timestamp  time.Time
}

// EndpointRegisteredEvent when a new API endpoint is added
type EndpointRegisteredEvent struct {
	EndpointID   uuid.UUID
	Path         string
	Method       string
	Service      string
	RegisteredBy uuid.UUID
	Timestamp    time.Time
}

// AuthService defines the application service for authentication and authorization
type AuthService interface {
	RegisterUser(ctx context.Context, cmd RegisterUserCommand) (*models.User, error)
	LoginUser(ctx context.Context, cmd LoginCommand, ipAddress, userAgent string) (accessToken string, refreshToken string, expiresIn int, err error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken string, newRefreshToken string, expiresIn int, err error)
	LogoutUser(ctx context.Context, userID uuid.UUID, accessToken string) error
	RequestPasswordReset(ctx context.Context, cmd RequestPasswordResetCommand) error
	ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error

	CreateRole(ctx context.Context, cmd CreateRoleCommand) (*models.Role, error)
	AssignRoleToUser(ctx context.Context, cmd AssignRoleCommand) (*models.User, error)
	GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error)

	RegisterAPIEndpoint(ctx context.Context, cmd RegisterAPIEndpointCommand) (*models.APIEndpoint, error)
	UpdateEndpointPermissions(ctx context.Context, endpointID uuid.UUID, roleNames []string) (*models.APIEndpoint, error)
	GetAllEndpoints(ctx context.Context) ([]models.APIEndpoint, error)

	ValidateAccessToken(ctx context.Context, token string) (*models.JWTClaims, error)
	AuthorizeRequest(ctx context.Context, query AuthorizationQuery) (bool, error)
}

type authService struct {
	userRepo      repository.UserRepository
	roleRepo      repository.RoleRepository
	authTokenRepo repository.AuthTokenRepository // For saving refresh tokens
	cacheRepo     repository.CacheRepository     // For blacklisting, rate limiting
	endpointRepo  repository.EndpointRepository
	jwtSecret     []byte // This field is used in generateTokens
	// Add other dependencies as needed, e.g., emailService, eventPublisher
}

// authApplicationService implements AuthService
type authApplicationService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	endpointRepo   repository.EndpointRepository
	authTokenRepo  repository.AuthTokenRepository
	cacheRepo      repository.CacheRepository
	eventPublisher repository.EventPublisher
	jwtSecret      []byte
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	endpointRepo repository.EndpointRepository,
	authTokenRepo repository.AuthTokenRepository,
	cacheRepo repository.CacheRepository,
	eventPublisher repository.EventPublisher,
	jwtSecret []byte,
) AuthService {
	return &authApplicationService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		endpointRepo:   endpointRepo,
		authTokenRepo:  authTokenRepo,
		cacheRepo:      cacheRepo,
		eventPublisher: eventPublisher,
		jwtSecret:      jwtSecret,
	}
}
func (s *authApplicationService) RegisterUser(ctx context.Context, cmd RegisterUserCommand) (*models.User, error) {
	existingUser, err := s.repository.FindByUsernameOrEmail(ctx, cmd.Username, cmd.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return nil, errors.New("username or email already exists")
	}

	passwordVO, err := models.NewPassword(cmd.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password: %w", err)
	}
	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return nil, errors.New("could not hash password")
	}

	user := &models.User{
		ID:           uuid.New(),
		Username:     cmd.Username,
		Email:        cmd.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, fmt.Errorf("could not create user: %w", err)
	}

	s.eventPublisher.Publish(ctx, UserRegisteredEvent{
		UserID:    user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Timestamp: time.Now(),
	})

	return user, nil
}
func (s *authApplicationService) LoginUser(ctx context.Context, cmd LoginCommand, ipAddress, userAgent string) (accessToken string, refreshToken string, expiresIn int, err error) {
	key := fmt.Sprintf("login_fail:%s", cmd.Username)
	attempts, err := s.cacheRepo.Increment(ctx, key)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to increment login attempts: %w", err)
	}
	s.cacheRepo.Expire(ctx, key, 1*time.Hour) // Fire and forget for expire

	if attempts > 5 {
		return "", "", 0, errors.New("too many login attempts, please try again later")
	}

	user, err := s.userRepo.FindByUsernameOrEmail(ctx, cmd.Username, cmd.Username)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return "", "", 0, errors.New("invalid credentials")
	}

	passwordVO, err := models.NewPassword(cmd.Password)
	if err != nil {
		return "", "", 0, errors.New("invalid credentials")
	}
	match, err := passwordVO.MatchesHash(user.PasswordHash)
	if err != nil || !match {
		return "", "", 0, errors.New("invalid credentials")
	}

	accessToken, refreshToken, err = s.generateTokens(ctx, *user)
	if err != nil {
		return "", "", 0, fmt.Errorf("could not generate tokens: %w", err)
	}

	s.userRepo.UpdateLastLogin(ctx, user) // Update last login, fire and forget

	s.eventPublisher.Publish(ctx, UserLoggedInEvent{
		UserID:    user.ID,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		Timestamp: time.Now(),
	})

	s.cacheRepo.Set(ctx, key, 0, 0) // Reset login attempts on success

	return accessToken, refreshToken, 15 * 60, nil // 15 minutes
}
func (s *authApplicationService) RefreshToken(ctx context.Context, oldRefreshToken string) (accessToken string, newRefreshToken string, expiresIn int, err error) {
	token, err := s.authTokenRepo.FindByRefreshToken(ctx, oldRefreshToken)
	if err != nil {
		return "", "", 0, errors.New("invalid refresh token")
	}
	if token == nil || token.ExpiresAt.Before(time.Now()) {
		return "", "", 0, errors.New("invalid or expired refresh token")
	}

	user, err := s.userRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to find user for refresh: %w", err)
	}
	if user == nil {
		return "", "", 0, errors.New("user not found for refresh token")
	}

	accessToken, newRefreshToken, err = s.generateTokens(ctx, *user)
	if err != nil {
		return "", "", 0, fmt.Errorf("could not generate tokens: %w", err)
	}

	// Delete old refresh token
	if err := s.repository.DeleteUserTokens(ctx, token.UserID); err != nil {
		// Log this, but don't fail the refresh if deletion fails
		fmt.Printf("Warning: failed to delete old refresh token for user %s: %v\n", token.UserID, err)
	}

	return accessToken, newRefreshToken, 15 * 60, nil
}
func (s *authApplicationService) LogoutUser(ctx context.Context, userID uuid.UUID, accessToken string) error {
	claims, err := middleware.ValidateJWT(accessToken)
	if err != nil {
		return fmt.Errorf("invalid access token for logout: %w", err)
	}

	// Voeg token toe aan blacklist
	blacklisted := &models.TokenBlacklist{
		Token:     accessToken,
		ExpiresAt: time.Unix(claims.ExpiresAt, 0),
	}
	if err := repository.AddBlacklistedAccessToken(ctx, blacklisted); err != nil {
		return fmt.Errorf("could not blacklist token: %w", err)
	}

	// Verwijder refresh token
	if err := repository.DeleteRefreshTokensByUserID(ctx, userID); err != nil {
		return fmt.Errorf("could not delete refresh token: %w", err)
	}
	return nil
}

func (s *authApplicationService) RequestPasswordReset(ctx context.Context, cmd RequestPasswordResetCommand) error {
	user, err := repository.FindByUsernameOrEmail(ctx, cmd.Email, cmd.Email)
	if err != nil {
		return fmt.Errorf("failed to check user for password reset: %w", err)
	}
	if user == nil {
		// Liever geen foutmelding geven of het email bestaat (security)
		fmt.Println("Password reset requested for non-existent email (security measure).")
		return nil // Meld success om user enumeration te voorkomen
	}

	token := uuid.New().String()
	resetToken := &models.PasswordResetToken{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1 * time.Hour),
		Used:      false,
	}

	if err := repository.SavePasswordResetToken(ctx, resetToken); err != nil {
		return fmt.Errorf("could not create reset token: %w", err)
	}

	// In productie: Stuur email met reset link
	resetLink := fmt.Sprintf("https://arjanslab.nl/reset-password?token=%s", token)
	fmt.Printf("DEV: Password reset link for %s: %s\n", user.Email, resetLink)

	return nil
}

func (s *authApplicationService) ResetPassword(ctx context.Context, cmd ResetPasswordCommand) error {
	resetToken, err := repository.FindValidPasswordResetToken(ctx, cmd.Token)
	if err != nil {
		return fmt.Errorf("failed to find reset token: %w", err)
	}
	if resetToken == nil || resetToken.Used || resetToken.ExpiresAt.Before(time.Now()) {
		return errors.New("invalid or expired token")
	}

	passwordVO, err := models.NewPassword(cmd.Password)
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}
	hashedPassword, err := passwordVO.Hash()
	if err != nil {
		return errors.New("could not hash password")
	}

	user, err := s.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user for password reset: %w", err)
	}
	if user == nil {
		return errors.New("user not found for password reset")
	}
	user.PasswordHash = hashedPassword

	if err := s.userRepo.Save(ctx, user); err != nil {
		return fmt.Errorf("could not update password: %w", err)
	}

	if err := repository.MarkPasswordResetTokenUsed(ctx, resetToken); err != nil {
		fmt.Printf("Warning: failed to mark password reset token used: %v\n", err) // Log but don't fail
	}

	// Invalidate all sessions for this user
	if err := s.repository.DeleteRefreshTokensByUserID(ctx, user.ID); err != nil {
		fmt.Printf("Warning: failed to invalidate old refresh tokens after password reset: %v\n", err) // Log but don't fail
	}

	return nil
}

func (s *authApplicationService) CreateRole(ctx context.Context, cmd CreateRoleCommand) (*models.Role, error) {
	existingRole, err := s.repository.FindByName(ctx, cmd.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existingRole != nil {
		return nil, errors.New("role already exists")
	}

	// Fetch permissions by IDs
	var permissions []models.Permission

	for _, id := range cmd.Permissions {
		permissions = append(permissions, models.Permission{ID: id})
	}

	role := &models.Role{
		ID:          uuid.New(),
		Name:        cmd.Name,
		Description: cmd.Description,
		Permissions: permissions,
		CreatedAt:   time.Now(),
	}

	if err := s.roleRepo.Save(ctx, role); err != nil {
		return nil, fmt.Errorf("could not create role: %w", err)
	}

	s.eventPublisher.Publish(ctx, RoleAssignedEvent{
		// ?
	})

	return role, nil
}

func (s *authApplicationService) AssignRoleToUser(ctx context.Context, cmd AssignRoleCommand) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, cmd.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	role, err := s.roleRepo.FindByIds(ctx, []uuid.UUID{cmd.RoleID}) // FindByIds returneert een slice, dus aanpassen
	if err != nil || len(role) == 0 {
		return nil, errors.New("role not found")
	}

	if err := s.roleRepo.AddUserRole(ctx, cmd.UserID, cmd.RoleID, cmd.AdminID); err != nil {
		return nil, fmt.Errorf("could not assign role: %w", err)
	}

	s.eventPublisher.Publish(ctx, RoleAssignedEvent{
		UserID:     cmd.UserID,
		RoleID:     cmd.RoleID,
		AssignedBy: cmd.AdminID,
		Timestamp:  time.Now(),
	})

	return user, nil
}

func (s *authApplicationService) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]models.Role, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	roles, err := s.roleRepo.FindUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user roles: %w", err)
	}

	// Laad permissies voor elke rol
	for i := range roles {
		if err := s.userRepo.LoadRolePermissions(ctx, &roles[i]); err != nil {
			return nil, fmt.Errorf("failed to load permissions for role %s: %w", roles[i].Name, err)
		}
	}

	return roles, nil
}

func (s *authApplicationService) RegisterAPIEndpoint(ctx context.Context, cmd RegisterAPIEndpointCommand) (*models.APIEndpoint, error) {
	roles, err := s.repository.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all roles: %w", err)
	}

	var allowedRoles []models.Role
	for _, role := range roles {
		for _, name := range cmd.RoleNames {
			if role.Name == name {
				allowedRoles = append(allowedRoles, role)
				break
			}
		}
	}
	if len(allowedRoles) != len(cmd.RoleNames) {
		return nil, errors.New("one or more role names are invalid")
	}

	endpoint := &models.APIEndpoint{
		ID:          uuid.New(),
		ServiceName: cmd.ServiceName,
		Path:        cmd.Path,
		Method:      cmd.Method,
		Description: cmd.Description,
		Version:     cmd.Version,
		Roles:       allowedRoles,
		CreatedAt:   time.Now(),
	}

	if err := s.endpointRepo.Save(ctx, endpoint); err != nil {
		return nil, fmt.Errorf("could not register endpoint: %w", err)
	}

	s.eventPublisher.Publish(ctx, EndpointRegisteredEvent{
		EndpointID:   endpoint.ID,
		Path:         endpoint.Path,
		Method:       endpoint.Method,
		Service:      endpoint.ServiceName,
		RegisteredBy: uuid.Nil,
		Timestamp:    time.Now(),
	})

	return endpoint, nil
}

func (s *authApplicationService) UpdateEndpointPermissions(ctx context.Context, endpointID uuid.UUID, roleNames []string) (*models.APIEndpoint, error) {
	endpoint, err := s.endpointRepo.FindByID(ctx, endpointID)
	if err != nil {
		return nil, fmt.Errorf("failed to find endpoint: %w", err)
	}
	if endpoint == nil {
		return nil, errors.New("endpoint not found")
	}

	roles, err := s.roleRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all roles: %w", err)
	}

	var newRoles []models.Role
	for _, role := range roles {
		for _, name := range roleNames {
			if role.Name == name {
				newRoles = append(newRoles, role)
				break
			}
		}
	}
	if len(newRoles) != len(roleNames) {
		return nil, errors.New("one or more new role names are invalid")
	}

	if err := s.endpointRepo.UpdateRoles(ctx, endpoint, newRoles); err != nil {
		return nil, fmt.Errorf("could not update endpoint permissions: %w", err)
	}

	// Reload endpoint with updated roles
	if err := s.endpointRepo.LoadEndpointRoles(ctx, endpoint); err != nil {
		return nil, fmt.Errorf("failed to reload endpoint roles: %w", err)
	}

	return endpoint, nil
}

func (s *authApplicationService) GetAllEndpoints(ctx context.Context) ([]models.APIEndpoint, error) {
	endpoints, err := s.endpointRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not fetch endpoints: %w", err)
	}
	// Preload roles for all endpoints (dit kan in de repository geoptimaliseerd worden met Preload)
	for i := range endpoints {
		if err := s.endpointRepo.LoadEndpointRoles(ctx, &endpoints[i]); err != nil {
			return nil, fmt.Errorf("failed to load roles for endpoint %s: %w", endpoints[i].ID, err)
		}
	}
	return endpoints, nil
}

func (s *authApplicationService) ValidateAccessToken(ctx context.Context, token string) (*models.JWTClaims, error) {
	if blacklisted, err := s.authTokenRepo.IsBlacklisted(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %w", err)
	} else if blacklisted {
		return nil, errors.New("token is blacklisted")
	}

	claims, err := jwt.ValidateJWT(token)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	if claims.Issuer != "your-auth-service" {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

func (s *authApplicationService) AuthorizeRequest(ctx context.Context, query AuthorizationQuery) (bool, error) {
	cacheKey := fmt.Sprintf("authz:%s:%s:%s", query.UserID, query.Resource, query.Action)

	if cached, err := s.cacheRepo.Get(ctx, cacheKey); err == nil {
		return cached == "true", nil
	}

	user, err := s.userRepo.FindByID(ctx, query.UserID)
	if err != nil {
		return false, fmt.Errorf("failed to find user for authorization: %w", err)
	}
	if user == nil {
		return false, errors.New("user not found")
	}

	// Load user roles and their permissions
	if err := s.userRepo.LoadUserRoles(ctx, user); err != nil {
		return false, fmt.Errorf("failed to load user roles for authorization: %w", err)
	}

	authorized := false
	for _, role := range user.Roles {
		if err := s.userRepo.LoadRolePermissions(ctx, &role); err != nil {
			return false, fmt.Errorf("failed to load role permissions for authorization: %w", err)
		}
		for _, perm := range role.Permissions {
			if perm.Resource == query.Resource && perm.Action == query.Action {
				authorized = true
				break
			}
		}
		if authorized {
			break
		}
	}

	if err := s.cacheRepo.Set(ctx, cacheKey, authorized, 5*time.Minute); err != nil {
		fmt.Printf("Warning: failed to cache authorization result: %v\n", err) // Log, maar niet fatal
	}

	return authorized, nil
}
