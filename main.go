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
	shortenURL("https://www.facebook.com")
}

func shortenURL(str string) {
	data := []byte(str)
	hash := sha256.Sum256(data)
	yamlFile, err := ioutil.ReadFile("urls-storage.yaml")
	if err != nil {
		fmt.Printf("An Error occured in reading the yaml file %s", err.Error())
	}

	urlStoreObject := make(map[string]string)
	fmt.Println(string(yamlFile))
	err = yaml.Unmarshal(yamlFile, &urlStoreObject)
	if err != nil {
		fmt.Printf("An Error occured in un marshalling the yaml file %s", err.Error())
	}

	if _, exists := urlStoreObject[str]; exists {
		fmt.Println(urlStoreObject[str])
	} else {
		hashString := hex.EncodeToString(hash[0:7])
		urlStoreObject[str] = hashString
		data, err = yaml.Marshal(urlStoreObject)
		if err != nil {
			fmt.Printf("An Error occured in marshalling the map object %s", err.Error())
		}

		err = ioutil.WriteFile("urls-storage.yaml", data, 0644)
		if err != nil {
			fmt.Printf("An Error occured in writing the map data to yaml file %s", err.Error())
		}
	}

}
