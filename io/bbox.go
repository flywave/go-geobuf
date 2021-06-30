package io

import (
	"github.com/flywave/go-geom"
)

func BBoxPoints(pts [][]float64) []float64 {
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

func MergeBoundingBoxs(bb1 []float64, bb2 []float64) []float64 {
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

func ExpandBBoxs(bboxs [][]float64) []float64 {
	bbox := bboxs[0]
	for _, temp_bbox := range bboxs[1:] {
		bbox = MergeBoundingBoxs(bbox, temp_bbox)
	}
	return bbox
}

func BBoxPointGeometry(pt []float64) []float64 {
	return []float64{pt[0], pt[1], pt[0], pt[1]}
}

func BBoxMultiPointGeometry(pts [][]float64) []float64 {
	return BBoxPoints(pts)
}

func BBoxLineStringGeometry(line [][]float64) []float64 {
	return BBoxPoints(line)
}

func BBoxMultiLineStringGeometry(multiline [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, line := range multiline {
		bboxs = append(bboxs, BBoxPoints(line))
	}
	return ExpandBBoxs(bboxs)
}

func BBoxPolygonGeometry(polygon [][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, cont := range polygon {
		bboxs = append(bboxs, BBoxPoints(cont))
	}
	return ExpandBBoxs(bboxs)
}

func BBoxMultiPolygonGeometry(multipolygon [][][][]float64) []float64 {
	bboxs := [][]float64{}
	for _, polygon := range multipolygon {
		for _, cont := range polygon {
			bboxs = append(bboxs, BBoxPoints(cont))
		}
	}
	return ExpandBBoxs(bboxs)
}

func CaclBoundingBox(g geom.Geometry) []float64 {
	switch v := g.(type) {
	case geom.Point:
		return BBoxPointGeometry(v.Data())
	case geom.Point3:
		return BBoxPointGeometry(v.Data())
	case geom.MultiPoint:
		return BBoxMultiPointGeometry(v.Data())
	case geom.MultiPoint3:
		return BBoxMultiPointGeometry(v.Data())
	case geom.LineString:
		return BBoxLineStringGeometry(v.Data())
	case geom.LineString3:
		return BBoxLineStringGeometry(v.Data())
	case geom.MultiLine:
		return BBoxMultiLineStringGeometry(v.Data())
	case geom.MultiLine3:
		return BBoxMultiLineStringGeometry(v.Data())
	case geom.Polygon:
		return BBoxPolygonGeometry(v.Data())
	case geom.Polygon3:
		return BBoxPolygonGeometry(v.Data())
	case geom.MultiPolygon:
		return BBoxMultiPolygonGeometry(v.Data())
	case geom.MultiPolygon3:
		return BBoxMultiPolygonGeometry(v.Data())
	}
	return []float64{}
}

func BBoxGeometryCollection(gs []geom.Geometry) []float64 {
	bboxs := [][]float64{}
	for _, g := range gs {
		bboxs = append(bboxs, CaclBoundingBox(g))
	}
	return ExpandBBoxs(bboxs)
}
