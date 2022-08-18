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
)

type PodChaosMonkey struct {
	Client               kubernetes.Interface
	Namespace            string
	IntervalInSeconds    int
	GracePeriodInSeconds int64
	Labels               labels.Selector
	IncludeFinalizers    bool
}

func NewPodChaosMonkey(client kubernetes.Interface) *PodChaosMonkey {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
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
		IncludeFinalizers:    viper.GetBool("INCLUDE_FINALIZERS"),
	}
}

func (p *PodChaosMonkey) filterPods(pods []v1.Pod) []v1.Pod {
	results := filterPodsByLabels(pods, p.Labels)
	results = filterPodsByState(results, v1.PodRunning, p.IncludeFinalizers)
	return results
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

func filterPodsByState(pods []v1.Pod, phase v1.PodPhase, includeFinalizers bool) []v1.Pod {
	results := []v1.Pod{}

	for _, pod := range pods {
		if pod.Status.Phase == phase {
			// Pods that have finalizers should not be deleted, unless we force it with the includeVariables config
			if pod.ObjectMeta.Finalizers == nil || includeFinalizers {
				results = append(results, pod)
			}
		}
	}
	return results
}

func (p *PodChaosMonkey) DeletePods(pods []v1.Pod) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	if len(pods) > 0 {
		log.Infof("Successfully retrieved %d matching pods", len(pods))
		randomPodIndex := rand.Intn(len(pods))
		selectedPod, err := p.Client.CoreV1().Pods(p.Namespace).Get(context.TODO(), pods[randomPodIndex].ObjectMeta.Name, metav1.GetOptions{})
		
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
func (p *PodChaosMonkey) FilterPods(ctx context.Context) []v1.Pod {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	podList, err := p.Client.CoreV1().Pods(p.Namespace).List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &p.GracePeriodInSeconds})
	if err != nil {
		log.Warnf("No pods found in namespace %s", p.Namespace)
	}

	filteredPods := p.filterPods(podList.Items)
	log.Debugf("Total %d pods, %d match the label selector", len(podList.Items), len(filteredPods))
	return filteredPods
}

func (p *PodChaosMonkey) Start(ctx context.Context) {
	for {
		elegiblePods := p.FilterPods(ctx)
		p.DeletePods(elegiblePods)
		time.Sleep(time.Duration(p.IntervalInSeconds) * time.Second)
	}
}
