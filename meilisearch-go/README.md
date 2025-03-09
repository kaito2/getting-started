# Meilisearch

Make friends with Meilisearch.

[Getting started with Meilisearch Cloud — Meilisearch documentation](https://www.meilisearch.com/docs/learn/getting_started/cloud_quick_start#getting-started-with-meilisearch-cloud)

## Scripts

### Add Document

Add a document to an index.

```bash
go run scripts/add-document/main.go \
   -host=MEILISEARCH_HOST \
   -api-key=MEILISEARCH_API_KEY \
   -index=MEILISEARCH_INDEX \
   -file=./asset/add-movies.json
```

### Search

Search for documents in an index.

```bash
go run scripts/search/main.go \
   -host=MEILISEARCH_HOST \
   -api-key=MEILISEARCH_API_KEY \
   -index=MEILISEARCH_INDEX \
   -query=kaito2
```
