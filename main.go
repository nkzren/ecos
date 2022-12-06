package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/nkzren/ecoscheduler/config"
	"github.com/nkzren/ecoscheduler/kube"
	"github.com/nkzren/ecoscheduler/score"
)

func main() {
	sigs, done, config := setup()
	defer waitForExit(sigs, done)

	go startScheduling(config.Scheduler, func() {
		labelNodes(config.Kube)
	})
}

func startScheduling(c config.SchedulerConf, task func()) {
	s := gocron.NewScheduler(time.UTC)
	s.Every(c.Interval).Do(task)
	s.StartAsync()
}

func labelNodes(kubeConf config.KubeConf) {
	nodes, err := kube.GetNodes()
	for _, node := range nodes.Items {
		label := node.Labels["ecos"]
		if label == "" {
			fmt.Printf("Ecos' label not found for node (%s), setting to neutral", node.Name)
			kube.UpdateLabel(&node, "neutral")
		} else {
			result := score.GetResult("weather")
			kube.UpdateLabel(&node, result)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}
}

func setup() (chan os.Signal, chan bool, config.Configurations) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan bool, 1)
	config := config.Root
	return sigs, done, config
}

func waitForExit(sigs chan os.Signal, done chan bool) {
	go func() {
		sig := <-sigs
		fmt.Println(sig)
		done <- true
	}()
	fmt.Println("wait for exit signal")
	<-done
	fmt.Println("exiting")
}
