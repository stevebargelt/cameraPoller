package camera

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Token : an auth token for the camera API
type Token struct {
	Username       string
	Password       string
	CameraLoginURL string
}

// LoginTokenRequest is the model for building a login token request
type loginTokenRequest struct {
	Cmd    string `json:"cmd"`
	Action int    `json:"action"`
	Param  struct {
		User struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		} `json:"User"`
	} `json:"param"`
}

// LoginTokenResponse is the model for the Token response from camera
type loginTokenResponse struct {
	Cmd   string `json:"cmd"`
	Code  int    `json:"code"`
	Value struct {
		Token struct {
			LeaseTime int    `json:"leaseTime"`
			Name      string `json:"name"`
		} `json:"Token"`
	} `json:"value"`
}

// GetToken returns an auth token from the camera API
func (c *Token) GetToken() string {

	var loginTokenRequest [1]loginTokenRequest
	loginTokenRequest[0].Cmd = "Login"
	loginTokenRequest[0].Action = 0
	loginTokenRequest[0].Param.User.UserName = c.Username
	loginTokenRequest[0].Param.User.Password = c.Password

	var jsonStr, err = json.Marshal(loginTokenRequest)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("%s", string(jsonStr))

	client := http.Client{}

	request, _ := http.NewRequest("GET", c.CameraLoginURL, bytes.NewBuffer(jsonStr))
	request.Header.Set("Content-type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var l []loginTokenResponse
	err = json.Unmarshal(body, &l)
	return l[0].Value.Token.Name

}
