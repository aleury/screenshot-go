package main

import (
	"io/ioutil"
	"log"

	"github.com/mxschmitt/playwright-go"
)

func main() {
	content, err := ioutil.ReadFile("test.html")
	if err != nil {
		log.Fatalf("could not load html file: %v", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Fatalf("could not launch playwright: %v", err)
	}

	browser, err := pw.Chromium.Launch()
	if err != nil {
		log.Fatalf("could not launch Chromium: %v", err)
	}

	page, err := browser.NewPage()
	if err != nil {
		log.Fatalf("could not creae page: %v", err)
	}

	scOpts := playwright.PageSetContentOptions{WaitUntil: playwright.String("networkidle")}
	if err = page.SetContent(string(content), scOpts); err != nil {
		log.Fatalf("could not set page content: %v", err)
	}

	targetHandle, err := page.QuerySelector("#screenshot-target")
	if err != nil {
		log.Fatalf("could not target handle: %v", err)
	}

	scrollWidthHandle, err := targetHandle.GetProperty("scrollWidth")
	if err != nil {
		log.Fatalf("could not get scrollWidth handle: %v", err)
	}
	scrollWidth, err := scrollWidthHandle.JSONValue()
	if err != nil {
		log.Fatalf("could not get scrollWidth value: %v", err)
	}

	scrollHeightHandle, err := targetHandle.GetProperty("scrollHeight")
	if err != nil {
		log.Fatalf("could not get scrollHeight handle: %v", err)
	}
	scrollHeight, err := scrollHeightHandle.JSONValue()
	if err != nil {
		log.Fatalf("could not get scrollHeight value: %v", err)
	}

	sOpts := playwright.PageScreenshotOptions{
		Path:     playwright.String("test.png"),
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
		log.Fatalf("could not create screenshot: %v", err)
	}

	if err = browser.Close(); err != nil {
		log.Fatalf("could not close browser: %v", err)
	}

	if err = pw.Stop(); err != nil {
		log.Fatalf("could not stop playwright: %v", err)
	}
}
