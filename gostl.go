package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"strconv"

	"github.com/pmmaga/gostl/model"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func usage() {
	fmt.Println("usage: gostl [flags] [pathtofile]")
	flag.PrintDefaults()
	os.Exit(1)
}

var (
	// Command Flag Declaration
	info = flag.Bool("i", true, "Print gathered information about the file")
	draw = flag.Bool("d", true, "Draw the model from a direction on a size x size grid (-d [front|side|top] size)")

	// Option Flags
	preLoad = flag.Bool("pl", false, "Preload the file into memory (binary STL only)")

	// Debugging
	cpuprofile = flag.String("cpuprofile", "", "Write cpu profile to this file")
)

func main() {
	//Read the flags
	flag.Usage = usage
	flag.Parse()

	//Create a CPU profile for debugging
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "creating cpu profile: %s\n", err)
			check(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	//for runs := 0; runs < 50; runs++ {

	//Create the model
	var aModel model.Model

	//File path
	filePath := flag.Arg(flag.NArg() - 1)
	if filePath == "" {
		usage()
	}

	//If we want to preload the model in memory
	if *preLoad {
		//Load the whole file to memory
		fileSlice, err := ioutil.ReadFile(filePath)
		check(err)
		//Create the model from it
		aModel, err = model.CreateFromByteSlice(fileSlice)
		check(err)
	} else {
		//Open the passed file for reading
		fileHandle, err := os.Open(filePath)
		check(err)
		defer fileHandle.Close()

		//Create the reader
		fileReader := bufio.NewReaderSize(fileHandle, 50*1000)

		//Check if it is ASCII
		asciiCheck, err := fileReader.Peek(5)
		check(err)

		//Try ASCII if it looks like it
		if string(asciiCheck) == "solid" {
			//Discard the error so we try binary if this fails
			aModel, _ = model.CreateFromASCIISTL(fileReader)
		}
		//If it failed, try binary
		if len(aModel.Triangles) == 0 {
			//Reset the reader in case it tried ASCII
			if string(asciiCheck) == "solid" {
				fileHandle.Seek(0, 0)
				fileReader = bufio.NewReader(fileHandle)
			}
			//Read the Binary STL
			aModel, err = model.CreateFromBinarySTL(fileReader)
			check(err)
		}
	}

	if *info {
		//Print the Model Info
		fmt.Println(&aModel)
	}

	if *draw {
		var perspective model.ProjectFrom
		var size int64

		if flag.NArg() == 1 {
			perspective = model.ProjectFromFront
			size = 160
		} else {
			//Check the perspective and size params
			switch flag.Arg(0) {
			case "front":
				perspective = model.ProjectFromFront
			case "side":
				perspective = model.ProjectFromSide
			case "top":
				perspective = model.ProjectFromTop
			default:
				usage()
			}
			var err error
			size, err = strconv.ParseInt(flag.Arg(1), 10, 0)
			if err != nil {
				usage()
			}
		}
		//Paint the model
		fmt.Println(model.DrawMatrix(model.ProjectModelVertices(&aModel, int(size), perspective)))
	}
	// }
}
