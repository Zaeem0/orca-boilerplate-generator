package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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
	template []interface{}
	config   []interface{}
}

func main() {
	//Octal value 0700 for user
	os.Mkdir("./generated-template", 0700)
	os.Mkdir("./generated-template/CreativeTemplates", 0700)
	os.Mkdir("./generated-template/FrameTemplates", 0700)
	os.Mkdir("./generated-template/GlobalTemplates", 0700)
	os.Mkdir("./generated-template/ThumbnailImages", 0700)

	var tempJSON []CreativeTemplates
	var obj CreativeTemplates

	//hardcoded values, use user input flags/os args for now then form inputs
	templateName := "statefarm"
	name := "statefarm-business"
	sizes := []string{"160x600", "300x250", "300x600"}
	start := []string{}
	middle := []string{"1", "2", "3"}
	end := []string{}
	limit := 5
	min := 1
	base := 6000

	for _, size := range sizes {
		obj = CreativeTemplates{}
		obj.Name = fmt.Sprintf("%s-%s", name, size)
		obj.Size = size
		//set width and height properties
		wxh := strings.Split(size, "x")
		width, _ := strconv.Atoi(wxh[0])
		height, _ := strconv.Atoi(wxh[1])
		obj.Width = width
		obj.Height = height

		obj.FrameLimit = limit
		obj.FrameMinCount = min

		//frames
		obj.Start = start
		obj.Middle = middle
		obj.End = end
		obj.BaseSize = base

		tempJSON = append(tempJSON, obj)
	}

	//indented JSON with 2 spaces
	data, _ := json.MarshalIndent(tempJSON, "", "  ")
	//need to create a file and folder rather than just writing to a test file
	_ = ioutil.WriteFile(fmt.Sprintf("./generated-template/CreativeTemplates/%s.json", templateName), data, 0644)
}
