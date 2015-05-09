package githubstat

type Metrics interface {
	GetMetrics()
}

type MetricsRequest interface {
	express() // express the meaning of request
	SetParameters(param *MetricsParameters)
	FetchMetrics() Metrics
}
