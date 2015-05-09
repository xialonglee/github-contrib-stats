package githubstat

import "fmt"

type IssueMetrics struct {
	issue int
}

func (m *IssueMetrics) GetMetrics() {
}

type IssueMetricsRequest struct {
}

func (m *IssueMetricsRequest) express() {
	fmt.Println("issue selected")
}

func (m *IssueMetricsRequest) FetchMetrics() {
}

func (m *IssueMetricsRequest) GetMetrics() Metrics {
	return &IssueMetrics{issue: 0} // temporary returning. @todo
}
