package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func EmptyCreativeTemplateObject() CreativeTemplates {
	obj := CreativeTemplates{}
	obj.Start = []string{}
	obj.Middle = []string{}
	obj.End = []string{}
	obj.BaseSize = 8000
	return obj
}

func main() {
	file, _ := ioutil.ReadFile("creativetemplate.json")
	fmt.Println(string(file))

	var tempJSON []CreativeTemplates
	json.Unmarshal(file, &tempJSON)

	// name := "statefarm-business"
	// sizes := []string{"160x600", "300x250", "300x600"}
	// limit := 5
	// min := 1

	// fmt.Printf("%#v", CreativeTemplate[0])
	// CreativeTemplate[0].Name = "Changed name"
	// frames := []string{}
	// CreativeTemplate[0].Start = frames
	tempJSON = append(tempJSON, EmptyCreativeTemplateObject())
	// fmt.Println("\n", CreativeTemplate[0].Start)
	// fmt.Println("\n", CreativeTemplate[0])

	file, _ = json.MarshalIndent(tempJSON, "", "  ")
	_ = ioutil.WriteFile("creativetemplate.json", file, 0644)
}
