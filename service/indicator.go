package service

import (
	"elb2c/gh-action/api/checkrun"
	"elb2c/gh-action/config"
	"errors"
	"log"
)

type TestFailureAnnotator interface {
	Annotate() error
}

type TestFailureAnnotateService struct {
	config          *config.Config
	parser          TestResultParser
	checkRunCreator checkrun.Creator
	checkRunUpdater checkrun.Updater
}

func NewTestFailureAnnotator(cfg *config.Config, parser TestResultParser,
	creator checkrun.Creator, updater checkrun.Updater) TestFailureAnnotator {

	return &TestFailureAnnotateService{
		config:          cfg,
		parser:          parser,
		checkRunCreator: creator,
		checkRunUpdater: updater,
	}
}

func (self *TestFailureAnnotateService) Annotate() error {
	if self.config == nil {
		return errors.New("Config must not be nil")
	}

	// Create a check run
	ID, err := self.checkRunCreator.Create()
	if err != nil {
		return err
	}

	// Parser test results
	failures, err := self.parser.Parse(self.config.TestResult())
	if err != nil {
		log.Printf("Failed to parser the test report because: %s\n", err)
	}

	// Covert test failures to GitHub annotations
	annotations := make([]checkrun.Annotation, 0)
	for _, failure := range failures {
		annotation := checkrun.Annotation{
			Title:     failure.Name,
			Path:      failure.File,
			StartLine: failure.Line,
			EndLine:   failure.Line,
			Level:     "failure",
			Message:   failure.Reason,
		}

		annotations = append(annotations, annotation)
	}

	// Complete the check run
	if err := self.checkRunUpdater.Update(ID, annotations); err != nil {
		return err
	}

	return nil
}
