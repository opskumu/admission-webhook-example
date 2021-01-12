package main

import (
	"flag"
	"os"

	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

var (
	scheme = runtime.NewScheme()
	log    = logf.Log.WithName("pod-admission-webhook")
)

func init() {
	logf.SetLogger(zap.New())
}

func main() {
	entryLog := log.WithName("entrypoint")

	var (
		metricsAddr, certDir string
		port                 int
	)

	pflag.IntVar(&port, "port", 9443, "pod-admission-webhook listen port.")
	pflag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	pflag.StringVar(&certDir, "cert-dir", "", "CertDir is the directory that contains the server key and certificate. "+
		"if not set, webhook server would look up the server key and certificate in "+
		"{TempDir}/k8s-webhook-server/serving-certs. The server key and certificate "+
		"must be named tls.key and tls.crt, respectively.")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// Setup a Manager
	entryLog.Info("setting up manager")
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               port,
		CertDir:            certDir,
	})
	if err != nil {
		entryLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Setup webhooks
	entryLog.Info("setting up webhook server")
	hookServer := mgr.GetWebhookServer()

	entryLog.Info("registering webhooks to the webhook server")
	hookServer.Register("/mutate-pod", &webhook.Admission{Handler: &podMutate{Client: mgr.GetClient()}})

	entryLog.Info("starting manager")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
