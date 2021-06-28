package geobuf

import (
	raw "github.com/flywave/go-geobuf/geobuf_raw"
	geojson "github.com/paulmach/go.geojson"
)

type Concurrent struct {
	Reader       *Reader
	C            chan *geojson.Feature
	Count        int
	Limit        int
	FeatureCount int
}

func NewConcurrent(buf *Reader, limit int) *Concurrent {
	return &Concurrent{Reader: buf, Limit: limit, Count: limit}
}

func (con *Concurrent) StartProcesses() {
	i := 0
	for con.Reader.Next() && i < con.Limit {
		bytevals := con.Reader.Bytes()
		go func(bytevals []byte) {
			con.C <- raw.ReadFeature(bytevals)
		}(bytevals)
		i++
	}
	con.Reader.FeatureCount--
}

func (con *Concurrent) Next() bool {
	if con.Count == con.Limit || con.Reader.Reader.EndBool {
		if con.Reader.Reader.EndBool && con.Reader.FeatureCount > con.FeatureCount {
			return true
		} else if con.Reader.Reader.EndBool {
			return false
		} else {
			con.Count = 0
			con.C = make(chan *geojson.Feature)

			go con.StartProcesses()
			return true
		}

	} else {
		return true
	}
	return false
}

func (con *Concurrent) Feature() *geojson.Feature {
	con.Count++
	con.FeatureCount++
	feature := <-con.C
	return feature
}
