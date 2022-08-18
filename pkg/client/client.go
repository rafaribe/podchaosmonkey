package client

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func InitKubernetesClient() kubernetes.Interface {
	kubeConfigEnvPath := viper.GetString("KUBECONFIG")
	if kubeConfigEnvPath != "" {
		return getLocalKubernetesClient(&kubeConfigEnvPath)
	}
	return getInClusterKubernetesClient()
}

func getLocalKubernetesClient(path *string) kubernetes.Interface {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *path)
	if err != nil {
		log.Error("Failed to get kubeconfig: %s", err.Error())
		return nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func getInClusterKubernetesClient() kubernetes.Interface {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error("Failed to get kubeconfig: %s", err.Error())
		return nil
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}
