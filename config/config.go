package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Configurations struct {
	Kube       KubeConf
	Scheduler  SchedulerConf
	Weatherbit WeatherbitConf
}

type KubeConf struct {
	ConfPath string
}

type SchedulerConf struct {
	Interval string
}

type WeatherbitConf struct {
	Path string
	Key  string
}

var Root Configurations

func setup() Configurations {
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %s", err)
	}

	var c Configurations
	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Printf("Unable to decode config: %s", err)
	}
	return c
}

func init() {
	Root = setup()
}
