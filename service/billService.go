package service

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

//保存票据
//return string ：为链码层返回的context
func (setup *FabricSetupService) SaveBill(bill Bill) (string, error) {
	//优化
	var args []string
	args = append(args, "issue")
	//将票据对象转换为字节数组 应用程序之间的数据交换是通过字节来实现的（数据传输）
	b, err := json.Marshal(bill)
	if err != nil {
		return "", fmt.Errorf("指定的票据对象序列化发生错误")
	}
	//设置调用链码执行交易请求参数
	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{b},
	} //结构体大秀
	//调用链码执行交易
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf("发布票据失败： %v", err)
	}
	return response.TransactionID.ID, nil

}

//当前持有人的票据列表
func (setup *FabricSetupService) FindBills(holdCmID string) ([]byte, error) {
	var args []string
	args = append(args, "queryMyBills")
	args = append(args, holdCmID)
	//设置请求参数
	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1])},
	}
	//调用链码
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x000}, fmt.Errorf(err.Error())
	}
	//payload:查询链码层后的信息
	b := response.Payload

	return b[:], nil

}

//发起背书请求
//args:billNO waitEndorseCmID waitendorseAccet
func (setup *FabricSetupService) Endorse(billNo string, WaitEndorseCmID string, WaitEndorseAcct string) (string, error) {
	var args []string
	args = append(args, "endorse")
	args = append(args, billNo)
	args = append(args, WaitEndorseCmID)
	args = append(args, WaitEndorseAcct)
	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
	}
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return string(response.Payload), nil

}

//背书列表
func (setup *FabricSetupService) FindWaitBills(WaitEndorseCmID string) ([]byte, error) {
	var args []string
	args = append(args, "queryMyWaitBills")
	args = append(args, WaitEndorseCmID)

	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1])},
	}
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x00}, fmt.Errorf(err.Error())
	}
	b := response.Payload
	return b[:], nil
}

//签收票据 args:billNo waitEndorseID waitEndorseAccet
func (setup *FabricSetupService) Accept(billNo string, WaitEndorseCmID string, WaitEndorseAcct string) (string, error) {
	var args []string
	args = append(args, "accept")
	args = append(args, billNo)
	args = append(args, WaitEndorseCmID)
	args = append(args, WaitEndorseAcct)

	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
	}
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return string(response.Payload), nil
}

//根据票据no查询详情
func (setup *FabricSetupService) FindBillByNo(billNo string) ([]byte, error) {
	var args []string
	args = append(args, "queryBillByNo")
	args = append(args, billNo)
	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1])},
	}
	response, err := setup.Fabric.Client.Query(req)
	if err != nil {
		return []byte{0x00}, fmt.Errorf("!!!!!!!查询指定票据信息失败：=>"+err.Error())
	}
	return response.Payload, nil
}
//票据拒签 args:billNo waitEndorseID waitEndorseAccet
func (setup *FabricSetupService) Reject(billNo string, WaitEndorseCmID string, WaitEndorseAcct string) (string, error)  {
	var args []string
	args = append(args, "reject")
	args = append(args, billNo)
	args = append(args, WaitEndorseCmID)
	args = append(args, WaitEndorseAcct)
	req := chclient.Request{
		ChaincodeID: setup.Fabric.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])},
	}
	response, err := setup.Fabric.Client.Execute(req)
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return string(response.Payload), nil
}

