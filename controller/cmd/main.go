package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michalschott/kubebootstrapper/pkg/version"
	log "github.com/sirupsen/logrus"
)

func main() {
	logLevel := flag.String("l", "info", "Log level [default: info]")
	syncPeriod := flag.String("s", "5s", "Control loop sync period [default: 5s]")

	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	switch *logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	log.Info("Starting kubebootstrapper controller, version: ", version.Version)
	log.Debug("-l ", *logLevel)
	log.Debug("-s ", *syncPeriod)

	syncPeriodParsed, err := time.ParseDuration(*syncPeriod)
	if err != nil {
		log.Fatal(err.Error())
	}

	go controlLoop(syncPeriodParsed)

	<-stop
	log.Info("Shutting down")
	os.Exit(0)
}

func controlLoop(syncPeriod time.Duration) {
	for {
		log.Info("Starting control loop")
		log.Debug("Finished control loop")
		time.Sleep(syncPeriod)
	}
}
