package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	"time"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"

	 "github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	 "github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	 "github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
)

// FabricSetup implementation
type FabricSetup struct {
	ConfigFile    string
	ChannelID     string
	initialized   bool
	ChannelConfig string
	OrgAdmin      string
	OrgName       string
	admin         resmgmtclient.ResourceMgmtClient
	sdk           *fabsdk.FabricSDK

	 // 链码相关    
	ChaincodeID     string    
	ChaincodeGoPath string
	ChaincodePath   string    
	UserName        string    
	Client          chclient.ChannelClient


}
const chaincodeVersion = "1.0"
// Initialize reads the configuration file and setsup the client, chain and event hub
func (t *FabricSetup) Initialize() error {

	fmt.Println("开始初始化......")

	if t.initialized {
		return fmt.Errorf("sdk已初始化完毕\n")
	}

	// 使用指定的配置文件创建SDK
	sdk, err := fabsdk.New(config.FromFile(t.ConfigFile))
	if err != nil {
		return fmt.Errorf("SDK创建失败: %v\n", err)
	}
	t.sdk = sdk

	// 根据指定的具有特权的用户创建用于管理通道的客户端API
	chMgmtClient, err := t.sdk.NewClient(fabsdk.WithUser(t.OrgAdmin),
	 fabsdk.WithOrg(t.OrgName)).ChannelMgmt()
	if err != nil {
		return fmt.Errorf("SDK添加管理用户失败: %v\n", err)
	}

	// 获取客户端的会话用户(目前只有session方法能够获取)
	session, err := t.sdk.NewClient(fabsdk.WithUser(t.OrgAdmin),
		fabsdk.WithOrg(t.OrgName)).Session()
	if err != nil {
		return fmt.Errorf("获取会话用户失败 %s, %s: %s\n", t.OrgName, t.OrgAdmin, err)
	}
	orgAdminUser := session

	// 指定用于创建或更新通道的参数
	req := chmgmtclient.SaveChannelRequest{
		ChannelID: t.ChannelID,
		ChannelConfig: t.ChannelConfig,
		SigningIdentity: orgAdminUser}
	// 使用指定的参数创建或更新通道
	err = chMgmtClient.SaveChannel(req)
	if err != nil {
		return fmt.Errorf("创建通道失败: %v\n", err)
	}

	time.Sleep(time.Second * 10)

	// 创建一个用于管理系统资源的客户端API。
	t.admin, err = t.sdk.NewClient(fabsdk.WithUser(t.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("创建资源管理客户端失败: %v\n", err)
	}
	// 将peer加入通道
	 err = t.admin.JoinChannel(t.ChannelID)
	 if err != nil {
		return fmt.Errorf("peer加入通道失败: %v\n", err)
	}

	fmt.Println("初始化成功")
	t.initialized = true
	return nil
}

// 安装并实例化链码
func (setup *FabricSetup) InstallAndInstantiateCC() error {
	fmt.Println("开始安装链码......")

	// 创建一个新的链码包并使用我们的链码初始化
	ccPkg, err := gopackager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v\n", err)
	}

	// 指定要安装链码的各项参数
	installCCReq := resmgmtclient.InstallCCRequest{
		Name:    setup.ChaincodeID,
		Path:    setup.ChaincodePath,
		Version: chaincodeVersion,
		Package: ccPkg,
	}

	// 在Org Peer上安装链码
	_, err = setup.admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("安装链码失败: %v\n", err)
	}
	fmt.Println("链码安装成功!")
	fmt.Println("开始实例化链码......")

	// 设置链码策略
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"Org1MSP"})
	fmt.Println(ccPolicy)

	// 指定实例化链码相关参数
	instantiateCCReq := resmgmtclient.InstantiateCCRequest{
		Name:    setup.ChaincodeID,
		Path:    setup.ChaincodePath,
		Version: chaincodeVersion,
		Args:    [][]byte{[]byte("init")},
		Policy:  ccPolicy,
	}
	fmt.Println(instantiateCCReq)
	// 实例化链码
	err = setup.admin.InstantiateCC(setup.ChannelID, instantiateCCReq)
	if err != nil {
		return fmt.Errorf("实例化链码失败: %v\n", err)
	}
	fmt.Println("链码实例化成功")


	// 创建通道客户端用于查询与执行事务
	setup.Client, err = setup.sdk.NewClient(fabsdk.WithUser(setup.UserName)).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("创建新的通道客户端失败: %v\n", err)
	}


	fmt.Println("链码实例化成功!")
	return nil


}
// Creation of the client which will enables access to our channel events

