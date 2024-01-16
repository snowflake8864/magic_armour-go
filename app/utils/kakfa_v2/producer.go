package kakfa_v2

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/yml_config"
	"goskeleton/app/utils/yml_config/ymlconfig_interf"
	"log"
	"os"
	"strings"
	"time"
)

var kafkaProducerPool *ProducerPool
var configYml ymlconfig_interf.YmlConfigInterf

// ProducerPool 生产者连接池
type ProducerPool struct {
	bootstrapServers string
	pool             chan *kafka.Writer
}

// 处于程序底层的包，init 初始化的代码段的执行会优先于上层代码，因此这里读取配置项不能使用全局配置项变量
func init() {
	fmt.Println("kakfa_v2 init---------------------------------------------------------")
	configYml = yml_config.CreateYamlFactory("kafka_v2")
	kafkaProducerPool = NewProducerPool()
}

// NewProducerPool 创建生产者连接池
func NewProducerPool() *ProducerPool {
	fmt.Println("NewProducerPool---------------------------------------------------------")
	addr := configYml.GetString("KafkaV2.Producer.Addr")
	brokers := strings.Split(addr, ",")
	poolSize := configYml.GetInt("KafkaV2.Producer.PoolSize")
	writeTimeout := configYml.GetDuration("KafkaV2.Producer.WriteTimeout")
	batchSize := configYml.GetInt("KafkaV2.Producer.BatchSize")
	batchTimeout := configYml.GetDuration("KafkaV2.Producer.BatchTimeout")
	requiredAcks := configYml.GetInt("KafkaV2.Producer.RequiredAcks")

	pool := make(chan *kafka.Writer, poolSize)
	for i := 0; i < poolSize; i++ {
		writer := &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: writeTimeout * time.Second,
			BatchSize:    batchSize,
			BatchTimeout: batchTimeout * time.Second,
			RequiredAcks: kafka.RequiredAcks(requiredAcks),
			Async:        false,
			Logger:       log.New(os.Stdout, "kafka-go: ", log.LstdFlags),
		}
		pool <- writer
	}
	return &ProducerPool{
		bootstrapServers: addr,
		pool:             pool,
	}
}

// Get 从连接池中获取生产者
func (p *ProducerPool) Get() (*kafka.Writer, error) {
	select {
	case producer := <-p.pool:
		return producer, nil
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("timeout waiting for producer")
	}
}

// Release 释放生产者回连接池中
func (p *ProducerPool) Release(producer *kafka.Writer) {
	p.pool <- producer
}

// 从连接池获取一个Kafka生产者
func GetOneKafkaProducer() *KafkaProducer {
	maxRetryTimes := 3
	var w *kafka.Writer
	var err error
	for i := 1; i <= maxRetryTimes; i++ {
		w, err = kafkaProducerPool.Get()
		// 首先通过执行一个获取时间的命令检测连接是否有效，如果已有的连接无法执行命令，则重新尝试连接到Kafka服务器获取的连接池地址
		if err != nil {
			// 如果已有的连接池获取连接出错(官方库的说法是连接不可用)，那么继续使用从新初始化连接池
			NewProducerPool()
			w, _ = kafkaProducerPool.Get()
		}

		if err != nil {
			//variable.ZapLog.Error("Redis：网络中断,开始重连进行中..." , zap.Error(oneConn.Err()))
			if i == maxRetryTimes {
				variable.ZapLog.Error(my_errors.ErrorsRedisGetConnFail, zap.Error(err))
				return nil
			}
			//如果出现网络短暂的抖动，短暂休眠后，支持自动重连
			time.Sleep(time.Second * 3)
		} else {
			break
		}
	}
	return &KafkaProducer{Producer: *w}
}
func sendMessage(pool *ProducerPool, message string) error {
	// 在连接池中获取生产者
	producer, err := pool.Get()
	if err != nil {
		return fmt.Errorf("failed to get producer: %w", err)
	}
	defer pool.Release(producer)

	// 构建消息
	kMessage := kafka.Message{
		Key:   []byte("key"),
		Value: []byte(message),
	}

	// 发送消息
	err = producer.WriteMessages(context.Background(), kMessage)
	if err != nil {
		return fmt.Errorf("failed to write messages: %w", err)
	}

	log.Printf("Message sent successfully: %s", message)

	return nil
}

// 定义一个生产者构体
type KafkaProducer struct {
	Producer kafka.Writer
}

func (w *KafkaProducer) produce(message kafka.Message) error {
	defer kafkaProducerPool.Release(&w.Producer)

	err := w.Producer.WriteMessages(context.Background(), message)
	if err != nil {
		variable.ZapLog.Error("Failed to produce message:" + err.Error())
		return err
	}
	return nil
}

type ProducerMessage struct {
	Type string
	Data string
}

func (w *KafkaProducer) PushCommon(message ProducerMessage) bool {
	content, err := json.Marshal(message)
	if err != nil {
		variable.ZapLog.Error(fmt.Sprintf("PushCommon json marshal err, message:%v,err:%v", message, err))
	}
	topic := configYml.GetString("KafkaV2.TopicList.TopicComm.Name")

	msg := kafka.Message{
		Topic: topic,
		Value: content,
	}
	result := w.produce(msg)
	if result != nil {
		variable.ZapLog.Error(fmt.Sprintf("PushCommon produce, ,err:%v", result.Error()))
		return false
	}
	return true
}
