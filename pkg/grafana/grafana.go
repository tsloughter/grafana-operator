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

type DashboardsInterface interface {
	Search() ([]GrafanaDashboard, error)
	Create(dashboardJson io.Reader) error
	Delete(slug string) error
}

type APIClient struct {
	DashboardsClient DashboardsClient
}

type DashboardsClient struct {
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

func NewAPIClient(baseUrl *url.URL, c *http.Client) APIClient {
	return APIClient{
		DashboardsClient{
			BaseUrl:    baseUrl,
			HTTPClient: c,
		}}
}

func NewDashboardsClient(baseUrl *url.URL, c *http.Client) DashboardsInterface {
	return &DashboardsClient{
		BaseUrl:    baseUrl,
		HTTPClient: c,
	}
}

func (c *DashboardsClient) Search() ([]GrafanaDashboard, error) {
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

func (c *DashboardsClient) Delete(slug string) error {
	deleteUrl := makeUrl(c.BaseUrl, "/api/dashboards/db/"+slug)
	req, err := http.NewRequest("DELETE", deleteUrl, nil)
	if err != nil {
		return err
	}

	return doRequest(c.HTTPClient, req)
}

func (c *DashboardsClient) Create(dashboardJson io.Reader) error {
	importDashboardUrl := makeUrl(c.BaseUrl, "/api/dashboards/import")
	req, err := http.NewRequest("POST", importDashboardUrl, dashboardJson)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	return doRequest(c.HTTPClient, req)
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

type Interface interface {
	Dashboards() DashboardsInterface
}

type Clientset struct {
	BaseUrl    *url.URL
	HTTPClient *http.Client
}

func New(baseUrl *url.URL) *DashboardsClient {
	return &DashboardsClient{
		BaseUrl:    baseUrl,
		HTTPClient: http.DefaultClient,
	}
}

func Newer(baseUrl *url.URL) *APIClient {
	return &APIClient{
		DashboardsClient{
			BaseUrl:    baseUrl,
			HTTPClient: http.DefaultClient,
		}}
}

func (c *Clientset) Dashboards() DashboardsInterface {
	return NewDashboardsClient(c.BaseUrl, c.HTTPClient)
}

func makeUrl(baseUrl *url.URL, endpoint string) string {
	result := &url.URL{}
	*result = *baseUrl

	result.Path = path.Join(result.Path, endpoint)

	return result.String()
}
