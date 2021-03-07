package main

import (
	_ "fmt"
	"net/http"
	_ "net/url"

	"github.com/gin-gonic/gin"
	"github.com/sreesa7144/url-shortener/controllers"
)

func main() {
	//fmt.Println(shortenURL("https://www.google.com"))
	r := gin.Default()
	welcomeMessage := "Welcome to the URL Shortener please send your url to the /url path"
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": welcomeMessage})
	})
	r.POST("/url", controllers.ShortenURL)
	r.Run()
}
