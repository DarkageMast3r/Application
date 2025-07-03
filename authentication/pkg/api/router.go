package api

import (
	"authentication/pkg/cache"
	"authentication/pkg/database"
	"authentication/pkg/middleware"
	"context"
	"time"

	docs "authentication/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

func ContextMiddleware(authRepo api.AuthRepository) gin.HandlerFunc { // Changed type to api.AuthRepository
	return func(c *gin.Context) {
		// This line is often discouraged for performance.
		// Handlers should preferably get their dependencies through closures during router setup.
		c.Set("authRepository", authRepo)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	// --- Initialize concrete implementations of your domain repositories ---
	// These concrete implementations wrap your `database.Database` and `cache.Cache` interfaces.

	userRepo := repository.NewGormUserRepository(db)
	roleRepo := repository.NewGormRoleRepository(db)
	endpointRepo := repository.NewGormEndpointRepository(db)
	authTokenRepo := repository.NewGormAuthTokenRepository(db)   // For DB-backed refresh tokens
	cacheRepo := repository.NewRedisCacheRepository(redisClient) // For Redis-backed cache/blacklist/rate-limiting

	authRepository := api.NewAuthRepository(userRepo, roleRepo, endpointRepo, authTokenRepo, cacheRepo)

	r := gin.Default()

	// Global Middlewares
	r.Use(middleware.Logger(logger, mongoCollection))
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.WithRedisClient(redisClient)) // This middleware sets the redis client in context, mainly for older patterns or for other middlewares to use.
	// For handlers using `cacheRepo`, this specific middleware isn't strictly necessary.

	// New rate limiter configuration
	apiRateLimiter := middleware.RateLimiter(rate.Every(1*time.Minute), 60)
	authRateLimiter := middleware.RateLimiter(rate.Every(30*time.Second), 5) // Strenger voor auth endpoints

	// Swagger setup
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		// Authenticatie endpoints
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authRateLimiter, authRepository.Login)
			auth.POST("/register", authRateLimiter, authRepository.Register)
			auth.POST("/refresh", authRepository.RefreshToken)
			auth.POST("/logout", middleware.JWTAuth(), authRepository.Logout)
			// Ensure these methods exist on your api.AuthRepository interface and its implementation
			auth.POST("/request-password-reset", authRateLimiter, authRepository.RequestPasswordReset)
			auth.POST("/reset-password", authRepository.ResetPassword)
			// Duplicate routes removed:
			// auth.POST("/request-password-reset", authRateLimiter, authRepository.RequestPasswordReset)
			// auth.POST("/reset-password", authRepository.ResetPassword)

			auth.GET("/validate", authRepository.ValidateToken)
			auth.POST("/authorize", authRepository.Authorize)
		}

		// Admin endpoints
		admin := v1.Group("/admin")
		admin.Use(middleware.JWTAuth(), middleware.RequireRole("admin"))
		{
			admin.POST("/roles", authRepository.CreateRole)
			admin.POST("/roles/assign", authRepository.AssignRole)
			admin.GET("/users/:id/roles", authRepository.GetUserRoles)
			admin.POST("/endpoints", authRepository.RegisterEndpoint)
			admin.PUT("/endpoints/:id/permissions", authRepository.UpdateEndpointPermissions)
			admin.GET("/endpoints", authRepository.GetAllEndpoints)
		}
	}

	// Health check en swagger
	r.GET("/health", authRepository.Healthcheck)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
