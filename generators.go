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

func generateBoilerplate(data CreativeTemplateData) error {
	data.Frames = append(data.Frames, data.Start...)
	data.Frames = append(data.Frames, data.Middle...)
	data.Frames = append(data.Frames, data.End...)
	os.RemoveAll("./tmp")
	var folderStructure = []string{"./tmp", "./tmp/CreativeTemplates", "./tmp/FrameTemplates", "./tmp/GlobalTemplates", "./tmp/ThumbnailImages"}

	//Octal value 0700 for user
	for _, folder := range folderStructure {
		if err := os.Mkdir(folder, 0700); err != nil {
			return err
		}
	}

	var err error
	err = generateCreativeTemplates(data)
	if err != nil {
		return err
	}
	err = generateFrameTemplates(data.Sizes, data.Frames)
	if err != nil {
		return err
	}
	err = generateGlobalTemplates(data.Sizes)
	if err != nil {
		return err
	}
	err = generateThumbnails(data.Sizes, data.Frames)
	if err != nil {
		return err
	}

	return nil
}

func emptyFrame(clickout bool) []byte {
	var obj frame
	obj.Template = []script{}
	if clickout {
		obj.Template = append(obj.Template, script{T: "j", S: "ZG9jdW1lbnQuZ2V0RWxlbWVudEJ5SWQoJ2V4aXQnKS5hZGRFdmVudExpc3RlbmVyKCdjbGljaycsIGV4aXRDbGlja0hhbmRsZXIpOw=="})
	}
	obj.Config = []string{}
	jsonData, _ := json.MarshalIndent(obj, "", "  ")
	return jsonData
}

func generateThumbnails(sizes []string, frames []string) error {
	for _, size := range sizes {
		//.../ThumbnailImages/[SIZE]
		os.Mkdir(fmt.Sprintf("./tmp/ThumbnailImages/%s", size), 0700)
		for idx := range frames {
			err := createThumbnail(size, frames[idx])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func createThumbnail(size string, frameName string) error {
	//.../ThumbnailImages/[SIZE]/[SIZE]-[Frame].svg
	thumbs, err := os.Create(fmt.Sprintf("./tmp/ThumbnailImages/%s/%s-%s.svg", size, size, frameName))
	if err != nil {
		return err
	}

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

	return nil
}

func generateGlobalTemplates(sizes []string) error {
	for _, size := range sizes {
		//.../GlobalTemplates/[SIZE].json
		err := ioutil.WriteFile(fmt.Sprintf("./tmp/GlobalTemplates/%s.json", size), emptyFrame(true), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func generateFrameTemplates(sizes []string, frames []string) error {
	for _, size := range sizes {
		os.Mkdir(fmt.Sprintf("./tmp/FrameTemplates/%s", size), 0700)
		for f := range frames {
			//.../FrameTemplates/[SIZE]/[SIZE]-[Frame].json
			err := ioutil.WriteFile(fmt.Sprintf("./tmp/FrameTemplates/%s/%s-%s.json", size, size, frames[f]), emptyFrame(false), 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func generateCreativeTemplates(data CreativeTemplateData) error {
	var tempJSON []CreativeTemplates
	var obj CreativeTemplates

	for _, size := range data.Sizes {
		obj = CreativeTemplates{}
		obj.Name = fmt.Sprintf("%s-%s", data.Name, size)
		obj.Size = size
		//set width and height properties
		wxh := strings.Split(size, "x")
		width, err := strconv.Atoi(wxh[0])
		if err != nil {
			return err
		}
		height, err := strconv.Atoi(wxh[1])
		if err != nil {
			return err
		}
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
	err := ioutil.WriteFile(fmt.Sprintf("./tmp/CreativeTemplates/%s.json", data.TemplateSet), jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
