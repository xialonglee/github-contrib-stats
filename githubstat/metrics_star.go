package githubstat

import "fmt"

type StarMetrics struct {
	star int
}

func (m *StarMetrics) GetMetrics() {
	fmt.Printf("star: %d\n", m.star)
}

type StarMetricsRequest struct {
	param *MetricsParameters
}

func (m *StarMetricsRequest) express() {
	fmt.Println("star selected")
}

func (m *StarMetricsRequest) SetParameters(param *MetricsParameters) {
	m.param = param
}

func (m *StarMetricsRequest) FetchMetrics() Metrics {
	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()
	repos, _, _ := client.Repositories.Get(*m.param.OwnerName, *m.param.RepoName)

	return &StarMetrics{
		star: *repos.StargazersCount,
	}
}
