package api

import (
	"ZorgTechCatalogus/pkg/cache"
	"ZorgTechCatalogus/pkg/database"
	"ZorgTechCatalogus/pkg/middleware"
	"context"
	"net/http"
	"time"

	docs "ZorgTechCatalogus/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

func ContextMiddleware(zorgtechRepo *zorgTechProductRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set("zorgTechProductRepository", zorgtechRepo)

		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	zorgTechProductRepository := NewZorgTechProductRepository(db, redisClient, ctx)

	r := gin.Default()
	r.Use(ContextMiddleware(zorgTechProductRepository))

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
		// Healthcheck
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy"})
		})
		// Catalogus endpoints
		catalogus := v1.Group("/catalogus")
		{
			// Commands
			catalogus.POST("/product", middleware.APIKeyAuth(), middleware.JWTAuth(), zorgTechProductRepository.MaakZorgTechProduct)
			catalogus.PUT("/product", middleware.APIKeyAuth(), middleware.JWTAuth(), zorgTechProductRepository.WijzigZorgTechProduct)
			catalogus.DELETE("/product", middleware.APIKeyAuth(), middleware.JWTAuth(), zorgTechProductRepository.VerwijderZorgTechProduct)
			catalogus.POST("/product/technisch-detail", middleware.APIKeyAuth(), middleware.JWTAuth(), zorgTechProductRepository.VoegTechnischDetailToe)
			catalogus.DELETE("/product/technisch-detail", middleware.APIKeyAuth(), middleware.JWTAuth(), zorgTechProductRepository.VerwijderTechnischDetail)

			// Queries
			catalogus.GET("/product/:zorgtechId", middleware.APIKeyAuth(), zorgTechProductRepository.GetProductById)
			catalogus.GET("/categorie/:categorie", middleware.APIKeyAuth(), zorgTechProductRepository.FindByCategorie)
			catalogus.GET("/producten", middleware.APIKeyAuth(), zorgTechProductRepository.ListAlleProducten)
			catalogus.GET("/zoek", middleware.APIKeyAuth(), zorgTechProductRepository.ZoekOpNaam)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
