package requester

import (
	"io"
	"net/http"
	"time"
)

func SendHttpRequest(url string, method string, json io.Reader) (time.Duration, error) {
	start := time.Now()

	var response *http.Response
	var err error

	if method == "POST" {
		response, err = http.Post(url, "application/json", json)
	} else {
		response, err = http.Get(url)
	}
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	return time.Since(start), nil
}
