Обновляет существующий товар по ID.

**Реализация:**

<pre><code>func handleUpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	var p Products
	if err := c.ShouldBindJSON(&amp;p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated := serviceUpdate(id, p)
	if updated == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, updated)
}</code></pre>
