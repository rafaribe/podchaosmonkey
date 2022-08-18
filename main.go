// Note: the example only works with the code within the same release/branch.
package main

import (
	"context"

	"go.uber.org/zap"

	chaos "podchaosmonkey/pkg/chaos"
	env "podchaosmonkey/pkg/environment"
)

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	log.Info("Starting pod chaos monkey")
	env.LoadEnv()
	client := chaos.InitKubernetesClient()
	podChaosMonkey := chaos.NewPodChaosMonkey(client)
	podChaosMonkey.Start(context.Background())
}
