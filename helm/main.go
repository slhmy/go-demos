package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/kube"
)

func main() {
	actionConfig := new(action.Configuration)

	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	context := "minikube"
	namespace := "default"
	helmDriver := ""

	clientgetter := kube.GetConfig(kubeconfig, context, namespace)
	if err := actionConfig.Init(clientgetter, namespace, helmDriver, log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = namespace
	client.ReleaseName = "go-demo"
	client.Wait = true
	client.Atomic = true
	client.Timeout = 1 * time.Minute

	chartPath := "go-demo-helm-chart"
	chrt, err := loader.Load(chartPath)
	if err != nil {
		log.Printf("Failed to load chart: %v", err)
		os.Exit(1)
	}

	if _, err := client.Run(chrt, nil); err != nil {
		log.Printf("Failed to install chart: %v", err)
		os.Exit(1)
	}

	// Uninstall the chart
	uninstall := action.NewUninstall(actionConfig)
	if _, err := uninstall.Run("go-demo"); err != nil {
		log.Printf("Failed to uninstall chart: %v", err)
		os.Exit(1)
	}
}
