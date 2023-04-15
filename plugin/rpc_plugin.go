package main

import "C"
import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
	"strconv"
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
	"google.golang.org/grpc"
)

const rpcGffVarNameSeparator = "!"

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
const rpcCallActionParam2Default int32 = -1
const rpcCallActionParam2ResetCallAction int32 = 1

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

func (p *rpcPlugin) resetCallAction() {
	p.globalCallActionRequest = newCallActionRequest()
	p.globalCallActionResponse = newCallActionResponse()
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

	return client.NWNXGetInt(sFunction_, sParam1_, nParam2_, p.config.getTimeout())
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
		if nParam2_ == rpcCallActionParam2ResetCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_NValue{
				NValue: nValue_,
			},
		}

		return
	case rpcSetBool:
		if nParam2_ == rpcCallActionParam2ResetCallAction {
			p.resetCallAction()
		}

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

	client.NWNXSetInt(sFunction_, sParam1_, nParam2_, nValue_, p.config.getTimeout())
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

	return client.NWNXGetFloat(sFunction_, sParam1_, nParam2_, p.config.getTimeout())
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
		if nParam2_ == rpcCallActionParam2ResetCallAction {
			p.resetCallAction()
		}

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

	client.NWNXSetFloat(sFunction_, sParam1_, nParam2_, fValue_, p.config.getTimeout())
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

	return client.NWNXGetString(sFunction_, sParam1_, nParam2_, p.config.getTimeout())
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
		if nParam2_ == rpcCallActionParam2ResetCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1_] = &pb.Value{
			ValueType: &pb.Value_SValue{
				SValue: sValue_,
			},
		}

		return
	case rpcResetCallAction:
		p.resetCallAction()

		return
	case rpcCallAction:
		// sParam1_ holds the client identifier
		client, ok := p.getRpcClient(sParam1_)
		if !ok {
			return
		}

		// sValue_ contains the action
		p.globalCallActionRequest.Action = sValue_
		response, err := client.callAction(p.globalCallActionRequest, p.config.getTimeout())
		p.resetCallAction()

		if err == nil {
			*p.globalCallActionResponse = *response
		}

		return
	}

	// NWNXSetString()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	client.NWNXSetString(sFunction_, sParam1_, nParam2_, sValue_, p.config.getTimeout())
}

// getGffSize called at the start of RetrieveCampaignObject
func (p *rpcPlugin) getGffSize(sVarName *C.char) uint32 {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	var sFunction_, sVarName_ string
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else {
		return 0
	}
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

	return client.SCORCOGetGFFSize(sVarName_, p.config.getTimeout())
}

// getGff called at the end of RetrieveCampaignObject
func (p *rpcPlugin) getGff(sVarName *C.char, _ *C.uint8_t, _ C.size_t) []byte {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	var sFunction_, sVarName_ string
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else {
		return nil
	}
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

	return client.SCORCOGetGFF(sVarName_, p.config.getTimeout())
}

// setGff called during StoreCampaignObject
func (p *rpcPlugin) setGff(sVarName *C.char, gffData *C.uint8_t, gffDataSize C.size_t) {
	gffDataSize_ := uint32(gffDataSize)
	gffData_ := C.GoBytes(unsafe.Pointer(gffData), C.int(gffDataSize))
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 3)
	var sFunction_, sParam2_, sVarName_ string
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else if len(splits) == 3 {
		sFunction_, sParam2_, sVarName_ = splits[0], splits[1], splits[2]
	} else {
		return
	}
	log.Debugf("SetGFFSize(%s, %x, %d)", sVarName, gffData, gffDataSize)

	// CallAction()
	switch sFunction_ {
	case rpcSetGff:
		if nParam2_, err := strconv.Atoi(sParam2_); err == nil && int32(nParam2_) == rpcCallActionParam2ResetCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sVarName_] = &pb.Value{
			ValueType: &pb.Value_GffValue{
				GffValue: gffData_,
			},
		}

		return
	}

	// StoreCampaignObject()
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	client.SCORCOSetGFF(sVarName_, gffData_, gffDataSize_, p.config.getTimeout())
}
