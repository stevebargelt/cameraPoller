package vision

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.1/customvision/prediction"
	uuid "github.com/satori/go.uuid"
	"github.com/stevebargelt/cameraPoller/metrics"
)

// ImagePredictor : predicts what an image contains
type ImagePredictor struct {
	Predictor   prediction.BaseClient
	ProjectID   uuid.UUID
	IterationID uuid.UUID
	FilePath    string
}

// Predict - takes info and returns a prediction
func (p *ImagePredictor) Predict(mService metrics.UseCase, predictionType string) prediction.ImagePrediction {

	ctx := context.Background()
	var testImageData []byte
	var err error
	retryCount := 0
	//this is really UGLY. Finding that we error on the first file because it's not done writing when we read it.
	for ok := true; ok; ok = (len(testImageData) == 0) {
		testImageData, err = ioutil.ReadFile(p.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Length %v\n", len(testImageData))
		if len(testImageData) == 0 {
			retryCount++
		}
	}
	log.Printf("RetryCount = %v\n", retryCount)
	appMetric := metrics.NewVision("prediciton", predictionType)
	appMetric.Started()
	results, err := p.Predictor.PredictImage(ctx, p.ProjectID, ioutil.NopCloser(bytes.NewReader(testImageData)), &p.IterationID, "")
	appMetric.Finished()
	appMetric.StatusCode = strconv.Itoa(results.StatusCode)
	mService.SaveVision(appMetric)
	if err != nil {
		fmt.Println("\n\npredictor.PredictImage Failed.")
		log.Fatal(err)
	}
	return results
}
