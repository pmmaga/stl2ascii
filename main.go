package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
)

type Triangle struct {
	normal        [3]float32
	vertices      [3][3]float32
	attrByteCount uint16
}

type Model struct {
	header       string
	numTriangles uint32
	triangles    []Triangle
}

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

	//Create the reader
	fileReader := bufio.NewReader(fileHandle)

	//Check if it is ASCII
	asciiCheck, err := fileReader.Peek(5)
	check(err)

	if string(asciiCheck) == "solid" {
		//And blow up if it is!
		fmt.Println("ASCII reader not implemented yet")
		os.Exit(1)
	} else {
		fmt.Printf("Reading %v\n", cmdArgs[0])
		var model Model
		//Read the header
		byteHeader := make([]byte, 80)
		_, err = fileReader.Read(byteHeader)
		check(err)
		model.header = strings.Trim(string(byteHeader), "\x00")
		fmt.Printf("Header: %v\n", model.header)
		//Read the number of triangles
		err = binary.Read(fileReader, binary.LittleEndian, &model.numTriangles)
		check(err)
		//Print it
		fmt.Printf("Triangles: %v\n", model.numTriangles)
		//Initialize arrays for min x y z and max x y z
		mins := [...]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
		maxs := [...]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
		//Read the triangles
		for tri := uint32(0); tri < model.numTriangles; tri++ {
			var triangle Triangle
			//Read the normal
			for k := range triangle.normal {
				err = binary.Read(fileReader, binary.LittleEndian, &triangle.normal[k])
				check(err)
			}
			//Read the vertices
			for i := range triangle.vertices {
				for j := range triangle.vertices[i] {
					err = binary.Read(fileReader, binary.LittleEndian, &triangle.vertices[i][j])
					check(err)
					//Update min and max
					if triangle.vertices[i][j] < mins[j] {
						mins[j] = triangle.vertices[i][j]
					}
					if triangle.vertices[i][j] > maxs[j] {
						maxs[j] = triangle.vertices[i][j]
					}
				}
			}
			//Read the attribute byte count (which should be 0)
			err = binary.Read(fileReader, binary.LittleEndian, &triangle.attrByteCount)
			check(err)
			//If it isn't skip those bytes
			if triangle.attrByteCount != uint16(0) {
				_, err = fileHandle.Seek(int64(triangle.attrByteCount), 1)
				check(err)
			}
			//Apend the created Triangle to the Model
			model.triangles = append(model.triangles, triangle)
		}
		//Print the dimensions
		fmt.Printf("%v\n%v\n", mins, maxs)
		fmt.Printf("%v %v %v", maxs[0]-mins[0], maxs[1]-mins[1], maxs[2]-mins[2])
		//fmt.Printf("%v", model)
	}
}
