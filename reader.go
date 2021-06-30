package geobuf

import (
	"bufio"
	"os"

	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-mapbox/tileid"

	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-pbf"
)

type Reader struct {
	FileBool     bool
	Reader       *pbf.ProtobufScanner
	Filename     string
	File         *os.File
	Buf          []byte
	MetaData     MetaData
	MetaDataBool bool
	SubFileEnd   int
	FeatureCount int
}

type SubFile struct {
	Positions      [2]int
	NumberFeatures int
	Size           int
}

type MetaData struct {
	FileSize       int
	NumberFeatures int
	Files          map[string]*SubFile
	Bounds         tileid.Extrema
}

func (metadata *MetaData) LintMetaData(pos int) {
	for _, v := range metadata.Files {
		v.Positions = [2]int{v.Positions[0] + pos, v.Positions[1] + pos}
		v.Size = v.Positions[1] - v.Positions[0]
	}
}

func ReaderBuf(bytevals []byte) *Reader {
	buffer := bytes.NewReader(bytevals)
	buf := &Reader{Reader: pbf.NewProtobufScanner(buffer), Buf: bytevals, FileBool: false}
	buf.CheckMetaData()
	buf.FeatureCount = 0

	return buf
}

func ReaderFile(filename string) *Reader {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	reader := bufio.NewReader(file)

	buf := &Reader{
		Reader:   pbf.NewProtobufScanner(reader),
		Filename: filename,
		FileBool: true,
		File:     file,
	}
	buf.CheckMetaData()
	buf.FeatureCount = 0
	return buf
}

func (reader *Reader) Next() bool {
	reader.FeatureCount++
	return reader.Reader.Scan()
}

func (reader *Reader) Bytes() []byte {
	return reader.Reader.Protobuf()
}

func (reader *Reader) BytesIndicies() ([]byte, [2]int) {
	return reader.Reader.ProtobufIndicies()
}

func (reader *Reader) Feature() *geom.Feature {
	return io.ReadFeature(reader.Bytes())
}

func (reader *Reader) FeatureIndicies() (*geom.Feature, [2]int) {
	bytevals, indicies := reader.BytesIndicies()
	return io.ReadFeature(bytevals), indicies
}

func ReadFeature(bytevals []byte) *geom.Feature {
	return io.ReadFeature(bytevals)
}

func ReadKeys(bytevals []byte) []string {
	pbfval := pbf.Reader{Pbf: bytevals, Length: len(bytevals)}
	keys := []string{}
	key, val := pbfval.ReadTag()
	if key == 1 && val == 0 {
		pbfval.ReadVarint()
		key, val = pbfval.ReadTag()
	}
	for key == 2 && val == 2 {
		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size
		pbfval.Pos += 1
		keys = append(keys, pbfval.ReadString())

		pbfval.Pos = endpos
		key, val = pbfval.ReadTag()
	}

	return keys
}

func ReadBoundingBox(bytevals []byte, factor float64) []float64 {
	pos := len(bytevals) - 1
	alloc := make([]byte, 32)
	allocpos := 31
	boolval := true
	for boolval {
		alloc[allocpos] = bytevals[pos]
		if bytevals[pos] == 42 {
			boolval = false
		}
		pos--
		allocpos--
	}

	bb := make([]float64, 4)
	pbfval := pbf.NewReader(alloc[allocpos+3:])
	bb[0] = float64(io.ReadSVarintPower(pbfval, factor))
	bb[1] = float64(io.ReadSVarintPower(pbfval, factor))
	bb[2] = float64(io.ReadSVarintPower(pbfval, factor))
	bb[3] = float64(io.ReadSVarintPower(pbfval, factor))
	return bb
}

func (reader *Reader) ReadAll() []*geom.Feature {
	feats := []*geom.Feature{}
	for reader.Next() {
		feats = append(feats, reader.Feature())
	}
	return feats
}

func (reader *Reader) Reset() {
	if reader.FileBool {
		file, err := os.Open(reader.Filename)
		if err != nil {
			fmt.Println(err)
		}
		read := bufio.NewReader(file)
		reader.Reader = pbf.NewProtobufScanner(read)
	} else {
		buffer := bytes.NewReader(reader.Buf)
		reader.Reader = pbf.NewProtobufScanner(buffer)
	}
	if reader.MetaDataBool {
		reader.Next()
		reader.Bytes()
	}
	reader.FeatureCount = 0
}

func (reader *Reader) ReadIndAppend(inds [2]int) []byte {
	inds[0] = inds[0] - len(pbf.EncodeVarint(uint64(inds[1]-inds[0]))) - 1
	bytevals := make([]byte, inds[1]-inds[0])
	reader.File.ReadAt(bytevals, int64(inds[0]))
	return bytevals
}

func (reader *Reader) ReadIndFeature(inds [2]int) *geom.Feature {
	bytevals := make([]byte, inds[1]-inds[0])
	reader.File.ReadAt(bytevals, int64(inds[0]))
	return ReadFeature(bytevals)
}

func (reader *Reader) ReadIndicies(inds [2]int) []byte {
	bytevals := make([]byte, inds[1]-inds[0])
	reader.File.ReadAt(bytevals, int64(inds[0]))
	return bytevals
}

func (reader *Reader) Seek(pos int) {
	if reader.FileBool {
		reader.File.Seek(int64(pos), 0)
		myreader := bufio.NewReader(reader.File)
		reader.Reader = pbf.NewProtobufScanner(myreader)
		reader.Reader.TotalPosition = pos
	} else {
		buffer := bytes.NewReader(reader.Buf)
		buffer.Seek(int64(pos), 0)
		reader.Reader = pbf.NewProtobufScanner(buffer)
		reader.Reader.TotalPosition = pos
	}
}

func WriteMetaData(meta MetaData) interface{} {
	bb := bytes.NewBuffer([]byte{})
	dec := gob.NewEncoder(bb)
	err := dec.Encode(meta)
	if err != nil {
		fmt.Println(err)
	}
	return string(bb.Bytes())
}

func ReadMetaData(bytevals []byte) MetaData {
	dec := gob.NewDecoder(bytes.NewBuffer(bytevals))
	var q MetaData
	err := dec.Decode(&q)
	if err != nil {
		fmt.Println(err)
	}
	return q
}

func (reader *Reader) CheckMetaData() {
	reader.Next()
	feature := reader.Feature()
	_, boolval := feature.Properties["metadata"]
	if len(feature.Properties) == 1 && boolval {
		bytevals := []byte(feature.Properties["metadata"].(string))
		reader.MetaData = ReadMetaData(bytevals)
		reader.MetaData.LintMetaData(reader.Reader.TotalPosition)
		reader.MetaDataBool = true
		reader.FeatureCount = 0

	} else {
		reader.Reset()
	}
}

func (reader *Reader) SubFileSeek(key string) {
	subfile := reader.MetaData.Files[key]
	reader.Seek(subfile.Positions[0])
	reader.SubFileEnd = subfile.Positions[1]
}

func (reader *Reader) SubFileBytes(key string) *Reader {
	subfile, boolval := reader.MetaData.Files[key]
	if boolval {
		return ReaderBuf(reader.ReadIndicies(subfile.Positions))
	}
	return ReaderBuf([]byte{})
}

func (reader *Reader) SubFileNext() bool {
	return reader.Reader.Scan() && reader.Reader.TotalPosition < reader.SubFileEnd
}

func (reader *Reader) Close() {
	reader.File.Close()
}
