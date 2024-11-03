package data

type ShortUrlResponse struct {
	Value interface{} `json:"value"`
}

type HealthStatus struct {
	Status string `json:"status"`
	Environment string `json:"environment"`
	Version string `json:"version"`
}