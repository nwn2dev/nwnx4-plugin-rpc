package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/pkcs12"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	pb "nwnx4.org/nwn2dev/xp_rpc/proto"
	"os"
	"strings"
	"time"
)

import (
	"google.golang.org/grpc"
)

const rpcResetGeneric string = "RPC_RESET_GENERIC_"
const rpcBuildGeneric string = "RPC_BUILD_GENERIC_"
const rpcBuildGenericStream string = "RPC_BUILD_GENERIC_STREAM_"
const rpcPullGenericStream string = "RPC_PULL_GENERIC_STREAM_"

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

const rpcParam2Default int32 = 0
const rpcStartBuildGeneric int32 = 1
const rpcEndBuildGeneric int32 = 2

type rpcPlugin struct {
	config                      rpcConfig
	creds                       credentials.TransportCredentials
	clients                     map[string]*rpcClient
	globalExGenericRequest      *pb.ExBuildGenericRequest
	globalExGenericResponse     *pb.ExBuildGenericResponse
	globalExGenericStreamClient pb.ExService_ExBuildGenericStreamClient
}

// newRpcPlugin builds and returns a new RPC plugin
func newRpcPlugin() *rpcPlugin {
	return &rpcPlugin{
		config:                      rpcConfig{},
		creds:                       insecure.NewCredentials(),
		clients:                     make(map[string]*rpcClient),
		globalExGenericRequest:      newExGenericRequest(),
		globalExGenericResponse:     newExGenericResponse(),
		globalExGenericStreamClient: nil,
	}
}

func newExGenericRequest() *pb.ExBuildGenericRequest {
	return &pb.ExBuildGenericRequest{
		Action: "",
		Params: make(map[string]*pb.Value),
	}
}

func newExGenericResponse() *pb.ExBuildGenericResponse {
	return &pb.ExBuildGenericResponse{
		Data: make(map[string]*pb.Value),
	}
}

func (p *rpcPlugin) closeExGenericStreamClient() error {
	if p.globalExGenericStreamClient == nil {
		return nil
	}

	err := p.globalExGenericStreamClient.CloseSend()
	p.globalExGenericStreamClient = nil

	return err
}

func (p *rpcPlugin) buildExBuildGenericStream(client *rpcClient, request *pb.ExBuildGenericRequest) error {
	if err := p.closeExGenericStreamClient(); err != nil {
		return err
	}

	// Create stream client
	ctx := context.Background()
	var err error
	if p.globalExGenericStreamClient, err = client.exServiceClient.ExBuildGenericStream(ctx, request); err != nil {
		client.isValid = false
		log.Errorf("Call to ExBuildGenericStream returned error: %s", err)

		return err
	}

	return nil
}

func (p *rpcPlugin) pullExGenericStreamClient() (*pb.ExBuildGenericResponse, error) {
	if p.globalExGenericStreamClient == nil {
		return nil, errors.New("stream client not available")
	}

	response, err := p.globalExGenericStreamClient.Recv()

	if err != nil {
		// If you are at EOF, then return nothing
		if err.Error() == "EOF" {
			return nil, nil
		} else if s, ok := status.FromError(err); ok && s.Code() == codes.OutOfRange {
			return nil, nil
		}

		return nil, err
	}

	if response == nil {
		return nil, p.closeExGenericStreamClient()
	}

	return response, nil
}

// init initializes the RPC plugin
func (p *rpcPlugin) init() {
	log.Info("Initializing RPC plugin")

	// Add a certificate
	getCredentials := func() {
		if p.config.Auth.PfxFilePath == nil && p.config.Auth.PfxPassword == nil {
			log.Info("Using insecure auth. settings")

			return
		}

		pfxFilePath, pfxPassword := *p.config.Auth.PfxFilePath, ""

		if p.config.Auth.PfxPassword != nil {
			pfxPassword = *p.config.Auth.PfxPassword
		}

		// Load the PFX file
		pfxData, err := os.ReadFile(pfxFilePath)
		if err != nil {
			log.Fatalf("Error reading PFX file: %v", err)
		}

		// Parse the PFX data to get the certificate
		_, cert, err := pkcs12.Decode(pfxData, pfxPassword)
		if err != nil {
			log.Fatalf("Error decoding PFX file: %v", err)
		}

		// Create a new certificate pool and add the certificate
		caCertPool := x509.NewCertPool()
		caCertPool.AddCert(cert)

		// Create a TLS configuration using the parsed certificate
		tlsConfig := &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		}

		// Create a credentials object from the TLS configuration
		p.creds = credentials.NewTLS(tlsConfig)

		log.Info("Using secure auth. settings")
	}
	getCredentials()

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

	// Load the certificate
	var conn *grpc.ClientConn
	var err error
	conn, err = grpc.Dial(url, grpc.WithTransportCredentials(p.creds))

	// Dial with the loaded certificate
	if err != nil {
		log.Errorf("Unable to attach client: %s@%s", name, url)
		p.clients[name] = newRpcClient(name, url)

		return
	}

	// Create gRPC clients with the connection
	p.clients[name] = &rpcClient{
		isValid:             true,
		name:                name,
		url:                 url,
		exServiceClient:     pb.NewExServiceClient(conn),
		nwnxServiceClient:   pb.NewNWNXServiceClient(conn),
		scorcoServiceClient: pb.NewSCORCOServiceClient(conn),
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

// Reset the call action request and response
func (p *rpcPlugin) resetBuildGeneric() {
	p.globalExGenericRequest = newExGenericRequest()
	p.globalExGenericResponse = newExGenericResponse()
}

// getInt the body of the NWNXGetInt() call
func (p *rpcPlugin) getInt(sFunction, sParam1 string, nParam2 int32) int32 {
	log.Info("In NWNXGetInt")

	// ExBuildGeneric()
	switch sFunction {
	case rpcGetInt:
		var value int32 = 0
		v, found := p.globalExGenericResponse.Data[sParam1]
		if found {
			value = v.GetNValue()
		}

		if nParam2 == rpcEndBuildGeneric {
			p.resetBuildGeneric()
		}

		return value
	case rpcGetBool:
		var value int32 = 0
		v, found := p.globalExGenericResponse.Data[sParam1]
		if found {
			if v.GetBValue() {
				value = 1
			}
		}

		if nParam2 == rpcEndBuildGeneric {
			p.resetBuildGeneric()
		}

		return value
	case rpcPullGenericStream:
		p.resetBuildGeneric()
		response, err := p.pullExGenericStreamClient()
		p.globalExGenericResponse = response

		if response == nil || err != nil {
			if err != nil {
				log.Errorf("Response stream to ExBuildGenericStream returned error: %s", err)
			}

			return 0
		}

		return 1
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
	log.Info("In NWNXSetInt")

	// ExBuildGeneric()
	switch sFunction {
	case rpcSetInt:
		if nParam2 == rpcStartBuildGeneric {
			p.resetBuildGeneric()
		}

		p.globalExGenericRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_NValue{
				NValue: nValue,
			},
		}

		return
	case rpcSetBool:
		if nParam2 == rpcStartBuildGeneric {
			p.resetBuildGeneric()
		}

		p.globalExGenericRequest.Params[sParam1] = &pb.Value{
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
	log.Info("In NWNXGetFloat")

	// ExBuildGeneric()
	switch sFunction {
	case rpcGetFloat:
		var value float32 = 0.0
		v, found := p.globalExGenericResponse.Data[sParam1]
		if found {
			value = v.GetFValue()
		}

		if nParam2 == rpcEndBuildGeneric {
			p.resetBuildGeneric()
		}

		return value
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
	log.Info("In NWNXSetFloat")

	// ExBuildGeneric()
	switch sFunction {
	case rpcSetFloat:
		if nParam2 == rpcStartBuildGeneric {
			p.resetBuildGeneric()
		}

		p.globalExGenericRequest.Params[sParam1] = &pb.Value{
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
	log.Info("In NWNXGetString")

	// ExBuildGeneric()
	switch sFunction {
	case rpcGetString:
		var value string = ""
		v, found := p.globalExGenericResponse.Data[sParam1]
		if found {
			value = v.GetSValue()
		}

		if nParam2 == rpcEndBuildGeneric {
			p.resetBuildGeneric()
		}

		return value
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
	log.Info("In NWNXSetString")

	// ExBuildGeneric()
	switch sFunction {
	case rpcSetString:
		if nParam2 == rpcStartBuildGeneric {
			p.resetBuildGeneric()
		}

		p.globalExGenericRequest.Params[sParam1] = &pb.Value{
			ValueType: &pb.Value_SValue{
				SValue: sValue,
			},
		}

		return
	case rpcResetGeneric:
		p.resetBuildGeneric()

		return
	case rpcBuildGeneric:
		// sParam1_ holds the client identifier
		client, ok := p.getRpcClient(sParam1)
		if !ok {
			return
		}

		// sValue_ contains the action
		p.globalExGenericRequest.Action = sValue
		response, err := client.buildGeneric(p.globalExGenericRequest, p.config.getTimeout())
		p.resetBuildGeneric()

		if err == nil {
			*p.globalExGenericResponse = *response
		}

		return
	case rpcBuildGenericStream:
		// sParam1_ holds the client identifier
		client, ok := p.getRpcClient(sParam1)
		if !ok {
			return
		}

		// sValue_ contains the action
		p.globalExGenericRequest.Action = sValue
		if err := p.buildExBuildGenericStream(client, p.globalExGenericRequest); err != nil {
			log.Errorf("Call to ExBuildGenericStream returned error: %s", err)
		}
		p.resetBuildGeneric()

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
	log.Info("In SCORCOGetGFFSize")

	// ExBuildGeneric()
	switch sFunction {
	case rpcGetGff:
		if v, found := p.globalExGenericResponse.Data[sVarName]; found {
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
	log.Info("In SCORCOGetGFF")

	// ExBuildGeneric()
	switch sFunction {
	case rpcGetGff:
		var value []byte = nil
		v, found := p.globalExGenericResponse.Data[sVarName]
		if found {
			value = v.GetGffValue()
		}

		if nParam2 == rpcEndBuildGeneric {
			p.resetBuildGeneric()
		}

		return value
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
	log.Info("In SCORCOSetGFF")

	// ExBuildGeneric()
	switch sFunction {
	case rpcSetGff:
		if nParam2 == rpcStartBuildGeneric {
			p.resetBuildGeneric()
		}

		p.globalExGenericRequest.Params[sVarName] = &pb.Value{
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
