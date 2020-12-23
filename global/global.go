package global

import (
	"distributed_file/config"
	"encoding/json"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	//"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	NacosConfig *config.NacosInit = &config.NacosInit{}
	ServerConfig *config.ServerConfig = &config.ServerConfig{}
)

func init(){
	configFile := "./config.yaml"
	v := viper.New()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil{
		panic(err)
	}
	if err := v.Unmarshal(NacosConfig);err != nil{
		panic(err)
	}
	readFromnacos()

}

func readFromnacos(){
	nacosInfo := NacosConfig.Nacos
	sc := []constant.ServerConfig{
		{
			IpAddr: nacosInfo.Host,
			Port: uint64(nacosInfo.Port),
		},
	}
	cc := constant.ClientConfig{
		NamespaceId: nacosInfo.Namespace,
		TimeoutMs: 5000,
		NotLoadCacheAtStart: true,
		LogDir: "tmp/nacos/log",
		CacheDir: "tmp/nacos/cache",
		RotateTime: "1h",
		MaxAge: 3,
		LogLevel: "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":  cc,
	})
	if err != nil{
		panic(err)
	}
	content, err := configClient.GetConfig(
		vo.ConfigParam{
			DataId: nacosInfo.DataId,
			Group: nacosInfo.Group,
		})
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(content), &ServerConfig)
	if err != nil {
		panic(err)
	}
}
