package main

import (
	"fmt"

	"github.com/mxschmitt/playwright-go"
)

// App ...
type App struct {
	pw      *playwright.Playwright
	browser *playwright.Browser
	cache   map[string][]byte
}

// NewApp creates an application
func NewApp() App {
	return App{cache: make(map[string][]byte)}
}

// Start boots the app dependencies
func (app *App) Start() error {
	pw, err := playwright.Run()
	if err != nil {
		return fmt.Errorf("could not launch playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		return fmt.Errorf("could not launch chromium: %v", err)
	}

	app.pw = pw
	app.browser = browser
	return nil
}

// Stop shuts down the app dependencies.
func (app *App) Stop() error {
	if err := app.browser.Close(); err != nil {
		return fmt.Errorf("could not close browser: %v", err)
	}

	if err := app.Stop(); err != nil {
		return fmt.Errorf("could not stop playwright: %v", err)
	}

	return nil
}

// Get retrieves a rendered image from the application cache
func (app *App) Get(key string) ([]byte, bool) {
	image, ok := app.cache[key]
	return image, ok
}

// Store caches the rendered image
func (app *App) Store(key string, image []byte) {
	app.cache[key] = image
}

// RenderContent generates a PNG image from a string of html
func (app *App) RenderContent(filename, content string) error {
	page, err := app.browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not creae page: %v", err)
	}

	opts := playwright.PageSetContentOptions{WaitUntil: playwright.String("networkidle")}
	if err = page.SetContent(content, opts); err != nil {
		return fmt.Errorf("could not set page content: %v", err)
	}

	return app.screenshot(filename, page)
}

// RenderPage generates a PNG image from a web page
func (app *App) RenderPage(filename, url string) error {
	page, err := app.browser.NewPage()
	if err != nil {
		return fmt.Errorf("could not creae page: %v", err)
	}

	opts := playwright.PageGotoOptions{WaitUntil: playwright.String("networkidle")}
	if _, err = page.Goto(url, opts); err != nil {
		return fmt.Errorf("could not goto url: %v", err)
	}

	return app.screenshot(filename, page)
}

func (app *App) screenshot(filename string, page *playwright.Page) error {
	targetHandle, err := page.QuerySelector("#screenshot-target")
	if err != nil {
		return fmt.Errorf("could not get target handle: %v", err)
	}

	scrollWidthHandle, err := targetHandle.GetProperty("scrollWidth")
	if err != nil {
		return fmt.Errorf("could not get scrollWidth handle: %v", err)
	}
	scrollWidth, err := scrollWidthHandle.JSONValue()
	if err != nil {
		return fmt.Errorf("could not get scrollWidth value: %v", err)
	}

	scrollHeightHandle, err := targetHandle.GetProperty("scrollHeight")
	if err != nil {
		return fmt.Errorf("could not get scrollHeight handle: %v", err)
	}
	scrollHeight, err := scrollHeightHandle.JSONValue()
	if err != nil {
		return fmt.Errorf("could not get scrollHeight value: %v", err)
	}

	sOpts := playwright.PageScreenshotOptions{
		Path:     playwright.String(filename),
		FullPage: playwright.Bool(true),
		Clip: &playwright.PageScreenshotClip{
			X:      playwright.Int(0),
			Y:      playwright.Int(0),
			Width:  playwright.Int(scrollWidth.(int)),
			Height: playwright.Int(scrollHeight.(int)),
		},
	}
	_, err = page.Screenshot(sOpts)
	if err != nil {
		return fmt.Errorf("could not create screenshot: %v", err)
	}

	return nil
}
