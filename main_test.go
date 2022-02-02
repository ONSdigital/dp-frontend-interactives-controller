package main

import (
	"context"
	"flag"
	"github.com/ONSdigital/dp-frontend-interactives-controller/features/steps"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var (
	componentFlag = flag.Bool("component", false, "perform component tests")
	doNothing     = func() {}
)

type ComponentTest struct {
	testingT *testing.T
}

func (c *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	component, err := steps.NewComponent()
	if err != nil {
		panic(err)
	}

	ctx.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
		component.Reset()
		return ctx, nil
	})
	ctx.After(func(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
		e := component.Close()
		return ctx, e
	})

	component.RegisterSteps(ctx)
}

func (c *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(doNothing)
	ctx.AfterSuite(doNothing)
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{testingT: t}

		godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()
	} else {
		t.Skip("component flag required to run component tests")
	}
}
