package main

import (
	"fmt"
	"github.com/wuzhanfly.com/bill/blockchain"
	"github.com/wuzhanfly.com/bill/service"
	"os"
	"github.com/wuzhanfly.com/bill/web"
	"github.com/wuzhanfly.com/bill/web/controller"
)

func main() {
	// 定义SDK属性
	fSetup := blockchain.FabricSetup{
		// 公共参数
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		ConfigFile: "config.yaml",
		// 通道相关
		ChannelID:     "mychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/github.com/wuzhanfly.com/bill/fixtures/artifacts/channel.tx",
	// 链码相关参数
        ChaincodeID: "bill",
		ChaincodeGoPath:   os.Getenv("GOPATH"),
		ChaincodePath:  "github.com/wuzhanfly.com/bill/chaincode/",

	// 指定用户
		UserName:       "User1",
}
	// 初始化SDK

	err := fSetup.Initialize()
	if err != nil {
		fmt.Printf("无法初始化Fabric SDK: %v\n", err)
	}
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("无法安装及实例化链码: %\n", err)
	}

	//测试链码成
	fsService := new(service.FabricSetupService)
	fsService.Fabric = &fSetup
	//====================业务层开始测试=============================
	/*//票据
	bill := service.Bill{
		BIllInfoID:"BOC101",
		BillInfoAmt:"1000",
		BillInfoType:"111",
		BillInfoBgDate:"20180101",
		BillInfoEdDate:"20181010",

		DrwrAcct:"111",
		DrwrCmID:"111",

		AccptrAcct:"111",
		AccptrCmID:"111",

		PyeeAcct:"111",
		PyeeCmID:"111",

		HoldrAcct:"jack",
		HoldrCmID:"jackID",

	}
	bill2 := service.Bill{
		BIllInfoID:"BOC102",
		BillInfoAmt:"2000",
		BillInfoType:"111",
		BillInfoBgDate:"20180101",
		BillInfoEdDate:"20181010",

		DrwrAcct:"111",
		DrwrCmID:"111",

		AccptrAcct:"111",
		AccptrCmID:"111",

		PyeeAcct:"111",
		PyeeCmID:"111",

		HoldrAcct:"jack",
		HoldrCmID:"jackID",

	}
	// 发布票据 返回交易id
	msg, err := fsService.SaveBill(bill)
	if err != nil {
		fmt.Printf("发布票据失败: %v\n", err)
	}else {
		fmt.Println("发布票据成功 =========> " + msg)
	}
	// 发布票据 返回交易id
	msg, err = fsService.SaveBill(bill2)
	if err != nil {
		fmt.Printf("发布票据失败: %v\n", err)
	}else {
		fmt.Println("发布票据成功 =========> " + msg)
	}
	//查询持票人的票据列表
	result, err := fsService.FindBills("jackID")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("根据持有人查询票据成功")
		var bills = []service.Bill{}
		json.Unmarshal(result, &bills)

		for _,obj := range bills{
			fmt.Println(obj)
		}
	}

	//发起背书请求

	msg, err =fsService.Endorse("BOC101","aliceID","alice")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println(msg)
	}
	//查询待背书列表
	result, err = fsService.FindWaitBills("aliceID")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println("查询待背书列表成功")
		var bills = []service.Bill{}
		json.Unmarshal(result, &bills)
		for _,obj := range bills {
			fmt.Println(obj)
		}
	}

	//签收票据
	msg, err = fsService.Accept("BOC101","aliceID","alice")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println(msg)
	}
	//根据票据no查询详情
	result, err = fsService.FindBillByNo("BOC101")
	if err != nil {
		fmt.Errorf(err.Error())
	}else {
		var bill service.Bill
		json.Unmarshal(result,&bill)
		fmt.Println(bill)
	}

	//发起背书请求

	msg, err =fsService.Endorse("BOC102","aliceID","alice")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println(msg)
	}

	//票据拒签
	msg, err = fsService.Reject("BOC102","aliceID","alice")
	if err != nil {
		fmt.Println(err.Error())
	}else {
		fmt.Println(msg)
	}
	//根据票据no查询详情
	result, err = fsService.FindBillByNo("BOC102")
	if err != nil {
		fmt.Errorf(err.Error())
	}else {
		var bill service.Bill
		json.Unmarshal(result,&bill)
		fmt.Println(bill)
	}*/
	//====================业务层测试结束=============================

	//调用webservice启动web服务
	//web.WebStart()
	app := controller.Application{
		Setup: fsService,
	}
	web.WebStart(app)

}
