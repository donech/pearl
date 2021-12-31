package jwt

type Config struct {
	SingingAlgorithm string `yaml:"singingAlgorithm"`
	Key              string `yaml:"key"`
	PublicKeyFile    string `yaml:"publicKeyFile"`
	PrivateKeyFile   string `yaml:"privateKeyFile"`
	Timeout          string `yaml:"timeout"`
}
