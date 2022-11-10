package outbound

import "encoding/json"

func UnmarshalTableInfo(data []byte) (TableInfo, error) {
	var r TableInfo
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *TableInfo) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type TableInfo struct {
	ID int64 `json:"id"`
}
