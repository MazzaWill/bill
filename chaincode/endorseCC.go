package chaincode

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

//发起背书 args: billNO waitendorseCMID waitendorseAcct
func (t *BillChainCode) endorse(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		msg, _ := GetMsgString(1, "票据背书请求失败，票据号码，待背书人证件号码，待背书人名称")
		return shim.Error(msg)
	}
	//根据票据号码查询对应的票据状态
	bill, bl := GetBill(stub, args[0])
	if !bl {
		return shim.Error("票据背书失败，查询票据状态时发生错误")
	}
	if bill.HoldrCmID == args[1] {
		return shim.Error("票据背书请求失败，带背书人不能为当前持票人")
	}
	//历史持有人不能为背书人
	iterator, err := stub.GetHistoryForKey(args[0])
	if err != nil {
		return shim.Error("票据背书请求失败，查询历史信息发生错误")
	}
	defer iterator.Close()
	var hisBill Bill
	for iterator.HasNext() {
		hisDate, err := iterator.Next()
		if err != nil {
			return shim.Error("票据背书请求失败,获取具体历史流转信息发生错误")
		}
		var historyItem HistoryItem
		if hisDate.Value == nil {
			continue
		}

		err = json.Unmarshal(hisDate.Value, &historyItem)
		if err != nil {
			return shim.Error("反序列化历史记录发生错误")
		}

		//历史持有人不能为当前待背书人
		if hisBill.HoldrCmID == args[1] {
			return shim.Error("票据背书请求失败，/历史持有人不能为待背书人/")
		}
	}
	//更改票据状态 待背书人信息 删除拒绝背书人信息
	bill.State = BillInfo_State_EndorseWaitSign
	bill.WaitEndorseCmID = args[1]
	bill.WaitEndorseAcct = args[2]
	bill.RejectEndorseCmID = ""
	bill.RejectEndorseAcct = ""
	//保存票据
	_, bl = PutBill(stub, bill)
	if !bl {
		return shim.Error("票据背书请求失败，保存状态是发生错误")
	}
	//增加复合键，以便于批量查询(带背书人证件号码)
	waitEndorseCmIDBillInfoIDIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.WaitEndorseCmID, bill.BillInfoID})
	if err != nil {
		return shim.Error("票据背书请求失败，创建复合键时失败")
	}
	err = stub.PutState(waitEndorseCmIDBillInfoIDIndexKey, []byte{0x00})
	if err != nil {
		return shim.Error("票据背书请求失败，保存复合键时发生错误")
	}
	msg, _ := GetMsgByte(0, "票据背书请求成功")
	return shim.Success(msg)
}

//查询带背书票据列表
//args：waitEndorseCmID
func (t *BillChainCode) queryMyWaitBills(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 1 {
		return shim.Error("只能指定待背书人的证件号码")
	}
	//根据待背书人证件号码查询复合键
	iterator, err := stub.GetStateByPartialCompositeKey(IndexName, []string{args[0]})
	if err != nil {
		return shim.Error("查询待背书票据失败，查询对应的复合键发生错误")
	}

	defer iterator.Close()

	var bills []Bill //对象数组

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return shim.Error("查询待背书票据失败，获取复合键发生错误")
		}

		_, composites, err := stub.SplitCompositeKey(kv.Key)
		if err != nil {
			return shim.Error("查询待背书票据失败,分割复合键是发生错误")
		}
		bill, bl := GetBill(stub, composites[1])
		if !bl {
			return shim.Error("查询待背书票据失败,查询具体的待背书票据是发生错误")
		}
		if bill.State == BillInfo_State_EndorseWaitSign && bill.WaitEndorseCmID == args[0] {
			bills = append(bills, bill)
		}
	}
	//查询结果
	b, err := json.Marshal(bills)
	if err != nil {
		return shim.Error("查询待背书票据失败,序列化查询结果时发生错误")
	}
	return shim.Success(b)
}

/*
*背书签收
args:billNO endorseID endorseAcct
*/
func (t *BillChainCode) accept(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 3 {
		return shim.Error("票据签收失败，必须只能指定票据号码，签收人证件号码，签收人名称")
	}
	bill, bl := GetBill(stub, args[0])
	if !bl {
		return shim.Error("票据签收失败,根据票据号码查询票据信息是发生错误")
	}
	//创建复合键（根据前持票人证件号码，以及票据号码）
	HoldrCmIDBillInfoIDIndexKey, err := stub.CreateCompositeKey(IndexName, []string{bill.HoldrCmID, bill.BillInfoID})
	if err != nil {
		return shim.Error("背书签收失败，创建复合键时发生错误")
	}
	//根据复合键的key从账本中删除信息以便于前持票人无法查询到该结果 查询前持有人
	err = stub.DelState(HoldrCmIDBillInfoIDIndexKey)
	if err != nil {
		return shim.Error("背书签收失败，删除复合键时发生错误")
	}

	//更改票据信息，票据状态，当前持票人信息，待背书人信息
	bill.State = BillInfo_State_EndorseSigned
	bill.HoldrCmID = args[1]
	bill.HoldrAcct = args[2]

	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	_, bl = PutBill(stub, bill)
	if !bl {
		return shim.Error("票据背书签收失败，保存票据是发生错误")
	}
	return shim.Success([]byte("票据背书签收成功"))
}

/**
*拒签reject
*args:billNo,rejectCmID,rejectacct
*/

func (t *BillChainCode) reject(stup shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 3 {
		return shim.Error("票据背书拒签失败，参数有误")
	}

	bill, bl := GetBill(stup, args[0])
	if !bl {
		return shim.Error("票据背书拒签失败,g根据票据号码查询信息发生错误")
	}

	//以待背书人证件号码以及票据号码创建复合键，以便当前用户无法再从待背书票据列表总查询该票据信息
	waitEndorseCmIDBillInfoIDIndexKey, err := stup.CreateCompositeKey(IndexName, []string{args[1], bill.BillInfoID})

	if err != nil {
		return shim.Error("票据背书拒签失败,创建复合键时发生错误")
	}
	err = stup.DelState(waitEndorseCmIDBillInfoIDIndexKey)
	if err != nil {
		return shim.Error("票据背书拒签失败,删除复合键时发生错误")
	}

	//修改票据状态信息 背书状态 带倍数人信息 拒绝背书人信息
	bill.State = BillInfo_State_EndorseReject
	bill.RejectEndorseAcct = args[2]
	bill.RejectEndorseCmID = args[1]
	bill.WaitEndorseCmID = ""
	bill.WaitEndorseAcct = ""

	//保存票据状态
	_, bl = PutBill(stup, bill)
	if !bl {
		return shim.Error("票据背书拒签失败,保存票据状态发生错误")
	}
	msg, _ := GetMsgByte(0, "票据拒签成功")
	return shim.Success(msg)


}
