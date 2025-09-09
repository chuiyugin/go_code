package job

import (
	"context"
	"encoding/json"
	"errors"
	"review-job/internal/conf"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	kafka "github.com/segmentio/kafka-go"
)

// 评价数据流处理

// JobWorker 自定义执行job的结构体，实现 transport.Server
type JobWorker struct {
	kafkaReader *kafka.Reader
	esClient    *ESClient
	log         *log.Helper
}

type ESClient struct {
	client *elasticsearch.TypedClient
	Index  string
}

func NewJobWorker(kafkaReader *kafka.Reader, esClient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafkaReader,
		esClient:    esClient,
		log:         log.NewHelper(logger),
	}
}

func NewKafkaReader(cfg *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupId,
		Topic:    cfg.Topic,
		MaxBytes: 10e6, // 10MB
	})
}

func NewESClient(cfg *conf.Elasticsearch) (*ESClient, error) {
	// 连接es
	c := elasticsearch.Config{
		Addresses: cfg.Addresses,
	}
	esclient, err := elasticsearch.NewTypedClient(c)
	if err != nil {
		return nil, err
	}
	return &ESClient{
		client: esclient,
		Index:  cfg.Index,
	}, nil
}

// Msg 定义kafka中接收到的数据
type Msg struct {
	Type     string                   `json:"type"`
	Database string                   `json:"database"`
	Table    string                   `json:"table"`
	IsDdl    bool                     `json:"isDdl"`
	Data     []map[string]interface{} `json:"data"`
}

// Start kratos程序启动之后会调用的方法
// ctx 是kratos框架启动的时候传入ctx，是带有退出取消的
func (jw JobWorker) Start(ctx context.Context) error {
	jw.log.Debugf("JobWorker start....")
	// 1、从kafka中获取mysql中的数据变更消息
	// 接收消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			jw.log.Errorf("readMessage from kafkaa failed, err:%v\n", err)
			break
		}
		jw.log.Debugf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 2、将完整的评价数据写入ES
		msg := new(Msg)
		err = json.Unmarshal(m.Value, msg)
		if err != nil {
			jw.log.Errorf("Unmarshal msg from kafka failed, err:%v", err)
			continue
		}
		if msg.Type == "INSERT" {
			// 往ES中新增文档
			for idx := range msg.Data {
				jw.indexDocument(msg.Data[idx])
			}
		} else {
			// 往ES中更新文档
			for idx := range msg.Data {
				jw.updateDocument(msg.Data[idx])
			}
		}
	}
	return nil
}

// Stop kratos结束之后会调用的
func (jw JobWorker) Stop(context.Context) error {
	jw.log.Debugf("JobWorker stop....")
	// 程序退出前关闭Reader
	return jw.kafkaReader.Close()
}

// indexDocument 索引文档
func (jw JobWorker) indexDocument(data map[string]interface{}) {
	// 添加文档
	reviewID := data["review_id"].(string)
	resp, err := jw.esClient.client.Index(jw.esClient.Index).
		Id(reviewID).
		Document(data).
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("indexing document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWorker) updateDocument(data map[string]interface{}) {
	reviewID := data["review_id"].(string)
	resp, err := jw.esClient.client.Update(jw.esClient.Index, reviewID).
		Doc(data). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("update document failed, err:%v\n", err)
		return
	}
	jw.log.Debugf("result:%v\n", resp.Result)
}
