// Copyright 2016 The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grafana

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type APIInterface interface {
	SearchDashboard() ([]GrafanaDashboard, error)
	CreateDashboard(dashboardJson io.Reader) error
	DeleteDashboard(slug string) error
	CreateDatasource(datasourceJson io.Reader) error
}

type APIClient struct {
	BaseUrl    *url.URL
	HTTPClient *http.Client
}

type GrafanaDashboard struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Uri   string `json:"uri"`
}

func (d *GrafanaDashboard) Slug() string {
	// The uri in the search result contains the slug.
	// http://docs.grafana.org/v3.1/http_api/dashboard/#search-dashboards
	return strings.TrimPrefix(d.Uri, "db/")
}

func NewAPIClient(baseURL *url.URL, c *http.Client) APIInterface {
	return &APIClient{
		BaseUrl:    baseURL,
		HTTPClient: c,
	}
}

func (c *APIClient) SearchDashboard() ([]GrafanaDashboard, error) {
	searchUrl := makeUrl(c.BaseUrl, "/api/search")
	resp, err := c.HTTPClient.Get(searchUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	searchResult := make([]GrafanaDashboard, 0)
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func (c *APIClient) DeleteDashboard(slug string) error {
	deleteUrl := makeUrl(c.BaseUrl, "/api/dashboards/db/"+slug)
	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}

	return doRequest(c.HTTPClient, req)
}

func (c *APIClient) CreateDashboard(dashboardJSON io.Reader) error {
	return doPost(makeUrl(c.BaseUrl, "/api/dashboards/import"), dashboardJSON, c.HTTPClient)
}

func (c *APIClient) CreateDatasource(datasourceJSON io.Reader) error {
	return doPost(makeUrl(c.BaseUrl, "/api/datasources"), datasourceJSON, c.HTTPClient)
}

func doPost(url string, dataJSON io.Reader, c *http.Client) error {
	req, err := http.NewRequest("POST", url, dataJSON)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	return doRequest(c, req)
}

func doRequest(c *http.Client, req *http.Request) error {
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Unexpected status code returned from Grafana API (got: %d, expected: 200, msg:%s)", resp.StatusCode, resp.Status)
	}
	return nil
}

type Clientset struct {
	BaseUrl    *url.URL
	HTTPClient *http.Client
}

func New(baseUrl *url.URL) *APIClient {
	return &APIClient{
		BaseUrl:    baseUrl,
		HTTPClient: http.DefaultClient,
	}
}

func makeUrl(baseURL *url.URL, endpoint string) string {
	result := &url.URL{}
	*result = *baseURL

	result.Path = path.Join(result.Path, endpoint)

	return result.String()
}
