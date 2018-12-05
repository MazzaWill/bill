package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)
/*保存票据
*args:bill
*/
func PutBill(stub shim.ChaincodeStubInterface, bill Bill) ([]byte, bool)  {
	//将票据对象序列化
	b, err :=json.Marshal(bill)
	if err != nil {
		return nil, false
	}
	err = stub.PutState(Bill_Prefix + bill.BillInfoID, b)
	if err != nil {
		return nil, false
	}
	return b, true

}
/*
去重（GetBill）
*/
func GetBill(stub shim.ChaincodeStubInterface, billNo string) (Bill,bool)  {
	var bill Bill
	//根据票据号码查询票据状态
	b, err := stub.GetState(Bill_Prefix + billNo)
	if err != nil{
		return bill,false
	}
	if b == nil{
		return bill, false
	}
	//根据擦寻的票据状态进行反序列化
	err = json.Unmarshal(b, &bill)
	if err != nil{
		return bill, false
	}
	//返回结果
	return bill, true
}
//发布票据 args：billObject
func (t *BillChainCode) issue(stub shim.ChaincodeStubInterface, args []string)peer.Response  {

	if len(args) != 1 {
		msg, _ := GetMsgString(1,"接受的票据参数必须位票据对象")
		return shim.Error(msg)
	}




	var bill Bill
	err := json.Unmarshal([]byte(args[0]),&bill)
	if err != nil{
		msg, _ := GetMsgString(1,"反序列化票据对象是发生错误")
		return shim.Error(msg)
	}
	//票据号码必须唯一
	//查重(根据票据号码进行查询)
	_, exist := GetBill(stub, bill.BillInfoID)
	if exist {
		msg, _ := GetMsgString(1,"要发布的票据号码重复")
		return shim.Error(msg)
	}
	//更改当前发布票据的状态
	bill.State = BillInfo_State_NewPublish
	//保存
	_, bl := PutBill(stub, bill)
	if !bl {
		msg, _ := GetMsgString(1,"保存票据信息发生错误")
		return shim.Error(msg)
	}
	//创建复合键（持票人+票据号码），用于后期批量查询
	/**
	qw = holderCmID+billnoaaa+boc10001
	qw1= aaa+boc10002
	*/
	holderCmIDBillInfoNoIndexKey, err :=stub.CreateCompositeKey(IndexName, []string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil {
		msg, _ := GetMsgString(1,"保存票据状态后，创建对应的复合键发生错误")//用户无法根据票据号码得到票据对象,重新保存
		return shim.Error(msg)
	}
	//如果保存复合key时指定的value为nil，会导致后期查询下到相应的信息
	err = stub.PutState(holderCmIDBillInfoNoIndexKey, []byte{0x00})
	if err != nil {
		msg, _ := GetMsgString(1,"保存复合键发生错误")
		return shim.Error(msg)
	}
	msg, _ := GetMsgByte(0,"票据发布成功")
	return shim.Success(msg)
}


//指定批量查询持票人的票据列表queryMyBills
//args:HoldCmID

func (t *BillChainCode) queryMyBills(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if len(args) != 1 {
		msg, _ := GetMsgString(1,"请指定持票人的证件号码")
		return shim.Error(msg)
	}
	//根据制定的持票人的证件号码批量查询
	//根据持票人的证件号码创建的复合键中查询所有票据号码
	/*No1 查询相应的复合键*/
	billIterator, err := stub.GetStateByPartialCompositeKey(IndexName, []string{args[0]})
	if err != nil {
		msg, _ := GetMsgString(1,"查询列表失败，根据持票人证件号码查询所持有的票据号码时发生错误")
		return shim.Error(msg)

	}
	defer billIterator.Close()
	var bills []Bill
	//迭代处理
	for billIterator.HasNext()  {
		//k:bill.HoldrCmID, bill.BIllInfoID v:[]byte{0x00}
		kv, err := billIterator.Next()
		if err != nil {
			return shim.Error("查询票据失败，获取复合键时发生错误")
		}
		//k:[bill.HoldrCmID, bill.BIllInfoID]
		_, composites, err := stub.SplitCompositeKey(kv.Key)
		if err != nil{
			return shim.Error("查询票据失败，分割复合键是发生错误")
		}
		//从分割后的复合键中获取对应的票据id后，然后查询相应的票据信息
		bill, bl :=GetBill(stub, composites[1])
		if !bl {
			return shim.Error("根据获取到的票据号码查询相应的票据状态是发生错误")
		}
		//判断待背书人id
		if bill.WaitEndorseCmID == args[0] {
			continue
		}

		bills = append(bills, bill)
	}

	b, err := json.Marshal(bills)
	if err != nil {
		return shim.Error("查询票据失败，序列化票据信息发生错误")
	}
	return shim.Success(b)
}
//通过票据号查询 args:billno

func (t *BillChainCode) queryBillByNo(stub shim.ChaincodeStubInterface, args []string) peer.Response  {
	if  len(args) != 1{
		return shim.Error("请输入指定的票据号码")
	}
	//查询
	bill, bl := GetBill(stub, args[0])
	if !bl {
		return shim.Error("根据指定的票据号码查询票据失败")

	}
	//获取历史更变数据
	iterator, err := stub.GetHistoryForKey(Bill_Prefix + bill.BillInfoID)
	if err != nil {
		return shim.Error("根据指的票据号码查询对应的历史变更记录失败")

	}
	defer iterator.Close()
	//迭代处理
	var historys []HistoryItem

	var hisBill Bill
	for iterator.HasNext()  {
		hisData, err := iterator.Next()
		if err != nil {
			return shim.Error("获取票据历史变更数据失败")
		}

		var historyItem HistoryItem

		historyItem.TxId = hisData.TxId

		json.Unmarshal(hisData.Value, &hisBill)

		if hisData.Value == nil {
			var empty Bill
			historyItem.Bill = empty
		}else {
			historyItem.Bill = hisBill
		}
		historys = append(historys, historyItem)
	}
	bill.Historys = historys
	//返回 反序列化
	b, err := json.Marshal(bill)
	if err != nil {
		return shim.Error("获取票据状态记背书历史失败，序列化票据时发生错误")
	}
	return shim.Success(b)
}

