package main

import (
	"time"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/newrelic/infra-integrations-sdk/log"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Resolve bool   `default:"false" help:"Try to resolve numeric addresses. Default is false."`
	Filter  string `default:"" help:"Properly formatted ss filter string. Default is none."`
	SSArgs  string `default:"-iot" help:"ss command line args."`
}

const (
	integrationName    = "com.newrelic.ss"
	integrationVersion = "1.0.0"
	defaultHTTPTimeout = time.Second * 1
	eventType          = "NetstatV2Sample"
)

var Args argumentList

func main() {
	integration, err := integration.New(integrationName, integrationVersion, integration.Args(&Args))
	fatalIfErr(err)
	log.SetupLogging(Args.Verbose)
	log.Debug("Starting integration %s", integrationName)
	defer log.Debug("Exiting integration %s", integrationName)

	entity := integration.LocalEntity()

	if Args.All() || Args.Metrics {
		log.Debug("Fetching metrics for integration %s", integrationName)
		fatalIfErr(getMetrics(entity, Args))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
