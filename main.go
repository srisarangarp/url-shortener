package main

import (
	_ "fmt"
	"net/http"
	_ "net/url"

	"github.com/gin-gonic/gin"
	"github.com/srisarangarp/url-shortener/controllers"
	"github.com/srisarangarp/url-shortener/foundationdb"
)

//Lets Integrate foundationDB to replace yaml files

func main() {

	dbObject := &foundationdb.FdbWrapper{}
	dbObject.CreateDb("urlMatch", "hashCodeMatch")
	controllers.SetDbObject(dbObject)

	r := gin.Default()
	welcomeMessage := "Welcome to the URL Shortener please send your url to the /url path"
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": welcomeMessage})
	})
	r.POST("/url", controllers.ShortenURL)
	r.Run()
}

/*
What will be the data model
1. I will create a directory namely urlspace
2.Then i will create two subspaces a) hashcodeMatch b)urlMatch
3.In URL match store URL and hashcode as KV pair
4. In hashcode match save hashcode and url as the KV pair

*/
