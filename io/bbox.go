package io

import (
	"github.com/flywave/go-geom"
)

func BoundingBox_Points(pts [][]float64) []float64 {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	for pos, pt := range pts {
		if pos == 0 {
			west, east = pt[0], pt[0]
			south, north = pt[1], pt[1]
		}
		x, y := pt[0], pt[1]
		if x < west {
			west = x
		}
		if x > east {
			east = x
		}

		if y < south {
			south = y
		}
		if y > north {
			north = y
		}
	}
	return []float64{west, south, east, north}
}

func Push_Two_BoundingBoxs(bb1 []float64, bb2 []float64) []float64 {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0

	west1, south1, east1, north1 := bb1[0], bb1[1], bb1[2], bb1[3]
	west2, south2, east2, north2 := bb2[0], bb2[1], bb2[2], bb2[3]

	if west1 < west2 {
		west = west1
	} else {
		west = west2
	}

	if south1 < south2 {
		south = south1
	} else {
		south = south2
	}

	if east1 > east2 {
		east = east1
	} else {
		east = east2
	}

	if north1 > north2 {
		north = north1
	} else {
		north = north2
	}

	return []float64{west, south, east, north}
}

func Expand_BoundingBoxs(bboxs [][]float64) []float64 {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = Push_Two_BoundingBoxs(bbox, temp_bbox)
	}
	return bbox
}

func BoundingBox_PointGeometry(pt []float64) []float64 {
	return []float64{pt[0], pt[1], pt[0], pt[1]}
}

func BoundingBox_MultiPointGeometry(pts [][]float64) []float64 {
	return BoundingBox_Points(pts)
}

func BoundingBox_LineStringGeometry(line [][]float64) []float64 {
	return BoundingBox_Points(line)
}

func BoundingBox_MultiLineStringGeometry(multiline [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, line := range multiline {
		bboxs = append(bboxs, BoundingBox_Points(line))
	}
	return Expand_BoundingBoxs(bboxs)
}

func BoundingBox_PolygonGeometry(polygon [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BoundingBox_Points(cont))
	}
	return Expand_BoundingBoxs(bboxs)
}

func BoundingBox_MultiPolygonGeometry(multipolygon [][][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BoundingBox_Points(cont))
		}
	}
	return Expand_BoundingBoxs(bboxs)
}

func Get_BoundingBox(g geom.Geometry) []float64 {
	switch v := g.(type) {
	case geom.Point:
		return BoundingBox_PointGeometry(v.Data())
	case geom.Point3:
		return BoundingBox_PointGeometry(v.Data())
	case geom.MultiPoint:
		return BoundingBox_MultiPointGeometry(v.Data())
	case geom.MultiPoint3:
		return BoundingBox_MultiPointGeometry(v.Data())
	case geom.LineString:
		return BoundingBox_LineStringGeometry(v.Data())
	case geom.LineString3:
		return BoundingBox_LineStringGeometry(v.Data())
	case geom.MultiLine:
		return BoundingBox_MultiLineStringGeometry(v.Data())
	case geom.MultiLine3:
		return BoundingBox_MultiLineStringGeometry(v.Data())
	case geom.Polygon:
		return BoundingBox_PolygonGeometry(v.Data())
	case geom.Polygon3:
		return BoundingBox_PolygonGeometry(v.Data())
	case geom.MultiPolygon:
		return BoundingBox_MultiPolygonGeometry(v.Data())
	case geom.MultiPolygon3:
		return BoundingBox_MultiPolygonGeometry(v.Data())
	}
	return []float64{}
}

func BoundingBox_GeometryCollection(gs []geom.Geometry) []float64 {
	bboxs := [][]float64{}
	for _, g := range gs {
		bboxs = append(bboxs, Get_BoundingBox(g))
	}
	return Expand_BoundingBoxs(bboxs)
}
