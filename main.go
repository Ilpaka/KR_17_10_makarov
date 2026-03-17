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
		news.GET("/products")
		news.POST("/add_news")
		news.DELETE("/new/:id")
	}

	r.Run(":8090")
}

func get_products(c *gin.Context) {
	{
		c.JSON(http.StatusOK, gin.H{"data": products})
	}
}
