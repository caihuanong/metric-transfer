package sender

import (
	"metric-transfer/g"
	"metric-transfer/g/log"
	"metric-transfer/models"
	"metric-transfer/sender/http"
)

const (
	HTTP = "http"
)

type Sender interface {
	Send(metrics []models.DataPoint)
	Run()
}

type SenderManager struct {
	Senders  map[string]Sender
	MetricCh chan g.TransMessage
}

func NewSenderManager(metricCh chan g.TransMessage) *SenderManager {
	s := &SenderManager{
		MetricCh: metricCh,
		Senders:  map[string]Sender{},
	}
	s.MetricCh = metricCh
	config := g.GetSenderConfig()
	if config.HttpSenderEnable {
		if sender, err := http.NewHttpSender(config.HttpSenderConfig); err == nil {
			s.Senders[HTTP] = sender
		}
	}
	return s
}

func (s *SenderManager) DispatchMetric() {
	for _, sender := range s.Senders {
		go sender.Run()
	}
	for msg := range s.MetricCh {
		log.Info("SenderManager receive metrics: ", msg)
		if len(s.Senders) == 0 {
			log.Warning("SenderManager have no sender client, can't send data")
		}

		for senderType, sender := range s.Senders {
			datapoints := msg.Metrics
			if _, ok := msg.To[senderType]; ok {
				sender.Send(datapoints)
			}
		}
	}
}
