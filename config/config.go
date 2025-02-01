package config

import (
	"strconv"

	"github.com/spf13/viper"
)

var (
	DB_USER               = ""
	DB_PASS               = ""
	DB_NAME               = ""
	DB_HOST               = ""
	PORT                  = ""
	REDIS_URL             = ""
	ENVIRONMENT           = ""
	MIGRATE               = false
	FB_JSON               = ""
	AWS_ACCESS_KEY_ID     = ""
	AWS_SECRET_ACCESS_KEY = ""
	SES_AWS_REGION        = ""
	CORS_ALLOWED_ORIGINS  = ""
	SECRET_KEY            = ""
	ADMIN_SECRET_KEY      = ""
	SERVER_TYPE           = ""
	SENDER_EMAIL          = ""
	SENDER_PASS           = ""
)

func LoadConfig() {
	DB_USER = viper.GetString("DB_USER")
	DB_PASS = viper.GetString("DB_PASS")
	DB_NAME = viper.GetString("DB_NAME")
	DB_HOST = viper.GetString("DB_HOST")

	PORT = strconv.Itoa(viper.GetInt("PORT"))
	REDIS_URL = viper.GetString("REDIS_URL")
	ENVIRONMENT = viper.GetString("ENVIRONMENT")
	MIGRATE = viper.GetBool("MIGRATE")

	AWS_ACCESS_KEY_ID = viper.GetString("AWS_ACCESS_KEY_ID")
	AWS_SECRET_ACCESS_KEY = viper.GetString("AWS_SECRET_ACCESS_KEY")
	SES_AWS_REGION = viper.GetString("SES_AWS_REGION")
	SECRET_KEY = viper.GetString("SECRET_KEY")
	SECRET_KEY = viper.GetString("ADMIN_SECRET_KEY")
	SERVER_TYPE = viper.GetString("SERVER_TYPE")
	SENDER_EMAIL = viper.GetString("SENDER_EMAIL")
	SENDER_PASS = viper.GetString("SENDER_PASS")

}
