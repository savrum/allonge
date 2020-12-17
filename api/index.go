package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// To return nil for strings you can use string pointer
func base64ify(url string) (*string, int, error) {
	res, httpErr := http.Get(url) // Get image from url

	// Check if URL is valid
	if httpErr != nil {
		return nil, 404, fmt.Errorf("Requested URL not found \"%s\"", url)
	}
	defer res.Body.Close() // Close body

	// Check if URL responds successfully
	if res.StatusCode != 200 {
		return nil, res.StatusCode, fmt.Errorf("HTTP image read error %d", res.StatusCode)
	}

	// Check if it's really an image
	if strings.Split(res.Header.Get("Content-type"), "/")[0] != "image" {
		return nil, 415, errors.New("Request URL is not an image")
	}

	// If all is well, create buffer from image
	bodyBytes, readErr := ioutil.ReadAll(res.Body)
	// Check if there's an internal error
	if readErr != nil {
		return nil, 500, errors.New("File read error")
	}
	// Convert image data to base64
	bodyString := string(base64.StdEncoding.EncodeToString(bodyBytes))
	// Return image base64 data using address of pointer
	return &bodyString, 200, nil
}

// Handler for serverless function.
func Handler(res http.ResponseWriter, req *http.Request) {
	image, _image := req.URL.Query()["image"] // Get image url

	// Check if image param available
	if !_image {
		res.WriteHeader(400)                             // Set HTTP 400 (Bad Request)
		fmt.Fprintf(res, "Missing image url parameter!") // Add message
	} else {
		imageB64, status, b64Err := base64ify(image[0]) // Get first param to process
		if b64Err != nil {
			res.WriteHeader(status)          // Set HTTP 400 (Bad Request)
			fmt.Fprintf(res, b64Err.Error()) // Add error message
		} else {
			res.WriteHeader(status)      // Set HTTP 200 (Success)
			res.Write([]byte(*imageB64)) // Return byte array raw
		}
	}

}
