package main

import "C"
import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
	"strings"
	"time"
	"unsafe"
)

/*
#include <stdlib.h>
#include <stdint.h>
*/
import "C"
import (
	"context"
	"google.golang.org/grpc"
)

const rpcGffVarNameSeparator = "///"

const rpcGetInt string = "RPC_GET_INT_"
const rpcSetInt string = "RPC_SET_INT_"
const rpcGetBool string = "RPC_GET_BOOL_"
const rpcSetBool string = "RPC_SET_BOOL_"
const rpcGetFloat string = "RPC_GET_FLOAT_"
const rpcSetFloat string = "RPC_SET_FLOAT_"
const rpcGetString string = "RPC_GET_STRING_"
const rpcSetString string = "RPC_SET_STRING_"
const rpcGetGff string = "RPC_GET_GFF_"
const rpcSetGff string = "RPC_SET_GFF_"
const rpcResetCallAction string = "RPC_RESET_CALL_ACTION_"
const rpcCallAction string = "RPC_CALL_ACTION_"

type rpcPlugin struct {
	config                   rpcConfig
	clients                  map[string]*rpcClient
	globalCallActionRequest  *pb.CallActionRequest
	globalCallActionResponse *pb.CallActionResponse
}

// newRpcPlugin builds and returns a new RPC plugin
func newRpcPlugin() *rpcPlugin {
	return &rpcPlugin{
		config:                   rpcConfig{},
		clients:                  make(map[string]*rpcClient),
		globalCallActionRequest:  newCallActionRequest(),
		globalCallActionResponse: newCallActionResponse(),
	}
}

func newCallActionRequest() *pb.CallActionRequest {
	return &pb.CallActionRequest{
		Action: "",
		Params: make(map[string]*pb.Value),
	}
}

func newCallActionResponse() *pb.CallActionResponse {
	return &pb.CallActionResponse{
		Data: make(map[string]*pb.Value),
	}
}

// init initializes the RPC plugin
func (p *rpcPlugin) init() {
	log.Info("Initializing RPC plugin")

	// Set the log level based on what was passed if it matches a level
	for _, logLevel := range log.AllLevels {
		if strings.EqualFold(logLevel.String(), p.config.Log.LogLevel) {
			log.Infof("Log level set as %s", logLevel)
			log.SetLevel(logLevel)
			break
		}
	}

	// Build out the clients
	for name, url := range p.config.Clients {
		plugin.addRpcClient(name, url)
	}

	log.Info("Initialized RPC plugin")
}

// addRpcClient adds an RPC client
// Runs on an rpcPlugin and adds a client by name and URL
func (p *rpcPlugin) addRpcClient(name, url string) {
	log.Infof("Adding client: %s@%s", name, url)
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Unable to attach client: %s@%s", name, url)

		p.clients[name] = &rpcClient{
			isValid:             false,
			name:                name,
			url:                 url,
			nwnxServiceClient:   nil,
			scorcoServiceClient: nil,
			actionServiceClient: nil,
		}

		return
	}

	p.clients[name] = &rpcClient{
		isValid:             true,
		name:                name,
		url:                 url,
		nwnxServiceClient:   pb.NewNWNXServiceClient(conn),
		scorcoServiceClient: pb.NewSCORCOServiceClient(conn),
		actionServiceClient: pb.NewActionServiceClient(conn),
	}

	log.Infof("Established client connection and stubs: %s@%s", name, url)
}

// getRpcClient will get an rpcClient by name
func (p *rpcPlugin) getRpcClient(name string) (*rpcClient, bool) {
	log.Infof("Getting client with name: %s", name)
	rpcClient, exists := p.clients[name]

	if !exists {
		log.Errorf("Client not declared: %s", name)

		return nil, false
	}

	// Handle invalid clients; try to recreate
	if !rpcClient.isValid {
		url := p.config.Clients[name]
		delay := p.config.getDelay()

		// Invalid client; try to recreate
		for i := 0; i < p.config.PerClient.Retries; i++ {
			log.Infof("Reading client: %s@%s", name, url)
			p.addRpcClient(name, url)
			rpcClient, exists = p.clients[name]

			if exists && rpcClient.isValid {
				return rpcClient, true
			}

			log.Infof("Adding client failed; delaying for %ds", p.config.getDelay())
			time.Sleep(delay)
		}

		log.Errorf("Client is still invalid; could not be setup: %s@#%s", name, url)
	}

	return rpcClient, true
}

// getInt the body of the NWNXGetInt() call
func (p *rpcPlugin) getInt(sFunction, sParam1 *C.char, nParam2 C.int) int32 {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.Debugf("NWNXGetInt(%s, %s, %d)", sFunction_, sParam1_, nParam2_)

	// CallAction()
	switch sFunction_ {
	case rpcGetInt:
		if v, found := p.globalCallActionResponse.Data[sParam1_]; found {
			return v.GetNValue()
		}

		return 0
	case rpcGetBool:
		if v, found := p.globalCallActionResponse.Data[sParam1_]; found {
			if v.GetBValue() {
				return 1
			}
		}

		return 0
	}

	// NWNXGetInt()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXGetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetInt(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Errorf("Call to GetInt returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return 0
	}

	return response.Value
}

// setInt the body of the NWNXSetInt() call
func (p *rpcPlugin) setInt(sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	nValue_ := int32(nValue)
	log.Debugf("NWNXSetInt(%s, %s, %d, %d)", sFunction_, sParam1_, nParam2_, nValue_)

	// CallAction()
	switch sFunction_ {
	case rpcSetInt:
		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_NValue{
				NValue: nValue_,
			},
		}

		return
	case rpcSetBool:
		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_BValue{
				BValue: nValue != 0,
			},
		}

		return
	}

	// NWNXSetInt()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXSetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		NValue:    nValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetInt(ctx, &request); err != nil {
		client.isValid = false
		log.Errorf("Call to SetInt returned error: %s; %s, %s, %d, %d",
			err, request.SFunction, request.SParam1, request.NParam2, request.NValue)
	}
}

// getFloat the body of the NWNXGetFloat() call
func (p *rpcPlugin) getFloat(sFunction, sParam1 *C.char, nParam2 C.int) float32 {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.Debugf("NWNXGetFloat(%s, %s, %d)", sFunction_, sParam1_, nParam2_)

	// CallAction()
	switch sFunction_ {
	case rpcGetFloat:
		if v, found := p.globalCallActionResponse.Data[sParam1_]; found {
			return v.GetFValue()
		}

		return 0.0
	}

	// NWNXGetFloat()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0.0
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXGetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetFloat(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Errorf("Call to GetFloat returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return 0.0
	}

	return response.Value
}

// setFloat the body of the NWNXSetFloat() call
func (p *rpcPlugin) setFloat(sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	fValue_ := float32(fValue)
	log.Debugf("NWNXSetFloat(%s, %s, %d, %d, %f)", sFunction_, sParam1_, nParam2_, fValue_)

	// CallAction()
	switch sFunction_ {
	case rpcSetFloat:
		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_FValue{
				FValue: fValue_,
			},
		}

		return
	}

	// NWNXSetFloat()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXSetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		FValue:    fValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetFloat(ctx, &request); err != nil {
		client.isValid = false
		log.Errorf("Call to SetFloat returned error: %s; %s, %s, %d, %f",
			err, request.SFunction, request.SParam1, request.NParam2, request.FValue)
	}
}

// getString the body of the NWNXGetString() call
func (p *rpcPlugin) getString(sFunction, sParam1 *C.char, nParam2 C.int) string {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.Debugf("NWNXGetString(%s, %s, %d)", sFunction_, sParam1_, nParam2_)

	// CallAction()
	switch sFunction_ {
	case rpcGetString:
		if v, found := p.globalCallActionResponse.Data[sParam1_]; found {
			return v.GetSValue()
		}

		return ""
	}

	// NWNXGetString()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return ""
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXGetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetString(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Errorf("Call to GetString returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2)

		return ""
	}

	return response.Value
}

// setString the body of the NWNXSetString() call
func (p *rpcPlugin) setString(sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	sValue_ := C.GoString(sValue)
	log.Debugf("NWNXSetString(%s, %s, %d, %s)", sFunction_, sParam1_, nParam2_, sValue_)

	// CallAction()
	switch sFunction_ {
	case rpcSetString:
		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_SValue{
				SValue: sValue_,
			},
		}

		return
	case rpcResetCallAction:
		p.globalCallActionRequest = newCallActionRequest()
		p.globalCallActionResponse = newCallActionResponse()

		return
	case rpcCallAction:
		// sParam1_ holds the client identifier
		client, ok := p.getRpcClient(sParam1_)
		if !ok {
			return
		}

		// sValue_ contains the action
		p.callAction(client, sValue_)

		return
	}

	// NWNXSetString()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.NWNXSetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		SValue:    sValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetString(ctx, &request); err != nil {
		client.isValid = false
		log.Errorf("Call to SetString returned error: %s; %s, %s, %d, %s",
			err, request.SFunction, request.SParam1, request.NParam2, request.SValue)
	}
}

// callAction call an action on the client specified
func (p *rpcPlugin) callAction(client *rpcClient, action string) {
	p.globalCallActionRequest.Action = action
	log.Infof("CallAction(): %s@%s, %+v", client.name, client.url, p.globalCallActionRequest)

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	response, err := client.actionServiceClient.CallAction(ctx, p.globalCallActionRequest)
	p.globalCallActionRequest = newCallActionRequest()
	p.globalCallActionResponse = newCallActionResponse()

	if err != nil {
		client.isValid = false
		log.Errorf("Error sending request: %s", err)

		return
	}

	*p.globalCallActionResponse = *response

	return
}

// getGffSize called at the start of RetrieveCampaignObject
func (p *rpcPlugin) getGffSize(sVarName *C.char) uint32 {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	if len(splits) != 2 {
		return 0
	}
	sFunction_, sVarName_ := splits[0], splits[1]
	log.Debugf("GetGFFSize(%s)", sVarName)

	// CallAction()
	switch sFunction_ {
	case rpcGetGff:
		if v, found := p.globalCallActionResponse.Data[sVarName_]; found {
			return uint32(len(v.GetGffValue()))
		}

		return 0
	}

	// RetrieveCampaignObject()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.SCORCOGetGFFSizeRequest{
		SVarName: sVarName_,
	}
	response, err := client.scorcoServiceClient.SCORCOGetGFFSize(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Errorf("Call to GetGFFSize returned error: %s; %s", err, request.SVarName)

		return 0
	}

	return response.Size
}

func (p *rpcPlugin) getGff(sVarName *C.char, _ *C.uint8_t, _ C.size_t) []byte {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	if len(splits) != 2 {
		return nil
	}
	sFunction_, sVarName_ := splits[0], splits[1]
	log.Debugf("GetGFF(%s)", sVarName)

	// CallAction()
	switch sFunction_ {
	case rpcGetGff:
		if v, found := p.globalCallActionResponse.Data[sVarName_]; found {
			return v.GetGffValue()
		}

		return nil
	}

	// RetrieveCampaignObject()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	request := pb.SCORCOGetGFFRequest{
		SVarName: sVarName_,
	}
	response, err := client.scorcoServiceClient.SCORCOGetGFF(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Errorf("Call to GetGFF returned error: %s; %s", err, request.SVarName)

		return nil
	}

	return response.GetGffData()
}

func (p *rpcPlugin) setGff(sVarName *C.char, gffData *C.uint8_t, gffDataSize C.size_t) {
	gffDataSize_ := uint32(gffDataSize)
	ptr := unsafe.Pointer(gffData)
	defer C.free(ptr)
	gffData_ := (*[1 << 30]byte)(ptr)[:gffDataSize_:gffDataSize_]
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	if len(splits) != 2 {
		return
	}
	sFunction_, sVarName_ := splits[0], splits[1]
	log.Debugf("SetGFFSize(%s, %x, %d)", sVarName, gffData, gffDataSize)

	// CallAction()
	switch sFunction_ {
	case rpcSetGff:
		gffValue := make([]byte, gffDataSize_)
		copy(gffValue, gffData_)

		p.globalCallActionRequest.Params[sVarName_] = &pb.Value{
			ValueType: &pb.Value_GffValue{
				GffValue: gffValue,
			},
		}

		return
	}

	// StoreCampaignObject()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.getTimeout())
	defer cancel()
	var request = pb.SCORCOSetGFFRequest{
		SVarName:    sVarName_,
		GffData:     gffData_,
		GffDataSize: gffDataSize_,
	}
	if _, err := client.scorcoServiceClient.SCORCOSetGFF(ctx, &request); err != nil {
		client.isValid = false
		log.Errorf("Call to SetGFF returned error: %s; %s", err, request.SVarName)
	}
}
