package service

import (
	"github.com/wuzhanfly.com/bill/blockchain"
)

type FabricSetupService struct {
	Fabric *blockchain.FabricSetup
} 

type Bill struct {
	BillInfoID   string `json:"BillInfoID"`    //票据号码
	BillInfoAmt  string `json:"Bill_info_amt"` //票据金额
	BillInfoType string `json:"BillInfoType"`  //票据类型BillInfoAmt

	BillInfoBgDate string `json:"BillInfoBgDate"` //出票日期
	BillInfoEdDate string `json:"BillInfoEdDate"` //到期时间

	DrwrAcct string `json:"DrwrAcct"` //出票人名称
	DrwrCmID string `json:"DrwrCmID"` //出票人证件号码

	AccptrAcct string `json:"AccptrAcct"` //承兑人名称
	AccptrCmID string `json:"AccptrCmID"` //承兑人证件号码

	PyeeAcct string `json:"PyeeAcct"` //收款人名称
	PyeeCmID string `json:"PyeeCmID"` //收款人证件号码

	HoldrAcct string `json:"HoldrAcct"` //当前持有人名称
	HoldrCmID string `json:"HoldrCmID"` //当前持有人证件ｉｄ

	WaitEndorseAcct string `json:"WaitEndorseAcct"` //带背书人名称
	WaitEndorseCmID string `json:"WaitEndorseCmID"` //带背书人证件ｉｄ

	RejectEndorseAcct string `json:"RejectEndorseAcct"` //拒绝背书人名称
	RejectEndorseCmID string `json:"RejectEndorseCmID"` //拒绝背书人证件号码

	State    string `json:" State"` //票据状态
	Historys []HistoryItem          //票据背书历史
}

type HistoryItem struct {
	TxId string
	Bill Bill
}
