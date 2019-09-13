package checkrun

import (
	"bytes"
	"elb2c/gh-action/config"
	"elb2c/gh-action/http/httpconst"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//go:generate mockgen -package=checkrun -self_package=elb2c/gh-action/api/checkrun -destination=mock_updater.go elb2c/gh-action/api/checkrun Updater

type Updater interface {
	Update(checkID int, annotations []Annotation) error
}

type UpdateAPI struct {
	client  *http.Client
	baseURL string
	config  *config.Config
}

type UpdateRequestBody struct {
	Name        string `json:"name"`
	SHA         string `json:"head_sha"`
	Status      string `json:"status"`
	CompletedAt string `json:"completed_at"`
	Conclusion  string `json:"conclusion"`
	Output      Output `json:"output"`
}

type Output struct {
	Title       string       `json:"title"`
	Summary     string       `json:"summary"`
	Annotations []Annotation `json:"annotations"`
}

type Annotation struct {
	Title     string `json:"title"`
	Path      string `json:"path"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
	Level     string `json:"annotation_level"`
	Message   string `json:"message"`
}

func NewUpdater(client *http.Client, URL string, cfg *config.Config) Updater {
	return &UpdateAPI{
		client:  client,
		baseURL: URL,
		config:  cfg,
	}
}

func (self *UpdateAPI) Update(checkID int, annotations []Annotation) error {
	if checkID == 0 {
		return errors.New("Invalid check ID")
	}

	if annotations == nil {
		return errors.New("Annotation array must not be nil")
	}

	if self.client == nil {
		return errors.New("HTTP client must not be nil")
	}

	if self.baseURL == "" {
		return errors.New("Base URL must not be empty")
	}

	if self.config == nil {
		return errors.New("Config must not be nil")
	}

	if self.config.GitHub == (config.GitHub{}) {
		return errors.New("GitHub config must not be empty")
	}

	if self.config.GitHub.Repository == "" {
		return errors.New("GitHub repository must not be empty")
	}

	if self.config.GitHub.Token == "" {
		return errors.New("GitHub token must not be empty")
	}

	if self.config.GitHub.SHA == "" {
		return errors.New("GitHub SHA must not be empty")
	}

	URL := fmt.Sprintf("%s/repos/%s/check-runs/%d", self.baseURL, self.config.GitHub.Repository, checkID)
	method := httpconst.MethodPatch
	header := makeHeaders(self.config.GitHub.Token)
	body := self.makeBody(annotations)
	OKStatusCode := 200

	_, err := submit(self.client, URL, method, header, body, OKStatusCode)

	return err
}

func (self *UpdateAPI) makeBody(annotations []Annotation) *bytes.Buffer {
	req := UpdateRequestBody{
		Name:        nameOfCheckRun,
		SHA:         self.config.GitHub.SHA,
		Status:      "completed",
		CompletedAt: time.Now().UTC().Format(checkRunDateFormat),
		Conclusion:  self.determineConclusion(annotations),
		Output: Output{
			Title:       "Test failure details",
			Summary:     fmt.Sprintf("%d test failure(s) found", len(annotations)),
			Annotations: annotations,
		},
	}
	reqJSON, _ := json.Marshal(req)

	return bytes.NewBuffer(reqJSON)
}

func (self *UpdateAPI) determineConclusion(annotations []Annotation) string {
	if len(annotations) > 0 {
		return "failure"
	}

	return "success"
}
