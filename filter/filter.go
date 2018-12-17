package filter

import (
	"fmt"
	"metric-transfer/filter/network"
	"metric-transfer/g"
	"metric-transfer/g/log"
	"metric-transfer/models"
	"strings"
	"sync"

	cluster "github.com/bsm/sarama-cluster"
)

const messagesize int = 100

type Filter interface {
	FilterMetric(metrics []models.DataPoint, metricCh chan g.TransMessage)
}

type FilterManager struct {
	Filters        []Filter
	KafkaConsumers []*cluster.Consumer
	ConsumeMsg     chan []models.DataPoint
	TransCh        chan g.TransMessage
}

func NewFilterManager(metricCh chan g.TransMessage) *FilterManager {
	filterConfig := g.GetFilterConfig()
	fm := new(FilterManager)
	if filterConfig.MetricAndTagsFilterEnabled {
		filter := network.NewMetricAndTagsFilter(filterConfig.MetricAndTagsFilterConfig)
		fm.Filters = append(fm.Filters, filter)
	}
	kafkaConfig := g.GetKafKaConfig()
	fm.KafkaConsumers = newConsumer(kafkaConfig)
	fm.TransCh = metricCh
	fm.ConsumeMsg = make(chan []models.DataPoint, messagesize)

	log.Infof("Init filterManager is finished: %v", *fm)
	return fm
}

func (fm *FilterManager) Filter() {
	go fm.consume()
	for {
		select {
		case datapoints := <-fm.ConsumeMsg:
			for _, filter := range fm.Filters {
				go filter.FilterMetric(datapoints, fm.TransCh)
			}
		}
	}
}

func (fm *FilterManager) close() {
}

func (fm *FilterManager) consume() {
	var wg sync.WaitGroup
	for i, consumer := range fm.KafkaConsumers {
		wg.Add(1)
		consumerClone := consumer
		j := i
		go func() {
			defer wg.Done()
			log.Info("Start kafkaConsumer ", j)
			for {
				select {
				case err := <-consumerClone.Errors():
					log.Errorf("Consumer %d consumme message error: %v", j, err)
				case note := <-consumerClone.Notifications():
					log.Infof("Consummer %d rebalanced: %+v", j, note)
				case msg := <-consumerClone.Messages():
					datapoints, err := parseMessage(msg.Value)
					if err != nil {
						log.Info("Filter parse kafka message error: ", err)
						continue
					}
					//log.Infof("Receive datapoints: %v", datapoints)
					fm.ConsumeMsg <- datapoints
				}
			}
		}()
	}
	wg.Wait()
	log.Info("FilterManager exit all consumer")
}

func newConsumer(config *g.KafkaConfig) []*cluster.Consumer {
	consumers := make([]*cluster.Consumer, config.ConsumerNum)
	for i := 0; i < config.ConsumerNum; i++ {
		kafkaConfig := cluster.NewConfig()
		kafkaConfig.ClientID = fmt.Sprintf("metric_consumer%02d", i)
		kafkaConfig.Consumer.Return.Errors = true
		kafkaConfig.Group.Return.Notifications = true
		kafkaConfig.Net.SASL.Enable = config.EnableAuth
		kafkaConfig.Net.SASL.User = config.User
		kafkaConfig.Net.SASL.Password = config.Password
		groupID := config.GroupID
		brokers := strings.Split(config.KafkaBrokers, ",")
		topic := config.KafkaTopic

		consumer, err := cluster.NewConsumer(brokers, groupID, []string{topic}, kafkaConfig)
		if err != nil {
			log.Panic("kafka consumer create failed: ", err)
		}
		consumers[i] = consumer
	}

	return consumers
}
