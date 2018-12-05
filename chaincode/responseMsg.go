package chaincode

import (
	"encoding/json"
	"fmt"
)

type ChaincodeResponseMsg struct {
	Code int    //0=>成功　１=>失败
	Dec  string //描述
}

func GetMsgByte(code int, dec string) ([]byte, error) {
	b, err := getMsg(code, dec)
	if err != nil {
		return nil, err
	}
	return b[:],nil

}

// 根据返回码和描述信息返回序列化后的字符串
func GetMsgString(code int, dec string) (string, error) {
	b, err := getMsg(code, dec)
	if err != nil {
		return "", err
	}
	return string(b[:]), nil
}

//根据返回码和描述信息进行系列化

func getMsg(code int, dec string) ([]byte, error) {
	var crm ChaincodeResponseMsg
	crm.Code = code
	crm.Dec = dec
	b, err := json.Marshal(crm)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return b, nil
}

