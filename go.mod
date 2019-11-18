module go-issued-service

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/go-redis/redis v6.15.6+incompatible
	github.com/gorilla/websocket v1.4.1
	github.com/onsi/ginkgo v1.10.3 // indirect
	github.com/onsi/gomega v1.7.1 // indirect
	github.com/spf13/cobra v0.0.5
)

replace (
	golang.org/x/crypto v0.0.0-20181203042331-505ab145d0a9 => github.com/golang/crypto v0.0.0-20181203042331-505ab145d0a9
	golang.org/x/net v0.0.0-20180925072008-f04abc6bdfa7 => github.com/golang/net v0.0.0-20180911220305-26e67e76b6c3
	golang.org/x/sync v0.0.0-20180314180146-1d60e4601c6f => github.com/golang/sync v0.0.0-20180314180146-1d60e4601c6f
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190602015325-4c4f7f33c9ed
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)
