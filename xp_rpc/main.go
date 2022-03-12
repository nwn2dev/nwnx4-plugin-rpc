package main

// #include <stdlib.h>
// #include <stdio.h>
// #include <string.h>
import "C"

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"path"
	"time"
	"unsafe"
)

import (
	pbCore "nwnx4.org/xp_rpc/proto"
	pbNWScript "nwnx4.org/xp_rpc/proto/nwscript"
)

const PluginName string = "RPC"
const PluginVersion string = "0.2.1"

type Config struct {
	Server  *ServerConfig
	Clients map[string]string
}

type ServerConfig struct {
	Url      string
	Services ServicesConfig
}

type ServicesConfig struct {
	Logger bool
}

type rpcPlugin struct {
	header      string
	description string
	subClass    string
	version     string
	rpcServer   *rpcServer
	rpcClients  map[string]rpcClient
}

type rpcServer struct {
	pbCore.UnimplementedLogServiceServer
}

func (s *rpcServer) Trace(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Trace(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

func (s *rpcServer) Debug(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Debug(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

func (s *rpcServer) Info(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Info(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

func (s *rpcServer) Warn(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Warn(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

func (s *rpcServer) Err(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Error(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

func (s *rpcServer) LogStr(ctx context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Printf(stringValue.Value)

	return &empty.Empty{}, status.New(codes.OK, "").Err()
}

type rpcClient struct {
	name                 string
	url                  string
	nwnxServiceClient    pbNWScript.NWNXServiceClient
	messageServiceClient pbCore.MessageServiceClient
}

func setupRpcPlugin() {
	plugin = rpcPlugin{
		header: fmt.Sprintf(
			"NWNX4 %s Plugin %s \n"+
				"(c) 2021-2022 by ihatemundays (scottmunday84@gmail.com) \n", PluginName, PluginVersion),
		description: "A better way to build service-oriented applications in NWN2",
		subClass:    PluginName,
		version:     PluginVersion,
		rpcServer:   nil,
		rpcClients:  make(map[string]rpcClient),
	}
}

func (p rpcPlugin) setupRpcServer(serverConfig *ServerConfig) {
	if serverConfig == nil {
		log.Info("Skipping server setup")

		return
	}

	// Build server
	listen, err := net.Listen("tcp", serverConfig.Url)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen: %v", err))

		return
	}

	log.Info("Adding server")
	s := grpc.NewServer()
	p.rpcServer = &rpcServer{}

	// Setup logger (if toggled)
	if serverConfig.Services.Logger {
		log.Info("Adding logger service to server")
		pbCore.RegisterLogServiceServer(s, p.rpcServer)
	}

	// Setup the server in an asynchronous goroutine
	serve := func(s *grpc.Server, listen net.Listener) {
		log.Info("Serving server")

		if err := s.Serve(listen); err != nil {
			log.Error(fmt.Sprintf("Could not serve server at %s", serverConfig.Url))
		}

		log.Info("Serve is closed")
	}
	go serve(s, listen)
}

func (p rpcPlugin) addRpcClient(name, url string) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Error(fmt.Sprintf("Unable to login to %s", url))

		return
	}

	p.rpcClients[name] = rpcClient{
		name:                 name,
		url:                  url,
		nwnxServiceClient:    pbNWScript.NewNWNXServiceClient(conn),
		messageServiceClient: pbCore.NewMessageServiceClient(conn),
	}

	log.Info(fmt.Sprintf("Established connection and stubs for %s@%s", name, url))
}

func (p rpcPlugin) getRpcClient(name string) *rpcClient {
	rpcClient, exists := p.rpcClients[name]
	if !exists {
		log.Error(fmt.Sprintf("Client not declared: %s", name))

		return nil
	}

	return &rpcClient
}

//export IsProtoPlugin
func IsProtoPlugin() C.char {
	return 1
}

//export GetHeaderDescriptor
func GetHeaderDescriptor() *C.char {
	return C.CString(plugin.header)
}

//export GetDescriptionDescriptor
func GetDescriptionDescriptor() *C.char {
	return C.CString(plugin.description)
}

//export GetSubClassDescriptor
func GetSubClassDescriptor() *C.char {
	return C.CString(plugin.subClass)
}

//export GetVersionDescriptor
func GetVersionDescriptor() *C.char {
	return C.CString(plugin.version)
}

//export Init
func Init(nwnxHome *C.char) C.char {
	// Setup the log file
	nwnxHome_ := C.GoString(nwnxHome)
	logFile, err := os.OpenFile(path.Join(nwnxHome_, "xp_rpc.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(logFile)
	log.Info(plugin.header)

	// Get YAML file with services
	configFile, err2 := ioutil.ReadFile(path.Join(nwnxHome_, "xp_rpc.yml"))
	if err2 != nil {
		log.Error(err2)

		return 0
	}
	config := Config{}
	err3 := yaml.Unmarshal(configFile, &config)
	if err3 != nil {
		log.Error(err3)

		return 0
	}

	log.Info("Processing configuration file")

	// Build out the server
	plugin.setupRpcServer(config.Server)

	// Build out the clients
	for name, url := range config.Clients {
		log.Info(fmt.Sprintf("Adding client %s@%s", name, url))
		plugin.addRpcClient(name, url)
	}

	log.Info("Initialized plugin")

	return 1
}

//export GetFunctionClass
func GetFunctionClass(fClass *C.char) {
	pluginName := C.CString(PluginName)
	C.strcpy(fClass, pluginName)
	defer C.free(unsafe.Pointer(pluginName))
}

//export GetInt
func GetInt(sFunction, sParam1 *C.char, nParam2 C.int) C.int {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXGetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err2 := client.nwnxServiceClient.NWNXGetInt(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to GetInt returned error: %s, %s, %d",
			request.SFunction, request.SParam1, request.NParam2))

		return 0
	}

	return C.int(response.Value)
}

//export SetInt
func SetInt(sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	nValue_ := int32(nValue)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXSetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		NValue:    nValue_,
	}
	_, err2 := client.nwnxServiceClient.NWNXSetInt(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to SetInt returned error: %s, %s, %d, %d",
			request.SFunction, request.SParam1, request.NParam2, request.NValue))
	}
}

//export GetFloat
func GetFloat(sFunction, sParam1 *C.char, nParam2 C.int) C.float {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return 0.0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXGetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err2 := client.nwnxServiceClient.NWNXGetFloat(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to GetFloat returned error: %s, %s, %d",
			request.SFunction, request.SParam1, request.NParam2))

		return 0.0
	}

	return C.float(response.Value)
}

//export SetFloat
func SetFloat(sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	fValue_ := float32(fValue)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXSetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		FValue:    fValue_,
	}
	_, err2 := client.nwnxServiceClient.NWNXSetFloat(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to SetFloat returned error: %s, %s, %d, %f",
			request.SFunction, request.SParam1, request.NParam2, request.FValue))
	}
}

//export GetString
func GetString(sFunction, sParam1 *C.char, nParam2 C.int) *C.char {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXGetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err2 := client.nwnxServiceClient.NWNXGetString(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to GetString returned error: %s, %s, %d",
			request.SFunction, request.SParam1, request.NParam2))

		return nil
	}

	return C.CString(response.Value)
}

//export SetString
func SetString(sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	sValue_ := C.GoString(sValue)
	client := plugin.getRpcClient(sFunction_)
	if client == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	request := pbNWScript.NWNXSetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		SValue:    sValue_,
	}
	_, err2 := client.nwnxServiceClient.NWNXSetString(ctx, &request)
	if err2 != nil {
		log.Error(fmt.Sprintf("Call to SetString returned error: %s, %s, %d, %s",
			request.SFunction, request.SParam1, request.NParam2, request.SValue))
	}
}

var plugin rpcPlugin

func main() {
	setupRpcPlugin()
}
