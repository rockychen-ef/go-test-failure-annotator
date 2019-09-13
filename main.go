package main

import (
	"elb2c/gh-action/api/checkrun"
	"elb2c/gh-action/config"
	"elb2c/gh-action/service"
	"fmt"
	"net/http"
	"os"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	parser := service.NewTestResultParser()

	httpClient := &http.Client{}
	checkRunCreator := checkrun.NewCreator(httpClient, cfg.GitHub.URL, &cfg)
	checkRunUpdator := checkrun.NewUpdater(httpClient, cfg.GitHub.URL, &cfg)

	annotator := service.NewTestFailureAnnotator(&cfg, parser, checkRunCreator, checkRunUpdator)
	if err := annotator.Annotate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
