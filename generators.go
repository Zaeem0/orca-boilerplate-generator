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

func emptyFrame(clickout bool) []byte {
	var obj frame
	obj.Template = []script{}
	if clickout {
		obj.Template = append(obj.Template, script{T: "s", S: "ZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2V4aXQnKS5hZGRFdmVudExpc3RlbmVyKCdjbGljaycsIGV4aXRDbGlja0hhbmRsZXIpOw=="})
	}
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

	wxh := strings.Split(size, "x")
	width, _ := strconv.Atoi(wxh[0])
	height, _ := strconv.Atoi(wxh[1])

	canvas.Start(width, height)
	canvas.Rect(0, 0, width, height, "fill:white")
	canvas.Text(width/2, height/2, fmt.Sprintf("%s", frameName), "text-anchor:middle;font-size:40px;fill:black;font-family: Helvetica;")
	canvas.End()
	w.Flush()
}

func generateGlobalTemplates(sizes []string) {
	for _, size := range sizes {
		//.../GlobalTemplates/[SIZE].json
		_ = ioutil.WriteFile(fmt.Sprintf("./tmp/GlobalTemplates/%s.json", size), emptyFrame(true), 0644)
	}
}

func generateFrameTemplates(sizes []string, frames []string) {
	for _, size := range sizes {
		os.Mkdir(fmt.Sprintf("./tmp/FrameTemplates/%s", size), 0700)
		for f := range frames {
			//.../FrameTemplates/[SIZE]/[SIZE]-[Frame].json
			_ = ioutil.WriteFile(fmt.Sprintf("./tmp/FrameTemplates/%s/%s-%s.json", size, size, frames[f]), emptyFrame(false), 0644)
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
