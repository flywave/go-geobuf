package io

import (
	"math"

	"github.com/flywave/go-pbf"
)

type Decode struct {
	Keys    []string
	Values  []interface{}
	Lengths []uint64
	Dim     int
	Factor  float64
	reader  *pbf.Reader
}

func NewDecode(reader *pbf.Reader) *Decode {
	return &Decode{Keys: make([]string, 0), Values: make([]interface{}, 0), Lengths: make([]uint64, 0), Dim: 2, Factor: math.Pow(10.0, 7.0), reader: reader}
}

func (d *Decode) readDataField() {
	key, _ := d.reader.ReadTag()
	if key == DATA_KEYS {
		d.Keys = append(d.Keys, d.reader.ReadString())
	} else if key == DIMENSIONS {
		d.Dim = d.reader.ReadVarint()
	} else if key == PRECISION {
		d.Factor = math.Pow(10, float64(d.reader.ReadVarint()))
	} else if key == DATA_TYPE_FEATURE_COLLECTION {
		d.readFeatureCollection()
	} else if key == DATA_TYPE_FEATURE {
		d.readFeature()
	} else if key == DATA_TYPE_GEOMETRY {
		d.readGeometry()
	}
}

func (d *Decode) readFeatureCollection() {

}

func (d *Decode) readFeature() {

}

func (d *Decode) readGeometry() {

}

func (d *Decode) readFeatureCollectionField() {
}

func (d *Decode) readFeatureField() {
}

func (d *Decode) readGeometryField() {
}

func (d *Decode) readCoords() {
}

func (d *Decode) readValue() {

}

func (d *Decode) readProps() {

}
