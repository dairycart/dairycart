package main

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/dairycart/dairyclient/v1"
)

func getCookieFromRequest(req *http.Request) (*http.Cookie, error) {
	cookies := req.Cookies()
	if len(cookies) == 0 {
		return nil, errors.New("No cookies found in request")
	}

	if debug {
		log.Printf("found %d cookies associated with this request\n", len(cookies))
	}

	for _, c := range cookies {
		if c.Name == cookieName {
			return c, nil
		}
	}
	return nil, errors.New("no dairycart cookie found in request")
}

func buildClientFromRequest(res http.ResponseWriter, req *http.Request) (*dairyclient.V1Client, error) {
	dairyCookie, err := getCookieFromRequest(req)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return nil, err
	}

	dairyClient, err := dairyclient.NewV1ClientFromCookie(apiURL, dairyCookie, http.DefaultClient)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return nil, err
	}
	return dairyClient, err
}
