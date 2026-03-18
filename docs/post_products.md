Возвращает один товар

**Реализация:**

<pre><code>var p Products
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	func serviceCreate(p Products) Products {
		return repoCreate(p)
	}

	func repoCreate(p Products) Products {
	p.Id = int64(len(products) + 1)
	products = append(products, p)
	return p
}
</code></pre>
