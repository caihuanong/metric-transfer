package main

import (
	"flag"
	"metric-transfer/filter"
	"metric-transfer/g"
	"metric-transfer/g/log"
	"metric-transfer/sender"

	"github.com/vharitonsky/iniflags"
)

var (
	sendBuffSize = flag.Int("send-buffer-size", 1000, "send buffer size")
)

func main() {
	log.Init()
	iniflags.Parse()

	// start senders
	metricCh := make(chan g.TransMessage, *sendBuffSize)
	sendManager := sender.NewSenderManager(metricCh)
	// start dispatch metric
	go sendManager.DispatchMetric()

	// start monitors
	filterManager := filter.NewFilterManager(metricCh)
	filterManager.Filter()

	done := make(chan struct{})
	<-done // block forever
}
