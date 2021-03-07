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

	if hashString, exists := lookForString(urlStoreObject, str); exists {
		return hashString
	}
	indexValue := 0
	var hashString string
	for {
		hashString := generateHashString(str, 0)
		urlStoreObjectInverse, err := getObjectMap("urls-storage-inverse.yaml")

		if _, exists := lookForString(urlStoreObjectInverse, str); exists {
			indexValue++
			continue
		}

		urlStoreObject[str] = hashString
		urlStoreObjectInverse[hashString] = str
		err = WriteObjectToFile(urlStoreObject, "urls-storage.yaml")
		if err != nil {
			fmt.Println("An Error Occured while writing the map objects to file for URLStore object", err.Error())
			return err.Error()
		}

		err = WriteObjectToFile(urlStoreObjectInverse, "urls-storage-inverse.yaml")
		if err != nil {
			fmt.Println("An Error Occured while writing the map objects to file for Inverse URLStore object", err.Error())
			return err.Error()
		}
		break

	}
	return hashString

}

func generateHashString(str string, startIndex int) string {
	data := []byte(str)
	hash := sha256.Sum256(data)
	if len(hash) >= 5 {
		hashString := hex.EncodeToString(hash[startIndex:5])
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
