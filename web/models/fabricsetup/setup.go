package fabricsetup


import (
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/pkg/errors"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

const rootPath = "/home/itcast"

type FabricSetup struct {
	//xxx其他字段，暂且忽略
	//xxx其他字段，暂且忽略
	ConfigFile string //初始化 SDK 对应的配置文件

	OrgID string // orderer节点
	OrdererID string //通道id
	ChannelID string //链代码id
	ChainCodeID string

	initialized bool
	ChannelConfig string // 与channel相关的配置文件
	ChaincodeGoPath string // GOPATH 环境变量
	ChaincodePath string // ChainCode 存放路径
	OrgAdmin string // Admin 用户名
	OrgName string //组织名字
	UserName string // 普通用户名
	EventID string //消息id
	Client *channel.Client //通道客户端句柄
	Admin *resmgmt.Client //管理通道
	Sdk *fabsdk.FabricSDK //官方sdk 对象
	Event *event.Client //消息客户端
}

func NewFabricSetup() *FabricSetup {
	fSetup := FabricSetup{
		// Network parameters
		OrdererID: "orderer.itcast.cn",

		// Channel parameters
		ChannelID:     "fgjorgschannel", //duke,配套就行，与config.yaml

		ChannelConfig: rootPath + "/trace0/service/channel-artifacts/fgjorgschannel.tx",
		// ChannelConfig: os.Getenv("GOPATH") + "/service/channel-artifacts/fgjchannel.tx",

		ChainCodeID:     "mycc",

		ChaincodeGoPath: rootPath + "/trace0/web", //自己拼的，没有使用标准目录的代价
		ChaincodePath:   "chaincode",
		// rootPath + "/trace0/web"	+ "src" + ChaincodePath

		// ChaincodeGoPath: os.Getenv("GOPATH"),
		// ChaincodePath:   "origins/chaincode",

		OrgAdmin:        "Admin",
		OrgName:         "ofgj",
		ConfigFile:      rootPath + "/trace0/service/config.yaml",

		// User parameters
		UserName: "User1",
		EventID:"eventInvoke",
	}

	return &fSetup
}

func (sdk *FabricSetup)Init() error{
	fmt.Println("初始化sdk begin")
	err := sdk.initializeSDK()
	if err !=nil{
		return err
	}
	fmt.Println("创建通道，加入通道begin")
	err = sdk.createAndJoinChannel()
	if err !=nil{
		return err
	}

	fmt.Println("安装通道并初始化begin")
	err = sdk.installAndInstantiateCC()
	if err !=nil{
		return err
	}
	return nil
}

//+++++++++++++ 以下的代码是对sdk进行操作，包括，创建通道，添加通道，安装链码，初始化链码++++++

// Initialize reads the configuration file and sets up the client, chain and event hub
//初始化 读取配置文件 并且设置客户端 ,链 和事件
func (setup *FabricSetup) initializeSDK() error {
	// Add parameters for the initialization
	//判断是否进行了初始化的操作
	if setup.initialized == true {
		return errors.New("sdk already initialized")
	}

	// Initialize the SDK with the configuration file
	//通过配置文件得到sdk对象
	//通过sdk的配置文件  来对sdk进行创建
	sdk, err := fabsdk.New(config.FromFile(setup.ConfigFile))
	if err != nil {
		return errors.WithMessage(err, "failed to create SDK")
	}
	//将FabricSetup 的 sdk进行赋值
	setup.Sdk = sdk
	fmt.Println("SDK created")
	return nil
}


func (setup *FabricSetup)createAndJoinChannel() error {
	// The resource management client is responsible for managing channels (create/update channel)
	//创建sdk实例给用户和组织创建上下文件
	resourceManagerClientContext := setup.Sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	if resourceManagerClientContext == nil {
		return errors.New("failed to load Admin identity")
	}
	//依赖于上下文 创建1个资源管理客户端
	resMgmtClient, err := resmgmt.New(resourceManagerClientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create channel management client from Admin identity")
	}
	//得到客户端
	setup.Admin = resMgmtClient
	fmt.Println("Ressource management client created")
	//配置相关信息
	req := resmgmt.SaveChannelRequest{
		//通道id
		ChannelID: setup.ChannelID,
		//连代码所在文件
		ChannelConfigPath: setup.ChannelConfig,
	}
	//保存起来
	txID, err := setup.Admin.SaveChannel(req)
	//判断是否成功
	if err != nil || txID.TransactionID == "" {
		return errors.WithMessage(err, "failed to save channel")
	}
	fmt.Println("Channel created")

	// Make admin user join the previously created channel
	//让管理员用户加入前面创建的通道
	err = setup.Admin.JoinChannel(
		setup.ChannelID,
		//重试机制
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(setup.OrdererID),
	)
	if err != nil {
		return errors.WithMessage(err, "failed to make admin join channel")
	}
	fmt.Println("Channel joined")

	fmt.Println("Initialization Successful")
	setup.initialized = true
	return nil
}


//安装与实例化连代码
func (setup *FabricSetup) installAndInstantiateCC() error {

	// Create the chaincode package that will be sent to the peers
	//创建1个 sspkg  依赖于连代码的 连代码地址   和连代码的go语言环境
	//func NewCCPackage(chaincodePath string, goPath string)
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("ccPkg created")

	// Install example cc to org peers

	//创建1个安装连代码参数
	installCCReq := resmgmt.InstallCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodePath, Version: "0", Package: ccPkg}
	//安装链代码
	_, err = setup.Admin.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return errors.WithMessage(err, "failed to install chaincode")
	}
	fmt.Println("Chaincode installed")

	// 设置背书策略，该参数依赖 configtx.yaml 文件中 Organizations -> ID
	// 背书规则只针对chaincode中写入数据的操作进行校验，对于查询类操作不背书
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"ofgj.itcast.cn"})
	//实例化链代码
	resp, err := setup.Admin.InstantiateCC(setup.ChannelID, resmgmt.InstantiateCCRequest{Name: setup.ChainCodeID, Path: setup.ChaincodeGoPath, Version: "0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil || resp.TransactionID == "" {
		return errors.WithMessage(err, "failed to instantiate the chaincode")
	}
	fmt.Println("Chaincode instantiated")

	// Channel client is used to query and execute transactions
	//创建用户访问channel的client的上下文
	clientContext := setup.Sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	//基于上下文获取客户端
	setup.Client, err = channel.New(clientContext)
	if err != nil {
		return errors.WithMessage(err, "failed to create new channel client")
	}
	fmt.Println("Channel client created")

	fmt.Println("Chaincode Installation & Instantiation Successful")
	setup.initialized = true
	return nil
}