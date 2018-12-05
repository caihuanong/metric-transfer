package network

import (
	"metric-transfer/g"
	"metric-transfer/g/log"
	"metric-transfer/models"
)

type NetWorkFilter struct {
	filter map[string](map[string]string) //a map :metric_name -> tagK:tagV, use to filter datapoints
	to     map[string]struct{}            //datapoints after filtering send to
}

func NewNetWorkFilter(config []g.MetricFilterConfig) *NetWorkFilter {
	filter := make(map[string](map[string]string))
	for _, metricConfig := range config {
		tags := metricConfig.Dimensions
		filter[metricConfig.Metric] = tags
	}
	to := make(map[string]struct{})
	to["http"] = struct{}{}

	nf := &NetWorkFilter{
		filter: filter,
		to:     to,
	}

	return nf
}

func (nf *NetWorkFilter) FilterMetric(datapoints []models.DataPoint, metricCh chan g.TransMessage) {
	var targetPoints []models.DataPoint
	for _, point := range datapoints {
		if nf.filterNameAndTags(point) {
			targetPoints = append(targetPoints, point)
		}
	}

	if targetPoints != nil && len(targetPoints) > 0 {
		msg := g.TransMessage{
			To:      nf.to,
			Metrics: targetPoints,
		}
		log.Infof("fitered points %v", targetPoints)
		metricCh <- msg
	}
}

func (nf *NetWorkFilter) filterNameAndTags(point models.DataPoint) bool {
	//filter by metricName
	if tags, ok := nf.filter[point.MetricName]; ok {
		//filter by tags
		for filterK, filterV := range tags {
			if V, ok := point.Dimensions[filterK]; ok && V == filterV {
				continue
			} else {
				return false
			}
		}
		return true
	}
	return false
}
