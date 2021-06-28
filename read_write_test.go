package geobuf

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"

	geojson "github.com/paulmach/go.geojson"
)

func DeltaPt(pt []float64, testpt []float64) float64 {
	deltax := math.Abs(pt[0] - testpt[0])
	deltay := math.Abs(pt[1] - testpt[1])
	return deltax + deltay
}

var PrecisionError = math.Pow(10.0, -6.0)

func TestReadWriteFile(t *testing.T) {
	bytevals, err := ioutil.ReadFile("test_data/county.geojson")
	if err != nil {
		fmt.Println(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	featuremap := map[int]*geojson.Feature{}
	for _, feature := range fc.Features {
		featuremap[int(feature.ID.(float64))] = feature
	}

	size_fc := len(fc.Features)

	buf := WriterFileNew("test.geobuf")

	for _, feature := range fc.Features {
		buf.WriteFeature(feature)
	}

	readbuf := buf.Reader()

	size_buf := 0
	for readbuf.Next() {
		feature := readbuf.Feature()
		id := int(feature.ID.(int))
		testfeature := featuremap[id]
		if testfeature.Geometry.Type == feature.Geometry.Type {
			for i := range testfeature.Geometry.Polygon {
				testring := testfeature.Geometry.Polygon[i]
				ring := feature.Geometry.Polygon[i]
				if len(ring) != len(testring) {
					t.Errorf("Different ring sizes expected %d got %d", len(testring), len(ring))
				} else {
					for j := range ring {
						pt := ring[j]
						testpt := testring[j]
						deltapt := DeltaPt(pt, testpt)
						if PrecisionError < deltapt {
							t.Errorf("Different Points expected %v %v", testpt, pt)
						}

					}

				}
			}
		} else {
			t.Errorf("Different Types")
		}

		size_buf++
	}

	if size_buf != size_fc {
		t.Errorf("Error ReadWrite File %d fc %d buf", size_fc, size_buf)
	}
	os.Remove("test.geobuf")
}

func TestReadWriteBuf(t *testing.T) {
	bytevals, err := ioutil.ReadFile("test_data/county.geojson")
	if err != nil {
		fmt.Println(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	featuremap := map[int]*geojson.Feature{}
	for _, feature := range fc.Features {
		featuremap[int(feature.ID.(float64))] = feature
	}

	size_fc := len(fc.Features)

	buf := WriterBufNew()

	for _, feature := range fc.Features {
		buf.WriteFeature(feature)
	}

	readbuf := buf.Reader()

	size_buf := 0
	for readbuf.Next() {
		feature := readbuf.Feature()
		id := int(feature.ID.(int))
		testfeature := featuremap[id]
		if testfeature.Geometry.Type == feature.Geometry.Type {
			for i := range testfeature.Geometry.Polygon {
				testring := testfeature.Geometry.Polygon[i]
				ring := feature.Geometry.Polygon[i]
				if len(ring) != len(testring) {
					t.Errorf("Different ring sizes expected %d got %d", len(testring), len(ring))
				} else {
					for j := range ring {
						pt := ring[j]
						testpt := testring[j]
						deltapt := DeltaPt(pt, testpt)
						if PrecisionError < deltapt {
							t.Errorf("Different Points expected %v %v", testpt, pt)
						}

					}

				}
			}
		} else {
			t.Errorf("Different Types")
		}

		size_buf++
	}

	if size_buf != size_fc {
		t.Errorf("Error ReadWrite File %d fc %d buf", size_fc, size_buf)
	}
}

func CreateInds(size int) [][2]int {
	delta := size / 10
	current := 0
	newlist := []int{current}
	for current < size {
		current += delta
		newlist = append(newlist, current)
	}

	if newlist[len(newlist)-1] >= size {
		newlist[len(newlist)-1] = size
	}

	oldi := newlist[0]
	totalinds := [][2]int{}
	for _, i := range newlist[1:] {
		totalinds = append(totalinds, [2]int{oldi, i})
		oldi = i
	}
	return totalinds
}

func TestReadWriteMultiBufFile(t *testing.T) {
	bytevals, err := ioutil.ReadFile("test_data/county.geojson")
	if err != nil {
		fmt.Println(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	featuremap := map[int]*geojson.Feature{}
	for _, feature := range fc.Features {
		featuremap[int(feature.ID.(float64))] = feature
	}

	size_fc := len(fc.Features)

	inds := CreateInds(size_fc)

	buffers := []*Writer{}
	total_split := 0
	for _, ind := range inds {
		i1, i2 := ind[0], ind[1]
		features := fc.Features[i1:i2]
		total_split += len(features)

		buf := WriterBufNew()

		for _, feature := range features {
			buf.WriteFeature(feature)
		}
		buffers = append(buffers, buf)
	}

	if total_split != size_fc {
		t.Errorf("Split function not workign %d %d", size_fc, total_split)
	}

	bigbuffer := WriterFileNew("test.geobuf")

	for _, buf := range buffers {
		bigbuffer.AddGeobuf(buf)
	}

	bigbufreader := bigbuffer.Reader()

	size_buf := 0
	for bigbufreader.Next() {
		feature := bigbufreader.Feature()
		id := int(feature.ID.(int))
		testfeature := featuremap[id]
		if testfeature.Geometry.Type == feature.Geometry.Type {
			for i := range testfeature.Geometry.Polygon {
				testring := testfeature.Geometry.Polygon[i]
				ring := feature.Geometry.Polygon[i]
				if len(ring) != len(testring) {
					t.Errorf("Different ring sizes expected %d got %d", len(testring), len(ring))
				} else {
					for j := range ring {
						pt := ring[j]
						testpt := testring[j]
						deltapt := DeltaPt(pt, testpt)
						if PrecisionError < deltapt {
							t.Errorf("Different Points expected %v %v", testpt, pt)
						}

					}

				}
			}
		} else {
			t.Errorf("Different Types")
		}

		size_buf++
	}

	if size_buf != size_fc {
		t.Errorf("Error ReadWrite File %d fc %d buf", size_fc, size_buf)
	}
	os.Remove("test.geobuf")
}

func TestReadWriteMultiBuf(t *testing.T) {
	bytevals, err := ioutil.ReadFile("test_data/county.geojson")
	if err != nil {
		fmt.Println(err)
	}

	fc, err := geojson.UnmarshalFeatureCollection(bytevals)
	if err != nil {
		fmt.Println(err)
	}

	featuremap := map[int]*geojson.Feature{}
	for _, feature := range fc.Features {
		featuremap[int(feature.ID.(float64))] = feature
	}

	size_fc := len(fc.Features)

	inds := CreateInds(size_fc)

	buffers := []*Writer{}
	total_split := 0
	for _, ind := range inds {
		i1, i2 := ind[0], ind[1]
		features := fc.Features[i1:i2]
		total_split += len(features)

		buf := WriterBufNew()

		for _, feature := range features {
			buf.WriteFeature(feature)
		}
		buffers = append(buffers, buf)
	}

	if total_split != size_fc {
		t.Errorf("Split function not workign %d %d", size_fc, total_split)
	}

	bigbuffer := WriterBufNew()

	for _, buf := range buffers {
		bigbuffer.AddGeobuf(buf)
	}

	bigbufreader := bigbuffer.Reader()

	size_buf := 0
	for bigbufreader.Next() {
		feature := bigbufreader.Feature()
		id := int(feature.ID.(int))
		testfeature := featuremap[id]
		if testfeature.Geometry.Type == feature.Geometry.Type {
			for i := range testfeature.Geometry.Polygon {
				testring := testfeature.Geometry.Polygon[i]
				ring := feature.Geometry.Polygon[i]
				if len(ring) != len(testring) {
					t.Errorf("Different ring sizes expected %d got %d", len(testring), len(ring))
				} else {
					for j := range ring {
						pt := ring[j]
						testpt := testring[j]
						deltapt := DeltaPt(pt, testpt)
						if PrecisionError < deltapt {
							t.Errorf("Different Points expected %v %v", testpt, pt)
						}

					}

				}
			}
		} else {
			t.Errorf("Different Types")
		}

		size_buf++
	}

	if size_buf != size_fc {
		t.Errorf("Error ReadWrite File %d fc %d buf", size_fc, size_buf)
	}
}
