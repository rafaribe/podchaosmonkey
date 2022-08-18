package chaos

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func TestGetAppConfiguration(t *testing.T) {
	{
		tests := map[string]struct {
			input           map[string]string
			want            *PodChaosMonkey
			wantNilConfig   bool
			inClusterConfig bool
		}{
			"outside cluster config": {input: map[string]string{"KUBECONFIG": "/Users/rafaribe/.kube/config", "NAMESPACE": "workloads"}, want: &PodChaosMonkey{Namespace: "workloads", GracePeriodInSeconds: 5, IntervalInSeconds: 10}, wantNilConfig: false},
			"in cluster config":      {input: map[string]string{"NAMESPACE": "workloads"}, want: &PodChaosMonkey{Namespace: "workloads", GracePeriodInSeconds: 5, IntervalInSeconds: 10}, wantNilConfig: true},
		}

		for name, tc := range tests {
			t.Run(name, func(t *testing.T) {
				setVals(tc.input)
				var client kubernetes.Interface
				if tc.inClusterConfig {
					client = fake.NewSimpleClientset()
				} else {
					client = InitKubernetesClient()
				}
				got := NewPodChaosMonkey(client)

				assert.Equal(t, tc.want.Namespace, got.Namespace)
				assert.Equal(t, tc.want.GracePeriodInSeconds, got.GracePeriodInSeconds)
				assert.Equal(t, tc.want.IntervalInSeconds, got.IntervalInSeconds)

				if tc.wantNilConfig {
					assert.Nil(t, got.Client)
				} else {
					assert.NotNil(t, got.Client)
				}
				clearVals(tc.input)
			})
		}
	}
}

// Utility Functuon to set the env vars during testing
func setVals(vals map[string]string) {
	for k, v := range vals {
		os.Setenv(k, v)
	}
}

// Utility Functuon just to remove the the env vars during testing
func clearVals(vals map[string]string) {
	for k := range vals {
		os.Setenv(k, "")
	}
}
