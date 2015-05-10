package githubstat

import (
	"fmt"

	"github.com/mgutz/ansi"
)

type Metrics interface {
	GetMetrics()
}

type MetricsRequest interface {
	express()       // express the meaning of request
	validate() bool // determine whether parameters are valid or not.
	SetParameters(param *MetricsParameters)
	FetchMetrics() Metrics
}

// metrics parameters
type MetricsParameters struct {
	OrgName   *string
	OwnerName *string
	RepoName  *string
}

type DefaultMetrics struct{}

type DefaultMetricsRequest struct{}

func (m *DefaultMetrics) GetMetrics() {
	// void
}

func (m *DefaultMetricsRequest) express() {
	fmt.Println(ansi.Color("you must select available metrics at least one.", "yellow"))
}

func (m *DefaultMetricsRequest) validate() bool {
	return true
}

func (m *DefaultMetricsRequest) SetParameters(param *MetricsParameters) {
	// do nothing
}

func (m *DefaultMetricsRequest) FetchMetrics() Metrics {
	m.express()
	return &DefaultMetrics{}
}
