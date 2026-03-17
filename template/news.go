package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

//Реализовать API описывающие новостную ленту.
// Класс должен содержать:
// ■ массив новостей; ✅
// ■ get-свойство, которое возвращает количество новостей;✅
// ■ метод для вывода на экран всех новостей;✅
// ■ метод для добавления новости;✅
// ■ метод для удаления новости;✅
// ■ метод для сортировки новостей по дате (от последних новостей до старых); ✅
// ■ метод для поиска новостей по тегу (возвращает массив новостей, в которых указан переданный в метод тег).✅
//
// Продемонстрировать работу написанных методов.

type News struct {
	Header string
	Date   time.Time
	Teg    string
	Id     int
}

var news []News

func main() {
	r := gin.Default()

	news := r.Group("")
	{
		news.GET("/news", get_news)
		news.GET("/count", get_quantity_of_news)
		news.GET("/sort_news", get_sort_news)
		news.GET("/sort_by_teg/:teg", get_news_by_teg)
		news.POST("/add_news", create_news)
		news.DELETE("/new/:id", delete_news_by_id)
	}

	r.Run(":8090")
}

func create_news(c *gin.Context) {
	var new News

	if err := c.ShouldBindJSON(&new); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error:": err.Error()})
		return
	}

	new.Id = len(news) + 1
	new.Date = time.Now()

	news = append(news, new)

	c.JSON(http.StatusCreated, new)
}

func get_news(c *gin.Context) {
	{
		c.JSON(http.StatusOK, gin.H{"data": news})
	}
}

func get_quantity_of_news(c *gin.Context) {
	{
		c.JSON(http.StatusOK, gin.H{"Count": len(news)})
	}
}

func delete_news_by_id(c *gin.Context) {
	idx := -1
	delete_id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		{
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Некорректный формат id"})
		}
		return
	}

	for i, u := range news {
		if u.Id == delete_id {
			idx = i
			break
		}
	}

	if idx == -1 {
		{
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Такой новости не сущетсвует"})
		}
		return
	}

	news = append(news[:idx], news[idx+1:]...)
	{
		c.JSON(http.StatusOK, gin.H{"Data": fmt.Sprintf("Новость %d удалена", delete_id)})
	}
}

func get_sort_news(c *gin.Context) {
	sort_news := news

	sort.Slice(sort_news, func(i, j int) bool {
		return news[i].Date.Before(news[j].Date)
	})

	{
		c.JSON(http.StatusOK, gin.H{"data": sort_news})
	}

}

func get_news_by_teg(c *gin.Context) {
	teg := c.Param("teg")

	news_by_tegs := []News{}

	for _, i := range news {
		if i.Teg == teg {
			news_by_tegs = append(news_by_tegs, i)
		}
	}
	{
		c.JSON(http.StatusOK, gin.H{"data": news_by_tegs})
	}
}
