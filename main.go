package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	svg "github.com/ajstarks/svgo"
)

type CreativeTemplates struct {
	Name          string
	Size          string
	Width         int
	Height        int
	FrameLimit    int
	FrameMinCount int
	Start         []string
	Middle        []string
	End           []string
	BaseSize      int
}

type frame struct {
	Template []string `json:"template"`
	Config   []string `json:"config"`
}
type CreativeTemplateData struct {
	templateGroupName string
	templateSet       string
	sizes             []string
	name              string
	width             int
	height            int
	limit             int
	min               int
	start             []string
	middle            []string
	end               []string
	frames            []string
	base              int
}

func main() {
	//hardcoded values, use user input flags/os args for now then form inputs
	data := CreativeTemplateData{
		templateGroupName: "statefarm-business-orca",
		templateSet:       "statefarm",
		name:              "some-business",
		sizes:             []string{"160x600", "300x250", "300x600"},
		start:             []string{"another one", "start"},
		middle:            []string{"1-mid", "copy-2", "3"},
		end:               []string{"endframe", "99"},
		limit:             5,
		min:               1,
		base:              6000,
	}
	data.frames = append(data.frames, data.start...)
	data.frames = append(data.frames, data.middle...)
	data.frames = append(data.frames, data.end...)

	//Octal value 0700 for user
	os.RemoveAll("./tmp")
	os.Mkdir("./tmp", 0700)
	os.Mkdir("./tmp/CreativeTemplates", 0700)
	os.Mkdir("./tmp/FrameTemplates", 0700)
	os.Mkdir("./tmp/GlobalTemplates", 0700)
	os.Mkdir("./tmp/ThumbnailImages", 0700)

	generateCreativeTemplates(data)
	generateFrameTemplates(data.sizes, data.frames)
	generateGlobalTemplates(data.sizes)
	generateThumbnails(data.sizes, data.frames)

	os.RemoveAll(fmt.Sprintf("./%s", data.templateGroupName))
	os.Rename("./tmp", fmt.Sprintf("./%s", data.templateGroupName))
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

	for _, size := range data.sizes {
		obj = CreativeTemplates{}
		obj.Name = fmt.Sprintf("%s-%s", data.name, size)
		obj.Size = size
		//set width and height properties
		wxh := strings.Split(size, "x")
		width, _ := strconv.Atoi(wxh[0])
		height, _ := strconv.Atoi(wxh[1])
		obj.Width = width
		obj.Height = height

		obj.FrameLimit = data.limit
		obj.FrameMinCount = data.min

		//frames
		obj.Start = data.start
		obj.Middle = data.middle
		obj.End = data.end
		obj.BaseSize = data.base

		tempJSON = append(tempJSON, obj)
	}

	//indented JSON with 2 spaces
	jsonData, _ := json.MarshalIndent(tempJSON, "", "  ")
	_ = ioutil.WriteFile(fmt.Sprintf("./tmp/CreativeTemplates/%s.json", data.templateSet), jsonData, 0644)
}
