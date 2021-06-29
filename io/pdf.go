package io

import (
	"math"

	"github.com/flywave/go-pbf"
)

var powerfactor = math.Pow(10.0, 7.0)

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func ReadSVarintPower(pbf *pbf.Reader) float64 {
	num := int(pbf.ReadVarint())
	if num%2 == 1 {
		return float64((num+1)/-2) / powerfactor
	} else {
		return float64(num/2) / powerfactor
	}
	return float64(0)
}

func ReadPoint(pbf *pbf.Reader, endpos int) []float64 {
	for pbf.Pos < endpos {
		x := ReadSVarintPower(pbf)
		y := ReadSVarintPower(pbf)
		return []float64{Round(x, .5, 7), Round(y, .5, 7)}
	}
	return []float64{}
}

func ReadLine(pbf *pbf.Reader, num int, endpos int) [][]float64 {
	var x, y float64
	if num == 0 {

		for startpos := pbf.Pos; startpos < endpos; startpos++ {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += ReadSVarintPower(pbf)
			y += ReadSVarintPower(pbf)
			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}
		}

		return newlist
	} else {
		newlist := make([][]float64, num/2)

		for i := 0; i < num/2; i++ {
			x += ReadSVarintPower(pbf)
			y += ReadSVarintPower(pbf)

			newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}

		}
		return newlist
	}
	return [][]float64{}
}

func ReadPolygon(pbf *pbf.Reader, endpos int) [][][]float64 {
	polygon := [][][]float64{}
	for pbf.Pos < endpos {
		num := pbf.ReadVarint()
		polygon = append(polygon, ReadLine(pbf, num, endpos))
	}
	return polygon
}

func ReadMultiPolygon(pbf *pbf.Reader, endpos int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	for pbf.Pos < endpos {
		num_rings := pbf.ReadVarint()
		polygon := make([][][]float64, num_rings)
		for i := 0; i < num_rings; i++ {
			num := pbf.ReadVarint()
			polygon[i] = ReadLine(pbf, num, endpos)
		}
		multipolygon = append(multipolygon, polygon)
	}
	return multipolygon
}

func ReadBoundingBox(pbf *pbf.Reader) []float64 {
	bb := make([]float64, 4)
	pbf.ReadVarint()
	bb[0] = float64(ReadSVarintPower(pbf))
	bb[1] = float64(ReadSVarintPower(pbf))
	bb[2] = float64(ReadSVarintPower(pbf))
	bb[3] = float64(ReadSVarintPower(pbf))
	return bb
}

func ConvertPt(pt []float64) []int64 {
	newpt := make([]int64, 2)
	newpt[0] = int64(pt[0] * math.Pow(10.0, 7.0))
	newpt[1] = int64(pt[1] * math.Pow(10.0, 7.0))
	return newpt
}

func paramEnc(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 31))
}

func WritePoint(pbf *pbf.Writer, tag pbf.TagType, pt []float64) {
	point := ConvertPt(pt)
	pbf.WritePackedUInt64(tag, []uint64{paramEnc(point[0]), paramEnc(point[1])})
}

func WriteLine(pbf *pbf.Writer, tag pbf.TagType, line [][]float64) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	newline := make([]uint64, len(line)*2)
	deltapt := make([]int64, 2)
	pt := make([]int64, 2)
	oldpt := make([]int64, 2)

	for i, point := range line {
		x, y := point[0], point[1]
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}

		pt = ConvertPt(point)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
		} else {
			deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1]}
			newline[i*2] = paramEnc(deltapt[0])
			newline[i*2+1] = paramEnc(deltapt[1])
		}
		oldpt = pt
	}
	pbf.WritePackedUInt64(tag, newline)
}

func MakeLine2(line [][]float64) ([]uint64, []int64) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	newline := make([]uint64, len(line)*2)
	deltapt := make([]int64, 2)
	pt := make([]int64, 2)
	oldpt := make([]int64, 2)

	for i, point := range line {
		x, y := point[0], point[1]
		if x < west {
			west = x
		} else if x > east {
			east = x
		}

		if y < south {
			south = y
		} else if y > north {
			north = y
		}

		pt = ConvertPt(point)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
		} else {
			deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1]}
			newline[i*2] = paramEnc(deltapt[0])
			newline[i*2+1] = paramEnc(deltapt[1])
		}
		oldpt = pt
	}

	return newline, []int64{int64(west * powerfactor),
		int64(south * powerfactor),
		int64(east * powerfactor),
		int64(north * powerfactor)}
}

func MakePolygon2(polygon [][][]float64) ([]uint64, []int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*2))

		tmpgeom, tmpbb := MakeLine2(cont)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	return geometry, bb
}

func WriteMultiPolygon(pbf *pbf.Writer, tag pbf.TagType, multipolygon [][][][]float64) []int64 {
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west, south, east, north = west*powerfactor, south*powerfactor, east*powerfactor, north*powerfactor
	bb := []int64{int64(west), int64(south), int64(east), int64(north)}

	for _, polygon := range multipolygon {
		geometry = append(geometry, uint64(len(polygon)))
		tempgeom, tempbb := MakePolygon2(polygon)
		geometry = append(geometry, tempgeom...)
		if bb[0] > tempbb[0] {
			bb[0] = tempbb[0]
		}
		if bb[1] > tempbb[1] {
			bb[1] = tempbb[1]
		}
		if bb[2] < tempbb[2] {
			bb[2] = tempbb[2]
		}
		if bb[3] < tempbb[3] {
			bb[3] = tempbb[3]
		}
	}
	pbf.WritePackedUInt64(tag, geometry)
	return bb
}

func WritePolygon(pbf *pbf.Writer, tag pbf.TagType, polygon [][][]float64) []int64 {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*2))

		tmpgeom, tmpbb := MakeLine2(cont)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	pbf.WritePackedUInt64(tag, geometry)
	return bb
}

func WriteBoundingBox(pbf *pbf.Writer, tag pbf.TagType, box []float64) {
	boxs := []uint64{
		paramEnc(int64(box[0] * math.Pow(10.0, 7.0))),
		paramEnc(int64(box[1] * math.Pow(10.0, 7.0))),
		paramEnc(int64(box[2] * math.Pow(10.0, 7.0))),
		paramEnc(int64(box[3] * math.Pow(10.0, 7.0))),
	}
	pbf.WritePackedUInt64(tag, boxs)
}

func WriteKeyValue(pbf *pbf.Writer, keytag pbf.TagType, key string, valtag pbf.TagType, value interface{}) {
	pbf.WriteString(keytag, key)

	switch valtag {
	case 1:
		if s, ok := value.(string); ok {
			pbf.WriteString(valtag, s)
		}
	case 2:
		if f, ok := value.(float32); ok {
			pbf.WriteFloat(valtag, f)
		}
	case 3:
		if d, ok := value.(float64); ok {
			pbf.WriteDouble(valtag, d)
		}
	case 4:
		if i, ok := value.(int64); ok {
			pbf.WriteInt64(valtag, i)
		}
	case 5:
		if i, ok := value.(uint64); ok {
			pbf.WriteUInt64(valtag, i)
		}
	case 6:
		if i, ok := value.(uint64); ok {
			pbf.WriteUInt64(valtag, i)
		}
	case 7:
		if b, ok := value.(bool); ok {
			pbf.WriteBool(valtag, b)
		}
	}
}
