package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Products struct {
	Id       int64
	Name     string
	Price    int
	in_stock int
}

var products []Products

func main() {
	r := gin.Default()

	news := r.Group("")
	{
		news.GET("/products", get_products)
		news.POST("/add_products", create_products)
		news.DELETE("/products/:id")
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
