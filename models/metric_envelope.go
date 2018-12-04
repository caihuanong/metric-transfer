package models

import (
	"encoding/json"
)

type MetricEnvelope struct {
	DataPoints []DataPoint       `json:"datapoints" binding:"required"`
	Meta       map[string]string `json:"meta" binding:"required"`

	encoded []byte
	err     error
}

func (me *MetricEnvelope) ensureEncoded() {
	if me.encoded == nil && me.err == nil {
		me.encoded, me.err = json.Marshal(me)
	}
}

func (me *MetricEnvelope) Encode() ([]byte, error) {
	me.ensureEncoded()
	return me.encoded, me.err
}

func (me *MetricEnvelope) Length() int {
	me.ensureEncoded()
	return len(me.encoded)
}
