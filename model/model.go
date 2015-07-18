package model

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
)

type triangle struct {
	normal        [3]float32
	vertices      [3][3]float32
	attrByteCount uint16
}

type Model struct {
	Header       string
	NumTriangles uint32
	Triangles    []triangle
}

func (m *Model) String() string {
	return fmt.Sprintf("Header: %v\nTriangles: %v\nDimensions: %v\n", m.Header, m.NumTriangles, GetDimensions(m))
}

func CreateFromBinarySTL(r *bufio.Reader) (m Model, err error) {
	//Read the header
	byteHeader := make([]byte, 80)
	_, err = r.Read(byteHeader)
	if err != nil {
		return m, err
	}

	m.Header = strings.Trim(string(byteHeader), "\x00")
	//Read the number of triangles
	err = binary.Read(r, binary.LittleEndian, &m.NumTriangles)
	if err != nil {
		return m, err
	}
	//Read the triangles
	for tri := uint32(0); tri < m.NumTriangles; tri++ {
		var aTriangle triangle
		//Read the normal
		for k := range aTriangle.normal {
			err = binary.Read(r, binary.LittleEndian, &aTriangle.normal[k])
			if err != nil {
				return m, err
			}
		}
		//Read the vertices
		for i := range aTriangle.vertices {
			for j := range aTriangle.vertices[i] {
				err = binary.Read(r, binary.LittleEndian, &aTriangle.vertices[i][j])
				if err != nil {
					return m, err
				}
			}
		}
		//Read the attribute byte count (which should be 0)
		err = binary.Read(r, binary.LittleEndian, &aTriangle.attrByteCount)
		if err != nil {
			return m, err
		}
		//If it isn't skip those bytes
		if aTriangle.attrByteCount != uint16(0) {
			attr := make([]byte, aTriangle.attrByteCount)
			err = binary.Read(r, binary.LittleEndian, &attr)
			if err != nil {
				return m, err
			}
		}
		//Apend the created Triangle to the Model
		m.Triangles = append(m.Triangles, aTriangle)
	}
	return m, nil
}

func CreateFromASCIISTL(r *bufio.Reader) (m Model, err error) {
	//Create the header
	header, err := r.ReadBytes('\n')
	if err != nil {
		return m, err
	}
	m.Header = fmt.Sprintf("Imported from ASCII STL by gostl - %v", strings.Trim(string(header[6:]), "\n"))
	for {
		line, err := r.ReadBytes('\n')
		if err != nil || len(line) < 5 {
			break
		}
		if string(line[:5]) != "facet" {
			break
		}
	}

	return m, nil
}

func GetDimensions(m *Model) [3]float32 {
	mins, maxs := getMinsMaxs(m)
	return [3]float32{maxs[0] - mins[0], maxs[1] - mins[1], maxs[2] - mins[2]}
}

func getMinsMaxs(m *Model) (mins [3]float32, maxs [3]float32) {
	//Initialize arrays for min x y z and max x y z
	mins = [3]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	maxs = [3]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
	//Run through the triangles
	for i := range m.Triangles {
		//Each vertice
		for j := range m.Triangles[i].vertices {
			//Each coordinate
			for k := range m.Triangles[i].vertices[j] {
				//Update min and max
				if m.Triangles[i].vertices[j][k] < mins[k] {
					mins[k] = m.Triangles[i].vertices[j][k]
				}
				if m.Triangles[i].vertices[j][k] > maxs[k] {
					maxs[k] = m.Triangles[i].vertices[j][k]
				}
			}
		}
	}
	return mins, maxs
}
