package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mjevans93308/avoxi-demo-app/config"
)

var a App

func TestAliveHandler(t *testing.T) {
	a.Initialize(true)
	req, _ := http.NewRequest(http.MethodGet, config.Alive, nil)

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, rr.Code)
	}

	logger.Info(rr.Body.String())
	expected := "It's...ALIVE!!!"

	if body := rr.Body.String(); body != expected {
		t.Errorf("Expected %s. Got %s", expected, body)
	}
}

func TestCheckGeoLocation(t *testing.T) {
	a.Initialize(true)

	testBody := payload{
		Ip_address: "8.8.8.8",
		Country_names: []string{
			"Australia",
			"Bosnia",
		},
	}

	requestByte, _ := json.Marshal(testBody)

	req, _ := http.NewRequest(http.MethodPost, config.CheckIPLocation, bytes.NewReader(requestByte))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	if http.StatusNotFound != rr.Code {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusNotFound, rr.Code)
	}
}
