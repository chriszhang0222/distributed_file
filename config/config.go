package config

type RabbitConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Vhost string `mapstructure:"vhost" json:"vhost"`
}

type MySQLConfig struct {
	User string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
	Db string `mapstructure:"db" json:"db"`
}

type RedisConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type NacosInit struct {
	Nacos NacosConfig `mapstructure: "nacos"`
}

type ServerConfig struct {
	MySQL MySQLConfig `mapstructure:"mysql" json:"mysql"`
	Redis RedisConfig `mapstructure:"redis" json:"redis"`
	Rabbit RabbitConfig `mapstructure:"rabbit" json:"rabbit"`
}