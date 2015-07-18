package model

import (
	"encoding/binary"
	"io"
	"math"
	"strings"
)

type Triangle struct {
	normal        [3]float32
	vertices      [3][3]float32
	attrByteCount uint16
}

type Model struct {
	Header       string
	NumTriangles uint32
	Triangles    []Triangle
}

func (m *Model) ReadBinarySTL(r io.Reader) (err error) {
	//Read the header
	byteHeader := make([]byte, 80)
	_, err = r.Read(byteHeader)
	if err != nil {
		return err
	}

	m.Header = strings.Trim(string(byteHeader), "\x00")
	//Read the number of triangles
	err = binary.Read(r, binary.LittleEndian, &m.NumTriangles)
	if err != nil {
		return err
	}
	//Read the triangles
	for tri := uint32(0); tri < m.NumTriangles; tri++ {
		var triangle Triangle
		//Read the normal
		for k := range triangle.normal {
			err = binary.Read(r, binary.LittleEndian, &triangle.normal[k])
			if err != nil {
				return err
			}
		}
		//Read the vertices
		for i := range triangle.vertices {
			for j := range triangle.vertices[i] {
				err = binary.Read(r, binary.LittleEndian, &triangle.vertices[i][j])
				if err != nil {
					return err
				}
			}
		}
		//Read the attribute byte count (which should be 0)
		err = binary.Read(r, binary.LittleEndian, &triangle.attrByteCount)
		if err != nil {
			return err
		}
		//If it isn't skip those bytes
		if triangle.attrByteCount != uint16(0) {
			attr := make([]byte, triangle.attrByteCount)
			err = binary.Read(r, binary.LittleEndian, &attr)
			if err != nil {
				return err
			}
		}
		//Apend the created Triangle to the Model
		m.Triangles = append(m.Triangles, triangle)
	}
	return nil
}

func GetDimensions(m *Model) [3]float32 {
	//Initialize arrays for min x y z and max x y z
	mins := [...]float32{math.MaxFloat32, math.MaxFloat32, math.MaxFloat32}
	maxs := [...]float32{-math.MaxFloat32, -math.MaxFloat32, -math.MaxFloat32}
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
	return [3]float32{maxs[0] - mins[0], maxs[1] - mins[1], maxs[2] - mins[2]}
}
