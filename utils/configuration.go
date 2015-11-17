package utils

import (
"encoding/json"
"log"
"os"
	"time"
)

// Configuration is used to store json information
type Configuration struct {
Port          	string
Gitfolder      	string
Updateinterval  time.Duration
Htmlfiles 		string
Ignorefolders 	[]string
}

//NewConfiguration create configuration from file
func NewConfiguration(path string) *Configuration {
	file, _ := os.Open(path)
	decoder := json.NewDecoder(file)
	c := Configuration{}
	err := decoder.Decode(&c)
	if err != nil {
	log.Fatal("malformed json file ", err)
	}
return &c
}