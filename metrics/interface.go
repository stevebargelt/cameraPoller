package metrics

import "time"

//Camera define a Camera call
type Camera struct {
	Name       string
	StatusCode string
	StartedAt  time.Time
	FinishedAt time.Time
	Duration   float64
}

// NewCamera create a new Cameraapp
func NewCamera(name string) *Camera {
	return &Camera{
		Name: name,
	}
}

//Started start monitoring the app
func (c *Camera) Started() {
	c.StartedAt = time.Now()
}

// Finished app finished
func (c *Camera) Finished() {
	c.FinishedAt = time.Now()
	c.Duration = time.Since(c.StartedAt).Seconds()
}

//Vision Call
type Vision struct {
	Name           string
	PredictionType string
	StatusCode     string
	StartedAt      time.Time
	FinishedAt     time.Time
	Duration       float64
}

//NewVision create a new Vision Call
func NewVision(name string, predictionType string) *Vision {
	return &Vision{
		Name:           name,
		PredictionType: predictionType,
	}
}

//Started start monitoring the app
func (v *Vision) Started() {
	v.StartedAt = time.Now()
}

// Finished app finished
func (v *Vision) Finished() {
	v.FinishedAt = time.Now()
	v.Duration = time.Since(v.StartedAt).Seconds()
}

//Cat Metrics
type Cat struct {
	Name      string
	Direction string
}

//NewCat create a new Cat Call
func NewCat(name string) *Cat {
	return &Cat{
		Name: name,
	}
}

//UseCase definition
type UseCase interface {
	SaveCamera(m *Camera)
	SaveVision(v *Vision)
	IncrementCat(c *Cat)
}
