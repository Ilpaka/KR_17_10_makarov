package service

import (
	"strings"

	"github.com/Ilpaka/go-products-api/internal/model"
	"github.com/Ilpaka/go-products-api/internal/repository"
)

type ProductFilter struct {
	MinPrice *int
	MaxPrice *int
	InStock  bool
}

type ProductService struct {
	repo *repository.ProductRepo
}

func NewProductService(repo *repository.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) GetAll(filter ProductFilter) []model.Products {
	all := s.repo.GetAll()
	result := make([]model.Products, 0, len(all))

	for _, p := range all {
		if filter.MinPrice != nil && p.Price < *filter.MinPrice {
			continue
		}
		if filter.MaxPrice != nil && p.Price > *filter.MaxPrice {
			continue
		}
		if filter.InStock && p.InStock <= 0 {
			continue
		}
		result = append(result, p)
	}

	return result
}

func (s *ProductService) GetByID(id int64) *model.Products {
	return s.repo.GetByID(id)
}

func (s *ProductService) Create(p model.Products) model.Products {
	return s.repo.Create(p)
}

func (s *ProductService) Update(id int64, p model.Products) *model.Products {
	return s.repo.Update(id, p)
}

func (s *ProductService) Delete(id int64) bool {
	return s.repo.Delete(id)
}

func ParseInStock(s string) bool {
	return s == "1" || strings.EqualFold(s, "true")
}
