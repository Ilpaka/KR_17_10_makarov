package main

import "strings"

type ProductFilter struct {
	MinPrice *int
	MaxPrice *int
	InStock  bool
}

func serviceGetAll(filter ProductFilter) []Products {
	all := repoGetAll()
	result := make([]Products, 0, len(all))

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

func serviceGetByID(id int64) *Products {
	return repoGetByID(id)
}

func serviceCreate(p Products) Products {
	return repoCreate(p)
}

func serviceUpdate(id int64, p Products) *Products {
	return repoUpdate(id, p)
}

func serviceDelete(id int64) bool {
	return repoDelete(id)
}

func parseInStock(s string) bool {
	return s == "1" || strings.EqualFold(s, "true")
}
