package geobuf

import (
	"bufio"
	"os"

	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-geom"
	"github.com/flywave/go-pbf"

	"bytes"
	"fmt"
	"io/ioutil"
)

type Writer struct {
	Filename  string
	Writer    *bufio.Writer
	FileBool  bool
	Buffer    *bytes.Buffer
	File      *os.File
	Bytesvals []byte
}

var writersize = 64 * 4096

func WriterFileNew(filename string) *Writer {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	return &Writer{Filename: filename, FileBool: true, File: file, Bytesvals: []byte{}}
}

func WriterFile(filename string) *Writer {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	return &Writer{Filename: filename, FileBool: true, File: file}
}

func WriterBufNew() *Writer {
	b := bytes.NewBuffer([]byte{})
	return &Writer{Writer: bufio.NewWriterSize(b, writersize), Buffer: b, FileBool: false}
}

func WriterBuf(bytevals []byte) *Writer {
	buffer := bytes.NewBuffer(bytevals)
	return &Writer{Writer: bufio.NewWriterSize(buffer, writersize), Buffer: buffer, FileBool: false}
}

func (writer *Writer) WriteFeature(feature *geom.Feature) {
	bytevals := io.WriteFeature(feature)

	bytevals = append(
		append(
			[]byte{10}, pbf.EncodeVarint(uint64(len(bytevals)))...,
		),
		bytevals...)
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
	}
}

func (writer *Writer) WriteRaw(bytevals []byte) {
	bytevals = append(
		append(
			[]byte{10}, pbf.EncodeVarint(uint64(len(bytevals)))...,
		),
		bytevals...)
	if writer.FileBool {
		writer.File.Write(bytevals)
	} else {
		writer.Writer.Write(bytevals)
	}
}

func (writer *Writer) AddGeobuf(buf *Writer) {
	if !writer.FileBool {
		writer.Writer.Flush()
	}
	if !buf.FileBool {
		buf.Writer.Flush()
		if writer.FileBool {
			writer.File.Write(buf.Buffer.Bytes())
		} else {
			writer.Writer.Write(buf.Buffer.Bytes())
		}
	}
}

func (writer *Writer) Bytes() []byte {
	writer.Writer.Flush()

	if !writer.FileBool {
		return writer.Buffer.Bytes()
	} else {
		writer.File.Close()
		bytevals, _ := ioutil.ReadFile(writer.Filename)
		return bytevals
	}
}

func (writer *Writer) Reader() *Reader {
	if !writer.FileBool {
		newreader := ReaderBuf(writer.Bytes())
		return newreader
	} else {
		writer.File.Close()
		newreader := ReaderFile(writer.Filename)
		return newreader
	}
}

func (writer *Writer) Close() {
	writer.File.Close()
}
