package service

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

//go:generate mockgen -package=service -self_package=elb2c/gh-action/service -destination=mock_parser.go elb2c/gh-action/service TestResultParser

type TestResultParser interface {
	Parse(testResult string) ([]TestFailure, error)
}

type TestFailure struct {
	Line   int
	File   string
	Name   string
	Reason string
}

type testSuites struct {
	XMLName    xml.Name    `xml:"testsuites"`
	TestSuites []testSuite `xml:"testsuite"`
}

type testSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	TotalTests int        `xml:"tests,attr"`
	Failures   int        `xml:"failures,attr"`
	Time       float64    `xml:"time,attr"`
	Name       string     `xml:"name,attr"`
	TestCases  []testCase `xml:"testcase"`
}

type testCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      float64  `xml:"time,attr"`
	Details   string   `xml:"failure"`
}

type TestResultParseService struct {
}

var (
	regexErrorTrace  = regexp.MustCompile(`Error Trace:(\s+)(\w+\.\w+)\:(\d+)`)
	regexErrorDtails = regexp.MustCompile(`Error:(\s+)(.*)(\s.*)(\s.*)`)
)

func NewTestResultParser() TestResultParser {
	return &TestResultParseService{}
}

func (self *TestResultParseService) Parse(testResult string) ([]TestFailure, error) {
	if testResult == "" {
		return nil, errors.New("TestResult must not be empty")
	}

	xmlFile, err := os.Open(testResult)
	if err != nil {
		return nil, err
	}
	log.Printf("Successful open the test result file: %s\n", testResult)
	defer xmlFile.Close()

	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var testsuites testSuites
	xml.Unmarshal(byteValue, &testsuites)

	var failures []TestFailure
	failedTestCases := self.filterFailedTestCases(testsuites.TestSuites)
	for _, testCase := range failedTestCases {
		failure, err := self.buildFailure(testCase.Details)
		if err != nil {
			return nil, err
		}
		failure.Name = testCase.Name
		failure.File = self.getDirectory(testCase.ClassName) + "/" + failure.File
		failures = append(failures, *failure)
	}

	return failures, nil
}

func (self *TestResultParseService) filterFailedTestCases(testsuites []testSuite) (result []testCase) {
	for i := 0; i < len(testsuites); i++ {
		if testsuites[i].Failures > 0 {
			for j := 0; j < len(testsuites[i].TestCases); j++ {
				if len(testsuites[i].TestCases[j].Details) > 0 {
					result = append(result, testsuites[i].TestCases[j])
				}
			}
		}
	}

	return
}

func (self *TestResultParseService) buildFailure(details string) (*TestFailure, error) {
	lineNumber, err := self.findLineNumber(details)
	if err != nil {
		return nil, err
	}

	fileName, err := self.findFileName(details)
	if err != nil {
		return nil, err
	}

	reason, err := self.findReason(details)
	if err != nil {
		return nil, err
	}

	return &TestFailure{
		Line:   lineNumber,
		File:   fileName,
		Reason: reason,
	}, nil
}

func (self *TestResultParseService) getDirectory(className string) string {
	array := strings.Split(className, "/")
	return array[len(array)-1]
}

func (self *TestResultParseService) findFileName(details string) (string, error) {
	const targetIndex = 2
	match := regexErrorTrace.FindStringSubmatch(details)
	if len(match) == 0 {
		return "", errors.New("No file name matches")
	}
	return match[targetIndex], nil
}

func (self *TestResultParseService) findLineNumber(details string) (int, error) {
	const targetIndex = 3
	match := regexErrorTrace.FindStringSubmatch(details)
	if len(match) == 0 {
		return 0, errors.New("No line number matches")
	}
	return strconv.Atoi(match[targetIndex])
}

func (self *TestResultParseService) findReason(details string) (string, error) {
	const targetIndex = 0
	match := regexErrorDtails.FindStringSubmatch(details)
	if len(match) == 0 {
		return "", errors.New("No reason matches")
	}
	return match[targetIndex], nil
}
