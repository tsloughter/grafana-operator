package test

import (
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tsloughter/grafana-operator/pkg/controller"
	"github.com/tsloughter/grafana-operator/pkg/grafana"
	"github.com/tsloughter/grafana-operator/pkg/kubernetes"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const dashboardJson = `
{
	"dashboard": {
		"description": "Dashboard to get an overview of one server",
		"editable": true,
	}
}
`
const datasourceJson = `
{
	"id": 2,
	"name": "prometheus",
	"type": "prometheus",
	"access": "direct",
	"basicAuth": false,
	"withCredentials": true,
}
`

type APIClientMock struct {
	mock.Mock
}

func (m *APIClientMock) CreateDashboard(dashboardJson io.Reader) error {
	m.Called(dashboardJson)
	return nil
}

func (m *APIClientMock) DeleteDashboard(slug string) error {
	m.Called(slug)
	return nil
}

func (m *APIClientMock) SearchDashboard() ([]grafana.GrafanaDashboard, error) {
	m.Called()
	return nil, nil
}

func (m *APIClientMock) CreateDatasource(datasourceJson io.Reader) error {
	m.Called(datasourceJson)
	return nil
}

func TestCreateDashboards(t *testing.T) {
	apiClient := new(APIClientMock)
	apiClient.On("CreateDashboard", strings.NewReader(dashboardJson)).Return(nil)
	c := newConfigMapController(apiClient)
	cm := grafanaConfigMap(dashboardJson, true, true)
	c.CreateDashboards(cm)
	apiClient.AssertExpectations(t)
}

func TestCreateDatasource(t *testing.T) {
	apiClient := new(APIClientMock)
	apiClient.On("CreateDatasource", strings.NewReader(datasourceJson)).Return(nil)
	c := newConfigMapController(apiClient)
	cm := grafanaConfigMap(datasourceJson, false, true)
	c.CreateDashboards(cm)
	apiClient.AssertExpectations(t)
}

func TestCreateDatasourceAnnotationFalse(t *testing.T) {
	apiClient := new(APIClientMock)
	c := newConfigMapController(apiClient)
	cm := grafanaConfigMap(datasourceJson, false, false)
	c.CreateDashboards(cm)
	apiClient.AssertNotCalled(t, "CreateDatasource")
}

func grafanaConfigMap(value string, isDashboard, enable bool) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "some-config-map",
			Annotations:  annotate(isDashboard, enable),
		},
		Data: map[string]string{
			"some_file.json": value,
		},
	}
}

func annotate(isDashboard, enable bool) map[string]string {
	if isDashboard {
		return map[string]string{"grafana.net/dashboards": strconv.FormatBool(enable)}
	} else {
		return map[string]string{"grafana.net/datasource": strconv.FormatBool(enable)}
	}
}

func newConfigMapController(g grafana.APIInterface) *controller.ConfigMapController {
	clientset, _ := kubernetes.NewClientSet(true)
	return controller.NewConfigMapController(clientset, g)
}
