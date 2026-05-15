package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/Ilpaka/go-products-api/internal/config"
	"github.com/Ilpaka/go-products-api/internal/middleware"
	"github.com/Ilpaka/go-products-api/internal/service"
)

func NewRouter(cfg *config.Config, svc *service.ProductService) *gin.Engine {
	r := gin.Default()

	auth := NewAuthHandler(cfg)
	products := NewProductHandler(svc)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/token", auth.Token)

	api := r.Group("").Use(middleware.AuthMiddleware(cfg))
	{
		api.GET("/products", products.List)
		api.GET("/products/:id", products.Get)
		api.POST("/add_products", products.Create)
		api.PUT("/products/:id", products.Update)
		api.DELETE("/products/:id", products.Delete)
	}

	return r
}
