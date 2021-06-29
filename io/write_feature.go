package io

import (
	"reflect"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-pbf"
)

func WriteFeature(feat *geom.Feature) []byte {
	newbytes := []byte{8}

	fwriter := pbf.NewWriter()

	if feat.ID != nil {
		vv := reflect.ValueOf(feat.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fwriter.WriteUInt64(layer.Proto.Feature.ID, uint64(vv.Int()))
		}
	}

	for k, v := range feat.Properties {
		WriteKeyValue(fwriter, k, v)
	}
	if feat.Geometry != nil {
		switch feat.Geometry.GetType() {
		case "Point":
			WritePoint(fwriter, feat.Geometry.Point)
		case "LineString":
			WriteLine(fwriter, feat.Geometry.LineString)
		case "Polygon":
			WritePolygon(fwriter, feat.Geometry.Polygon)
		case "MultiPoint":
			WriteLine(fwriter, feat.Geometry.MultiPoint)
		case "MultiLineString":
			WritePolygon(fwriter, feat.Geometry.MultiLineString)
		case "MultiPolygon":
			WriteMultiPolygon(fwriter, feat.Geometry.MultiPolygon)
		}
	}

	if feat.Geometry != nil {
		bb := Get_BoundingBox(feat.Geometry)
		WriteBoundingBox(fwriter, bb)
	}

	return newbytes
}
