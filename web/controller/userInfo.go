package controller

import "github.com/wuzhanfly.com/bill/service"

type Application struct {
	Setup *service.FabricSetupService
}

type User struct {
	LoginName    string  `json:"UserName"`
	Password    string  `json:"Password"`
	CmId        string  `json:"CmId"`
	Acct        string  `json:"Acct"`
}

var users []User

func init()  {

	admin := User{LoginName:"Admin",Password:"123456",CmId:"AdminID",Acct:"管理员"}
	alice := User{LoginName:"Alice",Password:"123456",CmId:"AliceID",Acct:"Alice"}
	bob := User{LoginName:"bob",Password:"123456",CmId:"BobID",Acct:"BOB"}
	jack := User{LoginName:"jack",Password:"123456",CmId:"JackID",Acct:"Jack"}
	users = append(users, admin)
	users = append(users, alice)
	users = append(users, bob)
	users = append(users, jack)

}

