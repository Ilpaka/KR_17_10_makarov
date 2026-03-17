package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

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

func handleCreateProduct(c *gin.Context) {
	var p Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceCreate(p))
}

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
