package llamahttp

import (
	"bytes"
	"encoding/json"
	"errors"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Map is an alias of map[string]interface{}.
type Map = map[string]interface{}

// Headers is an alias of map[string]string.
type Headers = map[string]string

// Get does a HTTP GET request
func Get(URL string, parameters, headers map[string]string) (statusCode int, responseBytes []byte, err error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		statusCode = -1
		return
	}

	// Pass parameters into URL
	query := request.URL.Query()
	for key, value := range parameters {
		query.Add(key, value)
	}
	request.URL.RawQuery = query.Encode()

	// Pass headers into request headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// Make client and do a request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		statusCode = -1
		return
	}

	// Read response body
	responseBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		statusCode = -1
		return
	}

	// Nothing happened? return the result!
	statusCode = response.StatusCode
	return
}

// JSONPayload serializes the payload as a JSON string and ignores errors
func JSONPayload(data interface{}) io.Reader {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	return buf
}

// FormURLEncodedPayload serializes the payload like the parameters
// you see in GET requests (form URL encoeded)
func FormURLEncodedPayload(data map[string]string) io.Reader {
	requestPayloads := url.Values{}
	for key, value := range data {
		requestPayloads.Set(key, value)
	}
	return bytes.NewReader([]byte(requestPayloads.Encode()))
}

// Post does a HTTP POST request
func Post(URL string, payload io.Reader, headers map[string]string) (statusCode int, responseBytes []byte, err error) {
	// Make request
	request, err := http.NewRequest("POST", URL, payload)
	if err != nil {
		return
	}

	// Pass headers into request headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// Make client and do a request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	// Read response body
	responseBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

// Delete does a HTTP DELETE request
func Delete(URL string, payload io.Reader, headers map[string]string) (statusCode int, responseBytes []byte, err error) {
	// Make request
	request, err := http.NewRequest("DELETE", URL, payload)
	if err != nil {
		return
	}

	// Pass headers into request headers
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	// Make client and do a request
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return
	}

	statusCode = response.StatusCode

	// Read response body
	responseBytes, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	return
}

// GetImage is a helper that downloads and parses an image (GET request)
func GetImage(URL string) (image.Image, error) {
	request, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusOK {
		responseBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		decodedImage, _, err := image.Decode(bytes.NewBuffer(responseBytes))
		return decodedImage, err
	}

	return nil, errors.New(strconv.Itoa(response.StatusCode))
}
