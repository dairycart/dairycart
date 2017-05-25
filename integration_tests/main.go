package dairytest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	maxAttempts = 10
	baseURL     = `http://dairycart`
)

func buildURL(parts ...string) string {
	return fmt.Sprintf("%s/%s", baseURL, strings.Join(parts, "/"))
}

func turnResponseBodyIntoString(res *http.Response) (string, error) {
	bodyBytes, err := ioutil.ReadAll(res.Body)
	return string(bodyBytes), err
}

func ensureThatDairycartIsAlive() error {
	url := buildURL("health")
	dairyCartIsDown := true
	numberOfAttempts := 0
	for dairyCartIsDown {
		_, err := http.Get(url)
		if err != nil {
			log.Printf("error encountered: %v", err)
			log.Printf("waiting half a second before pinging Dairycart again")
			time.Sleep(500 * time.Millisecond)
			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				return errors.New("Maximum number of attempts made, something's gone awry")
			}
		} else {
			dairyCartIsDown = false
		}
	}
	return nil
}

func checkProductExistence(sku string) (string, error) {
	url := buildURL("product", sku)
	resp, err := http.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return turnResponseBodyIntoString(resp)
}
