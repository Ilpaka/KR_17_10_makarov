package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Products struct {
	Id      int64  `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	InStock int    `json:"in_stock"`
}

var products []Products

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT токен не передан"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный формат заголовка Authorization. Ожидается: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неожиданный метод подписи: %v", token.Header["alg"])
			}
			return []byte(settings.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Недействительный или истёкший JWT токен"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	news := r.Group("").Use(authMiddleware())
	{
		news.GET("/products", get_products)
		news.GET("/products/:id", get_products_by_id)
		news.POST("/add_products", create_products)
		news.PUT("/products/:id", update_product_by_id)
		news.DELETE("/products/:id", delete_product_by_id)
	}

	r.Run(":8090")
}

func get_products(c *gin.Context) {
	result := make([]Products, 0, len(products))

	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	inStockStr := c.Query("in_stock")

	var minPrice, maxPrice *int
	if minPriceStr != "" {
		v, err := strconv.Atoi(minPriceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат min_price"})
			return
		}
		minPrice = &v
	}
	if maxPriceStr != "" {
		v, err := strconv.Atoi(maxPriceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат max_price"})
			return
		}
		maxPrice = &v
	}

	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price должен быть меньше или равен max_price"})
		return
	}

	filterInStock := inStockStr == "1" || strings.EqualFold(inStockStr, "true")

	for _, p := range products {
		if minPrice != nil && p.Price < *minPrice {
			continue
		}
		if maxPrice != nil && p.Price > *maxPrice {
			continue
		}
		if filterInStock && p.InStock <= 0 {
			continue
		}
		result = append(result, p)
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
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

func update_product_by_id(c *gin.Context) {
	updateID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	var updated Products
	if err := c.ShouldBindJSON(&updated); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range products {
		if products[i].Id == updateID {
			updated.Id = updateID
			products[i] = updated
			c.JSON(http.StatusOK, products[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
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
