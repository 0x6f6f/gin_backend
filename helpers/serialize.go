package helpers

import "encoding/json"

// 序列化为字符串
func SerializeUintListToStr(arr []uint) (string, error) {
	data, err := json.Marshal(arr)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 反序列化为数组
func DeserializeStrToUintList(serializedStr string) ([]uint, error) {
	var arr []uint
	err := json.Unmarshal([]byte(serializedStr), &arr)
	if err != nil {
		return nil, err
	}
	return arr, nil
}