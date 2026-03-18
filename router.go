package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "makarov_project/docs"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/auth/token", handleAuthToken)

	api := r.Group("").Use(authMiddleware())
	{
		api.GET("/products", handleGetProducts)
		api.GET("/products/:id", handleGetProductByID)
		api.POST("/add_products", handleCreateProduct)
		api.PUT("/products/:id", handleUpdateProduct)
		api.DELETE("/products/:id", handleDeleteProduct)
	}

	return r
}

// @Summary Get JWT token
// @Security BearerAuth
// @Description Returns a JWT token for API authentication
// @Tags Auth
// @Produce json
// @Success 200 {object} TokenResponse "JWT token"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /auth/token [post]
func handleAuthToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(settings.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed})
}

// @Summary Handler Get Products with filters: min, max, in_stock
// @Security BearerAuth
// @Description.markdown get_products
// @Tags Products
// @Produce json
// @Param min_price query int false "Minimum price filter"
// @Param max_price query int false "Maximum price filter"
// @Param in_stock query bool false "Filter by stock availability"
// @Success 200 {object} []main.Products "List of products"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Router /products [get]
func handleGetProducts(c *gin.Context) {
	var filter ProductFilter

	if s := c.Query("min_price"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат min_price"})
			return
		}
		filter.MinPrice = &v
	}
	if s := c.Query("max_price"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат max_price"})
			return
		}
		filter.MaxPrice = &v
	}
	if filter.MinPrice != nil && filter.MaxPrice != nil && *filter.MinPrice > *filter.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price должен быть меньше или равен max_price"})
		return
	}

	filter.InStock = parseInStock(c.Query("in_stock"))

	c.JSON(http.StatusOK, gin.H{"data": serviceGetAll(filter)})
}

// @Summary Handler Get Products by id
// @Security BearerAuth
// @Description.markdown get_products_by_id.md
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} main.Products "Product"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Router /products/{id} [get]
func handleGetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	p := serviceGetByID(id)
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": p})
}

// @Summary Post Handler add_products
// @Security BearerAuth
// @Description.markdown post_products.md
// @Tags Products
// @Accept json
// @Produce json
// @Param product body main.Products true "Product data"
// @Success 201 {object} main.Products "Created product"
// @Failure 404 {object} ErrorResponse "Not Found"
// @Router /add_products [post]
func handleCreateProduct(c *gin.Context) {
	var p Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceCreate(p))
}

// @Summary Update product by id
// @Security BearerAuth
// @Description Updates an existing product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body main.Products true "Product data"
// @Success 200 {object} main.Products "Updated product"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /products/{id} [put]
func handleUpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	var p Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := serviceUpdate(id, p)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Summary Delete product by id
// @Security BearerAuth
// @Description Deletes a product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Deletion confirmation"
// @Failure 400 {object} ErrorResponse "Bad request"
// @Failure 404 {object} ErrorResponse "Product not found"
// @Router /products/{id} [delete]
func handleDeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	if !serviceDelete(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Продукт %d удалён", id)})
}
