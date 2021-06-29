package io

import (
	"math"

	"github.com/flywave/go-geom/general"
)

var precision = math.Pow(10.0, -7.0)

var feature_s = `{"id":1000001,"type":"Feature","bbox":[-83.647031,33.698307,-83.275933,33.9659119],"geometry":{"type":"MultiPolygon","coordinates":[[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]],[[[-83.537385,33.9659119],[-83.5084519,33.931233],[-83.4155119,33.918541],[-83.275933,33.847977],[-83.306619,33.811444],[-83.28034,33.7617739],[-83.29145,33.7343149],[-83.406189,33.698307],[-83.479523,33.802265],[-83.505928,33.81776],[-83.533165,33.820923],[-83.647031,33.9061979],[-83.537385,33.9659119]]]]},"properties":{"AREA":"13219","COLORKEY":"#03E174","area":"13219","index":1109}}`
var feature, _ = general.UnmarshalFeature([]byte(feature_s))
var bytevals = WriteFeature(feature)

var polygon_s = `{"geometry": {"type": "Polygon", "coordinates": [[[-7.734374999999999, 25.799891182088334], [10.8984375, -34.016241889667015], [45.703125, 17.644022027872726], [-5.9765625, 26.43122806450644], [-7.734374999999999, 25.799891182088334]]]}, "type": "Feature", "properties": {}}`
var multipolygon_s = `{"type":"Feature","properties":{},"geometry":{"type":"MultiPolygon","coordinates":[[[[-71.71875,51.17934297928927],[-36.2109375,-49.15296965617039],[30.585937499999996,0.3515602939922709],[29.179687499999996,59.17592824927136],[-38.3203125,70.72897946208789],[-71.71875,51.17934297928927]]],[[[33.3984375,74.68325030051861],[75.234375,16.29905101458183],[76.2890625,64.77412531292873],[32.6953125,75.23066741281573],[33.3984375,74.68325030051861]]]]}}`
var linestring_s = `{"geometry": {"type": "LineString", "coordinates": [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]}, "type": "Feature", "properties": {}}`
var multilinestring_s = `{"geometry": {"type": "MultiLineString", "coordinates": [[[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651], [15.468749999999998, -9.102096738726443]], [[10.8984375, 56.17002298293205], [16.5234375, -2.108898659243126], [59.4140625, 42.032974332441405], [61.17187499999999, 42.293564192170095]]]}, "type": "Feature", "properties": {}}`
var point_s = `{"geometry": {"type": "Point", "coordinates": [-48.1640625, 47.754097979680026]}, "type": "Feature", "properties": {}}`
var multipoint_s = `{"geometry": {"type": "MultiPoint", "coordinates": [[-48.1640625, 47.754097979680026], [-9.140625, 4.214943141390651]]}, "type": "Feature", "properties": {}}`

var polygon, err = general.UnmarshalFeature([]byte(polygon_s))
var multipolygon, _ = general.UnmarshalFeature([]byte(multipolygon_s))
var linestring, _ = general.UnmarshalFeature([]byte(linestring_s))
var multilinestring, _ = general.UnmarshalFeature([]byte(multilinestring_s))
var point, _ = general.UnmarshalFeature([]byte(point_s))
var multipoint, _ = general.UnmarshalFeature([]byte(multipoint_s))
