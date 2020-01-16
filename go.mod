module github.com/HyperService-Consortium/go-ves

go 1.12

replace (
	github.com/HyperService-Consortium/go-hexutil => github.com/HyperService-Consortium/go-hexutil v1.0.1
	github.com/HyperService-Consortium/go-uip => github.com/HyperService-Consortium/go-uip v0.0.0-20200116083857-d2061bd0d6df
)

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/HyperService-Consortium/NSB v0.7.3-0.20191110033951-7a0fad01c5fb
	github.com/HyperService-Consortium/go-ethabi v0.9.1
	github.com/HyperService-Consortium/go-hexutil v1.0.1 // indirect
	github.com/HyperService-Consortium/go-mpt v1.1.1
	github.com/HyperService-Consortium/go-rlp v1.0.2
	github.com/HyperService-Consortium/go-uip v0.0.0-20200116083857-d2061bd0d6df
	github.com/Myriad-Dreamin/minimum-lib v0.0.0-20191109053555-ffc58e6d4591
	github.com/Myriad-Dreamin/mydrest v1.0.1
	github.com/Myriad-Dreamin/screenrus v1.0.0
	github.com/boltdb/bolt v1.3.1
	github.com/gin-contrib/cors v1.3.0
	github.com/gin-gonic/gin v1.4.1-0.20190628155452-f65018d7b1f1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-querystring v1.0.0
	github.com/gorilla/websocket v1.4.0
	github.com/imroc/req v0.2.4
	github.com/mattn/go-sqlite3 v1.11.1-0.20191008083825-3f45aefa8dc8
	github.com/prologic/bitcask v0.3.5
	github.com/sirupsen/logrus v1.4.2
	github.com/syndtr/goleveldb v1.0.1-0.20190318030020-c3a204f8e965
	github.com/tidwall/gjson v1.3.2
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf
	golang.org/x/net v0.0.0-20191028085509-fe3aa8a45271
	google.golang.org/grpc v1.23.0
)
