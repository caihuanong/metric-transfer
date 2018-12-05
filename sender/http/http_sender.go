package http

import (
	"bytes"
	"encoding/json"
	"io"
	"metric-transfer/g"
	"metric-transfer/g/log"
	"metric-transfer/models"
	"net/http"
	"sync"
	"time"
)

type HttpSender struct {
	mu          sync.Mutex
	method      string
	apiAddr     string
	ticker      *time.Ticker
	httpClient  *http.Client
	metricsBuff []models.DataPoint
}

func NewHttpSender(config g.HttpSenderConfig) (*HttpSender, error) {
	sender := new(HttpSender)
	sender.apiAddr = config.HttpSenderAddr
	sender.method = "POST"
	sender.ticker = time.NewTicker(time.Second * time.Duration(config.Interval))
	sender.httpClient = &http.Client{}
	sender.metricsBuff = []models.DataPoint{}
	log.Infof("New HttpSender finished: %+v", sender)
	return sender, nil
}

func (s *HttpSender) Send(metrics []models.DataPoint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.metricsBuff = append(s.metricsBuff, metrics...)
	log.Debug("Buff Size = ", len(s.metricsBuff))
}

func (s *HttpSender) Run() {
	for {
		<-s.ticker.C
		s.send()
	}
}

func (s *HttpSender) send() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.metricsBuff) == 0 {
		return
	}
	var rawRequest io.Reader
	jsonRequest, err := json.Marshal(s.metricsBuff)
	if err != nil {
		log.Info("HttpSender marshal metricsBuff error, ", err)
		return
	}
	rawRequest = bytes.NewReader(jsonRequest)
	// build request
	req, err := http.NewRequest(s.method, s.apiAddr, rawRequest)
	if err != nil {
		log.Info("HttpSender build request error, ", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Info("HttpSender do request error: ", resp)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Warn("HttpSender resp statusCode is not 200, code = ", resp.StatusCode)
	}
	s.metricsBuff = s.metricsBuff[:0]
}
