package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/nkzren/ecoscheduler/config"
	"github.com/nkzren/ecoscheduler/kube"
	"github.com/nkzren/ecoscheduler/score"
	"github.com/nkzren/ecoscheduler/weather"

	corev1 "k8s.io/api/core/v1"
)

func main() {
	sigs, done, config := setup()
	defer waitForExit(sigs, done)

	fmt.Println("Start scheduling for node labeling")
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	var wg sync.WaitGroup
	wg.Add(len(nodes.Items))
	for i := 0; i < len(nodes.Items); i++ {
		go func(i int) {
			node := nodes.Items[i]
			label, err := getLabelFor(&node)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			err = kube.UpdateLabel(&node, label)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func getLabelFor(node *corev1.Node) (string, error) {
	cityLabel := node.Labels["city"]
	countryLabel := node.Labels["country"]
	if cityLabel == "" || countryLabel == "" {
		return "neutral", errors.New(fmt.Sprintf("Location labels missing for node (%s), will set to neutral", node.Name))
	}
	result := score.GetResult("weather", weather.Location{City: cityLabel, Country: countryLabel})
	return result, nil
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
