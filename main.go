package main

import (
	"bufio"
	"fmt"
	"os"

	"bitbucket.org/pmmaga/gostl/model"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	//Read the arguments
	cmdArgs := os.Args[1:]
	//If the number of arguments is not right, explode!
	if len(cmdArgs) != 1 {
		fmt.Println("You need to pass the name of the STL file.")
		os.Exit(1)
	}

	//Open the passed file for reading
	fileHandle, err := os.Open(cmdArgs[0])
	check(err)
	defer fileHandle.Close()
	fmt.Printf("Reading %v\n", cmdArgs[0])

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

	//Print the Model
	fmt.Println(&aModel)
	fmt.Println(aModel.Paint(100, model.PaintFromFront))
}
