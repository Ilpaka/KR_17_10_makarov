package main

import (
	"net/http"
	"strconv"

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
		news.DELETE("/products/:id", get_products_by_id)
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
