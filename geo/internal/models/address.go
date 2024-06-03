package models

import "encoding/json"

type Address struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func (a Address) MarshalBinary() ([]byte, error) {
	return json.Marshal(a)
}

func (u *Address) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
