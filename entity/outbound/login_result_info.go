package outbound

import "encoding/json"

func UnmarshalLoginResultInfo(data []byte) (LoginResultInfo, error) {
	var r LoginResultInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginResultInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginResultInfo struct {
	Success  bool   `json:"success"`
	Token    string `json:"token"`
	PlayerId int64  `json:"player_id"`
}
