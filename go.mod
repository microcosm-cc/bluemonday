module github.com/microcosm-cc/bluemonday

go 1.21

require (
	github.com/aymerick/douceur v0.2.0
	golang.org/x/net v0.17.0
)

require (
	github.com/PuerkitoBio/goquery v1.8.1 // indirect
	github.com/andybalholm/cascadia v1.3.2 // indirect
	github.com/gorilla/css v1.0.0 // indirect
)

retract [v1.0.0, v1.0.24] // Retract older versions as only latest is to be depended upon
