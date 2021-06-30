package geobuf

import (
	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-geom"
)

type Concurrent struct {
	Reader       *Reader
	C            chan *geom.Feature
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
			con.C <- io.ReadFeature(bytevals)
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
			con.C = make(chan *geom.Feature)

			go con.StartProcesses()
			return true
		}
	} else {
		return true
	}
}

func (con *Concurrent) Feature() *geom.Feature {
	con.Count++
	con.FeatureCount++
	feature := <-con.C
	return feature
}
