package chaos

import (
	"context"
	"math/rand"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type PodChaosMonkey struct {
	Client               kubernetes.Interface
	Namespace            string
	IntervalInSeconds    int
	GracePeriodInSeconds int64
	Labels               labels.Selector
}

func NewPodChaosMonkey(client kubernetes.Interface) *PodChaosMonkey {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	viper.BindEnv("INTERVAL_IN_SECONDS")
	viper.BindEnv("NAMESPACE")
	viper.SetDefault("INTERVAL_IN_SECONDS", 10)
	viper.SetDefault("GRACE_PERIOD_SECONDS", 5)
	viper.SetDefault("NAMESPACE", "workloads")
	viper.SetDefault("LABELS", "podchaosmonkey=true")
	labels, err := labels.Parse(viper.GetString("LABELS"))
	if err != nil {
		log.Error("Failed to parse labels: %s", err.Error())
		return nil
	}
	return &PodChaosMonkey{
		Namespace:            viper.GetString("NAMESPACE"),
		Client:               client,
		IntervalInSeconds:    viper.GetInt("INTERVAL_IN_SECONDS"),
		GracePeriodInSeconds: int64(viper.GetInt("GRACE_PERIOD_SECONDS")),
		Labels:               labels,
	}
}

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

func (p *PodChaosMonkey) getAndFilterPods() []v1.Pod {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	podList, err := p.Client.CoreV1().Pods(p.Namespace).List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &p.GracePeriodInSeconds})
	if err != nil {
		panic(err.Error())
	}
	filteredPods := filterPodsByLabels(podList.Items, p.Labels)
	log.Debugf("Total %d pods, %d match the label selector", len(podList.Items), len(filteredPods))
	return filteredPods
}

func filterPodsByLabels(pods []v1.Pod, labelSelector labels.Selector) []v1.Pod {
	results := []v1.Pod{}

	for _, pod := range pods {
		selector := labels.Set(pod.Labels)
		if labelSelector.Matches(selector) {
			results = append(results, pod)
		}
	}
	return results
}

func (p *PodChaosMonkey) deletePods(pods []v1.Pod) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	if len(pods) > 0 {
		log.Infof("Successfully retrieved %d matching pods", len(pods))
		randomPodIndex := rand.Intn(len(pods))
		selectedPod, err := p.Client.CoreV1().Pods(p.Namespace).Get(context.TODO(), pods[randomPodIndex].Name, metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		log.Infof("Deleting pod %s", selectedPod.Name)
		err = p.Client.CoreV1().Pods(p.Namespace).Delete(context.TODO(), selectedPod.Name, metav1.DeleteOptions{GracePeriodSeconds: &p.GracePeriodInSeconds})
		if err != nil {
			panic(err.Error())
		}
	} else {
		log.Infof("No pods matched the label selector")
	}
}

func (p *PodChaosMonkey) Start(ctx context.Context) {
	for {
		pods := p.getAndFilterPods()
		p.deletePods(pods)
		time.Sleep(time.Duration(p.IntervalInSeconds) * time.Second)
	}
}
