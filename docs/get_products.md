Возвращает список товаров с фильтрами.

**Реализация:**

<pre><code>func handleGetProducts(c *gin.Context) {
	var filter ProductFilter

	if s := c.Query("min_price"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат min_price"})
			return
		}
		filter.MinPrice = &amp;v
	}
	if s := c.Query("max_price"); s != "" {
		v, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат max_price"})
			return
		}
		filter.MaxPrice = &amp;v
	}
	if filter.MinPrice != nil &amp;&amp; filter.MaxPrice != nil &amp;&amp; *filter.MinPrice &gt; *filter.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "min_price должен быть меньше или равен max_price"})
		return
	}

	filter.InStock = parseInStock(c.Query("in_stock"))

	c.JSON(http.StatusOK, gin.H{"data": serviceGetAll(filter)})
}</code></pre>
