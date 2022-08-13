package configuration

import (
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeconfig() *rest.Config {
	kubeConfigEnvPath := os.Getenv("KUBECONFIG")
	if kubeConfigEnvPath != "" {
		return getKubeconfigFromEnv(&kubeConfigEnvPath)
	}
	return getInClusterKubeconfig()
}

func getKubeconfigFromEnv(path *string) *rest.Config {

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *path)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func getInClusterKubeconfig() *rest.Config {
	// use the current context in kubeconfig
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	return config
}
