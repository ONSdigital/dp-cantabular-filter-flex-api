package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/steps"
	cmptest "github.com/ONSdigital/dp-component-test"
	dplogs "github.com/ONSdigital/log.go/v2/log"
)

const (
	mongoVersion     = "4.4.8"
	databaseName     = "filters"
	componentLogFile = "component-output.txt"
)

var componentFlag = flag.Bool("component", false, "perform component tests")
var quietComponentFlag = flag.Bool("quiet-component", false, "perform component tests with dp logging disabled")

type ComponentTest struct {
	t *testing.T

	mongoFeature *cmptest.MongoFeature
	authFeature  *cmptest.AuthorizationFeature

	component *steps.Component
}

func init() {
	dplogs.Namespace = "dp-cantabular-filter-flex-api"
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	f.authFeature.RegisterSteps(ctx)
	f.component.RegisterSteps(ctx)

	ctx.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		f.authFeature.Reset()
		e := f.component.Reset()
		if e != nil {
			log.Printf("failed to reset component: %s", err)
		}
		return ctx, e
	})
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	var err error

	f.mongoFeature = cmptest.NewMongoFeature(cmptest.MongoOptions{
		MongoVersion: mongoVersion,
		DatabaseName: databaseName,
	})
	f.authFeature = cmptest.NewAuthorizationFeature()

	f.component, err = steps.NewComponent(f.t, f.authFeature.FakeAuthService.ResolveURL(""), f.mongoFeature.Server.URI())
	if err != nil {
		log.Panicf("unable to create component: %s", err)
	}
	f.component.ApiFeature = cmptest.NewAPIFeature(f.component.Init)

	ctx.AfterSuite(func() {
		if err := f.mongoFeature.Close(); err != nil {
			log.Printf("failed to close mongo feature: %s", err)
		}
		f.authFeature.Close()
		f.component.Close()
	})
}

func TestComponent(t *testing.T) {
	if *componentFlag || *quietComponentFlag {
		status := 0

		cfg, err := config.Get()
		if err != nil {
			t.Fatalf("failed to get service config: %s", err)
		}

		var output io.Writer = os.Stdout

		if cfg.ComponentTestUseLogFile {
			logfile, err := os.Create(componentLogFile)
			if err != nil {
				t.Fatalf("could not create logs file: %s", err)
			}

			defer func() {
				if err := logfile.Close(); err != nil {
					log.Printf("failed to close log file: %s", err)
				}
			}()

			output = logfile
			dplogs.SetDestination(logfile, nil)
		}

		if *quietComponentFlag {
			dplogs.SetDestination(io.Discard, io.Discard)
		}

		var opts = godog.Options{
			Output:   colors.Colored(output),
			Format:   "pretty",
			Paths:    flag.Args(),
			TestingT: t,
		}

		f := &ComponentTest{
			t: t,
		}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
