package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	// "math/rand"
	// "time"
)

type BoundingBox struct {
	xmax float64
	xmin float64
	ymax float64
	ymin float64
}
type Shapefile struct {
	BBox          BoundingBox
	GeometryType  uint32
	recordNumber  uint32
	ContentLength uint32
	fileLength    int64
	geom          []interface{}
}

var shapefile Shapefile
var totalRecord, dataLength uint32

func main() {
	readFile()

}

func readFile() {
	var fileName string = "sampleData/polyline/line.shp"
	file, errFile := os.Open(fileName)
	stat, errStat := os.Stat(fileName)
	defer file.Close()
	if errFile != nil {
		log.Fatal(errStat)
	}

	m := []byte{}
	fileLength := stat.Size()
	for i := int64(0); i < fileLength; i++ {
		m = append(m, readNextBytes(file, 1)...)
	}
	shapefile.fileLength = fileLength
	GeomType := binary.LittleEndian.Uint32(m[32:36])
	shapefile.GeometryType = GeomType
	shapefile.ContentLength = binary.BigEndian.Uint32(m[24:28])
	shapefile.BBox = ReadExtent(m[36:60])

	fmt.Println(shapefile.BBox)
	var ContentLength int
	a, b := 0, 0

	for {
		headerContent := m[100+a : 108+a]
		// recordNumber := headerContent[:4]
		ContentLength = int(ContentLength) + int(binary.BigEndian.Uint32(headerContent[4:8])*2) + 8
		a = ContentLength
		geomData := m[108+b : 108+a-8][4:]
		b = a
		// fmt.Println(recordNumber)
		// ReadPoint(geomData)
		geometry := ReadLine(geomData)
		// geometry := ReadPoint(geomData)
		shapefile.geom = append(shapefile.geom, geometry)
		if ContentLength == int(fileLength)-100 {
			break
		}

	}
	fmt.Println(shapefile.geom)
}

// func ParseShape(GeomType int)
func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes
}

func Float64Parse(data []byte) float64 {
	var floatValue64 float64
	binary.Read(bytes.NewReader(data), binary.LittleEndian, &floatValue64)
	return floatValue64
}
func ReadExtent(data []byte) BoundingBox {
	var Target BoundingBox
	var xmin float64 = Float64Parse(data[0:8])
	var ymin float64 = Float64Parse(data[8:16])
	var xmax float64 = Float64Parse(data[16:24])
	var ymax float64 = Float64Parse(data[24:32])
	Target.xmax = xmax
	Target.xmin = xmin
	Target.ymax = ymax
	Target.ymin = ymin
	return Target
}

// Point Identifier
type Point struct {
	x float64
	y float64
}

// bakal ada structure untuk header isinya record number
func ReadPoint(data []byte) Point {
	var geom Point
	geom.x = Float64Parse(data[0:8])
	geom.y = Float64Parse(data[8:16])
	return geom
}

// ------------

// PolyLine Identifier
type PolyLine struct {
	Box       BoundingBox
	NumParts  uint32
	NumPoints uint32
	Part      uint32
	Point     []Point
}

func ReadLine(data []byte) PolyLine {
	var geom PolyLine

	geom.Box = ReadExtent(data[0:32])
	geom.NumParts = binary.LittleEndian.Uint32(data[32:36])
	geom.NumPoints = binary.LittleEndian.Uint32(data[36:40])
	geom.Point = make([]Point, int(len(data[40:])/16))

	for i := 0; i < int(geom.NumPoints); i++ {
		geom.Point[i] = ReadPoint(data[44+i*16 : 44+i*16+16])
	}
	return geom
}

type Polygon struct {
	Box       BoundingBox
	NumParts  uint32
	NumPoints uint32
	Parts     []uint32
	Point     []Point
}

// func CheckType(GeomType int64) string {
// 	var stringType
// 	switch GeomType{
// 		case
// 	}
// }
