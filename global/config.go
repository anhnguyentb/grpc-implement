package global

import (
	"github.com/spf13/viper"
)

func LoadConfig() error {

	viper.SetConfigName("configurations")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../.")
	return viper.ReadInConfig()
}
