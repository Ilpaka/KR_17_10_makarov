package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/Ilpaka/go-products-api/internal/model"
	"github.com/Ilpaka/go-products-api/internal/service"
)

type ProductHandler struct {
	svc *service.ProductService
}

func NewProductHandler(svc *service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

// @Summary Handler Get Products with filters: min, max, in_stock
// @Security BearerAuth
// @Description.markdown get_products
// @Tags Products
// @Produce json
// @Param min_price query int false "Minimum price filter"
// @Param max_price query int false "Maximum price filter"
// @Param in_stock query bool false "Filter by stock availability"
// @Success 200 {object} []model.Products "List of products"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Router /products [get]
func (h *ProductHandler) List(c *gin.Context) {
	var filter service.ProductFilter

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

	filter.InStock = service.ParseInStock(c.Query("in_stock"))

	c.JSON(http.StatusOK, gin.H{"data": h.svc.GetAll(filter)})
}

// @Summary Handler Get Products by id
// @Security BearerAuth
// @Description.markdown get_products_by_id
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} model.Products "Product"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	p := h.svc.GetByID(id)
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": p})
}

// @Summary Post Handler add_products
// @Security BearerAuth
// @Description.markdown post_products
// @Tags Products
// @Accept json
// @Produce json
// @Param product body model.Products true "Product data"
// @Success 201 {object} model.Products "Created product"
// @Failure 404 {object} model.ErrorResponse "Not Found"
// @Router /add_products [post]
func (h *ProductHandler) Create(c *gin.Context) {
	var p model.Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, h.svc.Create(p))
}

// @Summary Update product by id
// @Security BearerAuth
// @Description.markdown put_products
// @Tags Products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body model.Products true "Product data"
// @Success 200 {object} model.Products "Updated product"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Router /products/{id} [put]
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	var p model.Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := h.svc.Update(id, p)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Summary Delete product by id
// @Security BearerAuth
// @Description.markdown delete_products
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{} "Deletion confirmation"
// @Failure 400 {object} model.ErrorResponse "Bad request"
// @Failure 404 {object} model.ErrorResponse "Product not found"
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	if !h.svc.Delete(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Продукт %d удалён", id)})
}
