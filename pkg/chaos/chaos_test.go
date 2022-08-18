package chaos

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

func TestPodChaosMonkey_getPodList(t *testing.T) {
	client := fake.NewSimpleClientset()
	testingPods := createPodArray()
	matchingLabels, _ := labels.Parse("podchaosmonkey=true")

	tests := map[string]struct {
		configuration *PodChaosMonkey
		want          *v1.PodList
	}{
		"pod list without finalizers": {configuration: &PodChaosMonkey{
			Client:               client,
			Namespace:            "workloads",
			IntervalInSeconds:    10,
			GracePeriodInSeconds: 5,
			Labels:               matchingLabels,
			IncludeFinalizers:    false,
		}, want: &v1.PodList{Items: []v1.Pod{testingPods[0]}},
		}}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			for _, pod := range testingPods {
				tc.configuration.Client.CoreV1().Pods(tc.configuration.Namespace).Create(context.TODO(), &pod, metav1.CreateOptions{})
			}
			got := tc.configuration.getPodList()
			assert.NotEmpty(t, tc.want, got)
		})
	}
}

func createPodArray() []v1.Pod {
	var result []v1.Pod
	podWithoutFinalizers := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-without-finalizers",
			Namespace: "workloads",
			Labels: map[string]string{
				"podchaosmonkey": "true",
			},
			Finalizers: nil},
		Status: v1.PodStatus{Phase: v1.PodRunning},
	}
	podWithoutMatchingLabels := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-without-matching-labels",
			Namespace: "workloads",
			Labels: map[string]string{
				"podchaosmonkey": "false",
			},
			Finalizers: nil,
		},
		Status: v1.PodStatus{Phase: v1.PodRunning},
	}
	podWithFinalizers := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod-with-finalizers",
			Namespace: "workloads",
			Labels: map[string]string{
				"podchaosmonkey": "true",
			},
			Finalizers: []string{"kubernetes"},
		},
		Status: v1.PodStatus{Phase: v1.PodRunning},
	}

	result = append(result, podWithoutFinalizers)
	result = append(result, podWithoutMatchingLabels)
	result = append(result, podWithFinalizers)
	return result
}
