package types

type CfgEnv struct {
	SL_APP_NAME  string `env:"SL_APP_NAME"`
	SL_LOG_LEVEL string `env:"SL_LOG_LEVEL"`
	SL_HTTP_IP   string `env:"SL_HTTP_IP"`
	SL_HTTP_PORT string `env:"SL_HTTP_PORT"`
	SL_DB_IP     string `env:"SL_DB_IP"`
	SL_DB_PORT   string `env:"SL_DB_PORT"`
}
