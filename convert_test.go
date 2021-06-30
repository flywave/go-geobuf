package geobuf

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/flywave/go-geom"
)

func BenchmarkRead(filename_geojson string, filename_geobuf string) {
	if strings.Contains(filename_geojson, ".geobuf") {
		dummy := filename_geojson
		filename_geojson = filename_geobuf
		filename_geobuf = dummy
	}

	s := time.Now()
	bytevals, err := ioutil.ReadFile(filename_geojson)
	if err != nil {
		fmt.Println(err)
	}
	_, err = geom.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}
	end_geojson := time.Now().Sub(s)

	s = time.Now()
	geobuf := ReaderFile(filename_geobuf)
	for geobuf.Next() {
		geobuf.Feature()
	}
	end_geobuf := time.Now().Sub(s)

	fmt.Printf("Time to Read Geojson File: %s\nTime to Read Geobuf File: %s\n", end_geojson, end_geobuf)
}

func BenchmarkWrite(filename_geojson string, filename_geobuf string) {
	if strings.Contains(filename_geojson, ".geobuf") {
		dummy := filename_geojson
		filename_geojson = filename_geobuf
		filename_geobuf = dummy
	}

	bytevals, err := ioutil.ReadFile(filename_geojson)
	if err != nil {
		fmt.Println(err)
	}
	fc, err := geom.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	s := time.Now()
	_, err = fc.MarshalJSON()
	if err != nil {
		fmt.Println(err)
	}
	end_geojson := time.Now().Sub(s)

	s = time.Now()
	geobuf := WriterBufNew()
	for _, feature := range fc.Features {
		geobuf.WriteFeature(feature)
	}
	geobuf.Bytes()
	end_geobuf := time.Now().Sub(s)

	fmt.Printf("Time to Write Geojson File: %s\nTime to Write Geobuf File: %s\n", end_geojson, end_geobuf)
}
