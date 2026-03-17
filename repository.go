package main

var products []Products

func repoGetAll() []Products {
	return products
}

func repoGetByID(id int64) *Products {
	for i := range products {
		if products[i].Id == id {
			return &products[i]
		}
	}
	return nil
}

func repoCreate(p Products) Products {
	p.Id = int64(len(products) + 1)
	products = append(products, p)
	return p
}

func repoUpdate(id int64, p Products) *Products {
	for i := range products {
		if products[i].Id == id {
			p.Id = id
			products[i] = p
			return &products[i]
		}
	}
	return nil
}

func repoDelete(id int64) bool {
	for i := range products {
		if products[i].Id == id {
			products = append(products[:i], products[i+1:]...)
			return true
		}
	}
	return false
}
