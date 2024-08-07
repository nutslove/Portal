package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"portal/config"
	"portal/controllers"
	"portal/routers"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
	"github.com/opensearch-project/opensearch-go/v4/opensearchutil"
)

func main() {
	router := gin.Default()

	config.OtelInitialSetting()

	/////////////////////
	// OpenSearch Test //
	/////////////////////
	client, err := config.OpensearchNewClient()
	if err != nil {
		log.Fatal("cannot initialize", err)
	}

	ctx := context.Background()
	infoResp, err := client.Info(ctx, nil)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	fmt.Printf("Cluster INFO:\n  Cluster Name: %s\n  Cluster UUID: %s\n  Version Number: %s\n", infoResp.ClusterName, infoResp.ClusterUUID, infoResp.Version.Number)

	// Define index settings.
	IndexName := "go-test-index1"
	mapping := strings.NewReader(`{
		"settings": {
			"index": {
						"number_of_shards": 1,
						"number_of_replicas": 2
						}
					}
		}`)

	// Create an index with non-default settings.
	createIndex := opensearchapi.IndicesCreateReq{
		Index: IndexName,
		Body:  mapping,
	}
	createIndexResponse, err := client.Indices.Create(ctx, createIndex)

	var opensearchError *opensearch.StructError

	// Load err into opensearch.Error to access the fields and tolerate if the index already exists
	if err != nil {
		if errors.As(err, &opensearchError) {
			if opensearchError.Err.Type != "resource_already_exists_exception" {
				return
			}
		} else {
			return
		}
	}

	fmt.Printf("Created Index: %s\n  Shards Acknowledged: %t\n", createIndexResponse.Index, createIndexResponse.ShardsAcknowledged)

	// When using a structure, the conversion process to io.Reader can be omitted using utility functions.
	// Add a document to the index.
	document := struct {
		Title    string `json:"title"`
		Director string `json:"director"`
		Year     string `json:"year"`
	}{
		Title:    "Moneyball",
		Director: "Bennett Miller",
		Year:     "2011",
	}

	docId := "1"
	req := opensearchapi.IndexReq{
		Index:      IndexName,
		DocumentID: docId,
		Body:       opensearchutil.NewJSONReader(&document),
		Params: opensearchapi.IndexParams{
			Refresh: "true",
		},
	}
	insertResp, err := client.Index(ctx, req)
	if err != nil {
		log.Fatal("failed to insert document ", err)
	}
	fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

	// Search for the document.
	content := strings.NewReader(`{
		"size": 5,
		"query": {
				"multi_match": {
						"query": "miller",
						"fields": ["title^2", "director"]
				}
		}
	}`)

	searchResp, err := client.Search(
		ctx,
		&opensearchapi.SearchReq{
			Body: content,
		},
	)
	if err != nil {
		log.Fatal("failed to search document ", err)
	}
	fmt.Printf("Search hits: %v\n", searchResp.Hits.Total.Value)

	if searchResp.Hits.Total.Value > 0 {
		indices := make([]string, 0)
		for _, hit := range searchResp.Hits.Hits {
			var doc struct {
				Title    string `json:"title"`
				Director string `json:"director"`
				Year     string `json:"year"`
			}
			err := json.Unmarshal(hit.Source, &doc)
			if err != nil {
				log.Printf("Error unmarshaling document: %v", err)
				continue
			}

			fmt.Printf("Document ID: %s\n", hit.ID)
			fmt.Printf("  Title: %s\n", doc.Title)
			fmt.Printf("  Director: %s\n", doc.Director)
			fmt.Printf("  Year: %s\n", doc.Year)

			add := true
			for _, index := range indices {
				if index == hit.Index {
					add = false
				}
			}
			if add {
				indices = append(indices, hit.Index)
			}
		}
		fmt.Printf("Search indices: %s\n", strings.Join(indices, ","))
	}
	/////////////////////
	// OpenSearch Test //
	/////////////////////

	// 静的ファイルの設定
	router.Static("/static", "./static")

	// テンプレート内で使うカスタム関数を登録
	controllers.RegisterCustomFunction(router)

	// テンプレートファイルの読み込み
	router.LoadHTMLGlob("templates/*")

	// ルーティングの設定
	routers.SetupRouter(router)

	// サーバーの起動
	router.Run(":8080")
}
