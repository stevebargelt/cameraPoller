package integrationtests

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.1/customvision/prediction"
	uuid "github.com/satori/go.uuid"
)

type ImagePredictor struct {
	Predictor   prediction.BaseClient
	ProjectID   uuid.UUID
	IterationID uuid.UUID
	FilePath    string
}

type fileTest struct {
	file        string
	expectedCat string
	expectedDir string
}

var fileTests []fileTest

func TestVision(t *testing.T) {

	err := filepath.Walk(configuration.CaptureFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			_, filename := filepath.Split(path)
			fmt.Printf("Predicting file: %s\n", filename)
			tags := strings.Split(filename, "-")
			fileTests = append(fileTests, fileTest{path, tags[0], tags[1]})
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	for _, tt := range fileTests {
		catPredictor.FilePath = tt.file
		results := catPredictor.Predict()
		actual := getHighest(results)
		_, filename := filepath.Split(tt.file)
		if actual != tt.expectedCat {
			t.Errorf("Cat(%s): expected %s, actual %s", filename, tt.expectedCat, actual)
		}
	}
}

func getHighest(results prediction.ImagePrediction) string {

	var highestTag float64
	var highestName string

	for _, prediction := range *results.Predictions {
		if *prediction.Probability > highestTag {
			highestName = *prediction.TagName
			highestTag = *prediction.Probability
		}
	}
	return strings.ToLower(highestName)
}
