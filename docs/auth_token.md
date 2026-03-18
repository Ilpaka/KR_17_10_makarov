Возвращает JWT токен для аутентификации API.

**Реализация:**

<pre><code>func handleAuthToken(c *gin.Context) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	})
	signed, err := token.SignedString([]byte(settings.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": signed})
}</code></pre>
