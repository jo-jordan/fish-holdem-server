package inbound

import "encoding/json"

func UnmarshalBaseInbound(data []byte) (BaseInbound, error) {
	var r BaseInbound
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *BaseInbound) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type BaseInbound struct {
	ReqType string      `json:"req_type"`
	Data    interface{} `json:"data"`
}
