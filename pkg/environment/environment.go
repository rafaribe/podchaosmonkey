package environment

import (
	"os"

	"github.com/spf13/viper"
)

func LoadEnv() {
	viper.AutomaticEnv()
	viper.BindEnv("INTERVAL_IN_SECONDS")
	viper.BindEnv("NAMESPACE")
	viper.SetDefault("INTERVAL_IN_SECONDS", 10)
	viper.SetDefault("GRACE_PERIOD_SECONDS", 5)
	viper.SetDefault("NAMESPACE", "workloads")
	viper.SetDefault("LABELS", "podchaosmonkey=true")
	viper.SetDefault("INCLUDE_FINALIZERS", "false")
}

// Utility Functuon to set the env vars during testing
func SetVals(vals map[string]string) {
	for k, v := range vals {
		os.Setenv(k, v)
	}
}

// Utility Functuon to remove the the env vars during testing
func ClearVals(vals map[string]string) {
	for k := range vals {
		os.Setenv(k, "")
	}
}
