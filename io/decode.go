package io

import (
	"math"

	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
)

type Decode struct {
	Keys   []string
	Values []interface{}
	Dim    int
	Factor float64
	reader *pbf.Reader
}

func NewDecode(reader *pbf.Reader) *Decode {
	d := &Decode{Keys: make([]string, 0), Values: make([]interface{}, 0), Dim: 2, Factor: math.Pow(10.0, 7.0), reader: reader}
	d.readDataField()
	return d
}

func (d *Decode) readDataField() {
	key, _ := d.reader.ReadTag()
	for d.reader.Pos < d.reader.Length {
		if key == DATA_KEYS {
			d.Keys = append(d.Keys, d.reader.ReadString())
			key, _ = d.reader.ReadTag()
		} else if key == DIMENSIONS {
			d.Dim = d.reader.ReadVarint()
			key, _ = d.reader.ReadTag()
		} else if key == PRECISION {
			d.Factor = math.Pow(10, float64(d.reader.ReadVarint()))
			key, _ = d.reader.ReadTag()
		} else if key == DATA_TYPE_FEATURE_COLLECTION {
			d.readFeatureCollection()
			key, _ = d.reader.ReadTag()
		} else if key == DATA_TYPE_FEATURE {
			d.readFeature()
			key, _ = d.reader.ReadTag()
		} else if key == DATA_TYPE_GEOMETRY {
			d.readGeometry()
			key, _ = d.reader.ReadTag()
		}
	}
}

func (d *Decode) readFeatureCollection() *geom.FeatureCollection {
	fc := &geom.FeatureCollection{}
	fc.Type = "FeatureCollection"
	d.readFeatures(fc)
	return fc
}

func (d *Decode) readFeatures(fc *geom.FeatureCollection) {

}

func (d *Decode) readFeature() *geom.Feature {
	feature := &geom.Feature{Properties: map[string]interface{}{}}

	key, val := d.reader.ReadTag()
	if key == FEATURE_GEOMETRY && val == pbf.Bytes {
		feature.Geometry, feature.Properties = d.readGeometry()
		key, val = d.reader.ReadTag()
	}
	if key == FEATURE_ID {
		feature.ID = d.reader.ReadString()
		key, val = d.reader.ReadTag()
	}
	if key == FEATURE_INTID {
		feature.ID = d.reader.ReadVarint()
		key, val = d.reader.ReadTag()
	}
	for key == FEATURE_UNIQUE_VALUES && val == pbf.Bytes {
		size := d.reader.ReadVarint()
		endpos := d.reader.Pos + size

		for d.reader.Pos < endpos {
			newkey, _ := d.reader.ReadTag()
			switch newkey {
			case VALUES_STRING_VALUE:
			case VALUES_JSON_VALUE:
				d.Values = append(d.Values, d.reader.ReadString())
			case VALUES_DOUBLE_VALUE:
				d.Values = append(d.Values, d.reader.ReadDouble())
			case VALUES_POS_INT_VALUE:
			case VALUES_NEG_INT_VALUE:
				d.Values = append(d.Values, d.reader.ReadUInt64())
			case VALUES_BOOL_VALUE:
				d.Values = append(d.Values, d.reader.ReadBool())
			}
			d.reader.Pos = endpos
			key, val = d.reader.ReadTag()
		}
	}
	if key == FEATURE_PROPERTIES {
		if feature.Properties == nil {
			feature.Properties = make(map[string]interface{})
		}
		d.readProps(feature.Properties)
		key, val = d.reader.ReadTag()
	}
	if key == FEATURE_CUSTOM_PROPERTIES {
		d.readProps(feature.Properties)
	}
	return feature
}

func (d *Decode) readGeometry() (geom.Geometry, map[string]interface{}) {
	var geomtype string
	var lengths []uint64
	var geometry geom.Geometry
	var properties map[string]interface{}
	key, val := d.reader.ReadTag()
	if key == GEOMETRY_TYPES && val == pbf.Varint {
		geomtype = GeometryTypes[d.reader.ReadVarint()]
		key, val = d.reader.ReadTag()
	}
	if key == GEOMETRY_LENGTHS && val == pbf.Varint {
		lengths = d.reader.ReadPackedUInt64()
		key, val = d.reader.ReadTag()
	}
	if key == GEOMETRY_COORDS {
		size := d.reader.ReadVarint()
		endpos := d.reader.Pos + size

		switch geomtype {
		case "Point":
			geometry = general.NewPoint(ReadPoint(d.reader, endpos, d.Factor, d.Dim))
		case "LineString":
			geometry = general.NewLineString(ReadLine(d.reader, 0, endpos, d.Factor, d.Dim, false))
		case "Polygon":
			geometry = general.NewPolygon(ReadPolygon(d.reader, endpos, lengths, true, d.Factor, d.Dim))
		case "MultiPoint":
			geometry = general.NewMultiPoint(ReadLine(d.reader, 0, endpos, d.Factor, d.Dim, false))
		case "MultiLineString":
			geometry = general.NewMultiLineString(ReadPolygon(d.reader, endpos, lengths, false, d.Factor, d.Dim))
		case "MultiPolygon":
			geometry = general.NewMultiPolygon(ReadMultiPolygon(d.reader, endpos, lengths, d.Factor, d.Dim))
		}
	}
	if key == GEOMETRY_GEOMETRYS && val == pbf.Bytes {
	}
	if key == GEOMETRY_VALUES && val == pbf.Bytes {
		size := d.reader.ReadVarint()
		endpos := d.reader.Pos + size

		for d.reader.Pos < endpos {
			newkey, _ := d.reader.ReadTag()
			switch newkey {
			case VALUES_STRING_VALUE:
			case VALUES_JSON_VALUE:
				d.Values = append(d.Values, d.reader.ReadString())
			case VALUES_DOUBLE_VALUE:
				d.Values = append(d.Values, d.reader.ReadDouble())
			case VALUES_POS_INT_VALUE:
			case VALUES_NEG_INT_VALUE:
				d.Values = append(d.Values, d.reader.ReadUInt64())
			case VALUES_BOOL_VALUE:
				d.Values = append(d.Values, d.reader.ReadBool())
			}
			d.reader.Pos = endpos
			key, val = d.reader.ReadTag()
		}
	}
	if key == GEOMETRY_CUSTOM_PROPERTIES {
		properties = d.readProps(make(map[string]interface{}))
	}

	return geometry, properties
}

func (d *Decode) readFeatureCollectionField() {
}

func (d *Decode) readProps(props map[string]interface{}) map[string]interface{} {
	size := d.reader.ReadVarint()
	endpos := size + d.reader.Pos

	for d.reader.Pos < endpos {
		props[d.Keys[d.reader.ReadVarint()]] = d.Values[d.reader.ReadVarint()]
	}
	d.Values = make([]interface{}, 0)
	return props
}

func ReadFeature(bytevals []byte) *geom.Feature {
	return nil
}
