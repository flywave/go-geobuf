package io

import (
	"math"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
)

type Decode struct {
	Keys              []string
	Dim               int
	Factor            float64
	reader            *pbf.Reader
	featureCollection *geom.FeatureCollection
	feature           *geom.Feature
	geometry          geom.Geometry
}

func NewDecode(reader *pbf.Reader) *Decode {
	d := &Decode{Keys: make([]string, 0), Dim: 2, Factor: math.Pow(10.0, 7.0), reader: reader}
	d.reader.ReadFields(readDataField, d, -1)
	return d
}

func readDataField(key pbf.TagType, tp pbf.WireType, res interface{}, reader *pbf.Reader) {
	d := res.(*Decode)
	if key == DATA_KEYS {
		d.Keys = append(d.Keys, d.reader.ReadString())
	} else if key == DIMENSIONS {
		d.Dim = d.reader.ReadVarint()
	} else if key == PRECISION {
		d.Factor = math.Pow(10, float64(d.reader.ReadVarint()))
	} else if key == DATA_TYPE_FEATURE_COLLECTION {
		d.featureCollection = d.readFeatureCollection()
		bboxs := make([][]float64, 0)
		for _, feat := range d.featureCollection.Features {
			bboxs = append(bboxs, Get_BoundingBox(feat.Geometry))
		}
		d.featureCollection.BoundingBox = Expand_BoundingBoxs(bboxs)
	} else if key == DATA_TYPE_FEATURE {
		d.feature = d.readFeature()
		d.feature.BoundingBox = Get_BoundingBox(d.feature.Geometry)
	} else if key == DATA_TYPE_GEOMETRY {
		d.geometry, _ = d.readGeometry()
	}
}

func readFeature(reader *pbf.Reader, ctx *readerContext) *geom.Feature {
	ctx.feature = &geom.Feature{Properties: map[string]interface{}{}}
	reader.ReadMessage(readFeatureField, ctx)
	return ctx.feature
}

type readerContext struct {
	Keys              []string
	Values            []interface{}
	Dim               int
	Factor            float64
	featureCollection *geom.FeatureCollection
	feature           *geom.Feature
	geometry          geom.Geometry
	properties        map[string]interface{}
}

func readProps(reader *pbf.Reader, ctx *readerContext, props map[string]interface{}) map[string]interface{} {
	size := reader.ReadVarint()
	endpos := size + reader.Pos

	for reader.Pos < endpos {
		props[ctx.Keys[reader.ReadVarint()]] = ctx.Values[reader.ReadVarint()]
	}
	return props
}

func readValue(reader *pbf.Reader, values []interface{}) {
	size := reader.ReadVarint()
	endpos := reader.Pos + size

	for reader.Pos < endpos {
		newkey, _ := reader.ReadTag()
		switch newkey {
		case VALUES_STRING_VALUE:
		case VALUES_JSON_VALUE:
			values = append(values, reader.ReadString())
		case VALUES_DOUBLE_VALUE:
			values = append(values, reader.ReadDouble())
		case VALUES_POS_INT_VALUE:
			values = append(values, reader.ReadUInt64())
		case VALUES_NEG_INT_VALUE:
			values = append(values, -int64(reader.ReadUInt64()))
		case VALUES_BOOL_VALUE:
			values = append(values, reader.ReadBool())
		}
		reader.Pos = endpos
	}
}

func readFeatureCollectionField(tag pbf.TagType, tp pbf.WireType, result interface{}, reader *pbf.Reader) {
	ctx := result.(*readerContext)
	if tag == FEATURE_COLLECTION_FEATURES {
		fctx := *ctx
		fctx.feature = &geom.Feature{}
		ctx.featureCollection.Features = append(ctx.featureCollection.Features, readFeature(reader, &fctx))
	} else if tag == FEATURE_COLLECTION_VALUES {
		readValue(reader, ctx.Values)
	} else if tag == FEATURE_COLLECTION_CUSTOM_PROPERTIES {
		ctx.properties = readProps(reader, ctx, ctx.properties)
	}
}

func (d *Decode) getReaderContext() *readerContext {
	return &readerContext{Keys: d.Keys, Values: make([]interface{}, 0), Dim: d.Dim, Factor: d.Factor}
}

func (d *Decode) readFeatureCollection() *geom.FeatureCollection {
	ctx := d.getReaderContext()
	ctx.featureCollection = &geom.FeatureCollection{}
	ctx.featureCollection.Type = "FeatureCollection"
	d.reader.ReadMessage(readFeatureCollectionField, ctx.featureCollection)
	return ctx.featureCollection
}

func readGeometry(reader *pbf.Reader, ctx *readerContext) {
	reader.ReadMessage(readGeometryField, ctx)
}

func readFeatureField(key pbf.TagType, val pbf.WireType, result interface{}, reader *pbf.Reader) {
	ctx := result.(*readerContext)
	feature := ctx.feature
	if key == FEATURE_GEOMETRY && val == pbf.Bytes {
		gctx := *ctx
		gctx.feature = nil
		readGeometry(reader, &gctx)
		feature.Geometry, feature.Properties = gctx.geometry, gctx.properties
	}
	if key == FEATURE_ID {
		feature.ID = reader.ReadString()
	}
	if key == FEATURE_INTID {
		feature.ID = reader.ReadVarint()
	}
	for key == FEATURE_UNIQUE_VALUES && val == pbf.Bytes {
		readValue(reader, ctx.Values)
	}
	if key == FEATURE_PROPERTIES {
		if feature.Properties == nil {
			feature.Properties = make(map[string]interface{})
		}
		feature.Properties = readProps(reader, ctx, feature.Properties)
	}
	if key == FEATURE_CUSTOM_PROPERTIES {
		feature.Properties = readProps(reader, ctx, feature.Properties)
	}
}

func (d *Decode) readFeature() *geom.Feature {
	ctx := d.getReaderContext()
	ctx.feature = &geom.Feature{Properties: map[string]interface{}{}}
	d.reader.ReadMessage(readFeatureField, ctx)
	return ctx.feature
}

func readGeometryField(key pbf.TagType, val pbf.WireType, result interface{}, reader *pbf.Reader) {
	ctx := result.(*readerContext)
	var geomtype string
	var lengths []uint64
	var geometry geom.Geometry
	if key == GEOMETRY_TYPES && val == pbf.Varint {
		geomtype = GeometryTypes[reader.ReadVarint()]
	}
	if key == GEOMETRY_LENGTHS && val == pbf.Varint {
		lengths = reader.ReadPackedUInt64()
	}
	if key == GEOMETRY_COORDS {
		size := reader.ReadVarint()
		endpos := reader.Pos + size

		switch geomtype {
		case "Point":
			geometry = general.NewPoint(ReadPoint(reader, endpos, ctx.Factor, ctx.Dim))
		case "LineString":
			geometry = general.NewLineString(ReadLine(reader, 0, endpos, ctx.Factor, ctx.Dim, false))
		case "Polygon":
			geometry = general.NewPolygon(ReadPolygon(reader, endpos, lengths, true, ctx.Factor, ctx.Dim))
		case "MultiPoint":
			geometry = general.NewMultiPoint(ReadLine(reader, 0, endpos, ctx.Factor, ctx.Dim, false))
		case "MultiLineString":
			geometry = general.NewMultiLineString(ReadPolygon(reader, endpos, lengths, false, ctx.Factor, ctx.Dim))
		case "MultiPolygon":
			geometry = general.NewMultiPolygon(ReadMultiPolygon(reader, endpos, lengths, ctx.Factor, ctx.Dim))
		}
	}
	if key == GEOMETRY_GEOMETRYS && val == pbf.Bytes {
		readGeometry(reader, ctx)
	}
	if key == GEOMETRY_VALUES && val == pbf.Bytes {
		readValue(reader, ctx.Values)
	}
	if key == GEOMETRY_CUSTOM_PROPERTIES {
		ctx.properties = readProps(reader, ctx, make(map[string]interface{}))
	}
	if ctx.geometry != nil {
		gc, ok := ctx.geometry.(geom.Collection)
		if ok {
			gc = append(gc, geometry)
		} else {
			temp := ctx.geometry
			gc = geom.Collection{temp, geometry}
			ctx.geometry = gc
		}
	} else {
		ctx.geometry = geometry
	}
}

func (d *Decode) readGeometry() (geom.Geometry, map[string]interface{}) {
	ctx := d.getReaderContext()
	d.reader.ReadMessage(readFeatureCollectionField, ctx.featureCollection)
	return ctx.geometry, ctx.properties
}

func ReadFeature(bytevals []byte) *geom.Feature {
	return nil
}

func ReadFeatureCollection(bytevals []byte) *geom.FeatureCollection {
	return nil
}

func ReadGeometry(bytevals []byte) geom.Geometry {
	return nil
}