module github.com/microcosm-cc/bluemonday

go 1.19

require (
	github.com/aymerick/douceur v0.2.0
	golang.org/x/net v0.0.0-20221002022538-bcab6841153b
)

require github.com/gorilla/css v1.0.0 // indirect

retract [v1.0.0, v1.0.18] // Retract older versions as only latest is to be depended upon

retract v1.0.19 // Uses older version of golang.org/x/net
