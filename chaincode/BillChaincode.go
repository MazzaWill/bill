package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric/protos/peer"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type BillChainCode struct {

}

func (t *BillChainCode) Init(stub shim.ChaincodeStubInterface) peer.Response  {
	return shim.Success(nil)
}
//获取用户意图
func (t *BillChainCode) Invoke(stub shim.ChaincodeStubInterface) peer.Response  {
	fun, args := stub.GetFunctionAndParameters()
		//发布票据
	if fun == "issue" {
		return t.issue(stub, args)
		//查询所在用户票据列表
	}else if fun == "queryMyBills" {

		return t.queryMyBills(stub, args)
		//通过票据号查询
	}else if fun== "queryBillByNo" {

		return t.queryBillByNo(stub, args)
		//查询待签收票据
	}else if fun == "queryMyWaitBills" {

		return t.queryMyWaitBills(stub, args)
		//发起背书
	}else if fun == "endorse" {

		return t.endorse(stub, args)
		//背书接受
	}else if fun == "accept" {

		return t.accept(stub, args)
		//拒绝背书
	}else if fun== "reject" {

		return t.reject(stub, args)
	}
		//return shim.Error("指定的函数名称错误")
	respMsg, _ := GetMsgString(1,"指定的函数名称错误")
	return shim.Error(respMsg)

}

func main() {
	err := shim.Start(new(BillChainCode))
	if err != nil {
		fmt.Printf("启动链码错误: %v", err)
	}
}
