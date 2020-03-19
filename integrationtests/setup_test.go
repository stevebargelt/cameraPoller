package integrationtests

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.1/customvision/prediction"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"github.com/stevebargelt/cameraPoller/config"
	"github.com/stevebargelt/cameraPoller/vision"
)

var (
	configuration      config.Configuration
	configTestFolder   string
	catPredictor       vision.ImagePredictor
	directionPredictor vision.ImagePredictor
)

func TestMain(m *testing.M) {

	setupConfig()
	os.Exit(m.Run())

}

func setupConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../") // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	viper.SetDefault("HTTP_RETRY_COUNT", 20)

	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	projectID, err := uuid.FromString(configuration.ProjectID)
	if err != nil {
		fmt.Printf("Something went wrong creating ProjectID UUID: %s", err)
	}

	projectIDDirection, err := uuid.FromString(configuration.ProjectIDDirection)
	if err != nil {
		fmt.Printf("Something went wrong creating ProjectID UUID: %s", err)
	}

	iterationID, err := uuid.FromString(configuration.IterationID)
	if err != nil {
		fmt.Printf("Something went wrong creating Iteration UUID: %s", err)
	}

	iterationIDDirection, err := uuid.FromString(configuration.IterationIDDirection)
	if err != nil {
		fmt.Printf("Something went wrong creating Iteration Direction UUID: %s", err)
	}

	//Replace the capture folder wtih the TEST folder
	configuration.CaptureFolder = viper.GetString("CAPTURE_FOLDER_TEST")
	//"/Users/stevebargelt/go/src/github.com/stevebargelt/cameraPoller/integrationtests/testdata"

	predictor := prediction.New(configuration.PredictionKey, configuration.EndpointURL)

	catPredictor.ProjectID = projectID
	catPredictor.IterationID = iterationID
	catPredictor.Predictor = predictor

	directionPredictor.ProjectID = projectIDDirection
	directionPredictor.IterationID = iterationIDDirection
	directionPredictor.Predictor = predictor

}
