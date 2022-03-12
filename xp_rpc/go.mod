module xp_rpc

go 1.18

replace nwnx4.org/xp_rpc/proto => ./proto

replace nwnx4.org/xp_rpc/proto/nwscript => ./proto/NWScript

require (
	github.com/sirupsen/logrus v1.8.1
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
	nwnx4.org/xp_rpc/proto v0.0.0-00010101000000-000000000000
	nwnx4.org/xp_rpc/proto/nwscript v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)
