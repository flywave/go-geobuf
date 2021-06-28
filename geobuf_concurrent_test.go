package geobuf

import (
	"io/ioutil"
	"os"
	"testing"

	ld "github.com/murphy214/ld-geojson"
	geojson "github.com/paulmach/go.geojson"
)

func I() int {
	ld.Convert_FeatureCollection("test_data/wv.geojson", "test_data/wv_ld.geojson")
	return 0
}

var _ = I()

func BenchmarkFeatureCollectionRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals, _ := ioutil.ReadFile("test_data/wv.geojson")
		_, _ = geojson.UnmarshalFeatureCollection(bytevals)
	}
}

func BenchmarkLineDelimitedGeojsonRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		ldjson := ld.Read_LD_Geojson("test_data/wv_ld.geojson")
		for ldjson.Next() {
			ldjson.Feature()
		}
	}
	os.Remove("test_data/wv_ld.geojson")
}

func BenchmarkGeobufRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		buf := ReaderFile("test_data/wv.geobuf")
		for buf.Next() {
			buf.Feature()
		}
	}
}

func BenchmarkGeobufConcurrentRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		buf := ReaderFile("test_data/wv.geobuf")
		buff := NewGeobufReaderConcurrent(buf)
		for buff.Next() {
			buff.Feature()
		}
	}
}
