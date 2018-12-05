package web

import (
	"net/http"
	"fmt"
	"github.com/wuzhanfly.com/bill/web/controller"
)


//启动web服务，及路由信息
func WebStart(app controller.Application)  {
	 //指定文件服务器
	 fs := http.FileServer(http.Dir("web/static"))
	  http.Handle("/static/", http.StripPrefix("/static/", fs))



	//	指定路由信息（根据客户端请求，做出匹配请求）
	http.HandleFunc("/",app.LoginView)
	http.HandleFunc("/login",app.Login)
	http.HandleFunc("/loginout",app.Loginout)

	//发布票据
	http.HandleFunc("/addBill",app.IssueShow)//显示发布票据页面
	http.HandleFunc("/issue",app.Issue)//提交发布票据请求

	http.HandleFunc("/bills",app.FindBills)//查询当前持票人票据？？？？
	http.HandleFunc("/billInfo",app.BillInfoByNo)

	http.HandleFunc("/endorse",app.Endorse)//发起背书

	http.HandleFunc("/waitEndorseBills",app.WaitEndorseBills)//查询当前用户等待背书
	http.HandleFunc("/waitEndorseInfo",app.WaitEndorseInfo)//根据票据号码查询当前用户等待背书详情

	http.HandleFunc("/accept",app.Accept)//背书签收处理
	http.HandleFunc("/reject",app.Reject)//背书签拒签

	fmt.Println("启动Web服务, 监听端口号: 9000")

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		fmt.Println("启动Web服务错误")
	}

}
