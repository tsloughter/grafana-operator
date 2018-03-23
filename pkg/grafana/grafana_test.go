package grafana

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeUrl(t *testing.T) {
	myUrl, _ := url.Parse("https://server:9000")
	result := makeUrl(myUrl, "api/grafana/datasource")
	assert.Equal(t, "https://server:9000/api/grafana/datasource", result)

}
