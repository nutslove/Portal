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
	IndexName := "career"
	mapping := strings.NewReader(`{
		"settings": {
			"index": {
						"number_of_shards": 1,
						"number_of_replicas": 1
						}
			},
			"mappings": {
					"properties": {
							"title": { "type": "text" },
							"author": { "type": "keyword" },
							"post": { "type": "text" }
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
	// document := struct {
	// 	Title  string `json:"title"`
	// 	Author string `json:"author"`
	// 	Post   string `json:"post"`
	// }{
	// 	Title:  "What a day!",
	// 	Author: "Lee Joonki",
	// 	Post:   "Today is very good!!! I am so happy~",
	// }
	type Posts struct {
		Title  string `json:"title"`
		Author string `json:"author"`
		Post   string `json:"post"`
	}
	post := Posts{
		Title:  "Are you happy?",
		Author: "Sena",
		Post:   "Yes! I am very happy with you!",
	}

	docId := "38"
	req := opensearchapi.IndexReq{
		Index:      IndexName,
		DocumentID: docId, // DocumentIDは内部的に `_id` フィールドとして扱われ、通常 keywordタイプと同様に動作する
		Body:       opensearchutil.NewJSONReader(&post),
		Params: opensearchapi.IndexParams{
			Refresh: "true", // ドキュメントをインデックスに追加した直後に検索可能にする（デフォルトは非同期更新）
		},
	}
	insertResp, err := client.Index(ctx, req)
	if err != nil {
		log.Fatal("failed to insert document ", err)
	}
	fmt.Printf("Created document in %s\n  ID: %s\n", insertResp.Index, insertResp.ID)

	// Search for the document.
	///////////////////////////////////////////////////////////////////////////
	//// size: 返す検索結果の最大数
	//// multi_match: 複数のフィールドで同じクエリを検索するためのクエリタイプ（単一フィールドでの検索はmatch）
	//// query: 検索するテキストを指定
	//// fields: 検索対象のフィールドを指定
	//// ^2: ブースト値。"title"フィールドのスコアを2倍にする
	///////////////////////////////////////////////////////////////////////////
	// content := strings.NewReader(`{
	// 	"size": 5,
	// 	"query": {
	// 			"multi_match": {
	// 					"query": "miller",
	// 					"fields": ["title^2", "director"]
	// 			}
	// 	}
	// }`)

	///////////////////////////////////////////////////////////////////////
	//// ids: DocumentID(内部的には`_id`フィールドとして扱われる)で検索
	//// values: 検索したいDocumentIDを指定
	///////////////////////////////////////////////////////////////////////
	content := strings.NewReader(`{
    "query": {
        "ids": {
            "values": ["37"]
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
			err := json.Unmarshal(hit.Source, &post) // post変数に検索内容で上書き
			if err != nil {
				log.Printf("Error unmarshaling document: %v", err)
				continue
			}

			fmt.Printf("Document ID: %s\n", hit.ID)
			fmt.Printf("  Title: %s\n", post.Title)
			fmt.Printf("  Author: %s\n", post.Author)
			fmt.Printf("  Post: %s\n", post.Post)

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
