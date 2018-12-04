package models

type DataPoint struct {
	MetricName string            `json:"name" binding:"required"`
	Dimensions map[string]string `json:"dimensions" binding:"required"`
	Timestamp  int64             `json:"timestamp" binding:"required"`
	Value      float64           `json:"value" binding:required"`
}
