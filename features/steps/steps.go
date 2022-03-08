package steps

import (
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

var elementMap = map[string]string{
	"title": "title",
}

func (c *Component) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^The browser loads the iframe$`, c.theBrowserLoadsTheIframe)
	ctx.Step(`^I should see the "([^"]*)" element$`, c.iShouldSeeTheElement)
}

func (c *Component) theBrowserLoadsTheIframe() error {
	err := chromedp.Run(c.chrome.ctx,
		chromedp.Navigate("http://localhost:27300/interactives/abcde123/index.html"),
		chromedp.WaitVisible(`#childt`),
	)

	if err != nil {
		return err
	}

	return nil
}

func (c *Component) iShouldSeeTheElement(elementKey string) error {
	elementSelector := elementMap[elementKey]
	if elementSelector == "" {
		return godog.ErrUndefined
	}

	var res []*cdp.Node
	err := chromedp.Run(c.chrome.ctx,
		chromedp.Nodes(elementSelector, &res, chromedp.AtLeast(0)),
	)

	if err != nil {
		return err
	}

	assert.NotEqual(c, len(res), 0)
	return c.StepError()
}
