package steps

import (
	"github.com/cucumber/godog"
)

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^I should receive a hello-world response$`, c.iShouldReceiveAHelloworldResponse)
}

func (c *Component) iShouldReceiveAHelloworldResponse() error {
	/*  TODO - implement correct tests once service is ready
	responseBody := c.apiFeature.HttpResponse.Body
	body, _ := ioutil.ReadAll(responseBody)

	assert.Equal(c, `{"message":"Hello, World!"}`, strings.TrimSpace(string(body)))
	*/
	return c.StepError()

}
