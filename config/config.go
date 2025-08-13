package config

import (
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var Config struct {
	AppProd                     bool
	AppMode                     string
	AppID                       string
	AppSecret                   string
	AppLanguage                 string
	MysqlHost                   string
	MysqlPort                   string
	MysqlName                   string
	MysqlUser                   string
	MysqlPass                   string
	AllowOrigins                string
	AllowHeaders                string
	LogLevel                    string
	AppSalt                     string
	TemplateID                  []string
	AlibabaCloudAccessKeyID     string
	AlibabaCloudAccessKeySecret string
}

func envOr(env string, or string) string {
	rt := os.Getenv(env)
	if rt != "" {
		return rt
	}
	return or
}

func initConfig() {
	Config.AppProd = os.Getenv("APP_PROD") != ""
	if Config.AppProd {
		Config.AppMode = "release"
	} else {
		Config.AppMode = "debug"
	}
	Config.AppID = envOr("APP_ID", "wx1234567890")
	Config.AppSecret = envOr("APP_SECRET", "gin-example:secret")
	Config.AppLanguage = envOr("APP_LANGUAGE", "en")
	Config.MysqlHost = envOr("APP_MYSQL_HOST", "127.0.0.1")
	Config.MysqlPort = envOr("APP_MYSQL_PORT", "3306")
	Config.MysqlName = envOr("APP_MYSQL_NAME", "static")
	Config.MysqlUser = envOr("APP_MYSQL_USER", "root")
	Config.MysqlPass = envOr("APP_MYSQL_PASS", "123456")
	Config.AllowOrigins = envOr("APP_ALLOW_ORIGINS", "*")
	Config.AllowHeaders = envOr("APP_ALLOW_HEADERS", "Origin|Content-Length|Content-Type|Authorization")
	Config.LogLevel = envOr("APP_LOG_LEVEL", "info")
	Config.AppSalt = envOr("APP_SALT", "saltwithlength16")
	Config.TemplateID = strings.Split(envOr("APP_TEMPLATE_ID", "wx1234567890"), ",")
	Config.AlibabaCloudAccessKeyID = envOr("ALIBABA_ACCESS_KEY_ID", "123456")
	Config.AlibabaCloudAccessKeySecret = envOr("ALIBABA_ACCESS_KEY_SECRET", "123456")
}
