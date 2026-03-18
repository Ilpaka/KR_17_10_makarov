Создаёт новый товар.

**Реализация:**

<pre><code>func handleCreateProduct(c *gin.Context) {
	var p Products
	if err := c.ShouldBindJSON(&amp;p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, serviceCreate(p))
}</code></pre>
