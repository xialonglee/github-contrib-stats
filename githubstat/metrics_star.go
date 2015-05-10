package githubstat

import (
	"fmt"

	"github.com/mgutz/ansi"
)

type StarMetrics struct {
	star int
}

func (m *StarMetrics) GetMetrics() {
	fmt.Printf(ansi.Color("star: %d\n", "green"), m.star)
}

type StarMetricsRequest struct {
	param *MetricsParameters
}

func (m *StarMetricsRequest) express() {
	fmt.Printf("target repository: %s/%s\n", *m.param.OwnerName, *m.param.RepoName)
	fmt.Printf("metrics: star\n")
}

func (m *StarMetricsRequest) SetParameters(param *MetricsParameters) {
	m.param = param
}

func (m *StarMetricsRequest) validate() bool {
	if *m.param.OwnerName == "" {
		return false
	} else if *m.param.RepoName == "" {
		return false
	}
	return true
}

func (m *StarMetricsRequest) FetchMetrics() Metrics {
	m.express()
	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()

	if m.validate() {
		repos, _, _ := client.Repositories.Get(*m.param.OwnerName, *m.param.RepoName)

		return &StarMetrics{
			star: *repos.StargazersCount,
		}
	} else {
		return &StarMetrics{
			star: -1,
		}
	}
}
