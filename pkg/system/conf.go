package system

type SystemConf struct {
	TrendCsv            string `yaml:"TrendCsv"`
	CointCsv            string `yaml:"CointCsv"`
	UmbrellaCsv         string `yaml:"UmbrellaCsv"`
	WeightCsv           string `yaml:"WeightCsv"`
	Platform            string `yaml:"PlatformCsv"`
	CointegrationSrcipt string `yaml:"CointegrationSrcipt"`
	LogPath             string `yaml:"LogPath"`
	DBPath              string `yaml:"DBPath"`
	DBType              string `yaml:"DBType"`
	Options             struct {
		Quantity  float32 `yaml:"quantity"`
		Pairing   string  `yaml:"pairing"`
		Test      bool    `yaml:"test"`
		Sl        float32 `yaml:"sl"`
		Tp        float32 `yaml:"tp"`
		EnableTsl bool    `yaml:"enable_tsl"`
		Tsl       float32 `yaml:"tsl"`
		Ttp       float32 `yaml:"ttp"`
	}
	Email struct {
		User     string   `yaml:"user"`
		Password string   `yaml:"password"`
		Host     string   `yaml:"host"`
		Port     string   `yaml:"port"`
		MailTo   []string `yaml:"mailTo"`
	}
	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database string `yaml:"Database"`
	}
}
