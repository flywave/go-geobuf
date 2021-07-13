package geobuf

import "testing"

func TestJsonConvert(t *testing.T) {
	src := "test_data/county.geojson"
	dst := "test_data/conty.geobuf"
	GeobufFrmCollection(src, dst)

	dst1 := "test_data/conty2.geojson"
	GeobufToCollection(dst, dst1)
}

func TestProper(t *testing.T) {
	js := "test_data/5_23_10.geojson"

	src := "test_data/5_23_10.geobuf"
	dest := "test_data/5_23_102.geojson"
	GeobufFrmCollection(js, src)
	GeobufToCollection(src, dest)
}
