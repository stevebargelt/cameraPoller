package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

//Service implements
type Service struct {
	cameraHistogram *prometheus.HistogramVec
	visionHistogram *prometheus.HistogramVec
}

//NewPrometheusService creates a new prometheus service
func NewPrometheusService() (*Service, error) {

	prediction := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "vision",
		Name:      "request_duration_seconds",
		Help:      "Histogram for the runtime of a call to Azure Custom Vision Prediction Service.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name", "predictiontype", "statuscode"})

	camera := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "camera",
		Name:      "request_duration_seconds",
		Help:      "Histogram for the runtime of a call to the camera service.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"name", "statuscode"})

	s := &Service{
		cameraHistogram: camera,
		visionHistogram: prediction,
	}

	err := prometheus.Register(s.cameraHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}

	err = prometheus.Register(s.visionHistogram)
	if err != nil && err.Error() != "duplicate metrics collector registration attempted" {
		return nil, err
	}
	return s, nil
}

//SaveVision send metrics to server
func (s *Service) SaveVision(v *Vision) {
	// gatewayURL := config.PROMETHEUS_PUSHGATEWAY
	// s.pHistogram.WithLabelValues(c.Name).Observe(c.Duration)
	// return push.New(gatewayURL, "cmd_job").Collector(s.pHistogram).Push()
	s.visionHistogram.WithLabelValues(v.Name, v.PredictionType, v.StatusCode).Observe(v.Duration)
}

//SaveCamera send metrics to server
func (s *Service) SaveCamera(c *Camera) {
	s.cameraHistogram.WithLabelValues(c.Name, c.StatusCode).Observe(c.Duration)
}
