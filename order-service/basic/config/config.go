package config

type AppConfig struct {
	Mysql struct {
		Host     string
		Port     int
		User     string
		Password string
		Database string
	}
	Redis struct {
		Host     string
		Port     int
		Password string
		Database int
	}
	Consul struct {
		Host        string
		Port        int
		ServiceName string
		ServicePort int
		TTL         int
	}
	AliPay struct {
		AppId     string
		Key       string
		NotifyPay string
		Return    string
	}
}
type NacosConfig struct {
	Addr      string
	Port      int
	Namespace string
	DataID    string
	Group     string
}
