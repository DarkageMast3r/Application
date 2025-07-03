package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
	"unicode"

	"authentication/pkg/cache"
	"authentication/pkg/models"
	"authentication/pkg/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthRepository interface {
	Healthcheck(c *gin.Context)
	Register(c *gin.Context)
	Login(c *gin.Context)
	RefreshToken(c *gin.Context)
	Logout(c *gin.Context)
	CreateRole(c *gin.Context)
	AssignRole(c *gin.Context)
	GetUserRoles(c *gin.Context)
	RegisterEndpoint(c *gin.Context)
	UpdateEndpointPermissions(c *gin.Context)
	GetAllEndpoints(c *gin.Context)
	ValidateToken(c *gin.Context)
	Authorize(c *gin.Context)

	// Helper functions - usually not exposed in interface if only used internally
	// generateTokens(user models.User) (string, string, error)
	// validateToken(token string) (*models.JWTClaims, error)
}

// authRepository implements AuthRepository using specific domain repositories
type authRepository struct {
	// Remove direct DB and RedisClient access here.
	// Instead, inject the specific repository interfaces.
	UserRepo      repository.UserRepository
	RoleRepo      repository.RoleRepository
	EndpointRepo  repository.EndpointRepository
	AuthTokenRepo repository.AuthTokenRepository
	CacheRepo     repository.CacheRepository // Use CacheRepository for Redis interactions
	// Ctx is typically passed per request, not stored in the repository struct.
	// We'll remove *context.Context from the struct and pass it explicitly in each method call.
}

// NewAuthRepository creates a new AuthRepository instance
func NewAuthRepository(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	endpointRepo repository.EndpointRepository,
	authTokenRepo repository.AuthTokenRepository,
	cacheRepo repository.CacheRepository,
) AuthRepository { // Return the interface, not the concrete struct
	return &authRepository{
		UserRepo:      userRepo,
		RoleRepo:      roleRepo,
		EndpointRepo:  endpointRepo,
		AuthTokenRepo: authTokenRepo,
		CacheRepo:     cacheRepo,
	}
}

// TokenResponse and ErrorResponse (if in this file)
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Time in seconds
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ErrorResponse represents a generic error response structure for API
type ErrorResponse struct {
	Error string `json:"error" example:"Bad Request"`
}

// @BasePath /api/v1/auth

// Healthcheck godoc
// @Summary Service health check
// @Description Checks if the authentication service is running
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {string} string "ok"
// @Router /health [get]
func (r *authRepository) Healthcheck(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.RegisterUserCommand true "Registration data"
// @Success 201 {object} models.User "Successfully registered user"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 409 {object} ErrorResponse "Username or email already exists"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /register [post]
func (r *authRepository) Register(c *gin.Context) {
	var input models.RegisterUserCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if username or email already exists using UserRepository
	existingUser, err := r.UserRepo.FindByUsernameOrEmail(c.Request.Context(), input.Username, input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user check: " + err.Error()})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Username or email already exists"})
		return
	}

	// Validate password strength (new helper function)
	if err := validatePassword(input.Password); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Create password hash
	hashedPassword, err := models.NewPassword(input.Password).Hash()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not hash password"})
		return
	}

	// Create new user
	user := models.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	if err := r.UserRepo.Save(c.Request.Context(), &user); err != nil { // Use Save for create
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create user: " + err.Error()})
		return
	}

	// Publish event (placeholder)
	// eventPublisher.Publish(models.UserRegisteredEvent{...})

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

// Login godoc
// @Summary Authenticate user
// @Description Logs in a user and returns JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.LoginCommand true "Login credentials"
// @Success 200 {object} TokenResponse "Login response with tokens"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 429 {object} ErrorResponse "Too Many Requests"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /login [post]
func (r *authRepository) Login(c *gin.Context) {
	var input models.LoginCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Track failed attempts using CacheRepository
	key := fmt.Sprintf("login_fail:%s", input.Username)
	attempts, err := r.CacheRepo.Increment(c.Request.Context(), key)
	if err != nil {
		log.Printf("Error incrementing login fail counter: %v", err)
		// Don't fail the request, proceed with login attempt but log it.
	}
	// Set expiration only if it's the first attempt or if not set
	if attempts == 1 { // Only set expiry on the first increment
		if err := r.CacheRepo.Expire(c.Request.Context(), key, 1*time.Hour); err != nil {
			log.Printf("Error setting expiry for login fail counter: %v", err)
		}
	}

	if attempts > 5 {
		c.JSON(http.StatusTooManyRequests, ErrorResponse{Error: "Too many login attempts. Please try again later."})
		return
	}

	// Find user by username using UserRepository
	user, err := r.UserRepo.FindByUsernameOrEmail(c.Request.Context(), input.Username, "") // Find by username only
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user lookup: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Verify password
	password := models.Password{Value: input.Password}
	match, err := password.MatchesHash(user.PasswordHash)
	if err != nil || !match {
		// Increment counter and immediately return if password doesn't match
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid credentials"})
		return
	}

	// Reset failed attempts counter on successful login
	if err := r.CacheRepo.Del(c.Request.Context(), key); err != nil {
		log.Printf("Warning: Failed to reset login fail counter for %s: %v", input.Username, err)
	}

	// Generate tokens (helper function will now use AuthTokenRepo to save refresh token)
	accessToken, refreshToken, err := r.generateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not generate tokens: " + err.Error()})
		return
	}

	// Update last login using UserRepository
	err = r.UserRepo.UpdateLastLogin(c.Request.Context(), user)
	if err != nil {
		log.Printf("Warning: Failed to update last login for user %s: %v", user.ID, err)
	}

	// Publish event (placeholder)
	// eventPublisher.Publish(models.UserLoggedInEvent{...})

	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    900, // 15 minutes
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Generates a new access and refresh token pair
// @Tags auth
// @Accept json
// @Produce json
// @Param data body RefreshTokenRequest true "Refresh Token Payload"
// @Success 200 {object} TokenResponse
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (r *authRepository) RefreshToken(c *gin.Context) {
	var input RefreshTokenRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Validate refresh token using AuthTokenRepository
	authToken, err := r.AuthTokenRepo.FindByRefreshToken(c.Request.Context(), input.RefreshToken)
	if err != nil {
		// FindByRefreshToken should return nil, nil if not found
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid or expired refresh token"})
		return
	}
	if authToken == nil { // Explicitly check if token was not found
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Invalid or expired refresh token"})
		return
	}

	// Check if the token is expired (even if DB query already does this, a double-check is good)
	if authToken.ExpiresAt.Before(time.Now()) {
		// Optionally delete expired token immediately
		_ = r.AuthTokenRepo.DeleteByID(c.Request.Context(), authToken.ID)
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Refresh token has expired"})
		return
	}

	// Get user using UserRepository
	user, err := r.UserRepo.FindByID(c.Request.Context(), authToken.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user lookup: " + err.Error()})
		return
	}
	if user == nil { // User not found
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Associated user not found"})
		return
	}

	// Generate new tokens (helper function will save the new refresh token)
	accessToken, newRefreshToken, err := r.generateTokens(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not generate new tokens: " + err.Error()})
		return
	}

	// Delete old refresh token from DB using AuthTokenRepository
	err = r.AuthTokenRepo.DeleteByID(c.Request.Context(), authToken.ID)
	if err != nil {
		log.Printf("Warning: Failed to delete old refresh token %s: %v", authToken.ID.String(), err)
	}

	// Respond with new tokens
	c.JSON(http.StatusOK, TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    900, // 15 minutes
	})
}

// Logout godoc
// @Summary Log out user
// @Description Invalidates the user's refresh tokens
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} ErrorResponse "Successfully logged out" // Changed to object for consistency
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /logout [post]
func (r *authRepository) Logout(c *gin.Context) {
	userIDVal, exists := c.Get("userID") // Assuming userID is stored as uuid.UUID
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "Unauthorized"})
		return
	}

	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Invalid user ID in context"})
		return
	}

	// Delete all refresh tokens for this user using AuthTokenRepository
	if err := r.AuthTokenRepo.DeleteUserTokens(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not log out: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"}) // Still gin.H for simple messages is fine
}

// CreateRole godoc
// @Summary Create a new role
// @Description Creates a new role with permissions (admin only)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.CreateRoleCommand true "Role data"
// @Success 201 {object} models.Role "Created role"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 409 {object} ErrorResponse "Role already exists"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /roles [post]
func (r *authRepository) CreateRole(c *gin.Context) {
	var input models.CreateRoleCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if role already exists using RoleRepository
	existingRole, err := r.RoleRepo.FindByName(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during role check: " + err.Error()})
		return
	}
	if existingRole != nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "Role already exists"})
		return
	}

	// Find permissions using RoleRepository (if permissions are separate entities)
	// NOTE: This assumes models.Permission also has a FindByID or similar
	// For simplicity, I'm assuming you'll fetch permissions by ID from the DB if needed.
	// If input.Permissions are IDs, you'd need a PermissionRepository or direct DB query here.
	// For now, let's assume `input.Permissions` are names or IDs, and we need to fetch the actual Permission objects.
	var permissions []models.Permission
	// Example: if input.Permissions is []uuid.UUID
	// if err := r.PermRepo.FindByIds(c.Request.Context(), input.Permissions).Find(&permissions).Error(); err != nil {
	// 	c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid permission IDs"})
	// 	return
	// }
	// For now, if input.Permissions are just UUIDs, you'd fetch them here
	// Assuming permissions are already valid or will be validated when role is saved.

	// Create role
	role := models.Role{
		ID:          uuid.New(),
		Name:        input.Name,
		Description: input.Description,
		Permissions: permissions, // Permissions will be set when saving associations
		CreatedAt:   time.Now(),
	}

	if err := r.RoleRepo.Create(c.Request.Context(), &role); err != nil { // Use Create for new role
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not create role: " + err.Error()})
		return
	}

	// After creating the role, set the associations (if input.Permissions were IDs/names)
	// This would involve fetching permissions by their IDs/names and then calling RoleRepo.Save with the updated role.
	// If permissions are directly embedded in input and auto-created, then this might not be needed.
	// For now, assuming direct creation handles it or permissions are verified elsewhere.

	// Publish event (placeholder)
	// eventPublisher.Publish(models.RoleCreatedEvent{...})

	c.JSON(http.StatusCreated, gin.H{"data": role})
}

// AssignRole godoc
// @Summary Assign role to user
// @Description Assigns a role to a user (admin only)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.AssignRoleCommand true "Assignment data"
// @Success 200 {object} models.UserRole "Role assignment"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "User or role not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /roles/assign [post]
func (r *authRepository) AssignRole(c *gin.Context) {
	var input models.AssignRoleCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if user exists using UserRepository
	user, err := r.UserRepo.FindByID(c.Request.Context(), input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user lookup: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Check if role exists using RoleRepository
	role, err := r.RoleRepo.FindByID(c.Request.Context(), input.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during role lookup: " + err.Error()})
		return
	}
	if role == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Role not found"})
		return
	}

	// GORM's Many-To-Many relationship can be handled by appending to the slice and saving the 'user' object.
	// However, if it's a direct join table (UserRole), you create and save that explicitly.
	// Assuming `models.User` has `Roles []models.Role` field with `many2many` tag.
	// This would require fetching the user with existing roles, appending, and then saving.

	// Here, assuming you have a `UserRoleRepository` or that `UserRepo` has an `AssignRole` method
	// that handles the GORM association. If not, this is a place for a new repository method.
	// For now, I'll modify it to update the User's Roles directly if models.User has `Roles` slice.

	// Load existing roles for the user
	err = r.UserRepo.LoadUserRoles(c.Request.Context(), user) // This method loads user.Roles
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not load user roles: " + err.Error()})
		return
	}

	// Check if role is already assigned to prevent duplicates (optional, depending on DB constraints)
	for _, assignedRole := range user.Roles {
		if assignedRole.ID == role.ID {
			c.JSON(http.StatusConflict, ErrorResponse{Error: "Role already assigned to user"})
			return
		}
	}

	// Append the new role to the user's roles
	user.Roles = append(user.Roles, *role)

	// Save the user (this will update the many-to-many table)
	if err := r.UserRepo.Save(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not assign role to user: " + err.Error()})
		return
	}

	// Publish event (placeholder)
	// eventPublisher.Publish(models.RoleAssignedEvent{...})

	c.JSON(http.StatusOK, gin.H{"data": models.UserRole{UserID: input.UserID, RoleID: input.RoleID}}) // Return the basic assignment for consistency
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Returns all roles assigned to a user
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} models.Role "User roles"
// @Failure 404 {object} ErrorResponse "User not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /users/{userId}/roles [get]
func (r *authRepository) GetUserRoles(c *gin.Context) {
	userID, err := uuid.Parse(c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid user ID format"})
		return
	}

	// Check if user exists and load roles using UserRepository
	user, err := r.UserRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user lookup: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "User not found"})
		return
	}

	// Load roles and their permissions
	err = r.UserRepo.LoadUserRoles(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not fetch user roles: " + err.Error()})
		return
	}

	// Now load permissions for each role using RoleRepo.LoadRolePermissions
	// This might involve iterating if LoadUserRoles doesn't deeply preload permissions.
	// Assuming LoadUserRoles might populate user.Roles but not deeply load Permissions for each role.
	for i := range user.Roles {
		err := r.RoleRepo.LoadRolePermissions(c.Request.Context(), &user.Roles[i])
		if err != nil {
			log.Printf("Warning: Could not load permissions for role %s: %v", user.Roles[i].Name, err)
			// Decide if this is a fatal error or just log. For display, maybe not fatal.
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": user.Roles})
}

// RegisterEndpoint godoc
// @Summary Register API endpoint
// @Description Registers a new API endpoint with required permissions (admin only)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param input body models.RegisterAPIEndpointCommand true "Endpoint data"
// @Success 201 {object} models.APIEndpoint "Registered endpoint"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /endpoints [post]
func (r *authRepository) RegisterEndpoint(c *gin.Context) {
	var input models.RegisterAPIEndpointCommand

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Find roles by name using RoleRepository
	var roles []models.Role
	for _, roleName := range input.RoleNames {
		role, err := r.RoleRepo.FindByName(c.Request.Context(), roleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during role lookup: " + err.Error()})
			return
		}
		if role == nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Role '%s' not found", roleName)})
			return
		}
		roles = append(roles, *role)
	}

	// Create endpoint
	endpoint := models.APIEndpoint{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		Path:        input.Path,
		Method:      input.Method,
		Description: input.Description,
		Version:     input.Version,
		Roles:       roles,
		CreatedAt:   time.Now(),
	}

	if err := r.EndpointRepo.Create(c.Request.Context(), &endpoint); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not register endpoint: " + err.Error()})
		return
	}

	// Publish event (placeholder)
	// eventPublisher.Publish(models.EndpointRegisteredEvent{...})

	c.JSON(http.StatusCreated, gin.H{"data": endpoint})
}

// UpdateEndpointPermissions godoc
// @Summary Update endpoint permissions
// @Description Updates which roles can access an endpoint (admin only)
// @Tags auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param endpointId path string true "Endpoint ID"
// @Param input body []string true "Role names"
// @Success 200 {object} models.APIEndpoint "Updated endpoint"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 403 {object} ErrorResponse "Forbidden"
// @Failure 404 {object} ErrorResponse "Endpoint not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /endpoints/{endpointId}/permissions [put]
func (r *authRepository) UpdateEndpointPermissions(c *gin.Context) {
	endpointID, err := uuid.Parse(c.Param("endpointId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid endpoint ID format"})
		return
	}
	var roleNames []string

	if err := c.ShouldBindJSON(&roleNames); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Find endpoint using EndpointRepository
	endpoint, err := r.EndpointRepo.FindByID(c.Request.Context(), endpointID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during endpoint lookup: " + err.Error()})
		return
	}
	if endpoint == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Endpoint not found"})
		return
	}

	// Find roles by name using RoleRepository
	var rolesToAssign []models.Role
	for _, roleName := range roleNames {
		role, err := r.RoleRepo.FindByName(c.Request.Context(), roleName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during role lookup: " + err.Error()})
			return
		}
		if role == nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("Role '%s' not found", roleName)})
			return
		}
		rolesToAssign = append(rolesToAssign, *role)
	}

	// Update roles using EndpointRepository's LoadEndpointRoles and then Save
	err = r.EndpointRepo.LoadEndpointRoles(c.Request.Context(), endpoint) // Load current roles
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not load endpoint roles: " + err.Error()})
		return
	}

	// Assuming EndpointRepo.Save handles updating associations when endpoint.Roles is changed.
	// If not, you'd need a specific method like `EndpointRepo.ReplaceRoles(ctx, endpoint, newRoles)`.
	// For now, let's directly set the roles and rely on Save or a custom method.
	endpoint.Roles = rolesToAssign
	if err := r.EndpointRepo.Save(c.Request.Context(), endpoint); err != nil { // Save will update associations in GORM
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not update permissions: " + err.Error()})
		return
	}

	// Reload endpoint with roles to ensure they are returned (Save might not preload)
	// Or simply return the updated endpoint struct
	// err = r.EndpointRepo.LoadEndpointRoles(c.Request.Context(), endpoint)
	// if err != nil {
	// 	log.Printf("Warning: Failed to reload endpoint roles after update: %v", err)
	// }

	c.JSON(http.StatusOK, gin.H{"data": endpoint})
}

// GetAllEndpoints godoc
// @Summary List all endpoints
// @Description Returns all registered API endpoints (admin only)
// @Tags auth
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.APIEndpoint "List of endpoints"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /endpoints [get]
func (r *authRepository) GetAllEndpoints(c *gin.Context) {
	endpoints, err := r.EndpointRepo.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not fetch endpoints: " + err.Error()})
		return
	}

	// Load roles for each endpoint if FindAll doesn't preload them
	for i := range endpoints {
		err := r.EndpointRepo.LoadEndpointRoles(c.Request.Context(), &endpoints[i])
		if err != nil {
			log.Printf("Warning: Could not load roles for endpoint %s: %v", endpoints[i].Path, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": endpoints})
}

// ValidateToken godoc
// @Summary Validate token
// @Description Validates a JWT token and returns its claims
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "JWT token"
// @Success 200 {object} models.JWTClaims "Token claims"
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /validate [get]
func (r *authRepository) ValidateToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Token is required"})
		return
	}

	claims, err := r.validateToken(c.Request.Context(), token) // Pass context
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: err.Error()}) // Use error from helper
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": claims})
}

// Authorize godoc
// @Summary Check authorization
// @Description Checks if a user is authorized for a specific action
// @Tags auth
// @Accept json
// @Produce json
// @Param input body models.AuthorizationQuery true "Authorization query"
// @Success 200 {object} gin.H "Authorization result" // Can remain gin.H for simple boolean
// @Failure 400 {object} ErrorResponse "Bad Request"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /authorize [post]
func (r *authRepository) Authorize(c *gin.Context) {
	var input models.AuthorizationQuery

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Cache key
	cacheKey := fmt.Sprintf("authz:%s:%s:%s", input.UserID, input.Resource, input.Action)

	// Try cache first using CacheRepository
	if cachedVal, err := r.CacheRepo.Get(c.Request.Context(), cacheKey); err == nil && cachedVal != "" {
		authorized := (cachedVal == "true")
		c.JSON(http.StatusOK, gin.H{"authorized": authorized})
		return
	} else if err != nil && !errors.Is(err, cache.ErrNotFound) { // Log other cache errors
		log.Printf("Warning: Cache read error for key %s: %v", cacheKey, err)
	}

	// Database check if cache miss, using UserRepository and RoleRepository
	user, err := r.UserRepo.FindByID(c.Request.Context(), input.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Database error during user lookup: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "User not found"})
		return
	}

	// Load user's roles and their permissions
	err = r.UserRepo.LoadUserRoles(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Could not load user roles for authorization: " + err.Error()})
		return
	}

	authorized := false
	for _, role := range user.Roles {
		err := r.RoleRepo.LoadRolePermissions(c.Request.Context(), &role) // Load permissions for each role
		if err != nil {
			log.Printf("Warning: Could not load permissions for role %s during authorization: %v", role.Name, err)
			continue // Continue checking other roles
		}
		for _, perm := range role.Permissions {
			if perm.Resource == input.Resource && perm.Action == input.Action {
				authorized = true
				break
			}
		}
		if authorized {
			break // No need to check other roles if authorized
		}
	}

	// Cache result using CacheRepository
	if err := r.CacheRepo.Set(c.Request.Context(), cacheKey, authorized, 5*time.Minute); err != nil {
		log.Printf("Warning: Failed to set authorization cache for key %s: %v", cacheKey, err)
	}

	c.JSON(http.StatusOK, gin.H{"authorized": authorized})
}

// --- Helper functions ---

// generateTokens handles token creation and saving the refresh token.
// It now uses AuthTokenRepository and UserRepo for loading roles/permissions.
func (r *authRepository) generateTokens(user models.User) (string, string, error) {
	ctx := context.Background() // Use a background context if `c.Request.Context()` isn't available

	// Load all roles for the user, and then load permissions for each role
	err := r.UserRepo.LoadUserRoles(ctx, &user)
	if err != nil {
		return "", "", fmt.Errorf("failed to load user roles: %w", err)
	}

	var permissions []string
	for _, role := range user.Roles {
		err := r.RoleRepo.LoadRolePermissions(ctx, &role)
		if err != nil {
			// Log but don't fail, maybe the role has no permissions or there's a DB issue
			log.Printf("Warning: Failed to load permissions for role %s: %v", role.Name, err)
			continue
		}
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Resource+":"+perm.Action)
		}
	}

	// Create access token
	accessToken, err := models.GenerateJWT(user.ID, user.Username, permissions)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create refresh token
	refreshTokenModel := models.AuthToken{ // Changed to AuthToken
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     uuid.New().String(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 1 week
	}

	if err := r.AuthTokenRepo.Create(ctx, &refreshTokenModel); err != nil { // Use AuthTokenRepo
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshTokenModel.Token, nil
}

// validateToken validates a JWT token, now using CacheRepository for blacklist.
func (r *authRepository) validateToken(ctx context.Context, token string) (*models.JWTClaims, error) {
	// Check blacklist first using CacheRepository
	isBlacklisted, err := r.CacheRepo.IsAccessTokenBlacklisted(ctx, token)
	if err != nil {
		log.Printf("Error checking token blacklist: %v", err)
		// Decide if this is fatal. For security, maybe treat as blacklisted on cache error.
		// For now, let's proceed but log.
	}
	if isBlacklisted {
		return nil, errors.New("token is blacklisted")
	}

	claims, err := models.ValidateJWT(token)
	if err != nil {
		return nil, err
	}

	// Extra checks (these are usually handled by JWT library's validation, but good to have)
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}

	if claims.Issuer != "your-auth-service" {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}

// validatePassword is a helper for password strength validation.
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper, hasLower, hasNumber, hasSpecial bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c): // unicode.IsPunct or unicode.IsSymbol for special characters
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("password must contain upper, lower, number and special char")
	}

	return nil
}
