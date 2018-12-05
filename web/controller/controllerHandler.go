package controller

import (
	"net/http"
	"github.com/wuzhanfly.com/bill/service"
	"encoding/json"
	"fmt"
)

var cuser User

func (app *Application)LoginView(w http.ResponseWriter, r *http.Request)  {
	ShowView(w,r,"login.html",nil)
}
func (app *Application)Login(w http.ResponseWriter, r *http.Request)  {

	loginName := r.FormValue("loginName")
	password := r.FormValue("password")
	data := &struct {
		Cuser User
		Flag  bool
	}{
		Flag:false,
	}

	var flag bool
	for _,user := range users {
		if user.LoginName == loginName && user.Password ==password {
			cuser = user
			flag = true
			break
		}
	}

	if flag {
		//登入成功
		fmt.Println("当前登入用户信息：", cuser)
		r.Form.Set("holdCmID", cuser.CmId)
		app.FindBills(w,r)
	}else {
		//导入失败
		data.Flag = true
		data.Cuser.LoginName = loginName
		ShowView(w,r,"login.html",data)
	}
	//fmt.Fprintf(w,"OK!")//输出到客户端
}
func (app *Application)Loginout(w http.ResponseWriter, r *http.Request)  {
	cuser = User{}
	ShowView(w,r,"login.html",nil)

}

//示发布票据页面
func (app *Application)IssueShow (w http.ResponseWriter, r *http.Request)  {
	data := &struct {
		Cuser User
		Flag bool
	}{
		Cuser:cuser,
		Flag:false,
	}
	ShowView(w,r,"issue.html",data)

}

//发布票据
func (app *Application)Issue(w http.ResponseWriter, r *http.Request)  {
	bill := service.Bill{
		BillInfoID:r.FormValue("BillInfoID"),
		BillInfoAmt:r.FormValue("BillInfoAmt"),
		BillInfoType:r.FormValue("BillInfoType"),
		BillInfoBgDate:r.FormValue("BillInfoBgDate"),
		BillInfoEdDate:r.FormValue("BillInfoEdDate"),
		DrwrAcct:r.FormValue("DrwrAcct"),
		DrwrCmID:r.FormValue("DrwrCmID"),
		AccptrAcct:r.FormValue("AccptrAcct"),
		AccptrCmID:r.FormValue("AccptrCmID"),
		PyeeAcct:r.FormValue("PyeeAcct"),
		PyeeCmID:r.FormValue("PyeeCmID"),
		HoldrAcct:r.FormValue("holdrAcct"),
		HoldrCmID:r.FormValue("holdrCmID"),
	}
	transactionID, err := app.Setup.SaveBill(bill)

	data:=&struct {
		Cuser User
		Msg string
		Flag bool
	}{
		Cuser:       cuser,
		Flag:        true,
		Msg:         "",
	}

	if err != nil {
		data.Msg = err.Error()
	}else {
		data.Msg = "票据发布成功:"+ transactionID
	}

	ShowView(w,r,"issue.html",data)

}
//查询持票人票据列表
func (app *Application)FindBills(w http.ResponseWriter,r * http.Request)  {
	holdeCmId := cuser.CmId
	b, err := app.Setup.FindBills(holdeCmId)
	if err != nil {
		fmt.Println("查询当前用户票据列表失败：",err.Error())
	}
	var bills = []service.Bill{}
	json.Unmarshal(b, &bills)
	data := &struct {
		Bills []service.Bill
		Cuser User
	}{
		Bills: bills,
		Cuser: cuser,
	}
	ShowView(w,r,"bills.html",data)
}
//根据票据号码查询票据详情
func (app *Application) BillInfoByNo(w http.ResponseWriter,r *http.Request) {
	billInfoNo:= r.FormValue("billInfoNo")
	result, err := app.Setup.FindBillByNo(billInfoNo)
	if err != nil {
		fmt.Println(err.Error())
	}
	var bill = service.Bill{}
	json.Unmarshal(result,&bill)
	data := &struct {
		Bill service.Bill
		Cuser User
		Flag bool
		Msg string
	}{
		Bill:bill,
		Cuser:cuser,
		Flag:false,
		Msg:"",
	}
	if r.FormValue("msg") != "" {
		data.Msg = r.FormValue("mag")
	}
	ShowView(w,r,"billInfo.html",data)
}
//发起背书请求详情
func (app *Application)Endorse(w http.ResponseWriter,r *http.Request)  {
	waitEndorseAcct := r.FormValue("WaitEndorseAcct")
	waitEndorseCmId := r.FormValue("WaitEndorseCmID")
	billNo := r.FormValue("billNo")
	fmt.Println(billNo,waitEndorseAcct,waitEndorseCmId)
	
	_, err := app.Setup.Endorse(billNo, waitEndorseCmId, waitEndorseAcct)
	if err != nil {
		r.Form.Set("Msg",err.Error())
	}else {
		r.Form.Set("Msg","背书请求成功,待"+waitEndorseCmId +"签收")
	}

	r.Form.Set("billInfoNo",billNo)
	r.Form.Set("flag","t")

	app.BillInfoByNo(w,r)

}


//查询当前用户的带倍数票据列表
func (app *Application)WaitEndorseBills(w http.ResponseWriter,r *http.Request)  {
	waitEndorseCmID :=cuser.CmId
	result, err := app.Setup.FindWaitBills(waitEndorseCmID)
	if err != nil {
		fmt.Println(err.Error())
	}
	
	var bills = []service.Bill{}
	json.Unmarshal(result, &bills)
	
	data := &struct {
		Bills []service.Bill
		Cuser User
	}{
		Bills:bills,
		Cuser:cuser,
	}
	ShowView(w,r,"waitEndorse.html",data)
}


//根据票据号码查询待签收票据详情

func (app *Application)WaitEndorseInfo(w http.ResponseWriter,r *http.Request)  {
	billNo := r.FormValue("billNo")
	result, _ := app.Setup.FindBillByNo(billNo)
	var bill service.Bill
	json.Unmarshal(result, &bill)

	data:= &struct {
		Bill service.Bill
		Cuser User
		Flag bool
		Msg string
	}{
	Bill:bill,
	Cuser:cuser,
	Flag:false,
	Msg:"",
	}
	if r.FormValue("msg") != "" {
		data.Flag=true
		data.Msg = r.FormValue("mag")
	}

	ShowView(w,r,"waitEndorseInfo.html",data)
}

//背书签收
func (app *Application)Accept(w http.ResponseWriter,r *http.Request){
	billNo := r.FormValue("billNo")
	signeAcct := cuser.Acct
	signCmId := cuser.CmId
	result, err :=app.Setup.Accept(billNo,signCmId,signeAcct)

	r.Form.Set("billNo",billNo)
	if err != nil {
		r.Form.Set("msg",err.Error())
	}else {
		r.Form.Set("msg",result)
	}

	app.WaitEndorseInfo(w,r)
}

func (app *Application)Reject(w http.ResponseWriter,r *http.Request){
	billNo := r.FormValue("billNo")
	rejectAcct := cuser.Acct
	rejectCmId := cuser.CmId
	result ,err := app.Setup.Reject(billNo,rejectCmId,rejectAcct)
	r.Form.Set("billNo",billNo)
	if err !=nil {
		r.Form.Set("msg",err.Error())
	}else {
		r.Form.Set("msg",result)
	}
	app.WaitEndorseInfo(w,r)

}