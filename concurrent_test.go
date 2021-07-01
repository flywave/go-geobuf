package geobuf

import (
	"io/ioutil"
	"testing"

	"github.com/flywave/go-geom/general"
)

func BenchmarkFeatureCollectionRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals, _ := ioutil.ReadFile("test_data/wv.geojson")
		_, _ = general.UnmarshalFeatureCollection(bytevals)
	}
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
