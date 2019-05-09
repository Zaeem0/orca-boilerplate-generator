package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/handlers"
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
	router.
		Name("download zip").
		Methods("GET").
		Path("/download/{templateGroupName}").
		HandlerFunc(DownloadZip)

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedMethods(methods))(router)))
}

func DownloadZip(w http.ResponseWriter, r *http.Request) {
	zipfile := mux.Vars(r)["templateGroupName"]
	f, err := os.Open(fmt.Sprintf("./downloads/%s.zip", zipfile))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		ResponseJSON(w, []string{fmt.Sprintf("Failed to download: %s.zip", zipfile)})
		return
	}

	_, file := filepath.Split(f.Name())

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file))

	// copy the file contents to the http Writer
	_, err = io.Copy(w, f)
	if err != nil {
		log.Fatal(err)
	}
}

func ReceiveData(w http.ResponseWriter, r *http.Request) {
	var data CreativeTemplateData
	if DecodeBody(w, r, &data) != nil {
		return
	}
	data.cleanData()
	log.Println(data)
	generateBoilerplate(data)
	createZip(data.TemplateGroupName)

	w.WriteHeader(http.StatusCreated)
	ResponseJSON(w, data)
}

func (data *CreativeTemplateData) cleanData() {
	if data.TemplateGroupName == "" {
		data.TemplateGroupName = "unknown"
	}
	//More checking of data to be done...
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
