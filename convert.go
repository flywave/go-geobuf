package geobuf

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

type GeojsonFile struct {
	Features []*geom.Feature
	Count    int
	File     *os.File
	Pos      int64
	Feat_Pos int
}

func NewGeojson(filename string) GeojsonFile {
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

	return GeojsonFile{File: file, Pos: int64(startpos + 1)}
}

func (geojsonfile *GeojsonFile) ReadChunk(size int) []string {
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
				feat, err := general.UnmarshalFeature([]byte(i))
				if err != nil {
					fmt.Println(err, feat)
				} else {
					if feat.Geometry != nil {
						geobuf.WriteFeature(feat)
					}
				}
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	count += len(feats)
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

func GeobufFrmCollection(infile string, outfile string) {
	geobuf, _ := os.Create(outfile)
	f, _ := os.Open(infile)
	bt, _ := ioutil.ReadAll(f)
	fc, _ := general.UnmarshalFeatureCollection(bt)
	bt = io.WriteFeatureCollection(fc)
	geobuf.Write(bt)
	geobuf.Close()
}

func GeobufToCollection(infile string, outfile string) {
	f, _ := os.Open(infile)
	bts, _ := ioutil.ReadAll(f)
	fc := io.ReadFeatureCollection(bts)
	file, _ := os.Create(outfile)
	bt, _ := fc.MarshalJSON()
	file.Write(bt)
	file.Close()
}
