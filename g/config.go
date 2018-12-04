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

	//network filter
	networkFilterEnable = flag.Bool("network-filter-enable", false, "enable filter network metrics")
	networkMetrics      = flag.String("network-metrics", "test", "filter metrics and tags")

	//network sender
	networkSenderEnable   = flag.Bool("network-sender-enable", false, "enable send network metrics")
	netWorkSenderApi      = flag.String("network-sender-api", "test", "address of network sender send to")
	networkSenderInterval = flag.Int("network-sender-interval", 30, "network sender message's interval")
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
	NetworkSenderEnable bool
	NetworkSenderConfig
}

type MetricFilterConfig struct {
	Metric     string
	Dimensions map[string]string
}

type NetworkSenderConfig struct {
	NetWorkSenderApi string
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
	networkSenderConfig := NetworkSenderConfig{
		NetWorkSenderApi: *netWorkSenderApi,
		Interval:         *networkSenderInterval,
	}

	senderConfig.NetworkSenderEnable = *networkSenderEnable
	senderConfig.NetworkSenderConfig = networkSenderConfig
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
