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

func ContextMiddleware(zorgtechRepo *catalogusRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Set("catalogusRepository", zorgtechRepo)

		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	catalogusRepository := NewCatalogusRepository(db, redisClient, ctx)

	r := gin.Default()
	r.Use(ContextMiddleware(catalogusRepository))

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
			catalogus.POST("/product", catalogusRepository.MaakZorgTechProduct)
			catalogus.PUT("/product", catalogusRepository.WijzigZorgTechProduct)
			catalogus.DELETE("/product", catalogusRepository.VerwijderZorgTechProduct)
			catalogus.POST("/product/technisch-detail", catalogusRepository.VoegTechnischDetailToe)
			catalogus.DELETE("/product/technisch-detail", catalogusRepository.VerwijderTechnischDetail)

			// Queries
			catalogus.GET("/product/:zorgtechId", catalogusRepository.GetProductById)
			catalogus.GET("/categorie/:categorie", catalogusRepository.FindByCategorie)
			catalogus.GET("/producten", catalogusRepository.ListAlleProducten)
			catalogus.GET("/zoek", catalogusRepository.ZoekOpNaam)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
