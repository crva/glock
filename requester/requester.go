package requester

import (
	"net/http"
	"time"
)

func SendHttpRequest(url string) (time.Duration, error) {
	start := time.Now()

	response, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer response.Body.Close()
	return time.Since(start), nil
}
