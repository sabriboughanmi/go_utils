package remoteconfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{}
var etag string
var configBytes []byte

type ProjectConfig struct {
	token                *oauth2.Token
	baseURl              string
	remoteConfigEndPoint string
	remoteConfigURL      string
}

//ServiceAccountInit Initialize the Service Account using Parameters
func InitFromCredentials(email, privateKey, projectID string) (*ProjectConfig, error) {
	var projectConfig ProjectConfig
	config := &jwt.Config{
		Email:      email,
		PrivateKey: []byte(privateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/firebase.remoteconfig",
		},
		TokenURL: google.JWTTokenURL,
	}

	token, err := config.TokenSource(oauth2.NoContext).Token()
	if err != nil {
		return &projectConfig, err
	}

	projectConfig.baseURl = "https://firebaseremoteconfig.googleapis.com"
	projectConfig.remoteConfigEndPoint = "v1/projects/" + projectID + "/remoteConfig"
	projectConfig.remoteConfigURL = projectConfig.baseURl + "/" + projectConfig.remoteConfigEndPoint
	projectConfig.token = token

	return &projectConfig, nil
}

//ServiceAccount Initialize the Service Account using a ServiceAccount.Json path
func InitFromPath(serviceAccountPath string) (*ProjectConfig, error) {
	var projectConfig ProjectConfig
	b, err := ioutil.ReadFile(serviceAccountPath)
	if err != nil {
		return &projectConfig, err
	}
	var c = struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
		ProjectID  string `json:"project_id"`
	}{}

	json.Unmarshal(b, &c)
	config := &jwt.Config{
		Email:      c.Email,
		PrivateKey: []byte(c.PrivateKey),
		Scopes: []string{
			"https://www.googleapis.com/auth/firebase.remoteconfig",
		},
		TokenURL: google.JWTTokenURL,
	}

	token, err := config.TokenSource(oauth2.NoContext).Token()
	if err != nil {
		return &projectConfig, err
	}

	projectConfig.baseURl = "https://firebaseremoteconfig.googleapis.com"
	projectConfig.remoteConfigEndPoint = "v1/projects/" + c.ProjectID + "/remoteConfig"
	projectConfig.remoteConfigURL = projectConfig.baseURl + "/" + projectConfig.remoteConfigEndPoint
	projectConfig.token = token

	return &projectConfig, nil
}

func (pc *ProjectConfig) Publish(Etag string) (string, error) {

	if configBytes == nil || len(configBytes) == 0 {
		return "", fmt.Errorf("UnInitialized Config Buffer")
	}

	request, err := http.NewRequest("PUT", pc.remoteConfigURL, bytes.NewReader(configBytes))
	if err != nil {
		return "", err
	}

	request.Header.Set("Authorization", "Bearer "+pc.token.AccessToken)
	request.Header.Add("Content-Type", "application/json; UTF-8")
	request.Header.Add("If-Match", Etag)

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	// if resp.Status is 200
	if response.StatusCode == http.StatusOK {
		return response.Header["Etag"][0], nil
	}
	return "", fmt.Errorf("resp StatusCode : %d", response.StatusCode)
}

//RollbackVersion Rollback the RemoteConfig to an older Version
func (pc *ProjectConfig) RollbackVersion(version int) error {
	respJSON := map[string]int{"version_number": version}
	b, _ := json.Marshal(respJSON)

	// Create new request
	req, err := http.NewRequest(http.MethodPost, pc.remoteConfigURL+":rollback", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+pc.token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {

		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Convert Response Body to String
		bodyString := string(bodyBytes)
		print(bodyString)

		// get latest etag
		_, err = pc.GetRemoteConfig()
		if err != nil {
			return err
		}

		writeEtag(etag)
	}
	return nil
}

//GetRemoteConfig get the Remote Config Value
func (pc *ProjectConfig) GetRemoteConfig() (*RemoteConfig, error) {
	var remoteConfig RemoteConfig

	req, err := http.NewRequest("GET", pc.remoteConfigURL, nil)
	if err != nil {
		return &remoteConfig, err
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+pc.token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return &remoteConfig, err
	}

	// if resp.Status is 200
	if resp.StatusCode == http.StatusOK {

		//Save Etag
		etag = resp.Header["Etag"][0]

		// Read response body
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return &remoteConfig, err
		}

		err = json.Unmarshal(bytes, &remoteConfig)
		// Convert Response Body to String
		return &remoteConfig, err
	}
	return &remoteConfig, fmt.Errorf("resp StatusCode : %d", resp.StatusCode)
}

func (pc *ProjectConfig) ListVersion(size int) (string, error) {

	req, err := http.NewRequest(http.MethodGet, pc.remoteConfigURL+":listVersions?pageSize="+strconv.Itoa(size), nil)
	if err != nil {
		return "", err
	}

	// Set Authorization Header
	req.Header.Set("Authorization", "Bearer "+pc.token.AccessToken)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {

		// Read response body
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		// Convert Response Body to String
		return string(bodyBytes), nil
	}

	return "", fmt.Errorf("resp StatusCode : %d", resp.StatusCode)
}

func readEtag() string {
	return etag
}

func writeEtag(etagV string) {
	etag = etagV
}

//Remote Config Version
type Version struct {
	VersionNumber string    `json:"versionNumber"`
	UpdateTime    time.Time `json:"updateTime"`
	UpdateUser    struct {
		Email string `json:"email"`
	} `json:"updateUser"`
	UpdateOrigin string `json:"updateOrigin"`
	UpdateType   string `json:"updateType"`
}

type ConditionName string

const (
	DefaultValue ConditionName = "defaultValue"
)

type Condition map[string]ConditionName
type Parameters map[string]map[ConditionName]map[string]string

type RemoteConfig struct {
	Conditions *[]Condition `json:"conditions,omitempty"`
	Parameters *Parameters  `json:"parameters,omitempty"`
	Version    *Version     `json:"Version,omitempty"`
}

//GetValue safely returns a value for a configKey by ConditionName
func (p Parameters) GetValue(configKey string, condition ConditionName) (string, bool) {
	if p == nil {
		return "", false
	}

	if conditionsMap, ok := p[configKey]; ok {
		if valMap, ok := conditionsMap[condition]; ok {
			if val, ok := valMap["value"]; ok {
				return val, ok
			}
		}
		return "", false
	}
	return "", false
}

//GetDefaultValue safely returns a the Default Value for a configKey
func (p Parameters) GetDefaultValue(configKey string) (string, bool) {
	if p == nil {
		return "", false
	}

	if conditionsMap, ok := p[configKey]; ok {
		if valMap, ok := conditionsMap[DefaultValue]; ok {
			if val, ok := valMap["value"]; ok {
				return val, ok
			}
		}
		return "", false
	}
	return "", false
}
