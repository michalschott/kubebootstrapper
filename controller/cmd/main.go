package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/michalschott/kubebootstrapper/pkg/version"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type controllerConfig struct {
	syncPeriod time.Duration
	server     string
	channel    string
}

type channelInfo struct {
	server  string
	Channel struct {
		Path   string
		Shasum string
		Update string
	}
}

func main() {
	logLevel := flag.String("l", "info", "Log level [default: info]")
	syncPeriod := flag.String("s", "15s", "Control loop sync period [default: 15s]")
	server := flag.String("server", "http://http.kubebootstrapperserver:8000", "Where to fetch manifests from")
	channel := flag.String("channel", "stable", "Channel to use [default: stable]")

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
	log.Debug("-server ", *server)
	log.Debug("-channel ", *channel)

	syncPeriodParsed, err := time.ParseDuration(*syncPeriod)
	if err != nil {
		log.Fatal(err.Error())
	}

	config := &controllerConfig{
		syncPeriod: syncPeriodParsed,
		server:     *server,
		channel:    *channel,
	}

	go syncLoop(config)
	go controlLoop(config)

	<-stop
	log.Info("Shutting down")
	os.Exit(0)
}

func controlLoop(config *controllerConfig) {
	for {
		log.Debug("Starting control loop")
		log.Debug("Finished control loop")
		time.Sleep(config.syncPeriod)
	}
}

func syncLoop(config *controllerConfig) {
	for {
		log.Debug("Starting sync loop")
		channelInfo, err := fetchChannelManifest(config)
		if err != nil {
			log.Info(err.Error())
		} else {
			channelData, err := fetchChannelData(channelInfo)
			if err != nil {
				log.Info(err.Error())
			} else {
				log.Info(channelData)
			}
		}
		log.Debug("Finished sync loop")
		time.Sleep(config.syncPeriod)
	}
}

func fetchChannelManifest(config *controllerConfig) (*channelInfo, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(config.server + "/" + config.channel + ".yaml")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	channel := channelInfo{}

	err = yaml.Unmarshal([]byte(body), &channel)
	if err != nil {
		return nil, err
	}

	channel.server = config.server

	return &channel, nil
}

func fetchChannelData(info *channelInfo) ([]byte, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get(info.server + info.Channel.Path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	return body, nil
}
