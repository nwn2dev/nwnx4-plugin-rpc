package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdint.h>

typedef struct {
	const char* dllPath;
	const char* nwnInstallHomePath;
	const char* nwnxHomePath;
} CPluginInitInfo;
*/
import "C"

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	empty "google.golang.org/protobuf/types/known/emptypb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net"
	"os"
	"path"
	"reflect"
	"time"
	// Protobuf
	pbCore "nwnx4.org/xp_rpc/proto"
	pbNWScript "nwnx4.org/xp_rpc/proto/nwscript"
)

const pluginName string = "RPC"           // Plugin name passed to hook
const pluginVersion string = "0.2.3"      // Plugin version passed to hook
const pluginInstance string = "Singleton" // Plugin instance passed to hook

// YAML configuration for xp_rpc
type Config struct {
	Server  *ServerConfig
	Clients map[string]string
}

// YAML server configuration for xp_rpc
type ServerConfig struct {
	Url      string
	Services ServerServicesConfig
}

// YAML server services configuration for xp_rpc
type ServerServicesConfig struct {
	Logger bool
}

// Core plugin class; singleton per DLL
type rpcPlugin struct {
	rpcServer  *rpcServer
	rpcClients map[string]rpcClient
}

// initRpcServer initializes the RPC server
// Runs on an rpcPlugin and accepts a ServerConfig
func (p *rpcPlugin) initRpcServer(serverConfig *ServerConfig) {
	if serverConfig == nil {
		log.Info("No server configuration; skipping")

		return
	}

	// Build server
	listen, err := net.Listen("tcp", serverConfig.Url)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to listen at server URL: %v", err))

		return
	}

	log.Info(fmt.Sprintf("Adding server: @%s", serverConfig.Url))
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
			log.Error(fmt.Sprintf("Could not serve server: @%s", serverConfig.Url))
		}

		log.Info("Server is closed")
	}
	go serve(s, listen)
}

// addRpcClient adds an RPC client
// Runs on an rpcPlugin and adds a client by name and URL
func (p *rpcPlugin) addRpcClient(name, url string) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Sprintf("Unable to attach client: @%s", url))

		return
	}

	p.rpcClients[name] = rpcClient{
		name:                 name,
		url:                  url,
		nwnxServiceClient:    pbNWScript.NewNWNXServiceClient(conn),
		messageServiceClient: pbCore.NewMessageServiceClient(conn),
	}

	log.Info(fmt.Sprintf("Established client connection and stubs: %s@%s", name, url))
}

// getRpcClient will get an rpcClient by name
func (p *rpcPlugin) getRpcClient(name string) (*rpcClient, bool) {
	rpcClient, exists := p.rpcClients[name]
	if !exists {
		log.Error(fmt.Sprintf("Client not declared: %s", name))

		return nil, false
	}

	return &rpcClient, true
}

// getInt the body of the NWNXGetInt() call
func (p *rpcPlugin) getInt(sFunction, sParam1 *C.char, nParam2 C.int) C.int {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXGetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetInt(ctx, &request)
	if err != nil {
		log.Error(fmt.Sprintf("Call to GetInt returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2))

		return 0
	}

	return C.int(response.Value)
}

// setInt the body of the NWNXSetInt() call
func (p *rpcPlugin) setInt(sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	nValue_ := int32(nValue)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXSetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		NValue:    nValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetInt(ctx, &request); err != nil {
		log.Error(fmt.Sprintf("Call to SetInt returned error: %s; %s, %s, %d, %d",
			err, request.SFunction, request.SParam1, request.NParam2, request.NValue))
	}
}

// getFloat the body of the NWNXGetFloat() call
func (p *rpcPlugin) getFloat(sFunction, sParam1 *C.char, nParam2 C.int) C.float {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0.0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXGetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetFloat(ctx, &request)
	if err != nil {
		log.Error(fmt.Sprintf("Call to GetFloat returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2))

		return 0.0
	}

	return C.float(response.Value)
}

func (p *rpcPlugin) setFloat(sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	fValue_ := float32(fValue)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXSetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		FValue:    fValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetFloat(ctx, &request); err != nil {
		log.Error(fmt.Sprintf("Call to SetFloat returned error: %s; %s, %s, %d, %f",
			err, request.SFunction, request.SParam1, request.NParam2, request.FValue))
	}
}

// getString the body of the NWNXGetString() call
func (p *rpcPlugin) getString(sFunction, sParam1 *C.char, nParam2 C.int) *C.char {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXGetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetString(ctx, &request)
	if err != nil {
		log.Error(fmt.Sprintf("Call to GetString returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2))

		return nil
	}

	return C.CString(response.Value)
}

// setString the body of the NWNXSetString() call
func (p *rpcPlugin) setString(sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	sValue_ := C.GoString(sValue)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	request := pbNWScript.NWNXSetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		SValue:    sValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetString(ctx, &request); err != nil {
		log.Error(fmt.Sprintf("Call to SetString returned error: %s; %s, %s, %d, %s",
			err, request.SFunction, request.SParam1, request.NParam2, request.SValue))
	}
}

// rpcServer contains the interfaces to the RPC server
type rpcServer struct {
	pbCore.UnimplementedLogServiceServer
}

// Trace is the method call equivalent on the logger
func (s *rpcServer) Trace(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Trace(stringValue.Value)

	return &empty.Empty{}, nil
}

// Debug is the method call equivalent on the logger
func (s *rpcServer) Debug(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Debug(stringValue.Value)

	return &empty.Empty{}, nil
}

// Info is the method call equivalent on the logger
func (s *rpcServer) Info(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Info(stringValue.Value)

	return &empty.Empty{}, nil
}

// Warn is the method call equivalent on the logger
func (s *rpcServer) Warn(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Warn(stringValue.Value)

	return &empty.Empty{}, nil
}

// Err is the method call equivalent on the logger
func (s *rpcServer) Err(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Error(stringValue.Value)

	return &empty.Empty{}, nil
}

// LogStr is the method call equivalent on the logger without a log level
func (s *rpcServer) LogStr(_ context.Context, stringValue *wrappers.StringValue) (*empty.Empty, error) {
	log.Printf(stringValue.Value)

	return &empty.Empty{}, nil
}

// rpcClient contains the clients to RPC microservices
type rpcClient struct {
	name                 string
	url                  string
	nwnxServiceClient    pbNWScript.NWNXServiceClient
	messageServiceClient pbCore.MessageServiceClient
}

// newRpcPlugin builds and returns a new RPC plugin
func newRpcPlugin() *rpcPlugin {
	return &rpcPlugin{
		rpcServer:  nil,
		rpcClients: make(map[string]rpcClient),
	}
}

var plugin *rpcPlugin // Singleton

// All exports to C library

//export NWNXCPlugin_GetAbiVersion
func NWNXCPlugin_GetAbiVersion() C.uint32_t {
	return 1
}

//export NWNXCPlugin_GetPluginName
func NWNXCPlugin_GetPluginName() *C.char {
	return C.CString(pluginName)
}

//export NWNXCPlugin_GetPluginVersion
func NWNXCPlugin_GetPluginVersion() *C.char {
	return C.CString(pluginVersion)
}

//export NWNXCPlugin_GetPluginInstance
func NWNXCPlugin_GetPluginInstance() *C.char {
	return C.CString(pluginInstance)
}

//export NWNXCPlugin_New
func NWNXCPlugin_New(initInfo *C.CPluginInitInfo) C.uint32_t {
	plugin = newRpcPlugin()

	// Setup the log file
	nwnxHomePath_ := C.GoString(initInfo.nwnxHomePath)
	logFile, err := os.OpenFile(path.Join(nwnxHomePath_, "xp_rpc.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0
	}

	// Adding the header/description to the log
	header := fmt.Sprintf(
		"NWNX4 %s Plugin %s \n"+
			"(c) 2021-2022 by ihatemundays (scottmunday84@gmail.com) \n", pluginName, pluginVersion)
	description := "A better way to build service-oriented applications in NWN2"

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(logFile)
	log.Info(header)
	log.Info(description)

	// Get YAML file with services
	configFile, err2 := ioutil.ReadFile(path.Join(nwnxHomePath_, "xp_rpc.yml"))
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

	// Initialize the server
	plugin.initRpcServer(config.Server)

	// Build out the clients
	for name, url := range config.Clients {
		log.Info(fmt.Sprintf("Adding client %s@%s", name, url))
		plugin.addRpcClient(name, url)
	}

	log.Info("Initialized plugin")

	return C.uint32_t(reflect.ValueOf(plugin).Pointer())
}

//export NWNXCPlugin_Delete
func NWNXCPlugin_Delete(ptr *C.void) C.char {
	// No pointer, then can't simulate a delete
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return 0
	}

	plugin = nil // If you got this far, the plugin should be the one stored on the DLL

	return 1
}

//export NWNXCPlugin_GetInt
func NWNXCPlugin_GetInt(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.int {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return 0
	}

	return plugin.getInt(sFunction, sParam1, nParam2)
}

//export NWNXCPlugin_SetInt
func NWNXCPlugin_SetInt(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return
	}

	plugin.setInt(sFunction, sParam1, nParam2, nValue)
}

//export NWNXCPlugin_GetFloat
func NWNXCPlugin_GetFloat(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.float {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return 0.0
	}

	return plugin.getFloat(sFunction, sParam1, nParam2)
}

//export NWNXCPlugin_SetFloat
func NWNXCPlugin_SetFloat(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return
	}

	plugin.setFloat(sFunction, sParam1, nParam2, fValue)
}

//export NWNXCPlugin_GetString
func NWNXCPlugin_GetString(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int) *C.char {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return nil
	}

	return plugin.getString(sFunction, sParam1, nParam2)
}

//export NWNXCPlugin_SetString
func NWNXCPlugin_SetString(ptr *C.void, sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	_, err := getRpcPlugin(ptr)
	if err != nil {
		log.Error(err)

		return
	}

	plugin.setString(sFunction, sParam1, nParam2, sValue)
}

// Get the RPC plugin from the pointer
func getRpcPlugin(ptr interface{}) (*rpcPlugin, error) {
	plugin, ok := ptr.(*rpcPlugin)
	if !ok {
		return nil, errors.New("plugin missing")
	}

	return plugin, nil
}

func main() {}
