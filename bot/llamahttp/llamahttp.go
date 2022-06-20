package llamahttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/LlamaNite/llamaimage"
)

// Map is an alias of map[string]interface{}.
type Map map[string]any

// Headers is an alias of map[string]string.
type Headers map[string]string

type Options struct {
	Parameters Map
	Payloads   io.Reader
	Headers    Headers
}

// JSONPayload serializes the payload as a JSON string and ignores errors
func JSONPayload(data interface{}) io.Reader {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(data)
	return buf
}

// FormURLEncodedPayload serializes the payload like the parameters
// you see in GET requests (form URL encoded)
func FormURLEncodedPayload(data Headers) io.Reader {
	requestPayloads := url.Values{}
	for key, value := range data {
		requestPayloads.Set(key, value)
	}
	return bytes.NewReader([]byte(requestPayloads.Encode()))
}

func Do(method, URL string, options Options) (*http.Response, error) {
	req, err := http.NewRequest(method, URL, options.Payloads)
	if err != nil {
		return nil, err
	}

	// Pass parameters into URL
	query := req.URL.Query()
	for key, value := range options.Parameters {
		query.Add(key, fmt.Sprint(value))
	}
	req.URL.RawQuery = query.Encode()

	// Pass headers into request headers
	for key, value := range options.Headers {
		req.Header.Set(key, value)
	}

	return http.DefaultClient.Do(req)
}

func GetImage(id, url string) (image.Image, error) {
	filepath := filepath.Join("assets", id)

	// use local-storage cache if exists
	if f, err := os.Open(filepath); err == nil {
		if icon, err := llamaimage.OpenImage(f); err == nil {
			return icon, nil
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch > status %d", resp.StatusCode)
	}

	icon, err := llamaimage.OpenImage(resp.Body)
	if err != nil {
		return nil, err
	}
	go llamaimage.Save(icon, filepath)
	return icon, nil
}
