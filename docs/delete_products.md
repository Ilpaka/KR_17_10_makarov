Удаляет товар по ID.

**Реализация:**

<pre><code>func handleDeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат id"})
		return
	}

	if !serviceDelete(id) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Продукт %d удалён", id)})
}</code></pre>
