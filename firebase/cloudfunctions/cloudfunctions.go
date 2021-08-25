package cloudfunctions

import (
	"encoding/json"
	"net/http"
	s "strings"
)

//StandardRequestResponse : needed for Callable Cloud Functions
type StandardRequestResponse struct {
	Data interface{} `json:"data"`
}

//FormatForRequestResponse : format a response to be Function Callable
func FormatForRequestResponse(value interface{}) StandardRequestResponse {
	return StandardRequestResponse{Data: value}
}

//GetBodyData genetically reads a response body from a http.Request
func GetBodyData(r *http.Request, out interface{}) error {

	type Body struct {
		Data interface{} `json:"data"`
	}

	var body Body

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(body.Data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, out)
	if err != nil {
		return err
	}

	return nil
}

//SetupResponse .
func SetupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

//CORSEnabledFunction .
func CORSEnabledFunction(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers for the preflight request
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set CORS headers for the main request.
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

//GetIDToken Returns the ID Token From Callable Functions
func GetIDToken(r *http.Request) (string, bool) {
	IDToken := r.Header.Get("Authorization")
	if IDToken == "" {
		return "", false
	}
	if s.HasPrefix(IDToken, "Bearer ") {
		IDToken = s.Replace(IDToken, "Bearer ", "", 1)
	}
	return IDToken, true
}
