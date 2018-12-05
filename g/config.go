package g

import (
	"metric-transfer/g/log"

	"encoding/json"
	"flag"
)

var (
	//kafka
	kafkaBrokers = flag.String("kafka-brokers", "http://127.0.0.1:9092", "kafka broker list separated by ','")
	kafkaTopic   = flag.String("kafka-topic", "datapoints", "kafka topic")
	groupID      = flag.String("kafka-group", "metric-transfer", "Consumer group ID")
	user         = flag.String("kafka-user", "alice", "Kafka username")
	password     = flag.String("kafka-password", "alice", "Kafka password")
	enableAuth   = flag.Bool("kafka-sasl-auth", false, "Kafka enable SASL auth")
	consumerNum  = flag.Int("thread-num", 4, "Kafka consume thread num")

	//http filter
	networkFilterEnable = flag.Bool("http-filter-enable", false, "enable filter http metrics")
	networkMetrics      = flag.String("http-metrics", "test", "filter metrics and tags")

	//http sender
	httpSenderEnable   = flag.Bool("http-sender-enable", false, "enable send http metrics")
	httpWorkSenderAddr      = flag.String("http-sender-api", "test", "address of http sender send to")
	httpSenderInterval = flag.Int("http-sender-interval", 30, "http sender message's interval")
)

type KafkaConfig struct {
	KafkaBrokers string
	KafkaTopic   string
	GroupID      string
	User         string
	Password     string
	EnableAuth   bool
	ConsumerNum  int
}

type FilterConfig struct {
	NetworkFilterEnable bool
	NetworkFilterConfig []MetricFilterConfig
}

type SenderConfig struct {
	HttpSenderEnable bool
	HttpSenderConfig
}

type MetricFilterConfig struct {
	Metric     string
	Dimensions map[string]string
}

type HttpSenderConfig struct {
	HttpSenderAddr string
	Interval         int
}

var (
	kafkaConfig  *KafkaConfig
	filterConfig *FilterConfig
	senderConfig *SenderConfig
)

func GetKafKaConfig() *KafkaConfig {
	if kafkaConfig == nil {
		initKafkaConfig()
	}
	return kafkaConfig
}

func GetFilterConfig() *FilterConfig {
	if filterConfig == nil {
		initFilterConfig()
	}
	return filterConfig
}

func GetSenderConfig() *SenderConfig {
	if senderConfig == nil {
		initSenderConfig()
	}
	return senderConfig
}

func initKafkaConfig() {
	kafkaConfig = &KafkaConfig{
		KafkaBrokers: *kafkaBrokers,
		KafkaTopic:   *kafkaTopic,
		GroupID:      *groupID,
		User:         *user,
		Password:     *password,
		EnableAuth:   *enableAuth,
		ConsumerNum:  *consumerNum,
	}
}

func initSenderConfig() {
	senderConfig = new(SenderConfig)
	httpSenderConfig := HttpSenderConfig{
		HttpSenderAddr: *httpWorkSenderAddr,
		Interval:         *httpSenderInterval,
	}

	senderConfig.HttpSenderEnable = *httpSenderEnable
	senderConfig.HttpSenderConfig = httpSenderConfig
}

type Dimension struct {
	K string
	V string
}

type RawMetricFilter struct {
	Metric     string
	Dimensions []Dimension
}

func initFilterConfig() {
	filterConfig = new(FilterConfig)

	var rawFilters []RawMetricFilter
	if err := json.Unmarshal([]byte(*networkMetrics), &rawFilters); err != nil {
		log.Fatalln("Error: ", err)
	}
	log.Infof("rawFilters %+v", rawFilters)
	var metricFilters []MetricFilterConfig
	for _, filter := range rawFilters {
		var dims = make(map[string]string)
		for _, dim := range filter.Dimensions {
			dims[dim.K] = dim.V
		}
		metricFilters = append(metricFilters, MetricFilterConfig{filter.Metric, dims})
	}

	filterConfig.NetworkFilterEnable = *networkFilterEnable
	filterConfig.NetworkFilterConfig = metricFilters
}
