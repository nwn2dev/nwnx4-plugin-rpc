package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
	"time"
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

func (c rpcClient) callAction(request *pb.CallActionRequest, timeout time.Duration) (*pb.CallActionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	response, err := c.actionServiceClient.CallAction(ctx, request)

	if err != nil {
		c.isValid = false
		log.Errorf("Error sending request: %s", err)
	}

	return response, err
}

func (c rpcClient) NWNXGetInt(sFunction, sParam1 string, nParam2 int32, timeout time.Duration) int32 {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXGetIntRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
	}

	response, err := c.nwnxServiceClient.NWNXGetInt(ctx, &request)
	if err != nil {
		c.isValid = false
		log.Errorf("Call to GetInt returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return 0
	}

	return response.Value
}

func (c rpcClient) NWNXSetInt(sFunction, sParam1 string, nParam2, nValue int32, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXSetIntRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
		NValue:    nValue,
	}
	if _, err := c.nwnxServiceClient.NWNXSetInt(ctx, &request); err != nil {
		c.isValid = false
		log.Errorf("Call to SetInt returned error: %s; %s, %s, %d, %d",
			err, request.SFunction, request.SParam1, request.NParam2, request.NValue)
	}
}

func (c rpcClient) NWNXGetFloat(sFunction, sParam1 string, nParam2 int32, timeout time.Duration) float32 {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXGetFloatRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
	}
	response, err := c.nwnxServiceClient.NWNXGetFloat(ctx, &request)
	if err != nil {
		c.isValid = false
		log.Errorf("Call to GetFloat returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return 0.0
	}

	return response.Value
}

func (c rpcClient) NWNXSetFloat(sFunction, sParam1 string, nParam2 int32, fValue float32, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXSetFloatRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
		FValue:    fValue,
	}
	if _, err := c.nwnxServiceClient.NWNXSetFloat(ctx, &request); err != nil {
		c.isValid = false
		log.Errorf("Call to SetFloat returned error: %s; %s, %s, %d, %f",
			err, request.SFunction, request.SParam1, request.NParam2, request.FValue)
	}
}

func (c rpcClient) NWNXGetString(sFunction, sParam1 string, nParam2 int32, timeout time.Duration) string {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXGetStringRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
	}
	response, err := c.nwnxServiceClient.NWNXGetString(ctx, &request)
	if err != nil {
		c.isValid = false
		log.Errorf("Call to GetString returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return ""
	}

	return response.Value
}

func (c rpcClient) NWNXSetString(sFunction string, sParam1 string, nParam2 int32, sValue string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.NWNXSetStringRequest{
		SFunction: sFunction,
		SParam1:   sParam1,
		NParam2:   nParam2,
		SValue:    sValue,
	}
	if _, err := c.nwnxServiceClient.NWNXSetString(ctx, &request); err != nil {
		c.isValid = false
		log.Errorf("Call to SetString returned error: %s; %s, %s, %d, %s",
			err, request.SFunction, request.SParam1, request.NParam2, request.SValue)
	}
}

func (c rpcClient) SCORCOGetGFFSize(sVarName string, timeout time.Duration) uint32 {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.SCORCOGetGFFSizeRequest{
		SVarName: sVarName,
	}
	response, err := c.scorcoServiceClient.SCORCOGetGFFSize(ctx, &request)
	if err != nil {
		c.isValid = false
		log.Errorf("Call to GetGFFSize returned error: %s; %s", err, request.SVarName)

		return 0
	}

	return response.Size
}

func (c rpcClient) SCORCOGetGFF(sVarName string, timeout time.Duration) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	request := pb.SCORCOGetGFFRequest{
		SVarName: sVarName,
	}
	response, err := c.scorcoServiceClient.SCORCOGetGFF(ctx, &request)
	if err != nil {
		c.isValid = false
		log.Errorf("Call to GetGFFSize returned error: %s; %s", err, request.SVarName)

		return nil
	}

	return response.GetGffData()
}

func (c rpcClient) SCORCOSetGFF(sVarName string, gffData []byte, gffDataSize uint32, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var request = pb.SCORCOSetGFFRequest{
		SVarName:    sVarName,
		GffData:     gffData,
		GffDataSize: gffDataSize,
	}
	if _, err := c.scorcoServiceClient.SCORCOSetGFF(ctx, &request); err != nil {
		c.isValid = false
		log.Errorf("Call to SetGFF returned error: %s; %s", err, request.SVarName)
	}
}
