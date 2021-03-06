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

const (
	STATIC_DIR = "/client/"
)

var logToFile *log.Logger

func main() {
	f, err := os.OpenFile("text.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logToFile = log.New(f, "", log.LstdFlags)

	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.
		Name("index").
		Methods("GET").
		Path("/api/").
		HandlerFunc(Index)
	apiRouter.
		Name("posting data").
		Methods("POST").
		Path("/api/").
		HandlerFunc(ReceiveData)
	apiRouter.
		Name("download zip").
		Methods("GET").
		Path("/api/download/{templateGroupName}").
		HandlerFunc(DownloadZip)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("client")))
	mux.Handle("/api/", apiRouter)

	methods := []string{"GET", "POST", "PUT", "DELETE"}
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(handlers.AllowedMethods(methods))(mux)))
}

func Index(w http.ResponseWriter, r *http.Request) {
	ResponseJSON(w, []string{"Index Page"})
}

func DownloadZip(w http.ResponseWriter, r *http.Request) {
	zipfile := mux.Vars(r)["templateGroupName"]
	logToFile.Printf("Attempt made to download: %s\n", zipfile)
	f, err := os.Open(fmt.Sprintf("./downloads/%s.zip", zipfile))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		ResponseJSON(w, []string{fmt.Sprintf("%s.zip does not exist", zipfile)})
		return
	}

	_, file := filepath.Split(f.Name())

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file))

	// copy the file contents to the http Writer
	_, err = io.Copy(w, f)
	if err != nil {
		logToFile.Printf("Failed to download %s\n", zipfile)
		log.Fatal(err)
	}
	logToFile.Printf("Successfully downloaded %s\n", zipfile)
}

func ReceiveData(w http.ResponseWriter, r *http.Request) {
	var data CreativeTemplateData
	if DecodeBody(w, r, &data) != nil {
		return
	}
	data.cleanData()
	log.Println(data)
	err := generateBoilerplate(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ResponseJSON(w, err)
		logToFile.Printf("Failed to generate boilerplate for: %s\n", data.TemplateGroupName)
		return
	}
	logToFile.Printf("Succesfully generated boilerplate for: %s\n", data.TemplateGroupName)

	err = createZip(data.TemplateGroupName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ResponseJSON(w, err)
		return
	}
	logToFile.Printf("Succesfully generated ZIP for: %s\n", data.TemplateGroupName)
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
