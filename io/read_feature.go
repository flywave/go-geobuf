package io

import (
	"github.com/flywave/go-geom"
	"github.com/flywave/go-geom/general"
	"github.com/flywave/go-pbf"
)

func ReadFeature(bytevals []byte) *geom.Feature {
	pbfval := &pbf.Reader{Pbf: bytevals, Length: len(bytevals)}
	var geomtype string
	feature := &geom.Feature{Properties: map[string]interface{}{}}

	key, val := pbfval.ReadTag()
	if key == 1 && val == 0 {
		feature.ID = pbfval.ReadVarint()
		key, val = pbfval.ReadTag()
	}
	for key == 2 && val == 2 {
		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size

		pbfval.Pos += 1
		keyvalue := pbfval.ReadString()

		pbfval.Pos += 1
		pbfval.Varint()
		newkey, _ := pbfval.ReadTag()
		switch newkey {
		case 1:
			feature.Properties[keyvalue] = pbfval.ReadString()
		case 2:
			feature.Properties[keyvalue] = pbfval.ReadFloat()
		case 3:
			feature.Properties[keyvalue] = pbfval.ReadDouble()
		case 4:
			feature.Properties[keyvalue] = pbfval.ReadInt64()
		case 5:
			feature.Properties[keyvalue] = pbfval.ReadUInt64()
		case 6:
			feature.Properties[keyvalue] = pbfval.ReadUInt64()
		case 7:
			feature.Properties[keyvalue] = pbfval.ReadBool()
		}
		pbfval.Pos = endpos
		key, val = pbfval.ReadTag()
	}
	if key == 3 && val == 0 {
		switch int(pbfval.Pbf[pbfval.Pos]) {
		case 1:
			geomtype = "Point"
		case 2:
			geomtype = "LineString"
		case 3:
			geomtype = "Polygon"
		case 4:
			geomtype = "MultiPoint"
		case 5:
			geomtype = "MultiLineString"
		case 6:
			geomtype = "MultiPolygon"
		}
		pbfval.Pos += 1
		key, val = pbfval.ReadTag()
	}
	if key == 4 && val == 2 {
		size := pbfval.ReadVarint()
		endpos := pbfval.Pos + size

		switch geomtype {
		case "Point":
			feature.Geometry = general.NewPoint(ReadPoint(pbfval, endpos))
		case "LineString":
			feature.Geometry = general.NewLineString(ReadLine(pbfval, 0, endpos))
		case "Polygon":
			feature.Geometry = general.NewPolygon(ReadPolygon(pbfval, endpos))
		case "MultiPoint":
			feature.Geometry = general.NewMultiPoint(ReadLine(pbfval, 0, endpos))
		case "MultiLineString":
			feature.Geometry = general.NewMultiLineString(ReadPolygon(pbfval, endpos))
		case "MultiPolygon":
			feature.Geometry = general.NewMultiPolygon(ReadMultiPolygon(pbfval, endpos))
		}
		key, val = pbfval.ReadTag()
	}
	if key == 5 && val == 2 {
		feature.BoundingBox = ReadBoundingBox(pbfval)
	}
	return feature
}
