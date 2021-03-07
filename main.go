package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	_ "fmt"
	"io/ioutil"
	_ "net/url"

	"github.com/ghodss/yaml"
)

func main() {
	fmt.Println(shortenURL("https://www.google.com"))
}

func shortenURL(str string) string {
	urlStoreObject, err := getObjectMap("urls-storage.yaml")
	if err != nil {
		fmt.Printf("An Error occured in getting the mapObject from yaml file %s", err.Error())
		return err.Error()
	}

	if hashString, exists := lookForURL(urlStoreObject, str); exists {
		return hashString
	}

	hashString := generateHashString(str, 0)
	urlStoreObject[str] = hashString
	data, err := yaml.Marshal(urlStoreObject)
	if err != nil {
		fmt.Printf("An Error occured in marshalling the map object %s", err.Error())
	}

	err = ioutil.WriteFile("urls-storage.yaml", data, 0644)
	if err != nil {
		fmt.Printf("An Error occured in writing the map data to yaml file %s", err.Error())
	}
	return hashString

}

func generateHashString(str string, startIndex int) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[startIndex:7])
	return hashString
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

func lookForURL(urlStoreObject map[string]string, url string) (string, bool) {
	if val, exists := urlStoreObject[url]; exists {
		return val, exists
	}
	return "", false
}
