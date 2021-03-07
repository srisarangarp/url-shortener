package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"

	template "net/url"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func ShortenURL(c *gin.Context) {
	url := c.PostForm("url")

	urlStoreObject, err := getObjectMap("urls-storage.yaml")
	if err != nil {
		fmt.Printf("An Error occured in getting the mapObject from yaml file %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if hashString, exists := lookForString(urlStoreObject, url); exists {
		c.JSON(http.StatusOK, gin.H{"shorten_URL": returnCompleteShortenURL(hashString)})
		return
	}
	indexValue := 0
	var hashString string
	for {

		hashString = generateHashString(url, indexValue)
		urlStoreObjectInverse, err := getObjectMap("urls-storage-inverse.yaml")
		if err != nil {
			fmt.Printf("An Error occured in getting the mapObject from yaml file %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		//In case of collision
		if _, exists := lookForString(urlStoreObjectInverse, hashString); exists {
			indexValue++
			continue
		}

		urlStoreObject[url] = hashString
		urlStoreObjectInverse[hashString] = url
		err = WriteObjectToFile(urlStoreObject, "urls-storage.yaml")
		if err != nil {
			fmt.Println("An Error Occured while writing the map objects to file for URLStore object", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}

		err = WriteObjectToFile(urlStoreObjectInverse, "urls-storage-inverse.yaml")
		if err != nil {
			fmt.Println("An Error Occured while writing the map objects to file for Inverse URLStore object", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
			return
		}
		break

	}
	c.JSON(http.StatusOK, gin.H{"shorten_URL": returnCompleteShortenURL(hashString)})

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
	fmt.Println(string(yamlFile))

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
