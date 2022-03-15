package main

import (
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
	t            *testing.T
	MongoFeature *cmptest.MongoFeature
}

func init() {
	dplogs.Namespace = "dp-cantabular-filter-flex-api"
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	authFeature := cmptest.NewAuthorizationFeature()
	zebedeeURL := authFeature.FakeAuthService.ResolveURL("")
	mongoAddr := f.MongoFeature.Server.URI()

	component, err := steps.NewComponent(f.t, zebedeeURL, mongoAddr)
	if err != nil {
		log.Panicf("unable to create component: %s", err)
	}

	if _, err := component.Init(); err != nil {
		log.Panicf("unable to initialize component: %s", err)
	}

	apiFeature := cmptest.NewAPIFeature(component.Init)
	component.ApiFeature = apiFeature
	ctx.BeforeScenario(func(*godog.Scenario) {
		apiFeature.Reset()
		if err := f.MongoFeature.Reset(); err != nil {
			log.Panicf("failed to reset mongo feature: %s", err)
		}
		if err := component.Reset(); err != nil {
			log.Panicf("unable to initialise scenario: %s", err)
		}
		authFeature.Reset()
	})

	ctx.AfterScenario(func(*godog.Scenario, error) {
		component.Reset()
		component.Close()
		authFeature.Close()
	})

	authFeature.RegisterSteps(ctx)
	apiFeature.RegisterSteps(ctx)
	component.RegisterSteps(ctx)

}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		f.MongoFeature = cmptest.NewMongoFeature(cmptest.MongoOptions{
			MongoVersion: mongoVersion,
			DatabaseName: databaseName,
		})
	})
	ctx.AfterSuite(func() {
		if err := f.MongoFeature.Close(); err != nil {
			log.Printf("failed to close mongo feature: %s", err)
		}
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
