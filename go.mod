module gitlab.orayer.com/golang/pubsub

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.1
	github.com/mkevac/debugcharts v0.0.0-20180124214838-d3203a8fa926
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.1
	github.com/shirou/gopsutil v2.19.10+incompatible // indirect
	github.com/shirou/w32 v0.0.0-20160930032740-bb4de0191aa4 // indirect
	github.com/ugorji/go v1.1.7 // indirect
	gitlab.orayer.com/golang/errors v0.0.0-20191022062103-5ac945e04fd9
	gitlab.orayer.com/golang/server v0.0.0-20191119074007-4e60456b9561
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c
	golang.org/x/sys v0.0.0-20191119195528-f068ffe820e4 // indirect
	google.golang.org/grpc v1.21.0
)

replace (
	golang.org/x/crypto v0.0.0-20191117063200-497ca9f6d64f => github.com/golang/crypto v0.0.0-20191117063200-497ca9f6d64f
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c => github.com/golang/net v0.0.0-20190503192946-f4e77d36d62c
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e => github.com/golang/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190602015325-4c4f7f33c9ed
	golang.org/x/text v0.3.2 => github.com/golang/text v0.3.2
	google.golang.org/grpc => github.com/grpc/grpc-go v1.21.0
)
