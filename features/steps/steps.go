package steps

import (
	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	//c.apiFeature.RegisterSteps(ctx)
	// ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)
	// ctx.Step(`^the service starts`, c.theServiceStarts)
}

// theServiceStarts starts the service under test in a new go-routine
// note that this step should be called only after all dependencies have been setup,
// to prevent any race condition, specially during the first healthcheck iteration.
// func (c *Component) theServiceStarts() error {
// 	log.Info(c.ctx, "000000")
// 	c.wg.Add(1)
// 	go c.startService(c.ctx)
// 	return nil
// }

// func (c *Component) iShouldReceiveAHelloworldResponse() error {
// 	/*  TODO - implement correct tests once service is ready */
// 	responseBody := c.apiFeature.HttpResponse.Body
// 	body, _ := ioutil.ReadAll(responseBody)

// 	assert.Equal(c, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))
// 	return c.StepError()
// }
