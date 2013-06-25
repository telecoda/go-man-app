package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func addResponseHeaders(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	// allow cross origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

}

func httpErrorResponse(w http.ResponseWriter, err error) {

	http.Error(w, err.Error(), http.StatusInternalServerError)

}

func getRequestBody(request *http.Request) ([]byte, error) {

	jsonBody, err := ioutil.ReadAll(request.Body)

	return jsonBody, err
}

func OptionsHandler(w http.ResponseWriter, r *http.Request) {

	addResponseHeaders(w)
	w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT")

	fmt.Fprint(w, Response{"success": true, "message": "Welcome to go-man options", "method": r.Method})
}
