package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/youchuangcd/gopkg/common/utils"
	"github.com/youchuangcd/gopkg/mylog"
	"os"
)

var (
	// LogCategory 日志分类名称
	LogCategory = "nacos"
)

type Nacos struct {
	// 监听配置回调通知函数 ，接收content内容 返回 int
	// 业务侧创建一个回调函数  func xxxxx(data string) (flag int) {}
	// nacos.Nacos{Callback: xxxx}
	Callback     func(confUnmarshalMapValue map[string]UnmarshalMapValue, group, dataId, content string) error
	config       Config
	client       config_client.IConfigClient
	namingClient naming_client.INamingClient
}

type Config struct {
	IpAddr       string
	ContextPath  string
	Scheme       string
	Port         uint64
	NamespaceId  string
	Group        string
	AccessKey    string
	SecretKey    string
	UnmarshalMap map[string]UnmarshalMapValue
}
type UnmarshalMapValue struct {
	Conf           interface{}
	ChangeCallback func()
}

func InitNacos() {
	logger.InitLogger(logger.Config{CustomLogger: mylog.GetLogger()})
}

func New(config Config) *Nacos {
	nacos := &Nacos{
		config: config,
	}
	client, err := nacos.createConfig()
	if err != nil {
		panic("无法创建Nacos动态配置客户端")
	}
	nacos.client = client
	return nacos
}

func NewNaming(config Config) *Nacos {
	nacos := &Nacos{
		config: config,
	}
	namingClient, err := nacos.createNaming()
	if err != nil {
		panic("无法创建Nacos服务发现客户端")
	}
	nacos.namingClient = namingClient
	return nacos
}

// 创建动态配置客户端
func (p *Nacos) createConfig() (configClient config_client.IConfigClient, err error) {
	sc := []constant.ServerConfig{{
		IpAddr: p.config.IpAddr,
		Port:   p.config.Port,
		//ContextPath: p.config.ContextPath,
		//Scheme: p.config.Scheme,
	}}
	currentDir, err := os.Getwd()
	if currentDir == "/" {
		currentDir = ""
	}
	cc := constant.ClientConfig{
		NamespaceId:         p.config.NamespaceId,
		TimeoutMs:           5 * 1000,
		ListenInterval:      30 * 1000,
		NotLoadCacheAtStart: true,
		CacheDir:            currentDir,         //默认会把缓存下来的文件写入 currentDir/config
		AccessKey:           p.config.AccessKey, // ACM&KMS的AccessKey，用于配置中心的鉴权
		SecretKey:           p.config.SecretKey,
	}
	// 创建动态配置客户端的另一种方式 (推荐)
	configClient, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	return
}

// 读取配置内容
// (n Nacos)
func (p *Nacos) GetConfig(group, dataId string) (content string, err error) {
	content, err = p.client.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		mylog.Error(nil, LogCategory, fmt.Sprintf("send msg failed, err:%s", err.Error()))
		panic("无法连接Nacos，需要重新发布")
	}
	_ = p.Callback(p.config.UnmarshalMap, group, dataId, content)
	// 增加监听配置是否变化
	p.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, content string) {
			mylog.WithInfo(nil, LogCategory, map[string]interface{}{
				"namespace": namespace,
				"group":     group,
				"dataId":    dataId,
				"content":   content,
			}, "配置变动")
			// 避免出现运行中改错配置，导致监听配置恐慌
			utils.WithRecover(func() {
				_ = p.Callback(p.config.UnmarshalMap, group, dataId, content)
			})
		},
	})
	return
}

func (p *Nacos) createNaming() (namingClient naming_client.INamingClient, err error) {
	sc := []constant.ServerConfig{{
		IpAddr: p.config.IpAddr,
		Port:   p.config.Port,
		//ContextPath: p.config.ContextPath,
		//Scheme: p.config.Scheme,
	}}
	currentDir, err := os.Getwd()
	if currentDir == "/" {
		currentDir = ""
	}
	cc := constant.ClientConfig{
		NamespaceId:         p.config.NamespaceId,
		TimeoutMs:           5 * 1000,
		ListenInterval:      30 * 1000,
		NotLoadCacheAtStart: true,
		CacheDir:            currentDir,         //默认会把缓存下来的文件写入 currentDir/config
		AccessKey:           p.config.AccessKey, // ACM&KMS的AccessKey，用于配置中心的鉴权
		SecretKey:           p.config.SecretKey,
	}
	// 创建服务发现客户端的另一种方式 (推荐)
	namingClient, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	return
}

func (p *Nacos) SelectOneHealthyInstance(serviceName string, groupName string, clusters []string) (instance *model.Instance, err error) {
	instance, err = p.namingClient.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		ServiceName: serviceName,
		GroupName:   groupName, // 默认值DEFAULT_GROUP
		Clusters:    clusters,  // 默认值DEFAULT
	})
	return
}

type SubscribeParam struct {
	ServiceName       string                                       `param:"serviceName"` //required
	Clusters          []string                                     `param:"clusters"`    //optional,default:DEFAULT
	GroupName         string                                       `param:"groupName"`   //optional,default:DEFAULT_GROUP
	SubscribeCallback func(services []SubscribeService, err error) //required
}

type SubscribeService model.SubscribeService

func (p *Nacos) Subscribe(param *SubscribeParam) error {
	return p.namingClient.Subscribe(&vo.SubscribeParam{
		ServiceName: param.ServiceName,
		Clusters:    param.Clusters,
		GroupName:   param.GroupName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			var copyService = make([]SubscribeService, 0, len(services))
			for _, v := range services {
				copyService = append(copyService, SubscribeService{
					ClusterName: v.ClusterName,
					Enable:      v.Enable,
					InstanceId:  v.InstanceId,
					Ip:          v.Ip,
					Metadata:    v.Metadata,
					Port:        v.Port,
					ServiceName: v.ServiceName,
					Valid:       v.Valid,
					Weight:      v.Weight,
					Healthy:     v.Healthy,
				})
			}
			param.SubscribeCallback(copyService, err)
		},
	})
}
