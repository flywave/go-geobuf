package io

import (
	"math"
	"reflect"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-pbf"
)

const (
	maxPrecision = 1e6
)

type JSON string

type Encode struct {
	Keys    map[string]int
	KeysNum int
	KeysArr []string
	Dim     int
	Factor  float64
	writer  *pbf.Writer
}

func NewEncode(obj interface{}) *Encode {
	d := &Encode{Keys: make(map[string]int), KeysNum: 0, KeysArr: make([]string, 0), Dim: 2, Factor: math.Pow(10.0, 7.0), writer: pbf.NewWriter()}
	d.writeDataField(obj)
	return d
}

func (e *Encode) Bytes() []byte {
	return e.writer.Finish()
}

func (e *Encode) writeDataField(obj interface{}) {
	e.analyze(obj)

	e.Factor = math.Min(e.Factor, maxPrecision)

	precision := int(math.Ceil(math.Log(e.Factor) / math.Ln10))

	for _, str := range e.KeysArr {
		e.writer.WriteString(DATA_KEYS, str)
	}

	if e.Dim != 2 {
		e.writer.WriteVarint(DIMENSIONS, e.Dim)
	}

	if precision != 6 {
		e.writer.WriteVarint(PRECISION, precision)
	}
	switch g := obj.(type) {
	case *geom.Feature:
		e.writeFeature(g)
	case *geom.FeatureCollection:
		e.writeFeatureCollection(g)
	case geom.Geometry:
		e.writeGeometry(g)
	}
}

func saveKey(key string, e *Encode) {
	if _, ok := e.Keys[key]; !ok {
		e.KeysArr = append(e.KeysArr, key)
		e.Keys[key] = e.KeysNum
		e.KeysNum++
	}
}

func (e *Encode) analyzeMultiLine(coords [][][]float64) {
	for i := 0; i < len(coords); i++ {
		e.analyzePoints(coords[i])
	}
}

func (e *Encode) analyzePoints(coords [][]float64) {
	for i := 0; i < len(coords); i++ {
		e.analyzePoint(coords[i])
	}
}

func (e *Encode) analyzePoint(point []float64) {
	dim := e.Dim
	if len(point) > dim {
		dim = len(point)
	}

	for i := 0; i < len(point); i++ {
		for math.Round(point[i]*e.Factor)/e.Factor != point[i] && e.Factor < maxPrecision {
			e.Factor *= 10
		}
	}
}

func (e *Encode) analyze(obj interface{}) {
	switch o := obj.(type) {
	case *geom.Feature:
		e.analyze(o.Geometry)
		for key := range o.Properties {
			saveKey(key, e)
		}
		for key := range o.CRS {
			saveKey(key, e)
		}
	case *geom.FeatureCollection:
		for i := range o.Features {
			e.analyze(o.Features[i])
		}
		for key := range o.CRS {
			saveKey(key, e)
		}
	case geom.Point:
		e.analyzePoint(o.Data())
	case geom.Point3:
		e.analyzePoint(o.Data())
	case geom.MultiPoint:
		e.analyzePoints(o.Data())
	case geom.MultiPoint3:
		e.analyzePoints(o.Data())
	case geom.Collection:
		for i := 0; i < len(o.Geometries()); i++ {
			e.analyze(o[i])
		}
	case geom.LineString:
		e.analyzePoints(o.Data())
	case geom.LineString3:
		e.analyzePoints(o.Data())
	case geom.Polygon:
		e.analyzeMultiLine(o.Data())
	case geom.Polygon3:
		e.analyzeMultiLine(o.Data())
	case geom.MultiLine:
		e.analyzeMultiLine(o.Data())
	case geom.MultiLine3:
		e.analyzeMultiLine(o.Data())
	case geom.MultiPolygon:
		for _, ls := range o.Data() {
			e.analyzeMultiLine(ls)
		}
	case geom.MultiPolygon3:
		for _, ls := range o.Data() {
			e.analyzeMultiLine(ls)
		}
	}
}

func writeGeometry(geometry geom.Geometry, writer *pbf.Writer, factor float64, dim int) {
	writer.WriteVarint(GEOMETRY_TYPES, GeometryTagTypes[geometry.GetType()])

	switch g := geometry.(type) {
	case geom.Point:
		WritePoint(writer, g.Data(), factor, dim)
	case geom.Point3:
		WritePoint(writer, g.Data(), factor, dim)
	case geom.LineString:
		WriteLine(writer, g.Data(), factor, dim)
	case geom.LineString3:
		WriteLine(writer, g.Data(), factor, dim)
	case geom.Polygon:
		WritePolygon(writer, g.Data(), factor, dim, true)
	case geom.Polygon3:
		WritePolygon(writer, g.Data(), factor, dim, true)
	case geom.MultiPoint:
		WriteLine(writer, g.Data(), factor, dim)
	case geom.MultiPoint3:
		WriteLine(writer, g.Data(), factor, dim)
	case geom.MultiLine:
		WritePolygon(writer, g.Data(), factor, dim, false)
	case geom.MultiLine3:
		WritePolygon(writer, g.Data(), factor, dim, false)
	case geom.MultiPolygon:
		WriteMultiPolygon(writer, g.Data(), factor, dim)
	case geom.MultiPolygon3:
		WriteMultiPolygon(writer, g.Data(), factor, dim)
	case geom.Collection:
		for _, geom := range g {
			writer.WriteMessage(FEATURE_GEOMETRY, func(w *pbf.Writer) {
				writeGeometry(geom, w, factor, dim)
			})
		}
	}
}

func (e *Encode) writeFeature(feature *geom.Feature) {
	if feature.Geometry != nil {
		e.writer.WriteMessage(FEATURE_GEOMETRY, func(w *pbf.Writer) {
			writeGeometry(feature.Geometry, w, e.Factor, e.Dim)
		})
	}

	if feature.ID != nil {
		vv := reflect.ValueOf(feature.ID)
		kd := vv.Kind()
		switch kd {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			e.writer.WriteUInt64(FEATURE_INTID, uint64(vv.Int()))
		case reflect.String:
			e.writer.WriteString(FEATURE_ID, vv.String())
		}
	}

	if feature.Properties != nil {
		e.writeProps(feature.Properties, false)
	}
	if feature.CRS != nil {
		e.writeProps(feature.CRS, true)
	}
}

func (e *Encode) writeGeometry(geometry geom.Geometry) {
	writeGeometry(geometry, e.writer, e.Factor, e.Dim)
}

func (e *Encode) writeProps(props map[string]interface{}, isCustom bool) {
	indexes := make([]int, 0)
	valueIndex := 0

	for key := range props {
		e.writer.WriteMessage(pbf.TagType(13), func(w *pbf.Writer) {
			writeValue(props[key], w)
		})
		indexes = append(indexes, e.Keys[key])
		indexes = append(indexes, valueIndex)
		valueIndex++
	}
	if isCustom {
		e.writer.WritePackedVarint(pbf.TagType(15), indexes)
	} else {
		e.writer.WritePackedVarint(pbf.TagType(14), indexes)
	}
}

func writeValue(value interface{}, writer *pbf.Writer) {
	if value == nil {
		return
	}
	switch v := value.(type) {
	case string:
		writer.WriteValue(VALUES_STRING_VALUE, v)
	case float64:
		writer.WriteValue(VALUES_DOUBLE_VALUE, v)
	case uint64:
		writer.WriteValue(VALUES_POS_INT_VALUE, v)
	case int64:
		if v < 0 {
			writer.WriteValue(VALUES_NEG_INT_VALUE, uint64(-v))
		} else {
			writer.WriteValue(VALUES_POS_INT_VALUE, uint64(v))
		}
	case int:
		if v < 0 {
			writer.WriteValue(VALUES_NEG_INT_VALUE, uint64(-v))
		} else {
			writer.WriteValue(VALUES_POS_INT_VALUE, uint64(v))
		}
	case bool:
		writer.WriteValue(VALUES_BOOL_VALUE, v)
	case JSON:
		writer.WriteValue(VALUES_JSON_VALUE, string(v))
	}
}

func (e *Encode) writeFeatureCollection(obj *geom.FeatureCollection) {
	for _, feat := range obj.Features {
		e.writer.WriteMessage(FEATURE_COLLECTION_FEATURES, func(w *pbf.Writer) {
			e.writeFeature(feat)
		})
	}
	e.writeProps(obj.CRS, true)
}

func WriteFeature(feat *geom.Feature) []byte {
	encode := NewEncode(feat)
	return encode.Bytes()
}

func WriteFeatureCollection(featcol *geom.FeatureCollection) []byte {
	encode := NewEncode(featcol)
	return encode.Bytes()
}

func WriteGeometry(geo geom.Geometry) []byte {
	encode := NewEncode(geo)
	return encode.Bytes()
}
