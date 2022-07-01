module github.com/microcosm-cc/bluemonday

go 1.18

require (
	github.com/aymerick/douceur v0.2.0
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e
)

require github.com/gorilla/css v1.0.0 // indirect

retract [v1.0.0, v1.0.18] // Retract older versions as only latest is to be depended upon
