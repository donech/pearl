package grpc

type Config struct {
	Addr          string `yaml:"addr"`
	Port          string `yaml:"port"`
	WebPort       string `yaml:"webPort"`
	EnableReflect bool   `yaml:"enableReflect"`
	EnableGateWay bool   `yaml:"enableGateWay"`
}
