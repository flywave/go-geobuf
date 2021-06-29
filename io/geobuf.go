package io

import "github.com/flywave/go-pbf"

var geometryTypes = map[string]pbf.TagType{
	"Point":              0,
	"MultiPoint":         1,
	"LineString":         2,
	"MultiLineString":    3,
	"Polygon":            4,
	"MultiPolygon":       5,
	"GeometryCollection": 6,
}

const (
	DATA_KEYS                            pbf.TagType = 1
	DIMENSIONS                           pbf.TagType = 2
	PRECISION                            pbf.TagType = 3
	DATA_TYPE_FEATURE_COLLECTION         pbf.TagType = 4
	DATA_TYPE_FEATURE                    pbf.TagType = 5
	DATA_TYPE_GEOMETRY                   pbf.TagType = 6
	FEATURE_GEOMETRY                     pbf.TagType = 1
	FEATURE_ID                           pbf.TagType = 11
	FEATURE_INTID                        pbf.TagType = 12
	FEATURE_UNIQUE_VALUES                pbf.TagType = 13
	FEATURE_PROPERTIES                   pbf.TagType = 14
	FEATURE_CUSTOM_PROPERTIES            pbf.TagType = 15
	GEOMETRY_TYPES                       pbf.TagType = 1
	GEOMETRY_LENGTHS                     pbf.TagType = 2
	GEOMETRY_COORDS                      pbf.TagType = 3
	GEOMETRY_GEOMETRYS                   pbf.TagType = 4
	GEOMETRY_VALUES                      pbf.TagType = 5
	GEOMETRY_CUSTOM_PROPERTIES           pbf.TagType = 6
	FEATURE_COLLECTION_FEATURES          pbf.TagType = 1
	FEATURE_COLLECTION_VALUES            pbf.TagType = 13
	FEATURE_COLLECTION_CUSTOM_PROPERTIES pbf.TagType = 15
	VALUES_STRING_VALUE                  pbf.TagType = 1
	VALUES_DOUBLE_VALUE                  pbf.TagType = 2
	VALUES_POS_INT_VALUE                 pbf.TagType = 3
	VALUES_NEG_INT_VALUE                 pbf.TagType = 4
	VALUES_BOOL_VALUE                    pbf.TagType = 5
	VALUES_JSON_VALUE                    pbf.TagType = 6
)
