package filter

import (
	"encoding/json"

	"metric-transfer/g/log"
	"metric-transfer/models"
)

func parseMessage(metricMsg []byte) ([]models.DataPoint, error) {
	var metricEnvelope models.MetricEnvelope
	var err error

	if err = json.Unmarshal(metricMsg, &metricEnvelope); err != nil {
		log.Error("Error: ", err)
	}

	return metricEnvelope.DataPoints, err
}
