package geobuf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/flywave/go-geom"
)

type Geojson_File struct {
	Features []*geom.Feature
	Count    int
	File     *os.File
	Pos      int64
	Feat_Pos int
}

func NewGeojson(filename string) Geojson_File {
	file, _ := os.Open(filename)

	bytevals := make([]byte, 100)

	file.ReadAt(bytevals, int64(0))
	boolval := false
	var startpos int
	for ii, i := range string(bytevals) {
		if string(i) == "[" && boolval == false {
			startpos = ii
			boolval = true
		}
	}

	return Geojson_File{File: file, Pos: int64(startpos + 1)}
}

func (geojsonfile *Geojson_File) ReadChunk(size int) []string {
	var bytevals []byte
	if size > int(geojsonfile.Pos)+10000000 {
		bytevals = make([]byte, 10000000)
	} else {
		bytevals = make([]byte, size-int(geojsonfile.Pos))
	}

	geojsonfile.File.ReadAt(bytevals, geojsonfile.Pos)
	debt := 0
	newlist := []int{}
	boolval := false
	for i, run := range string(bytevals) {
		if "{" == string(run) {
			boolval = true
			if debt == 0 {
				newlist = append(newlist, i)
			}
			debt += 1
		} else if "}" == string(run) && boolval == true {
			debt -= 1
			if debt == 0 {
				newlist = append(newlist, i)
			}
		}
	}
	boolval = false
	row := []int{}
	geojsons := []string{}
	for _, i := range newlist {
		row = append(row, i)
		if boolval == false {
			boolval = true
		} else if boolval == true {
			vals := string(bytevals[row[0]:row[1]])
			geojsons = append(geojsons, vals+"}")

			row = []int{}
			boolval = false
		}

	}
	var newpos int64
	if len(newlist) > 0 {
		newpos = geojsonfile.Pos + int64(newlist[len(newlist)-1])
	} else {
		newpos = int64(size)
	}
	geojsonfile.Pos = newpos
	return geojsons
}

func AddFeatures(geobuf *Writer, feats []string, count int, s time.Time) int {
	var wg sync.WaitGroup
	for _, i := range feats {
		wg.Add(1)
		go func(i string) {
			if len(i) > 0 {
				feat, err := geom.UnmarshalFeature([]byte(i))
				if err != nil {
					fmt.Println(err, feat)
				} else {
					if feat.Geometry != nil {
						geobuf.WriteFeature(feat)
					} else {
						fmt.Println(feat)
					}
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	count += len(feats)
	fmt.Printf("\r%d features created from raw geojson string in %s", count, time.Now().Sub(s))

	return count
}

func GetFilesize(filename string) int {
	fi, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err)
	}

	size := fi.Size()
	return int(size)
}

func ConvertGeojson(infile string, outfile string) {
	s := time.Now()
	size := GetFilesize(infile)

	geobuf := WriterFileNew(outfile)
	geojsonfile := NewGeojson(infile)
	count := 0
	feats := []string{"d"}

	for len(feats) > 0 {
		feats = geojsonfile.ReadChunk(size)
		count = AddFeatures(geobuf, feats, count, s)
	}
}

func ConvertGeobuf(infile string, outfile string) {
	geobuf := ReaderFile(infile)

	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
	}
	file.WriteString(`{"type": "FeatureCollection", "features": [`)

	for geobuf.Next() {
		feature := geobuf.Feature()
		s, _ := feature.MarshalJSON()
		if geobuf.Next() {
			file.Write(append(s, 44))

		} else {
			file.Write(s)

		}
	}
	file.WriteString("]}")
}

type MapFunc func(feature *geom.Feature) *geom.Feature

func MapGeobuf(infile string, newfile string, mapfunc MapFunc) {
	geobuf := ReaderFile(infile)
	geobuf2 := WriterFileNew(newfile)
	for geobuf.Next() {
		feature := geobuf.Feature()
		feature = mapfunc(feature)
		geobuf2.WriteFeature(feature)
	}
}

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
