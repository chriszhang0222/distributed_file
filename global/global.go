package global
import (
	"distributed_file/config"
	"fmt"
	//"github.com/nacos-group/nacos-sdk-go/clients"
	//"github.com/nacos-group/nacos-sdk-go/common/constant"
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
	fmt.Println(NacosConfig.Nacos)
}
