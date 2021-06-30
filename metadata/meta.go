package metadata

import (
	"fmt"
	"os"
	"time"

	"github.com/flywave/go-geobuf"
	"github.com/flywave/go-geobuf/io"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

type Meta struct {
	Type          string
	Verticies     int
	Properties    int
	SizeJSON      int
	TimeReadJSON  int
	TimeWriteJSON int
	SizeBUF       int
	TimeReadBUF   int
	TimeWriteBUF  int
}

func (meta *Meta) MakeString() string {
	return fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d,%d,%d\n", meta.Type, meta.Verticies, meta.Properties,
		meta.SizeJSON, meta.TimeReadJSON, meta.TimeWriteJSON, meta.SizeBUF, meta.TimeReadBUF, meta.TimeWriteBUF,
	)
}

type MetaCSV struct {
	File     *os.File
	FileName string
}

func NewMetaCSV(filename string) *MetaCSV {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	metacsv := &MetaCSV{File: file, FileName: filename}
	metacsv.File.WriteString("type,verticies,number_properties,size_json,time_read_json,time_write_json,size_buf,time_read_buf,time_write_buf\n")
	return metacsv
}

func (metacsv *MetaCSV) AddMeta(feature *geom.Feature) {
	featurestring := MakeMeta(feature).MakeString()
	metacsv.File.WriteString(featurestring)
}

func MakeMeta(feature *geom.Feature) *Meta {
	meta := &Meta{}
	meta.Type = string(feature.Geometry.GetType())
	var total int
	switch g := feature.Geometry.(type) {
	case geom.Point:
	case geom.Point3:
		meta.Verticies = 1
	case geom.MultiPoint:
		meta.Verticies = len(g.Data())
	case geom.MultiPoint3:
		meta.Verticies = len(g.Data())
	case geom.LineString:
		meta.Verticies = len(g.Data())
	case geom.LineString3:
		meta.Verticies = len(g.Data())
	case geom.MultiLine:
		for _, line := range g.Data() {
			total += len(line)
		}
		meta.Verticies = total
	case geom.MultiLine3:
		for _, line := range g.Data() {
			total += len(line)
		}
		meta.Verticies = total
	case geom.Polygon:
		for _, line := range g.Data() {
			total += len(line)
		}
		meta.Verticies = total
	case geom.Polygon3:
		for _, line := range g.Data() {
			total += len(line)
		}
		meta.Verticies = total
	case geom.MultiPolygon:
		for _, polygon := range g.Data() {
			for _, line := range polygon {
				total += len(line)
			}
		}
		meta.Verticies = total
	case geom.MultiPolygon3:
		for _, polygon := range g.Data() {
			for _, line := range polygon {
				total += len(line)
			}
		}
		meta.Verticies = total
	}
	meta.Properties = len(feature.Properties)
	s := time.Now()
	bytevals, err := feature.MarshalJSON()
	meta.TimeWriteJSON = int(time.Now().Sub(s).Nanoseconds())
	if err != nil {
		fmt.Println(err)
	}
	meta.SizeJSON = len(bytevals)
	s = time.Now()
	_, err = general.UnmarshalFeature(bytevals)
	meta.TimeReadJSON = int(time.Now().Sub(s).Nanoseconds())
	if err != nil {
		fmt.Println(err)
	}
	s = time.Now()
	bytevals = io.WriteFeature(feature)
	meta.TimeWriteBUF = int(time.Now().Sub(s).Nanoseconds())
	meta.SizeBUF = len(bytevals)
	s = time.Now()
	feature = io.ReadFeature(bytevals)
	meta.TimeReadBUF = int(time.Now().Sub(s).Nanoseconds())
	return meta
}

func CreateMetaCSV(buf *geobuf.Reader, outfilecsv string) {
	outcsv := NewMetaCSV(outfilecsv)
	i := 0
	for buf.Next() {
		outcsv.AddMeta(buf.Feature())
		i++
		if i%1000 == 0 {
			fmt.Printf("\rTotal Number of Meta Data Features Completed %d", i)
		}
	}
}
