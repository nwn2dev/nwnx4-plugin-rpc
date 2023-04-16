package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
	"strings"
	"time"
)

import (
	"google.golang.org/grpc"
)

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
const rpcStartCallAction int32 = 1
const rpcEndCallAction int32 = 2

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
func (p *rpcPlugin) getInt(sFunction, sParam1 string, nParam2 int32) int32 {
	// CallAction()
	switch sFunction {
	case rpcGetInt:
		if v, found := p.globalCallActionResponse.Data[sParam1]; found {
			value := v.GetNValue()
			if nParam2 == rpcEndCallAction {
				p.resetCallAction()
			}

			return value
		}

		return 0
	case rpcGetBool:
		if v, found := p.globalCallActionResponse.Data[sParam1]; found {
			value := v.GetBValue()
			if nParam2 == rpcStartCallAction {
				p.resetCallAction()
			}
			if value {
				return 1
			}
		}

		return 0
	}

	// NWNXGetInt()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return 0
	}

	return client.NWNXGetInt(sFunction, sParam1, nParam2, p.config.getTimeout())
}

// setInt the body of the NWNXSetInt() call
func (p *rpcPlugin) setInt(sFunction, sParam1 string, nParam2, nValue int32) {
	// CallAction()
	switch sFunction {
	case rpcSetInt:
		if nParam2 == rpcStartCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_NValue{
				NValue: nValue,
			},
		}

		return
	case rpcSetBool:
		if nParam2 == rpcStartCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_BValue{
				BValue: nValue != 0,
			},
		}

		return
	}

	// NWNXSetInt()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return
	}

	client.NWNXSetInt(sFunction, sParam1, nParam2, nValue, p.config.getTimeout())
}

// getFloat the body of the NWNXGetFloat() call
func (p *rpcPlugin) getFloat(sFunction, sParam1 string, nParam2 int32) float32 {
	// CallAction()
	switch sFunction {
	case rpcGetFloat:
		if v, found := p.globalCallActionResponse.Data[sParam1]; found {
			value := v.GetFValue()
			if nParam2 == rpcEndCallAction {
				p.resetCallAction()
			}

			return value
		}

		return 0.0
	}

	// NWNXGetFloat()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return 0.0
	}

	return client.NWNXGetFloat(sFunction, sParam1, nParam2, p.config.getTimeout())
}

// setFloat the body of the NWNXSetFloat() call
func (p *rpcPlugin) setFloat(sFunction, sParam1 string, nParam2 int32, fValue float32) {
	// CallAction()
	switch sFunction {
	case rpcSetFloat:
		if nParam2 == rpcStartCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_FValue{
				FValue: fValue,
			},
		}

		return
	}

	// NWNXSetFloat()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return
	}

	client.NWNXSetFloat(sFunction, sParam1, nParam2, fValue, p.config.getTimeout())
}

// getString the body of the NWNXGetString() call
func (p *rpcPlugin) getString(sFunction, sParam1 string, nParam2 int32) string {
	// CallAction()
	switch sFunction {
	case rpcGetString:
		if v, found := p.globalCallActionResponse.Data[sParam1]; found {
			value := v.GetSValue()
			if nParam2 == rpcEndCallAction {
				p.resetCallAction()
			}

			return value
		}

		return ""
	}

	// NWNXGetString()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return ""
	}

	return client.NWNXGetString(sFunction, sParam1, nParam2, p.config.getTimeout())
}

// setString the body of the NWNXSetString() call
func (p *rpcPlugin) setString(sFunction, sParam1 string, nParam2 int32, sValue string) {
	// CallAction()
	switch sFunction {
	case rpcSetString:
		if nParam2 == rpcStartCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_SValue{
				SValue: sValue,
			},
		}

		return
	case rpcResetCallAction:
		p.resetCallAction()

		return
	case rpcCallAction:
		// sParam1_ holds the client identifier
		client, ok := p.getRpcClient(sParam1)
		if !ok {
			return
		}

		// sValue_ contains the action
		p.globalCallActionRequest.Action = sValue
		response, err := client.callAction(p.globalCallActionRequest, p.config.getTimeout())
		p.resetCallAction()

		if err == nil {
			*p.globalCallActionResponse = *response
		}

		return
	}

	// NWNXSetString()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return
	}

	client.NWNXSetString(sFunction, sParam1, nParam2, sValue, p.config.getTimeout())
}

// getGffSize called at the start of RetrieveCampaignObject
func (p *rpcPlugin) getGffSize(sFunction, sVarName string, _ int32) uint32 {
	// CallAction()
	switch sFunction {
	case rpcGetGff:
		if v, found := p.globalCallActionResponse.Data[sVarName]; found {
			return uint32(len(v.GetGffValue()))
		}

		return 0
	}

	// RetrieveCampaignObject()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return 0
	}

	return client.SCORCOGetGFFSize(sVarName, p.config.getTimeout())
}

// getGff called at the end of RetrieveCampaignObject
func (p *rpcPlugin) getGff(sFunction, sVarName string, nParam2 int32) []byte {
	// CallAction()
	switch sFunction {
	case rpcGetGff:
		if v, found := p.globalCallActionResponse.Data[sVarName]; found {
			value := v.GetGffValue()
			if nParam2 == rpcEndCallAction {
				p.resetCallAction()
			}

			return value
		}

		return nil
	}

	// RetrieveCampaignObject()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return nil
	}

	return client.SCORCOGetGFF(sVarName, p.config.getTimeout())
}

// setGff called during StoreCampaignObject
func (p *rpcPlugin) setGff(sFunction, sVarName string, nParam2 int32, gffData []byte, gffDataSize uint32) {
	// CallAction()
	switch sFunction {
	case rpcSetGff:
		if nParam2 == rpcStartCallAction {
			p.resetCallAction()
		}

		p.globalCallActionRequest.Params[sVarName] = &pb.Value{
			ValueType: &pb.Value_GffValue{
				GffValue: gffData,
			},
		}

		return
	}

	// StoreCampaignObject()
	client, ok := p.getRpcClient(sFunction)
	if !ok {
		return
	}

	client.SCORCOSetGFF(sVarName, gffData, gffDataSize, p.config.getTimeout())
}
