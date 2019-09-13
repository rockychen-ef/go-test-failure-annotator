package service

import (
	"elb2c/gh-action/api/checkrun"
	"elb2c/gh-action/config"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_passing_test_failure_report_returns_no_error_and_creates_and_updates_an_check_run(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	cfg := config.Config{
		TestResultFile: "test_report.xml",
	}

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(&cfg, parserMock, creatorMock, updaterMock)

	checkID := 1
	creatorMock.EXPECT().Create().Return(checkID, nil)

	var failures []TestFailure
	failures = append(failures, TestFailure{
		Line:   1,
		File:   "error_test.go",
		Name:   "Test_passing_an_error_returns_an_error",
		Reason: "Because of errors",
	})
	failures = append(failures, TestFailure{
		Line:   10,
		File:   "blender_test.go",
		Name:   "Test_passing_an_apple_returns_an_apple_juice",
		Reason: "Because of blender",
	})
	parserMock.EXPECT().Parse(gomock.Any()).Return(failures, nil)

	var annotations []checkrun.Annotation
	annotations = append(annotations, checkrun.Annotation{
		Title:     "Test_passing_an_error_returns_an_error",
		Path:      "error_test.go",
		StartLine: 1,
		EndLine:   1,
		Level:     "failure",
		Message:   "Because of errors",
	})
	annotations = append(annotations, checkrun.Annotation{
		Title:     "Test_passing_an_apple_returns_an_apple_juice",
		Path:      "blender_test.go",
		StartLine: 10,
		EndLine:   10,
		Level:     "failure",
		Message:   "Because of blender",
	})
	updaterMock.EXPECT().Update(checkID, annotations).Return(nil)

	err := svc.Annotate()

	assert.NoError(test, err)
}

func Test_passing_no_test_failure_report_returns_no_error_and_creates_and_updates_an_check_run(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	cfg := config.Config{
		TestResultFile: "test_report.xml",
	}

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(&cfg, parserMock, creatorMock, updaterMock)

	checkID := 1
	creatorMock.EXPECT().Create().Return(checkID, nil)

	var failures []TestFailure
	parserMock.EXPECT().Parse(gomock.Any()).Return(failures, nil)

	emptyAnnotations := make([]checkrun.Annotation, 0)
	updaterMock.EXPECT().Update(checkID, emptyAnnotations).Return(nil)

	err := svc.Annotate()

	assert.NoError(test, err)
}

func Test_missing_test_report_config_returns_an_error(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(nil, parserMock, creatorMock, updaterMock)

	err := svc.Annotate()

	assert.Error(test, err)
	creatorMock.EXPECT().Create().Times(0)
	parserMock.EXPECT().Parse(gomock.Any()).Times(0)
	updaterMock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
}

func Test_assuming_failed_to_parse_test_report_returns_no_error(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	cfg := config.Config{
		TestResultFile: "test_report.xml",
	}

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(&cfg, parserMock, creatorMock, updaterMock)

	checkID := 1
	creatorMock.EXPECT().Create().Return(checkID, nil)

	parseFailed := errors.New("Failed to parse the test report")
	parserMock.EXPECT().Parse(gomock.Any()).Return(nil, parseFailed)

	emptyAnnotations := make([]checkrun.Annotation, 0)
	updaterMock.EXPECT().Update(checkID, emptyAnnotations).Return(nil)

	err := svc.Annotate()

	assert.NoError(test, err)
}

func Test_assuming_failed_to_create_a_check_run_returns_an_error(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	cfg := config.Config{
		TestResultFile: "test_report.xml",
	}

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(&cfg, parserMock, creatorMock, updaterMock)

	createFailed := errors.New("Failed to create a check run")
	creatorMock.EXPECT().Create().Return(0, createFailed)

	err := svc.Annotate()

	assert.Error(test, err)
	parserMock.EXPECT().Parse(gomock.Any()).Times(0)
	updaterMock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
}

func Test_assuming_failed_to_update_a_check_run_returns_an_error(test *testing.T) {
	mockCtl := gomock.NewController(test)
	defer mockCtl.Finish()

	cfg := config.Config{
		TestResultFile: "test_report.xml",
	}

	parserMock := NewMockTestResultParser(mockCtl)
	creatorMock := checkrun.NewMockCreator(mockCtl)
	updaterMock := checkrun.NewMockUpdater(mockCtl)
	svc := NewTestFailureAnnotator(&cfg, parserMock, creatorMock, updaterMock)

	checkID := 1
	creatorMock.EXPECT().Create().Return(checkID, nil)

	var failures []TestFailure
	failures = append(failures, TestFailure{
		Line:   1,
		File:   "error_test.go",
		Name:   "Test_passing_an_error_returns_an_error",
		Reason: "Because of errors",
	})
	parserMock.EXPECT().Parse(gomock.Any()).Return(failures, nil)

	updateFailed := errors.New("Failed to update a check run")
	updaterMock.EXPECT().Update(checkID, gomock.Any()).Return(updateFailed)

	err := svc.Annotate()

	assert.Error(test, err)
}
