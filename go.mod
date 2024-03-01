module github.com/johnewart/freighter

go 1.22.0

require (
	github.com/google/go-containerregistry v0.19.0
	github.com/hanwen/go-fuse/v2 v2.5.0
	google.golang.org/appengine v1.6.8
	google.golang.org/grpc v1.62.0
	google.golang.org/protobuf v1.32.0
	zombiezen.com/go/log v1.1.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240227224415-6ceb2ff114de // indirect
	gorm.io/driver/sqlite v1.5.5 // indirect
	gorm.io/gorm v1.25.7 // indirect
)

replace github.com/google/go-containerregistry => ../go-containerregistry
