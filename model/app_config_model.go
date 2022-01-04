package model

type AppConfig struct {
	ApplicationName string
	JwtSignatureKey string
	RedisHost       string
	RedisPort       string
	ApiHost         string
	ApiPort         string
}
