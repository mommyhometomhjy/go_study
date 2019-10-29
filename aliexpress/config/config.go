package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	getConfig()
}

func getConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func GetDataBaseConfig() (driver, dbname string) {
	return viper.GetString("database.driver"), viper.GetString("database.dbname")
}
func GetAliexpress() (percent, exchnge float64) {
	return viper.GetFloat64("aliexpress.percent"), viper.GetFloat64("aliexpress.exchange")
}
