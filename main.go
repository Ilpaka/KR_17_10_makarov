package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Products struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	InStock int    `json:"in_stock"`
}

var products []Products

func main() {
	r := gin.Default()

	news := r.Group("")
	{
		news.GET("/products", get_products)
		news.GET("/products/:id", get_products_by_id)
		news.POST("/add_products", create_products)
		news.DELETE("/products/:id", delete_product_by_id)
	}

	r.Run(":8090")
}

func get_products(c *gin.Context) {
	{
		c.JSON(http.StatusOK, gin.H{"data": products})
	}
}

func create_products(c *gin.Context) {
	var product Products

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product.Id = int64(len(products) + 1)
	products = append(products, product)

	c.JSON(http.StatusCreated, product)
}

func get_products_by_id(c *gin.Context) {
	get_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	var get_new *Products
	for i := range products {
		if products[i].Id == get_id {
			get_new = &products[i]
			break
		}
	}

	if get_new == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": get_new})
}

func delete_product_by_id(c *gin.Context) {
	deleteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	idx := -1
	for i := range products {
		if products[i].Id == deleteID {
			idx = i
			break
		}
	}

	if idx == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	products = append(products[:idx], products[idx+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Продукт %d удалён", deleteID)})
}
