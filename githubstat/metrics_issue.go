package githubstat

import "fmt"

type IssueMetrics struct {
	issue int
}

func (m *IssueMetrics) Show() {
	fmt.Printf("this metrics has not been implemented yet.")
}

type IssueMetricsRequest struct {
	param *MetricsParameters
}

func (m *IssueMetricsRequest) express() {
	//fmt.Printf("target repository: %s/%s\n", *m.param.OwnerName, *m.param.RepoName)
	fmt.Println("metrics: issue count")
}

func (m *IssueMetricsRequest) SetParameters(param *MetricsParameters) {
	m.param = param
}

func (m *IssueMetricsRequest) validate() bool {
	return true
}

func (m *IssueMetricsRequest) FetchMetrics() Metrics {
	return &IssueMetrics{
		issue: 0,
	} // temporary returning. @todo
}
