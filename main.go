package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v1.1/customvision/prediction"
	"github.com/hashicorp/go-retryablehttp"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"google.golang.org/api/option"

	camera "github.com/stevebargelt/cameraPoller/camera"
	config "github.com/stevebargelt/cameraPoller/config"
	"github.com/stevebargelt/cameraPoller/vision"
)

// LitterboxUser = defines the attributes of a cat using the litterbox
type LitterboxUser struct {
	Name                 string
	NameProbability      float64
	Direction            string
	DirectionProbability float64
	Photo                string
}

func main() {

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".") // look for config in the working directory
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s ", err))
	}
	viper.SetDefault("HTTP_RETRY_COUNT", 20)

	var configuration config.Configuration
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

	predictor := prediction.New(configuration.PredictionKey, configuration.EndpointURL)

	client := retryablehttp.NewClient() // http.Client{}
	client.Backoff = retryablehttp.LinearJitterBackoff
	client.RetryWaitMin = 1 * time.Second
	client.RetryWaitMax = 5 * time.Second
	client.RetryMax = 50 // we know the most likely reason is that the network is down (rebooted hardware)

	motionCap := new(camera.Motion)
	motionCap.Client = client
	motionCap.CameraMotionURL = configuration.CameraMotionURL
	motionCap.CameraStillPicURL = configuration.CameraStillPicURL
	motionCap.CameraLoginURL = configuration.CameraLoginURL
	motionCap.CameraUsername = configuration.CameraUsername
	motionCap.CameraPassword = configuration.CameraPassword
	motionCap.CaptureFolder = configuration.CaptureFolder

	catPredictor := new(vision.ImagePredictor)
	catPredictor.ProjectID = projectID
	catPredictor.IterationID = iterationID
	catPredictor.Predictor = predictor

	directionPredictor := new(vision.ImagePredictor)
	directionPredictor.ProjectID = projectIDDirection
	directionPredictor.IterationID = iterationIDDirection
	directionPredictor.Predictor = predictor

	var photos []string
	var litterboxPicSet []LitterboxUser
	ticker := time.NewTicker(1 * time.Second)
	timeout := make(chan bool, 1)
	// haveIdentity := make(chan bool, 1)
	quit := make(chan struct{})
	// start := time.Now()
	go func() {
		for {
			select {
			case <-ticker.C:
				fileName := motionCap.MotionCap()
				if len(fileName) > 0 {
					photos = append(photos, fileName)
					catPredictor.FilePath = fileName
					results := catPredictor.Predict()
					highestProbabilityTag := processCatResults(results, fileName)
					litterboxPicSet = append(litterboxPicSet, highestProbabilityTag)
					// If this is the first photo then set a timer so we don't wait indef for n (configuration.PhotosInSet) photos...
					if len(litterboxPicSet) == 1 {
						go func() {
							time.Sleep(time.Duration(configuration.TimeoutValue) * time.Second)
							timeout <- true
						}()
					}
					if len(photos) == configuration.PhotosInSet {
						ticker.Stop()
						litterboxUser, weHaveCat := determineResults(litterboxPicSet)
						if weHaveCat {
							directionPredictor.FilePath = litterboxUser.Photo
							directionResults := directionPredictor.Predict()
							setDirection(directionResults, &litterboxUser)
						}
						doStuffWithResult(litterboxUser, configuration.FirebaseCredentials, configuration.FirestoreCollection, weHaveCat)
						moveProcessedFiles(litterboxPicSet, configuration.ProcessedFolder)
						litterboxPicSet = nil
					}
				}
			// case <-haveIdentity:
			// 	fmt.Println("We IDd the cat")
			// 	fmt.Println("Re-Starting Ticker")
			// 	ticker = time.NewTicker(1 * time.Second)
			case <-timeout:
				ticker.Stop()
				if len(photos) == 0 {
					fmt.Printf("We Good. Timeout called but we processed %v pics.\n", configuration.PhotosInSet)
				} else {
					fmt.Printf("Timeout called but we processed %v pics.\n", len(photos))
					litterboxUser, weHaveCat := determineResults(litterboxPicSet)
					if weHaveCat {
						directionPredictor.FilePath = litterboxUser.Photo
						directionResults := directionPredictor.Predict()
						setDirection(directionResults, &litterboxUser)
					}
					doStuffWithResult(litterboxUser, configuration.FirebaseCredentials, configuration.FirestoreCollection, weHaveCat)
					moveProcessedFiles(litterboxPicSet, configuration.ProcessedFolder)
					litterboxPicSet = nil
				}
				fmt.Println("Re-Starting Ticker")
				ticker = time.NewTicker(1 * time.Second)
			case <-quit:
				fmt.Println("In quit")
				ticker.Stop()
				return
			}
		}
	}()

	select {}

}

func doStuffWithResult(litterboxUser LitterboxUser, firebaseCredentials string, firestoreCollection string, weHaveCat bool) {

	if weHaveCat {
		fmt.Printf("I am %v%% sure that it was %s and I am ", litterboxUser.NameProbability*100, litterboxUser.Name)
		fmt.Printf("%v%% sure that they were headed %s the catbox!\n", litterboxUser.DirectionProbability*100, litterboxUser.Direction)
		addLitterBoxTripToFirestore(litterboxUser, firebaseCredentials, firestoreCollection)
	} else {
		fmt.Printf("I am %v%% sure that we had a false motion event!\n", litterboxUser.NameProbability*100)
	}
}

func processCatResults(results prediction.ImagePrediction, fileName string) LitterboxUser {

	//Process the results of ONE image... loop through the tag predictions in the image

	litterboxUser := LitterboxUser{"Negative", 0.00, "", 0.00, fileName}
	for _, prediction := range *results.Predictions {
		fmt.Printf("\t%s: %.2f%%\n", *prediction.TagName, *prediction.Probability*100)

		//of the tags in the model pick the highest probability
		// TODO: Use a slice or enum for the direction, no magic strings and decouple the directions
		if *prediction.Probability > litterboxUser.NameProbability {
			litterboxUser.Name = *prediction.TagName
			litterboxUser.NameProbability = *prediction.Probability
			litterboxUser.Photo = fileName
		}
	}
	return litterboxUser
}

func setDirection(directionResults prediction.ImagePrediction, litterboxUser *LitterboxUser) {

	for _, prediction := range *directionResults.Predictions {
		fmt.Printf("\t%s: %.2f%%\n", *prediction.TagName, *prediction.Probability*100)

		if *prediction.TagName == "in" || *prediction.TagName == "out" {
			if *prediction.Probability > litterboxUser.DirectionProbability {
				litterboxUser.Direction = *prediction.TagName
				litterboxUser.DirectionProbability = *prediction.Probability
			}
		}
	}
}

func determineResults(litterboxPicSet []LitterboxUser) (LitterboxUser, bool) {
	var highestCatIndex int
	var highestCatProbability = 0.0
	var highestNegProbability = 0.0
	var highestNegIndex int
	var weHaveCat = false
	//log.Printf("litterboxPicSet: %v\n", litterboxPicSet)
	for index, element := range litterboxPicSet {
		if element.Name != "Negative" {
			if element.NameProbability > highestCatProbability {
				highestCatProbability = element.NameProbability
				highestCatIndex = index
				weHaveCat = true
			}
		} else {
			if element.NameProbability > highestNegProbability {
				highestNegProbability = element.NameProbability
				highestNegIndex = index
			}
		}
	}
	if weHaveCat {
		return litterboxPicSet[highestCatIndex], weHaveCat
	}
	return litterboxPicSet[highestNegIndex], weHaveCat
}

// Next Steps?
func addLitterBoxTripToFirestore(user LitterboxUser, firebaseCredentials string, firestoreCollection string) {
	ctx := context.Background()
	sa := option.WithCredentialsFile(firebaseCredentials)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	//TODO: Find cat first. Add if not found?
	_, _, err = client.Collection(firestoreCollection).Doc(user.Name).Collection("LitterTrips").Add(ctx, map[string]interface{}{
		"Probability":          user.NameProbability,
		"Direction":            user.Direction,
		"DirectionProbability": user.DirectionProbability,
		"Photo":                user.Photo, // right now this is the local name. Could be the URL to the photo in Cloud Storage.
		"timestamp":            firestore.ServerTimestamp,
	})
	if err != nil {
		log.Fatalf("Failed adding litterbox trip: %v", err)
	}
}

func moveProcessedFiles(litterboxPicSet []LitterboxUser, processedFolder string) {

	for _, litterboxUser := range litterboxPicSet {
		//fmt.Printf("Photo: %s\n", litterboxUser.Photo)
		_, file := path.Split(litterboxUser.Photo)
		newpath := path.Join(processedFolder, file)
		err := os.Rename(litterboxUser.Photo, newpath)
		if err != nil {
			log.Fatal(err)
		}
	}
}
