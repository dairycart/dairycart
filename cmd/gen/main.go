package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	maxAttempts       = 10
	currentAPIVersion = `v1`
	validPassword     = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

var (
	baseURL string
)

func init() {
	baseURL = os.Getenv("DAIRYCART_API_URL")
}

func mapToQueryValues(in map[string]string) string {
	out := url.Values{}
	for k, v := range in {
		out.Set(k, v)
	}
	return out.Encode()
}

func buildPath(parts ...interface{}) string {
	stringParts := []string{}
	for _, p := range parts {
		switch part := p.(type) {
		case string:
			stringParts = append(stringParts, part)
		case uint64:
			stringParts = append(stringParts, fmt.Sprintf("%d", part))
		}
	}

	return fmt.Sprintf("%s/%s/%s", baseURL, currentAPIVersion, strings.Join(stringParts, "/"))
}

func buildURL(path string, queryParams map[string]string) string {
	u, _ := url.Parse(path)
	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString
	return u.String()
}

func ensureThatDairycartIsAlive() {
	path := buildPath("health")
	u := buildURL(path, nil)
	dairyCartIsDown := true
	numberOfAttempts := 0
	for dairyCartIsDown {
		_, err := http.Get(u)
		if err != nil {
			log.Printf("waiting half a second before pinging Dairycart again")
			time.Sleep(500 * time.Millisecond)
			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				log.Fatalf("Maximum number of attempts made, something's gone awry")
			}
		} else {
			dairyCartIsDown = false
		}
	}
}

func main() {
	ensureThatDairycartIsAlive()
	var (
		cmd *exec.Cmd
		err error
	)

	cmd = exec.Command("gnorm", "gen")
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Printf("%+v", err)
		panic(err)
	}
}
