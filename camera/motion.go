package camera

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stevebargelt/cameraPoller/metrics"
)

// Motion is an object for motion capture
type Motion struct {
	Client            *retryablehttp.Client
	CameraMotionURL   string
	CameraStillPicURL string
	CameraUsername    string
	CameraPassword    string
	CameraLoginURL    string
	Token             string
	CaptureFolder     string
}

// MotionRequest is the model for for building a motion  request
type motionRequest struct {
	Cmd    string `json:"cmd"`
	Action int    `json:"action"`
	Param  struct {
		Channel int `json:"channel"`
	} `json:"param"`
}

// MotionResponse is the model for the Motion response from the camera
type motionResponse struct {
	Cmd   string `json:"cmd"`
	Code  int    `json:"code"`
	Value struct {
		State int `json:"state"`
	} `json:"value"`
	Error struct {
		Detail  string `json:"detail"`
		RSPCode int    `json:"rspCode"`
	} `json:"error"`
}

//MotionCap - captures motion from camera and writes an image to file
func (m *Motion) MotionCap(mService metrics.UseCase) string {

	var motionRequest [1]motionRequest
	motionRequest[0].Cmd = "GetMdState"
	motionRequest[0].Action = 0
	motionRequest[0].Param.Channel = 0

	appMetric := metrics.NewCamera("motion")
	appMetric.Started()
	var motionRequestJSON, err2 = json.Marshal(motionRequest)
	if err2 != nil {
		panic(err2)
	}
	// fmt.Printf("This is m: %v\n\n", m)
	request, _ := retryablehttp.NewRequest("GET", m.CameraMotionURL+m.Token, bytes.NewBuffer(motionRequestJSON))
	request.Header.Set("Content-type", "application/json")

	response, err := m.Client.Do(request)
	if err != nil {
		panic(err)
	}
	appMetric.Finished()
	appMetric.StatusCode = strconv.Itoa(response.StatusCode)
	mService.SaveCamera(appMetric)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var motionResp []motionResponse
	err = json.Unmarshal(body, &motionResp)
	if err != nil {
		fmt.Println("Error Unmarshalling")
	}
	errorCode := motionResp[0].Error.RSPCode
	if errorCode != 0 {
		fmt.Printf("Error returned: %v\n", errorCode)
		if errorCode == -6 {
			m.getCamToken()
		}
	} else {
		state := motionResp[0].Value.State
		if state == 1 {
			fmt.Printf("Motion State: %v \n", state)
			file := createFile(m.CaptureFolder)
			putFile(file, m.CameraStillPicURL, httpClient())
			return file.Name()
		}
	}
	return ""
}

func (m *Motion) getCamToken() {
	camToken := Token{CameraLoginURL: m.CameraLoginURL, Username: m.CameraUsername, Password: m.CameraPassword}
	m.Token = camToken.GetToken()
}

func buildFileName() string {
	return time.Now().Format("20060102150405")
}

func putFile(file *os.File, stillPicURL string, client *http.Client) {
	resp, err := client.Get(stillPicURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("Just Downloaded a file with size %d", size)
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	return &client
}

func createFile(captureFolder string) *os.File {
	file, err := os.Create(captureFolder + buildFileName() + ".jpg")
	if err != nil {
		panic(err)
	}
	return file
}
