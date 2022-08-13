package configuration

import (
	"reflect"
	"testing"

	"k8s.io/client-go/rest"
)

func TestGetKubeconfig(t *testing.T) {
	tests := []struct {
		name string
		want *rest.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetKubeconfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKubeconfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getKubeconfigFromEnv(t *testing.T) {
	type args struct {
		path *string
	}
	tests := []struct {
		name string
		args args
		want *rest.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getKubeconfigFromEnv(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getKubeconfigFromEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getInClusterKubeconfig(t *testing.T) {
	tests := []struct {
		name string
		want *rest.Config
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getInClusterKubeconfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInClusterKubeconfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
