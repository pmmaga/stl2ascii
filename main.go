package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"bitbucket.org/pmmaga/gostl/model"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func usage() {
	fmt.Println("usage: gostl [info|paint] [front|side|top size] [pathtostl]")
	os.Exit(1)
}

func main() {
	//Read the arguments
	cmdArgs := os.Args[1:]
	//If the number of arguments is enough, show usage
	if len(cmdArgs) < 2 {
		usage()
	}

	//Open the passed file for reading
	fileHandle, err := os.Open(cmdArgs[len(cmdArgs)-1])
	check(err)
	defer fileHandle.Close()
	fmt.Printf("Reading %v\n", cmdArgs[len(cmdArgs)-1])

	//Create the reader
	fileReader := bufio.NewReader(fileHandle)

	//Check if it is ASCII
	asciiCheck, err := fileReader.Peek(5)
	check(err)

	var aModel model.Model
	if string(asciiCheck) == "solid" {
		aModel, err = model.CreateFromASCIISTL(fileReader)
		check(err)

	} else {
		aModel, err = model.CreateFromBinarySTL(fileReader)
		check(err)
	}

	switch cmdArgs[0] {
	case "info":
		//Print the Model Info
		fmt.Println(&aModel)
	case "paint":
		//Check the perspective and size params
		var perspective model.PaintFrom
		switch cmdArgs[1] {
		case "front":
			perspective = model.PaintFromFront
		case "side":
			perspective = model.PaintFromSide
		case "top":
			perspective = model.PaintFromTop
		default:
			usage()
		}
		size, err := strconv.ParseInt(cmdArgs[2], 10, 0)
		if err != nil {
			usage()
		}
		//Paint the model
		fmt.Println(aModel.Paint(int(size), perspective))
	default:
		usage()
	}
}
