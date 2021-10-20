module github.com/max-gui/logagent

go 1.15

require (
	github.com/antonfisher/nested-logrus-formatter v1.3.1
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/gomodule/redigo v1.8.4
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/mattn/go-colorable v0.1.11 // indirect
	github.com/max-gui/regagent v0.1.1
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/zsais/go-gin-prometheus v0.1.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect

)

replace github.com/max-gui/regagent => ../Regagent
