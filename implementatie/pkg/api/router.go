package api

import (
	"ZorgTechImplementatie/pkg/cache"
	"ZorgTechImplementatie/pkg/database"
	"ZorgTechImplementatie/pkg/middleware"
	"context"
	"time"

	docs "ZorgTechImplementatie/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

func ContextMiddleware(implRepo *implementatieRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("implementatieRepository", implRepo)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	implementatieRepository := NewImplementatieRepository(db, redisClient, ctx)

	r := gin.Default()
	r.Use(ContextMiddleware(implementatieRepository))

	//r.Use(gin.Logger())
	r.Use(middleware.Logger(logger, mongoCollection))
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{

		// Implementatie endpoints
		implementatie := v1.Group("/implementatie")
		{
			// Commands
			implementatie.POST("/aanvraag", implementatieRepository.AanvraagProduct)
			implementatie.POST("/ontvang", implementatieRepository.OntvangProduct)
			implementatie.POST("/installeer", implementatieRepository.InstalleerProduct)
			implementatie.POST("/personaliseer", implementatieRepository.PersonaliseerProduct)
			implementatie.POST("/lever", implementatieRepository.LeverProduct)
			implementatie.POST("/voltooi", implementatieRepository.MarkeerAlsGeimplementeerd)

			// Queries
			implementatie.GET("/product/:zorgtechId", implementatieRepository.GetProductInformatie)
			implementatie.GET("/status/:clientId", implementatieRepository.GetImplementatieStatus)
			implementatie.GET("/installatie-status/:clientId", implementatieRepository.GetInstallatieStatus)
			implementatie.GET("/instellingen/:clientId", implementatieRepository.GetPersoonlijkeInstellingen)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
