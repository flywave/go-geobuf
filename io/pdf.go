package io

import (
	"math"

	"github.com/flywave/go-pbf"
)

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

func ReadSVarintPower(pbf *pbf.Reader, factor float64) float64 {
	num := int(pbf.ReadVarint())
	if num%2 == 1 {
		return float64((num+1)/-2) / factor
	} else {
		return float64(num/2) / factor
	}
}

func ReadPoint(pbf *pbf.Reader, endpos int, factor float64, dim int) []float64 {
	for pbf.Pos < endpos {
		if dim == 2 {
			x := ReadSVarintPower(pbf, factor)
			y := ReadSVarintPower(pbf, factor)
			return []float64{Round(x, .5, 7), Round(y, .5, 7)}
		} else if dim == 3 {
			x := ReadSVarintPower(pbf, factor)
			y := ReadSVarintPower(pbf, factor)
			z := ReadSVarintPower(pbf, factor)
			return []float64{Round(x, .5, 7), Round(y, .5, 7), Round(z, .5, 7)}
		}
	}
	return []float64{}
}

func ReadLine(pbf *pbf.Reader, num int, endpos int, factor float64, dim int, closed bool) [][]float64 {
	var x, y, z float64
	var newlist [][]float64
	if num == 0 {
		for startpos := pbf.Pos; startpos < endpos; startpos++ {
			if pbf.Pbf[startpos] <= 127 {
				num += 1
			}
		}
		newlist = make([][]float64, num/dim)

		for i := 0; i < num/dim; i++ {
			if dim == 2 {
				x += ReadSVarintPower(pbf, factor)
				y += ReadSVarintPower(pbf, factor)
				newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}
			} else if dim == 3 {
				x += ReadSVarintPower(pbf, factor)
				y += ReadSVarintPower(pbf, factor)
				z += ReadSVarintPower(pbf, factor)
				newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7), Round(z, .5, 7)}
			}
		}
	} else {
		newlist = make([][]float64, num/dim)

		for i := 0; i < num/dim; i++ {
			if dim == 2 {
				x += ReadSVarintPower(pbf, factor)
				y += ReadSVarintPower(pbf, factor)
				newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7)}
			} else if dim == 3 {
				x += ReadSVarintPower(pbf, factor)
				y += ReadSVarintPower(pbf, factor)
				z += ReadSVarintPower(pbf, factor)
				newlist[i] = []float64{Round(x, .5, 7), Round(y, .5, 7), Round(z, .5, 7)}
			}
		}
	}

	if closed {
		newlist = append(newlist, newlist[0])
	}
	return newlist
}

func ReadPolygon(pbf *pbf.Reader, endpos int, lengths []uint64, closed bool, factor float64, dim int) [][][]float64 {
	polygon := [][][]float64{}
	if lengths == nil {
		for pbf.Pos < endpos {
			num := pbf.ReadVarint()
			polygon = append(polygon, ReadLine(pbf, num, endpos, factor, dim, closed))
		}
	} else {
		for i := 0; i < len(lengths); i++ {
			num := int(lengths[i])
			polygon = append(polygon, ReadLine(pbf, num, endpos, factor, dim, closed))
		}
	}

	return polygon
}

func ReadMultiPolygon(pbf *pbf.Reader, endpos int, lengths []uint64, factor float64, dim int) [][][][]float64 {
	multipolygon := [][][][]float64{}
	if lengths == nil {
		for pbf.Pos < endpos {
			num_rings := pbf.ReadVarint()
			polygon := make([][][]float64, num_rings)
			for i := 0; i < num_rings; i++ {
				num := pbf.ReadVarint()
				polygon[i] = ReadLine(pbf, num, endpos, factor, dim, true)
			}
			multipolygon = append(multipolygon, polygon)
		}
	} else {
		var j = 1
		for i := 0; i < int(lengths[0]); i++ {
			polygon := [][][]float64{}
			for k := 0; k < int(lengths[j]); k++ {
				polygon = append(polygon, ReadLine(pbf, int(lengths[j+1+k]), endpos, factor, dim, true))
			}
			j += int(lengths[j]) + 1
			multipolygon = append(multipolygon, polygon)
		}
	}
	return multipolygon
}

func ReadBoundingBox(pbf *pbf.Reader, factor float64) []float64 {
	bb := make([]float64, 4)
	pbf.ReadVarint()
	bb[0] = float64(ReadSVarintPower(pbf, factor))
	bb[1] = float64(ReadSVarintPower(pbf, factor))
	bb[2] = float64(ReadSVarintPower(pbf, factor))
	bb[3] = float64(ReadSVarintPower(pbf, factor))
	return bb
}

func ConvertPt(pt []float64, factor float64, dim int) []int64 {
	if dim == 2 {
		newpt := make([]int64, 2)
		newpt[0] = int64(pt[0] * factor)
		newpt[1] = int64(pt[1] * factor)
		return newpt
	} else if dim == 3 {
		newpt := make([]int64, 3)
		newpt[0] = int64(pt[0] * factor)
		newpt[1] = int64(pt[1] * factor)
		newpt[2] = int64(pt[2] * factor)
		return newpt
	}
	return []int64{}
}

func paramEnc(value int64) uint64 {
	return uint64((value << 1) ^ (value >> 31))
}

func WritePoint(pbf *pbf.Writer, tag pbf.TagType, pt []float64, factor float64, dim int) {
	point := ConvertPt(pt, factor, dim)
	if dim == 2 {
		pbf.WritePackedUInt64(tag, []uint64{paramEnc(point[0]), paramEnc(point[1])})
	} else if dim == 3 {
		pbf.WritePackedUInt64(tag, []uint64{paramEnc(point[0]), paramEnc(point[1]), paramEnc(point[2])})
	}
}

func WriteLine(pbf *pbf.Writer, tag pbf.TagType, line [][]float64, factor float64, dim int) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	newline := make([]uint64, len(line)*dim)
	deltapt := make([]int64, dim)
	pt := make([]int64, dim)
	oldpt := make([]int64, dim)

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

		pt = ConvertPt(point, factor, dim)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
			if dim == 3 {
				newline[2] = paramEnc(pt[2])
			}
		} else {
			if dim == 2 {
				deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1]}
				newline[i*dim] = paramEnc(deltapt[0])
				newline[i*dim+1] = paramEnc(deltapt[1])
			} else if dim == 3 {
				deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1], pt[2] - oldpt[2]}
				newline[i*dim] = paramEnc(deltapt[0])
				newline[i*dim+1] = paramEnc(deltapt[1])
				newline[i*dim+2] = paramEnc(deltapt[2])
			}
		}
		oldpt = pt
	}
	pbf.WritePackedUInt64(tag, newline)
}

func MakeLine2(line [][]float64, factor float64, dim int) ([]uint64, []int64) {
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	newline := make([]uint64, len(line)*dim)
	deltapt := make([]int64, dim)
	pt := make([]int64, dim)
	oldpt := make([]int64, dim)

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

		pt = ConvertPt(point, factor, dim)
		if i == 0 {
			newline[0] = paramEnc(pt[0])
			newline[1] = paramEnc(pt[1])
			if dim == 3 {
				newline[2] = paramEnc(pt[2])
			}
		} else {
			if dim == 2 {
				deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1]}
				newline[i*2] = paramEnc(deltapt[0])
				newline[i*2+1] = paramEnc(deltapt[1])
			} else if dim == 3 {
				deltapt = []int64{pt[0] - oldpt[0], pt[1] - oldpt[1], pt[2] - oldpt[2]}
				newline[i*dim] = paramEnc(deltapt[0])
				newline[i*dim+1] = paramEnc(deltapt[1])
				newline[i*dim+2] = paramEnc(deltapt[2])
			}
		}
		oldpt = pt
	}

	return newline, []int64{int64(west * factor),
		int64(south * factor),
		int64(east * factor),
		int64(north * factor)}
}

func MakePolygon2(polygon [][][]float64, factor float64, dim int) ([]uint64, []int64) {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*dim))

		tmpgeom, tmpbb := MakeLine2(cont, factor, dim)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	return geometry, bb
}

func WriteMultiPolygon(pbf *pbf.Writer, tag pbf.TagType, multipolygon [][][][]float64, factor float64, dim int) []int64 {
	geometry := []uint64{}
	west, south, east, north := 180.0, 90.0, -180.0, -90.0
	west, south, east, north = west*factor, south*factor, east*factor, north*factor
	bb := []int64{int64(west), int64(south), int64(east), int64(north)}

	for _, polygon := range multipolygon {
		geometry = append(geometry, uint64(len(polygon)))
		tempgeom, tempbb := MakePolygon2(polygon, factor, dim)
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

func WritePolygon(pbf *pbf.Writer, tag pbf.TagType, polygon [][][]float64, factor float64, dim int) []int64 {
	geometry := []uint64{}
	bb := []int64{}
	for i, cont := range polygon {
		geometry = append(geometry, uint64(len(cont)*dim))

		tmpgeom, tmpbb := MakeLine2(cont, factor, dim)
		geometry = append(geometry, tmpgeom...)
		if i == 0 {
			bb = tmpbb
		}
	}
	pbf.WritePackedUInt64(tag, geometry)
	return bb
}

func WriteBoundingBox(pbf *pbf.Writer, tag pbf.TagType, box []float64, factor float64) {
	boxs := []uint64{
		paramEnc(int64(box[0] * factor)),
		paramEnc(int64(box[1] * factor)),
		paramEnc(int64(box[2] * factor)),
		paramEnc(int64(box[3] * factor)),
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
