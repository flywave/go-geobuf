package geobuf

import "testing"

func TestJsonConvert(t *testing.T) {
	src := "test_data/county.geojson"
	dst := "test_data/conty.geobuf"
	GeobufFrmCollection(src, dst)

	dst1 := "test_data/conty2.geojson"
	GeobufToCollection(dst, dst1)
}
