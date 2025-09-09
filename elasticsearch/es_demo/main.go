package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

func main() {
	// 连接es
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // es服务地址
		},
	}
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("NewTypedClient failed, err:%v\n", err)
		return
	}
	// 连接成功
	fmt.Println(client)
	// 创建index
	// CreateIndex(client)

	// 创建 document
	//indexDocument(client)

	// 查询 document
	// getDocumentByID(client, "1")

	// searchDocument2 指定条件搜索文档
	// searchDocument2(client)

	// searchDocument 搜索所有文档
	searchDocument(client)
}

func CreateIndex(client *elasticsearch.TypedClient) {
	resp, err := client.Indices.Create("my-review-1").Do(context.Background())
	if err != nil {
		fmt.Printf("Indices.Create failed, err:%v\n", err)
		return
	}
	fmt.Printf("Acknowledged:%v\n", resp.Acknowledged)
}

// indexDocument 索引文档
func indexDocument(client *elasticsearch.TypedClient) {
	// 定义 document 结构体对象
	d1 := Review{
		ID:      1,
		UserID:  147982601,
		Score:   5,
		Content: "这是一个好评！",
		Tags: []Tag{
			{1000, "好评"},
			{1100, "物超所值"},
			{9000, "有图"},
		},
		Status:      2,
		PublishTime: time.Now(),
	}

	// 添加文档
	resp, err := client.Index("my-review-1").
		Id(strconv.FormatInt(d1.ID, 10)).
		Document(d1).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

func getDocumentByID(client *elasticsearch.TypedClient, id string) {
	resp, err := client.Get("my-review-1", id).Do(context.Background())
	if err != nil {
		fmt.Printf("getDocumentByID failed, err:%v\n", err)
	}
	fmt.Printf("result:%s\n", resp.Source_)
}

// searchDocument 搜索所有文档
func searchDocument(client *elasticsearch.TypedClient) {
	// 搜索文档
	resp, err := client.Search().
		Index("my-review-1").
		Query(&types.Query{
			MatchAll: &types.MatchAllQuery{},
		}).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		return
	}
	fmt.Printf("total: %d\n", resp.Hits.Total.Value)
	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
}

// searchDocument2 指定条件搜索文档
func searchDocument2(client *elasticsearch.TypedClient) {
	// 搜索content中包含好评的文档
	resp, err := client.Search().
		Index("my-review-1").
		Query(&types.Query{
			MatchPhrase: map[string]types.MatchPhraseQuery{
				"content": {Query: "好评"},
			},
		}).
		Do(context.Background())
	if err != nil {
		fmt.Printf("search document failed, err:%v\n", err)
		return
	}
	fmt.Printf("total: %d\n", resp.Hits.Total.Value)
	// 遍历所有结果
	for _, hit := range resp.Hits.Hits {
		fmt.Printf("%s\n", hit.Source_)
	}
}
