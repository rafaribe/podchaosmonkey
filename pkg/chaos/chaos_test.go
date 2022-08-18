package chaos

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes/fake"

	env "podchaosmonkey/pkg/environment"
)

func TestNewPodChaosMonkey(t *testing.T) {
	tests := map[string]struct {
		input           map[string]string
		want            *PodChaosMonkey
		inClusterConfig bool
	}{
		"outside cluster config": {input: map[string]string{"KUBECONFIG": "/Users/rafaribe/.kube/config", "NAMESPACE": "workloads"}, want: &PodChaosMonkey{Namespace: "workloads", GracePeriodInSeconds: 5, IntervalInSeconds: 10}, inClusterConfig: false},
		"in cluster config":      {input: map[string]string{"NAMESPACE": "workloads"}, want: &PodChaosMonkey{Namespace: "workloads", GracePeriodInSeconds: 5, IntervalInSeconds: 10}, inClusterConfig: true},
	}

	for name, tc := range tests {
		env.LoadEnv()
		t.Run(name, func(t *testing.T) {
			env.SetVals(tc.input)
			client := fake.NewSimpleClientset()
			if tc.inClusterConfig {
				assert.Empty(t, viper.GetString("KUBECONFIG"))
			}
			got := NewPodChaosMonkey(client)

			assert.Equal(t, tc.want.Namespace, got.Namespace)
			assert.Equal(t, tc.want.GracePeriodInSeconds, got.GracePeriodInSeconds)
			assert.Equal(t, tc.want.IntervalInSeconds, got.IntervalInSeconds)

		})
		env.ClearVals(tc.input)
	}
}
