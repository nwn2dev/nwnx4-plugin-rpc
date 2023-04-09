module src

go 1.18

replace nwnx4.org/plugin/proto => ./proto

replace nwnx4.org/plugin/proto/nwscript => ./proto/nwscript

require (
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.54.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)
