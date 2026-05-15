package repository

import "github.com/Ilpaka/go-products-api/internal/model"

type ProductRepo struct {
	items []model.Products
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{}
}

func (r *ProductRepo) GetAll() []model.Products {
	return r.items
}

func (r *ProductRepo) GetByID(id int64) *model.Products {
	for i := range r.items {
		if r.items[i].Id == id {
			return &r.items[i]
		}
	}
	return nil
}

func (r *ProductRepo) Create(p model.Products) model.Products {
	p.Id = int64(len(r.items) + 1)
	r.items = append(r.items, p)
	return p
}

func (r *ProductRepo) Update(id int64, p model.Products) *model.Products {
	for i := range r.items {
		if r.items[i].Id == id {
			p.Id = id
			r.items[i] = p
			return &r.items[i]
		}
	}
	return nil
}

func (r *ProductRepo) Delete(id int64) bool {
	for i := range r.items {
		if r.items[i].Id == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return true
		}
	}
	return false
}
