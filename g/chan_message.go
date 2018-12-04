package g

import (
	"metric-transfer/models"
)

type TransMessage struct {
	To      map[string]struct{} //to senders
	Metrics []models.DataPoint  //metrics
}
