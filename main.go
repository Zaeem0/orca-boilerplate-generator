package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
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
	Template []string `json:"template"`
	Config   []string `json:"config"`
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

func generateBoilerplate(data CreativeTemplateData) {
	data.Frames = append(data.Frames, data.Start...)
	data.Frames = append(data.Frames, data.Middle...)
	data.Frames = append(data.Frames, data.End...)

	//Octal value 0700 for user
	os.RemoveAll("./tmp")
	os.Mkdir("./tmp", 0700)
	os.Mkdir("./tmp/CreativeTemplates", 0700)
	os.Mkdir("./tmp/FrameTemplates", 0700)
	os.Mkdir("./tmp/GlobalTemplates", 0700)
	os.Mkdir("./tmp/ThumbnailImages", 0700)

	generateCreativeTemplates(data)
	generateFrameTemplates(data.Sizes, data.Frames)
	generateGlobalTemplates(data.Sizes)
	generateThumbnails(data.Sizes, data.Frames)

	os.RemoveAll(fmt.Sprintf("./%s", data.TemplateGroupName))
	os.Rename("./tmp", fmt.Sprintf("./%s", data.TemplateGroupName))
}

func ReceiveData(w http.ResponseWriter, r *http.Request) {
	var data CreativeTemplateData
	if DecodeBody(w, r, &data) != nil {
		return
	}
	log.Println(data)
	generateBoilerplate(data)
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

func emptyFrame() []byte {
	var obj frame
	obj.Template = []string{}
	obj.Config = []string{}
	jsonData, _ := json.MarshalIndent(obj, "", "  ")
	return jsonData
}

func generateThumbnails(sizes []string, frames []string) {
	for _, size := range sizes {
		//.../ThumbnailImages/[SIZE]
		os.Mkdir(fmt.Sprintf("./tmp/ThumbnailImages/%s", size), 0700)
		for idx := range frames {
			createThumbnail(size, frames[idx])
		}
	}
}

func createThumbnail(size string, frameName string) {
	//.../ThumbnailImages/[SIZE]/[SIZE]-[Frame].svg
	thumbs, _ := os.Create(fmt.Sprintf("./tmp/ThumbnailImages/%s/%s-%s.svg", size, size, frameName))

	w := bufio.NewWriter(thumbs)
	canvas := svg.New(w)
	width := 300
	height := 250
	canvas.Start(width, height)
	canvas.Text(width/2, height/2, fmt.Sprintf("%s", frameName), "text-anchor:middle;font-size:12px;fill:black")
	canvas.End()
	w.Flush()
}

func generateGlobalTemplates(sizes []string) {
	for _, size := range sizes {
		//.../GlobalTemplates/[SIZE].json
		_ = ioutil.WriteFile(fmt.Sprintf("./tmp/GlobalTemplates/%s.json", size), emptyFrame(), 0644)
	}
}

func generateFrameTemplates(sizes []string, frames []string) {
	for _, size := range sizes {
		os.Mkdir(fmt.Sprintf("./tmp/FrameTemplates/%s", size), 0700)
		for f := range frames {
			//.../FrameTemplates/[SIZE]/[SIZE]-[Frame].json
			_ = ioutil.WriteFile(fmt.Sprintf("./tmp/FrameTemplates/%s/%s-%s.json", size, size, frames[f]), emptyFrame(), 0644)
		}
	}
}

func generateCreativeTemplates(data CreativeTemplateData) {
	var tempJSON []CreativeTemplates
	var obj CreativeTemplates

	for _, size := range data.Sizes {
		obj = CreativeTemplates{}
		obj.Name = fmt.Sprintf("%s-%s", data.Name, size)
		obj.Size = size
		//set width and height properties
		wxh := strings.Split(size, "x")
		width, _ := strconv.Atoi(wxh[0])
		height, _ := strconv.Atoi(wxh[1])
		obj.Width = width
		obj.Height = height

		obj.FrameLimit = data.Limit
		obj.FrameMinCount = data.Min

		//frames
		obj.Start = data.Start
		obj.Middle = data.Middle
		obj.End = data.End
		obj.BaseSize = data.Base

		tempJSON = append(tempJSON, obj)
	}

	//indented JSON with 2 spaces
	jsonData, _ := json.MarshalIndent(tempJSON, "", "  ")
	_ = ioutil.WriteFile(fmt.Sprintf("./tmp/CreativeTemplates/%s.json", data.TemplateSet), jsonData, 0644)
}
