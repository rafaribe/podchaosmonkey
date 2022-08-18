package main

import (
	"context"

	"go.uber.org/zap"

	chaos "podchaosmonkey/pkg/chaos"
	client "podchaosmonkey/pkg/client"
	env "podchaosmonkey/pkg/environment"
)

func main() {

	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()
	log.Info("Starting pod chaos monkey")
	env.LoadEnv()
	client := client.InitKubernetesClient()
	podChaosMonkey := chaos.NewPodChaosMonkey(client)
	ctx := context.Background()
	podChaosMonkey.Start(ctx)
}
