package conf

// 配置文件的结构体 对应yaml文件中的 mysql logger system
type Config struct {
	Mysql     Mysql     `yaml:"mysql"`
	Logger    Logger    `yaml:"logger"`
	System    System    `yaml:"system"`
	Redis     Redis     `yaml:"redis"`
	JwtSecret JwtSecret `yaml:"jwt"`
}
