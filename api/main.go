package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = http.Client{}

func base64ify(url string) (*string, error) {
	res, httpErr := http.Get(url)
	bodyBytes, readErr := ioutil.ReadAll(res.Body)

	if httpErr != nil || res.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP image read error %d", res.StatusCode)
	}
	if readErr != nil {
		return nil, errors.New("File read error")
	}

	bodyString := string(base64.StdEncoding.EncodeToString(bodyBytes))
	return &bodyString, nil
}

// Handler for serverless function.
func Handler(res http.ResponseWriter, req *http.Request) {
	image, _image := req.URL.Query()["image"] // Get image url

	// Check if image param available
	if !_image {
		res.WriteHeader(400)                             // Set HTTP 400 (Bad Request)
		fmt.Fprintf(res, "Missing image url parameter!") // Add message
	} else {
		imageB64, b64Err := base64ify(image[0])
		if b64Err != nil {
			res.WriteHeader(400)             // Set HTTP 400 (Bad Request)
			fmt.Fprintf(res, b64Err.Error()) // Add error message
		} else {
			res.WriteHeader(200)         // Set HTTP 200 (Success)
			res.Write([]byte(*imageB64)) // Return byte array raw
		}
	}

}
