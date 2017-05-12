package main

import (
	"net/http"
	"strconv"
)

// DetermineRequestLimits determines requested limits
func DetermineRequestLimits(req *http.Request) int {
	var actualLimit int
	requestedLimit := req.URL.Query().Get("limit")
	actualLimit, err := strconv.Atoi(requestedLimit)

	if requestedLimit == "" || err != nil {
		actualLimit = 25
		req.URL.Query().Set("limit", "25")
	}

	return actualLimit
}
