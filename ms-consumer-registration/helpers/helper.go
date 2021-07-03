package helpers

import "encoding/json"

func ToByte(value interface{}) []byte {
	val, err := json.Marshal(value)
	if err != nil{
		return val
	}

	return val
}

func ToInterface(value []byte) (val interface{}){
	err := json.Unmarshal(value, &val)
	if err != nil {
		return val
	}

	return val
}
