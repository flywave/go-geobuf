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
	file, err := os.Open(filename)
	if err != nil {
		return GeojsonFile{}
	}

	bufsize := 4096
	bytevals := make([]byte, bufsize)

	n, err := file.ReadAt(bytevals, 0)
	if err != nil && n == 0 {
		return GeojsonFile{File: file}
	}
	bytevals = bytevals[:n]

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
	if geojsonfile.File == nil {
		return nil
	}

	remaining := size - int(geojsonfile.Pos)
	chunkSize := remaining
	if chunkSize <= 0 {
		return nil
	}
	if chunkSize > 10000000 {
		chunkSize = 10000000
	}

	bytevals := make([]byte, chunkSize)

	n, err := geojsonfile.File.ReadAt(bytevals, geojsonfile.Pos)
	if err != nil && n == 0 {
		return nil
	}
	bytevals = bytevals[:n]

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
	sem := make(chan struct{}, 100)
	for _, i := range feats {
		wg.Add(1)
		sem <- struct{}{}
		go func(i string) {
			defer wg.Done()
			defer func() { <-sem }()
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
		return 0
	}

	size := fi.Size()
	return int(size)
}

func ConvertGeojson(infile string, outfile string) {
	s := time.Now()
	size := GetFilesize(infile)
	if size == 0 {
		return
	}

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
		return
	}
	file.WriteString(`{"type": "FeatureCollection", "features": [`)

	first := true
	for geobuf.Next() {
		feature := geobuf.Feature()
		if feature == nil {
			continue
		}
		s, _ := feature.MarshalJSON()
		if !first {
			file.Write([]byte{','})
		}
		file.Write(s)
		first = false
	}
	file.WriteString("]}")
	file.Close()
}

type MapFunc func(feature *geom.Feature) *geom.Feature

func MapGeobuf(infile string, newfile string, mapfunc MapFunc) {
	geobuf := ReaderFile(infile)
	geobuf2 := WriterFileNew(newfile)
	for geobuf.Next() {
		feature := geobuf.Feature()
		if feature != nil {
			feature = mapfunc(feature)
			geobuf2.WriteFeature(feature)
		}
	}
}

func GeobufFrmCollection(infile string, outfile string) error {
	geobuf, err := os.Create(outfile)
	if err != nil {
		return err
	}
	defer geobuf.Close()

	f, err := os.Open(infile)
	if err != nil {
		return err
	}
	defer f.Close()

	bt, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fc, err := general.UnmarshalFeatureCollection(bt)
	if err != nil {
		return err
	}
	bt = io.WriteFeatureCollection(fc)
	_, err = geobuf.Write(bt)
	return err
}

func GeobufToCollection(infile string, outfile string) {
	f, err := os.Open(infile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	bts, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	fc := io.ReadFeatureCollection(bts)
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	bt, _ := fc.MarshalJSON()
	file.Write(bt)
}
