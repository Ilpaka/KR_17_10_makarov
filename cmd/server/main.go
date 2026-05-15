// @title Products API KR_17_03
// @version 1.0
// @description Simple CRUD API for Products with JWT auth.
// @host localhost:8090
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer {token}" to authorize.
package main

import (
	"github.com/Ilpaka/go-products-api/internal/config"
	"github.com/Ilpaka/go-products-api/internal/handler"
	"github.com/Ilpaka/go-products-api/internal/repository"
	"github.com/Ilpaka/go-products-api/internal/service"

	_ "github.com/Ilpaka/go-products-api/docs"
)

func main() {
	cfg := config.Load()
	repo := repository.NewProductRepo()
	svc := service.NewProductService(repo)

	r := handler.NewRouter(cfg, svc)
	r.Run(":8090")
}
