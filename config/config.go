package config

type APIConfig struct {
	Host string `hcl:"host"`
	Port int    `hcl:"port"`
}

type NutsDBCfg struct {
	Path        string `hcl:"path"`
	SegmentSize int64  `hcl:"segment_size"`
}

type RedisConf struct {
	DevMode  bool   `hcl:"dev_mode"`
	Password string `hcl:"password"`
	Host     string `hcl:"host"`
}

type CacheCfg struct {
	Type string `hcl:"type"`

	Redis  RedisConf `hcl:"redis"`
	NutsDB NutsDBCfg `hcl:"nutsdb"`
}

type CommonCfg struct {
	Exchange     string `hcl:"exchange"`
	ExchangeType string `hcl:"exchange_type"`
}

type RabbitMQCfg struct {
	Host        string    `hcl:"host"`
	User        string    `hcl:"user"`
	Password    string    `hcl:"password"`
	ConsumerTag string    `hcl:"consumer_tag"`
	Common      CommonCfg `hcl:"common"`
}

// Root config structure
type Config struct {
	API        APIConfig   `hcl:"api"`
	EnableAuth bool        `hcl:"enable_auth"`
	Cache      CacheCfg    `hcl:"cache"`
	Rabbit     RabbitMQCfg `hcl:"rabbitmq"`
}
