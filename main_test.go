package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"testing"

	"github.com/ONSdigital/dp-cantabular-filter-flex-api/config"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/features/steps"
	componenttest "github.com/ONSdigital/dp-component-test"
	dplogs "github.com/ONSdigital/log.go/v2/log"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

const componentLogFile = "component-output.txt"

var componentFlag = flag.Bool("component", false, "perform component tests")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func init() {
	dplogs.Namespace = "dp-cantabular-filter-flex-api"
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	component := steps.NewComponent()
	component.InitService()
	// component.InitProducer()
	component.StartService(context.Background())

	apiFeature := componenttest.NewAPIFeature(component.InitService)

	ctx.BeforeScenario(func(*godog.Scenario) {
		apiFeature.Reset()
		if err := component.Reset(); err != nil {
			log.Panicf("unable to initialise scenario: %s", err)
		}
	})

	ctx.AfterScenario(func(*godog.Scenario, error) {
		component.Close()
	})

	apiFeature.RegisterSteps(ctx)
	component.RegisterSteps(ctx)

}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {

}

func TestComponent(t *testing.T) {
	if *componentFlag {
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

			defer logfile.Close()
			output = logfile

			dplogs.SetDestination(logfile, nil)
		}

		var opts = godog.Options{
			Output: colors.Colored(output),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

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
