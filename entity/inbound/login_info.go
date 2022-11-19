package inbound

import "encoding/json"

func UnmarshalLoginInfo(data []byte) (LoginInfo, error) {
	var r LoginInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginInfo struct {
	Username string `json:"username"`
}
