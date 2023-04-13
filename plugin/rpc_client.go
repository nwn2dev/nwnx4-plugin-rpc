package main

import (
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
)

// rpcClient contains the clients to RPC microservices
type rpcClient struct {
	isValid             bool
	name                string
	url                 string
	nwnxServiceClient   pb.NWNXServiceClient
	scorcoServiceClient pb.SCORCOServiceClient
	actionServiceClient pb.ActionServiceClient
}
