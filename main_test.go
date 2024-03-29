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
	dplogs "github.com/ONSdigital/log.go/v2/log"
)

const componentLogFile = "component-output.txt"

var componentFlag = flag.Bool("component", false, "perform component tests")
var quietComponentFlag = flag.Bool("quiet-component", false, "perform component tests with dp logging disabled")

type ComponentTest struct {
	component *steps.Component
}

func (ct *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		ct.component.Reset()
		return ctx, nil
	})
}

func (ct *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ct.component.RegisterSteps(ctx.ScenarioContext())

	ctx.AfterSuite(func() {
		ct.component.Close()
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

		ct := &ComponentTest{
			component: steps.NewComponent(t),
		}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  ct.InitializeScenario,
			TestSuiteInitializer: ct.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
