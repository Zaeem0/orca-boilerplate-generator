package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type CreativeTemplates struct {
	Name          string   `json:"name"`
	Size          string   `json:"size"`
	Width         int      `json:"width"`
	Height        int      `json:"height"`
	FrameLimit    int      `json:"frameLimit"`
	FrameMinCount int      `json:"frameMinCount"`
	Start         []string `json:"start"`
	Middle        []string `json:"middle"`
	End           []string `json:"end"`
	BaseSize      int      `json:"baseSize"`
}

type frame struct {
	Template []script `json:"template"`
	Config   []string `json:"config"`
}

type script struct {
	T string `json:"t"`
	S string `json:"s"`
}

type CreativeTemplateData struct {
	TemplateGroupName string   `json:"templateGroupName"`
	TemplateSet       string   `json:"templateSet"`
	Sizes             []string `json:"sizes"`
	Name              string   `json:"templateName"`
	Limit             int      `json:"frameLimit"`
	Min               int      `json:"frameMinCount"`
	Start             []string `json:"start"`
	Middle            []string `json:"middle"`
	End               []string `json:"end"`
	Base              int      `json:"baseSize"`

	Frames []string
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.
		Name("posting data").
		Methods("POST").
		Path("/").
		HandlerFunc(ReceiveData)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func ReceiveData(w http.ResponseWriter, r *http.Request) {
	var data CreativeTemplateData
	if DecodeBody(w, r, &data) != nil {
		return
	}
	log.Println(data)

	generateBoilerplate(data)
	createZip(data.TemplateGroupName)

	w.WriteHeader(http.StatusCreated)
	ResponseJSON(w, data)
}

func DecodeBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		ResponseJSON(w, err)
		return err
	}
	return nil
}

func ResponseJSON(w http.ResponseWriter, data interface{}) {
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println(err)
	}
}
