module github.com/flywave/go-geobuf

go 1.23.0

toolchain go1.24.4

require (
	github.com/flywave/go-geom v0.0.0-20250607125323-f685bf20f12c
	github.com/flywave/go-mapbox v0.0.0-20250314092441-27874854ad1b
	github.com/flywave/go-pbf v0.0.0-20230306063816-5e5b0da27bbd
)

replace github.com/flywave/go-mapbox => ../go-mapbox

replace github.com/flywave/go-geom => ../go-geom
