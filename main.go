package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/go-co-op/gocron"
	"github.com/nkzren/ecos/config"
	"github.com/nkzren/ecos/kube"
	"github.com/nkzren/ecos/score"
	"github.com/nkzren/ecos/weather"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
}

func main() {
	sigs, done, config := setup()
	defer waitForExit(sigs, done)

	fmt.Println("Start scheduling for node labeling")
	go startScheduling(config.Scheduler, func() {
		labelNodes(config.Kube)
	})
	go startController()
}

func startController() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "a9806571.nkzren",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&kube.DeploymentReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Deployment")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

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
