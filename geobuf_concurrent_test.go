package geobuf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/flywave/go-geom"
)

func BenchmarkFeatureCollectionRead(b *testing.B) {
	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		bytevals, _ := ioutil.ReadFile("test_data/wv.geojson")
		_, _ = geom.UnmarshalFeatureCollection(bytevals)
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
