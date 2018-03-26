package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestHello(t *testing.T) {

	response, err := http.Get("http://127.0.0.1:8000/hello")
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	} else {
		defer response.Body.Close()

		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}

		// Check the status code is what we expect.
		if status := response.StatusCode; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		// test hello world api
		expected := `{"Hello":"world"}`
		contents, _ := ioutil.ReadAll(response.Body)
		if string(contents) != expected {
			t.Errorf("handler returned unexpected body: got '%v' want '%v'", string(contents), expected)
		}
	}
}
