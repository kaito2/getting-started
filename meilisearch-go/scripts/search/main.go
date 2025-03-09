package main

import (
	"flag"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	flag.String("host", "http://localhost:7700", "The host to use")
	flag.String("api-key", "masterKey", "The API key to use")
	flag.String("index", "movies", "The index to search")
	flag.String("query", "", "The query to search for")
	flag.Parse()

	host := flag.Lookup("host").Value.String()
	apiKey := flag.Lookup("api-key").Value.String()
	index := flag.Lookup("index").Value.String()
	query := flag.Lookup("query").Value.String()

	client := meilisearch.New(host, meilisearch.WithAPIKey(apiKey))
	results, err := client.Index(index).Search(query, &meilisearch.SearchRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Results: %#+v\n", results)
}
