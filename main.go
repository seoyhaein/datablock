package main

import (
	"encoding/json"
	blk "github.com/seoyhaein/datablock/watch"
	"log"
	"os"
)

func main() {

}

func loadConfig(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := make(map[string]interface{})
	if err := decoder.Decode(&config); err != nil {
		return err
	}

	if maxWatchCount, ok := config["MaxWatchCount"].(float64); ok {
		blk.MaxWatchCount = int(maxWatchCount)
	} else {
		log.Printf("Invalid or missing MaxWatchCount in configuration")
	}

	if rootDirValue, ok := config["rootDir"].(string); ok {
		blk.RootDir = rootDirValue
	} else {
		log.Printf("Invalid or missing rootDir in configuration")
	}

	// Initialize the eventQueue after configuration is loaded
	//eventQueue = make([]fsnotify.Event, 0)

	return nil
}
