package environment

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadEnv(t *testing.T) {
	{
		tests := map[string]struct {
			input  map[string]string
			output map[string]string
		}{
			"set different than default vals": {input: map[string]string{"KUBECONFIG": "/Users/myuser/.kube/config"}, output: map[string]string{"KUBECONFIG": "/Users/myuser/.kube/config"}},
			"test defaults":                   {input: map[string]string{}, output: map[string]string{"INTERVAL_IN_SECONDS": "10", "GRACE_PERIOD_SECONDS": "5", "NAMESPACE": "workloads", "LABELS": "podchaosmonkey=true"}},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				SetVals(tc.input)
				LoadEnv()
				for k, v := range tc.output {
					assert.NotNil(t, v, viper.GetString(k))
					assert.Equal(t, v, viper.GetString(k))
				}
				ClearVals(tc.input)
			})
		}
	}
}

func TestSetVals(t *testing.T) {
	{
		tests := map[string]struct {
			values map[string]string
		}{
			"set environment variables": {values: map[string]string{"KUBECONFIG": "/Users/myuser/.kube/config", "NAMESPACE": "workloads", "GRACE_PERIOD_IN_SECONDS": "5", "INTERVAL_IN_SECONDS": "10"}},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				SetVals(tc.values)
				for k, v := range tc.values {
					assert.Equal(t, v, os.Getenv(k))
				}
			})
		}
	}
}

func TestClearVals(t *testing.T) {
	{
		tests := map[string]struct {
			values map[string]string
		}{
			"set environment variables": {values: map[string]string{"KUBECONFIG": "/Users/myuser/.kube/config", "NAMESPACE": "workloads", "GRACE_PERIOD_IN_SECONDS": "5", "INTERVAL_IN_SECONDS": "10"}},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				SetVals(tc.values)
				for k, expected := range tc.values {
					actual := os.Getenv(k)
					assert.NotEmpty(t, expected, actual)
				}
				ClearVals(tc.values)
				for k := range tc.values {
					actual := os.Getenv(k)
					assert.Empty(t, actual)
				}
			})
		}
	}
}
