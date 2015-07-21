package config

type GlobalConfig struct {
	Version     string            `json:"version"`
	Credentials GlobalCredentials `json:"credentials"`
}

type GlobalCredentials struct {
	AccessKey string `json:"access-key"`
	SecretKey string `json:"secret-key"`
}

type Config struct {
	Name        string      `json:"name"`
	Status      int         `json:"status"`
	AwsTagName  string      `json:"aws-tag-name"`
	Credentials Credentials `json:"credentials"`
}

type Credentials struct {
	Credentials GlobalCredentials `json:"credentials"`
	FromGlobal  bool              `json:"from-global"`
}
