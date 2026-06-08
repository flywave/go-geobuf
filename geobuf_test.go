package geobuf

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/flywave/go-geobuf/io"
	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
)

func point(x, y float64) *geom.Feature {
	f, _ := general.UnmarshalFeature([]byte(fmt.Sprintf(
		`{"type":"Feature","geometry":{"type":"Point","coordinates":[%v,%v]},"properties":{"a":1}}`, x, y)))
	return f
}

func TestWriterBufNewClose(t *testing.T) {
	w := WriterBufNew()
	w.Close()
}

func TestWriterBufNewWriteFeature(t *testing.T) {
	w := WriterBufNew()
	w.WriteFeature(point(1, 2))
	data := w.Bytes()
	if len(data) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestWriterBufNewReader(t *testing.T) {
	w := WriterBufNew()
	f := point(1, 2)
	w.WriteFeature(f)
	f2 := point(3, 4)
	w.WriteFeature(f2)

	r := w.Reader()
	if r == nil {
		t.Fatal("Reader returned nil")
	}

	var count int
	for r.Next() {
		rf := r.Feature()
		if rf == nil {
			t.Fatal("Read feature is nil")
		}
		count++
	}
	r.Close()
	if count != 2 {
		t.Fatalf("expected 2 features, got %d", count)
	}
}

func TestWriterFileNewClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.geobuf")
	w := WriterFileNew(path)
	if w == nil {
		t.Fatal("WriterFileNew returned nil")
	}
	w.WriteFeature(point(1, 2))
	w.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty file")
	}
}

func TestWriterFileNewReader(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.geobuf")
	w := WriterFileNew(path)
	w.WriteFeature(point(1, 2))
	w.WriteFeature(point(3, 4))
	w.Close()

	r := ReaderFile(path)
	if r == nil {
		t.Fatal("ReaderFile returned nil")
	}
	defer r.Close()

	var count int
	for r.Next() {
		rf := r.Feature()
		if rf == nil {
			t.Fatal("Read feature is nil")
		}
		count++
	}
	if count != 2 {
		t.Fatalf("expected 2 features, got %d", count)
	}
}

func TestRoundTripPoint(t *testing.T) {
	orig := point(10.5, 20.7)
	data := io.WriteFeature(orig)
	if len(data) == 0 {
		t.Fatal("WriteFeature returned empty data")
	}
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	g, ok := read.Geometry.(geom.Point)
	if !ok {
		t.Fatal("expected Point geometry")
	}
	pt := g.Data()
	if math.Abs(pt[0]-10.5) > 1e-6 || math.Abs(pt[1]-20.7) > 1e-6 {
		t.Fatalf("point mismatch: got %v", pt)
	}
	if read.Properties["a"] != float64(1) {
		t.Fatalf("property mismatch: got %v", read.Properties["a"])
	}
}

func TestRoundTripLineString(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"LineString","coordinates":[[0,0],[1,1],[2,0]]},"properties":{}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripPolygon(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[0,0],[1,0],[1,1],[0,1],[0,0]]]},"properties":{"c":3.14}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripMultiPoint(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"MultiPoint","coordinates":[[0,0],[1,1],[2,2]]},"properties":{}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripMultiLineString(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"MultiLineString","coordinates":[[[0,0],[1,1]],[[2,2],[3,3]]]},"properties":{}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripMultiPolygon(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"MultiPolygon","coordinates":[[[[0,0],[1,0],[1,1],[0,1],[0,0]]],[[[2,2],[3,2],[3,3],[2,3],[2,2]]]]},"properties":{}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripPoint3D(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"Point","coordinates":[1,2,3]},"properties":{}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Geometry == nil {
		t.Fatal("expected non-nil geometry")
	}
}

func TestRoundTripFeatureCollection(t *testing.T) {
	fc := &geom.FeatureCollection{
		Features: []*geom.Feature{point(1, 2), point(3, 4)},
	}
	data := io.WriteFeatureCollection(fc)
	read := io.ReadFeatureCollection(data)
	if read == nil {
		t.Fatal("ReadFeatureCollection returned nil")
	}
	if len(read.Features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(read.Features))
	}
}

func TestRoundTripProperties(t *testing.T) {
	raw := `{"type":"Feature","geometry":{"type":"Point","coordinates":[0,0]},"properties":{"str":"hello","num":42,"float":3.14,"bool":true,"neg":-10}}`
	orig, _ := general.UnmarshalFeature([]byte(raw))
	data := io.WriteFeature(orig)
	read := ReadFeature(data)
	if read == nil {
		t.Fatal("ReadFeature returned nil")
	}
	if read.Properties["str"] != "hello" {
		t.Fatalf("string property: got %v", read.Properties["str"])
	}
	if read.Properties["num"] != float64(42) {
		t.Fatalf("int property: got %v (type %T)", read.Properties["num"], read.Properties["num"])
	}
	if read.Properties["bool"] != true {
		t.Fatalf("bool property: got %v", read.Properties["bool"])
	}
}

func TestReaderBufNil(t *testing.T) {
	r := ReaderBuf(nil)
	if r == nil {
		t.Fatal("ReaderBuf(nil) returned nil")
	}
	if r.Next() {
		t.Fatal("Next() on nil buf should return false")
	}
	r.Close()
}

func TestReaderBufEmpty(t *testing.T) {
	r := ReaderBuf([]byte{})
	if r == nil {
		t.Fatal("ReaderBuf returned nil")
	}
	if r.Next() {
		t.Fatal("Next() on empty buf should return false")
	}
	r.Close()
}

func TestReadFeatureNil(t *testing.T) {
	f := ReadFeature(nil)
	if f != nil {
		t.Fatal("ReadFeature(nil) should return nil")
	}
}

func TestWriterFileNewNonexistentDir(t *testing.T) {
	w := WriterFileNew("/nonexistent/path/test.geobuf")
	if w.Err() == nil {
		t.Fatal("expected error for nonexistent directory")
	}
	w.Close()
}

func TestReaderFileNonexistent(t *testing.T) {
	r := ReaderFile("/nonexistent/file.geobuf")
	if r.Err() == nil {
		t.Fatal("expected error for nonexistent file")
	}
	r.Close()
	if r.Next() {
		t.Fatal("Next() after error should return false")
	}
}

func TestWriterFileNewReadAfterClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.geobuf")
	w := WriterFileNew(path)
	w.WriteFeature(point(1, 2))
	w.Close()

	r := ReaderFile(path)
	if r == nil || r.Err() != nil {
		t.Fatal("ReaderFile failed")
	}
	defer r.Close()
	if !r.Next() {
		t.Fatal("expected at least one feature")
	}
}

func TestWriterBytes(t *testing.T) {
	w := WriterBufNew()
	w.WriteFeature(point(1, 2))
	data := w.Bytes()
	if len(data) == 0 {
		t.Fatal("Bytes() returned empty")
	}
}

func TestReaderReset(t *testing.T) {
	w := WriterBufNew()
	w.WriteFeature(point(1, 2))
	w.WriteFeature(point(3, 4))

	r := w.Reader()
	if r == nil {
		t.Fatal("Reader returned nil")
	}

	r.Next()
	r.Feature()

	r.Reset()

	var count int
	for r.Next() {
		r.Feature()
		count++
	}
	if count != 2 {
		t.Fatalf("after Reset, expected 2 features, got %d", count)
	}
}

func TestReadAll(t *testing.T) {
	w := WriterBufNew()
	w.WriteFeature(point(1, 2))
	w.WriteFeature(point(3, 4))

	r := w.Reader()
	feats := r.ReadAll()
	if len(feats) != 2 {
		t.Fatalf("ReadAll: expected 2 features, got %d", len(feats))
	}
}

func TestReadFeatureFromBuf(t *testing.T) {
	w := WriterBufNew()
	w.WriteFeature(point(1, 2))

	data := w.Bytes()
	r := ReaderBuf(data)
	if r == nil {
		t.Fatal("ReaderBuf returned nil")
	}
	defer r.Close()

	if !r.Next() {
		t.Fatal("expected at least one feature")
	}
	f := r.Feature()
	if f == nil {
		t.Fatal("Feature() returned nil")
	}
}
