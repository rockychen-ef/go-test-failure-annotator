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

//go:generate mockgen -package=checkrun -self_package=elb2c/gh-action/api/checkrun -destination=mock_creator.go elb2c/gh-action/api/checkrun Creator

type Creator interface {
	Create() (int, error)
}

type CreationAPI struct {
	client  *http.Client
	baseURL string
	config  *config.Config
}

type CreationRequestBody struct {
	Name      string `json:"name"`
	SHA       string `json:"head_sha"`
	Status    string `json:"status"`
	StartedAt string `json:"started_at"`
}

func NewCreator(client *http.Client, URL string, cfg *config.Config) Creator {
	return &CreationAPI{
		client:  client,
		baseURL: URL,
		config:  cfg,
	}
}

func (self *CreationAPI) Create() (int, error) {
	if self.client == nil {
		return 0, errors.New("HTTP client must not be nil")
	}

	if self.baseURL == "" {
		return 0, errors.New("Base URL must not be empty")
	}

	if self.config == nil {
		return 0, errors.New("Config must not be nil")
	}

	if self.config.GitHub == (config.GitHub{}) {
		return 0, errors.New("GitHub config must not be empty")
	}

	if self.config.GitHub.Repository == "" {
		return 0, errors.New("GitHub repository must not be empty")
	}

	if self.config.GitHub.Token == "" {
		return 0, errors.New("GitHub token must not be empty")
	}

	if self.config.GitHub.SHA == "" {
		return 0, errors.New("GitHub SHA must not be empty")
	}

	URL := fmt.Sprintf("%s/repos/%s/check-runs", self.baseURL, self.config.GitHub.Repository)
	method := httpconst.MethodPost
	header := makeHeaders(self.config.GitHub.Token)
	body := self.makeBody()
	OKStatusCode := 201

	resp, err := submit(self.client, URL, method, header, body, OKStatusCode)
	if err != nil {
		return 0, err
	}

	return int(resp["id"].(float64)), nil
}

func (self *CreationAPI) makeBody() *bytes.Buffer {
	req := CreationRequestBody{
		Name:      nameOfCheckRun,
		SHA:       self.config.GitHub.SHA,
		Status:    "in_progress",
		StartedAt: time.Now().UTC().Format(checkRunDateFormat),
	}
	reqJSON, _ := json.Marshal(req)

	return bytes.NewBuffer(reqJSON)
}
