module escapade

go 1.12

// +heroku goVersion go1.12
// +heroku install ./cmd/api/

require (
	github.com/bradfitz/gomemcache v0.0.0-20180710155616-bc664df96737
	github.com/garyburd/redigo v1.6.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/golang/protobuf v1.3.0
	github.com/gorilla/handlers v1.4.0
	github.com/gorilla/mux v1.7.0
	github.com/gorilla/websocket v1.4.0
	github.com/jinzhu/gorm v1.9.2
	github.com/jinzhu/inflection v0.0.0-20180308033659-04140366298a // indirect
	github.com/karrick/godirwalk v1.8.0 // indirect
	github.com/lib/pq v1.0.0
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/client_golang v0.9.2
	github.com/segmentio/ksuid v1.0.2
	github.com/streadway/amqp v0.0.0-20190225234609-30f8ed68076e
	github.com/uudashr/gopkgs v2.0.1+incompatible // indirect
	golang.org/x/crypto v0.0.0-20190228161510-8dd112bcdc25
	golang.org/x/net v0.0.0-20190301231341-16b79f2e4e95
	google.golang.org/grpc v1.19.0
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)
