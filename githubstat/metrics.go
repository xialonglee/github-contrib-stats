package githubstat

import "fmt"

type Metrics interface {
	Show()
}

type MetricsRequest interface {
	express()       // express the meaning of request
	validate() bool // determine whether parameters are valid or not.
	SetParameters(param *MetricsParameters)
	FetchMetrics() Metrics
}

// metrics parameters
type MetricsParameters struct {
	Repos     []*RepoParameters
	Dimension *string
}
type RepoParameters struct {
	OwnerName *string
	RepoName  *string
}
type DefaultMetrics struct{}

type DefaultMetricsRequest struct{}

func (m *DefaultMetrics) Show() {
	// void
}

func (m *DefaultMetricsRequest) express() {
	fmt.Println("you must select available metrics at least one.")
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
