package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	flag.String("host", "http://localhost:7700", "The host to use")
	flag.String("api-key", "masterKey", "The API key to use")
	flag.String("index", "movies", "The index to add documents to")
	flag.String("file", "", "The file to add documents from")
	flag.Parse()

	host := flag.Lookup("host").Value.String()
	apiKey := flag.Lookup("api-key").Value.String()
	index := flag.Lookup("index").Value.String()
	file := flag.Lookup("file").Value.String()

	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))

	if file == "" {
		panic("file is required")
	}
	jsonFile, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var movies []map[string]interface{}
	if err = json.Unmarshal(byteValue, &movies); err != nil {
		panic(err)
	}

	taskInfo, err := client.Index(index).AddDocuments(movies)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Task: %#+v\n", taskInfo)

	// wait for task to complete
	fmt.Println("Waiting for task to complete...")
	task, err := client.WaitForTask(taskInfo.TaskUID, 3*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Task: %#+v\n", task)
}
