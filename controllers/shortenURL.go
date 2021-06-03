package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	template "net/url"

	_ "github.com/apple/foundationdb/bindings/go/src/fdb"
	_ "github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/srisarangarp/url-shortener/foundationdb"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

var dbObject *foundationdb.FdbWrapper

func SetDbObject(object *foundationdb.FdbWrapper) {
	dbObject = object
}
func ShortenURL(c *gin.Context) {
	url := c.PostForm("url")

	if hashString, err, exists := dbObject.GetFromUrlMatch(url); exists && hashString != "" {
		if err != nil {
			fmt.Printf("An Error occured in getting the urlmatch subspace %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"shorten_URL": returnCompleteShortenURL(hashString)})
		return
	}
	indexValue := 0
	var hashedString string
	for {

		hashedString = generateHashString(url, indexValue)

		if exists, err := dbObject.LookIntoHasCodeMatch(hashedString); exists {
			if err != nil {
				fmt.Printf("An Error occured in getting the mapObject from yaml file %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
				return
			}
			indexValue++
			continue
		}

		dbObject.InsertIntoUrlMatch(url, hashedString)
		dbObject.InsertIntoHashCodeMatch(hashedString, url)
		log.Println("The Hashed string is ", hashedString)
		break

	}
	c.JSON(http.StatusOK, gin.H{"shorten_URL": returnCompleteShortenURL(hashedString)})

}

func generateHashString(str string, startIndex int) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	if len(hash) >= 5+startIndex {
		hashString := hex.EncodeToString(hash[startIndex : 5+startIndex])
		return hashString
	}
	return "Unable to hash with current configuration"

}

func getObjectMap(path string) (map[string]string, error) {
	urlStoreObject := make(map[string]string)
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &urlStoreObject)
	if err != nil {
		return nil, err
	}
	return urlStoreObject, nil

}

func lookForString(urlStoreObject map[string]string, url string) (string, bool) {
	if val, exists := urlStoreObject[url]; exists {
		return val, exists
	}
	return "", false
}

func WriteObjectToFile(urlStoreObject map[string]string, path string) error {
	data, err := yaml.Marshal(urlStoreObject)
	if err != nil {
		fmt.Printf("An Error occured in marshalling the map object %s", err.Error())
		return err
	}

	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		fmt.Printf("An Error occured in writing the map data to yaml file %s", err.Error())
		return err
	}
	return nil
}
func returnCompleteShortenURL(shortenURL string) string {
	myURLtemplate := "https://infc.com/here-is-the-shorten-url"
	url, _ := template.Parse(myURLtemplate)
	url.Path = shortenURL
	return url.String()
}
