module src

go 1.18

replace nwnx4.org/src/proto => ./proto

replace nwnx4.org/src/proto/nwscript => ./proto/nwscript

require (
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	nwnx4.org/src/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)
