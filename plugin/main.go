package main

/*
#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>

typedef struct {
	const char* dll_path;
	const char* nwnx_user_path;
	const char* nwn2_install_path;
	const char* nwn2_home_path;
	const char* nwn2_module_path;
	const char* nwnx_install_path;
} CPluginInitInfo;
*/
import "C"

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	pb "proto"
	"strings"
	"time"
	"unsafe"
)

const pluginName string = "NWNX RPC Plugin" // Plugin name passed to hook
const pluginVersion string = "0.3.1"        // Plugin version passed to hook
const pluginID string = "RPC"               // Plugin ID used for identification in the list

const rpcGffVarNameSeparator = "/"

/*const (
	rpcValueTypeInt = iota
	rpcValueTypeBool
	rpcValueTypeFloat
	rpcValueTypeString
	rpcValueTypeGff
)*/

type rpcConfig struct {
	host      *rpcHostConfig
	clients   map[string]string
	perClient *rpcPerClientConfig
}

type rpcPerClientConfig struct {
	retries int
	delay   int
	timeout int
}

func (p *rpcPerClientConfig) getDelay() time.Duration {
	return time.Second * time.Duration(p.delay)
}

func (p *rpcPerClientConfig) getTimeout() time.Duration {
	return time.Second * time.Duration(p.timeout)
}

type rpcHostConfig struct {
	log *rpcLogConfig
}

type rpcLogConfig struct {
	logLevel string
}

type rpcPlugin struct {
	config  *rpcConfig
	clients map[string]rpcClient
}

// initRpcServer initializes the RPC server
// Runs on an rpcPlugin and accepts a ServerConfig
func (p *rpcPlugin) initRpcServer(hostConfig *rpcHostConfig) {
	if hostConfig == nil {
		log.Info("No server configuration; skipping")

		return
	}

	if hostConfig.log != nil {
		// Set the log level based on what was passed if it matches a level
		for _, logLevel := range log.AllLevels {
			if strings.EqualFold(logLevel.String(), hostConfig.log.logLevel) {
				log.SetLevel(logLevel)
				break
			}
		}
	}
}

// addRpcClient adds an RPC client
// Runs on an rpcPlugin and adds a client by name and URL
func (p *rpcPlugin) addRpcClient(name, url string) {
	conn, err := grpc.Dial(url, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error(fmt.Sprintf("Unable to attach client: @%s", url))

		return
	}

	p.clients[name] = rpcClient{
		isValid:             true,
		name:                name,
		url:                 url,
		nwnxServiceClient:   pb.NewNWNXServiceClient(conn),
		scorcoServiceClient: pb.NewSCORCOServiceClient(conn),
		actionServiceClient: pb.NewActionServiceClient(conn),
		callActionRequest:   nil,
		callActionResponse:  nil,
		config:              p.config,
	}

	log.Info(fmt.Sprintf("Established client connection and stubs: %s@%s", name, url))
}

// getRpcClient will get an rpcClient by name
func (p *rpcPlugin) getRpcClient(name string) (*rpcClient, bool) {
	rpcClient, exists := p.clients[name]

	if !exists {
		log.Error(fmt.Sprintf("Client not declared: %s", name))

		return nil, false
	}

	// Handle invalid clients; try to recreate
	if !rpcClient.isValid {
		url := p.config.clients[name]

		// Invalid client; try to recreate
		for i := 0; i < p.config.perClient.retries; i++ {
			p.addRpcClient(name, url)
			rpcClient, exists = p.clients[name]

			if exists {
				return &rpcClient, true
			}

			time.Sleep(p.config.perClient.getDelay())
		}

		log.Error(fmt.Sprintf("Client is still invalid; could not be setup: %s@#%s", name, url))
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXGetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetInt(ctx, &request)
	if err != nil {
		client.isValid = false
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXSetIntRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		NValue:    nValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetInt(ctx, &request); err != nil {
		client.isValid = false
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXGetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetFloat(ctx, &request)
	if err != nil {
		client.isValid = false
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXSetFloatRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		FValue:    fValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetFloat(ctx, &request); err != nil {
		client.isValid = false
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
		return C.CString("")
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXGetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
	}
	response, err := client.nwnxServiceClient.NWNXGetString(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Error(fmt.Sprintf("Call to GetString returned error: %s; %s, %s, %d",
			err, request.SFunction, request.SParam1, request.NParam2))

		return C.CString("")
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

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.NWNXSetStringRequest{
		SFunction: sFunction_,
		SParam1:   sParam1_,
		NParam2:   nParam2_,
		SValue:    sValue_,
	}
	if _, err := client.nwnxServiceClient.NWNXSetString(ctx, &request); err != nil {
		client.isValid = false
		log.Error(fmt.Sprintf("Call to SetString returned error: %s; %s, %s, %d, %s",
			err, request.SFunction, request.SParam1, request.NParam2, request.SValue))
	}
}

func (p *rpcPlugin) getGffSize(sVarName *C.char) C.size_t {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	sVarName_, sFunction_ := splits[0], splits[1]
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.SCORCOGetGFFSizeRequest{
		SVarName: sVarName_,
	}
	response, err := client.scorcoServiceClient.SCORCOGetGFFSize(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Error(fmt.Sprintf("Call to GetGFFSize returned error: %s; %s",
			err, request.SVarName))

		return 0
	}

	return C.size_t(response.Size)
}

func (p *rpcPlugin) getGff(sVarName *C.char, result *C.uint8_t, resultSize C.size_t) {
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	sVarName_, sFunction_ := splits[0], splits[1]
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	request := pb.SCORCOGetGFFRequest{
		SVarName: sVarName_,
	}
	response, err := client.scorcoServiceClient.SCORCOGetGFF(ctx, &request)
	if err != nil {
		client.isValid = false
		log.Error(fmt.Sprintf("Call to GetGFF returned error: %s; %s",
			err, request.SVarName))

		return
	}

	// Do not need to free this memory; managed by the hook library
	C.memcpy(unsafe.Pointer(result), unsafe.Pointer(&response.GffData[0]), resultSize)
}

func (p *rpcPlugin) setGff(sVarName *C.char, gffData *C.uint8_t, gffDataSize C.size_t) {
	gffDataSize_ := uint32(gffDataSize)
	ptr := unsafe.Pointer(gffData)
	gffData_ := (*[1 << 30]byte)(ptr)[:gffDataSize_:gffDataSize_]
	splits := strings.SplitN(C.GoString(sVarName), rpcGffVarNameSeparator, 2)
	sVarName_, sFunction_ := splits[0], splits[1]
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		C.free(ptr)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), p.config.perClient.getTimeout())
	defer cancel()
	var request = pb.SCORCOSetGFFRequest{
		SVarName:    sVarName_,
		GffData:     gffData_,
		GffDataSize: gffDataSize_,
	}
	if _, err := client.scorcoServiceClient.SCORCOSetGFF(ctx, &request); err != nil {
		client.isValid = false
		log.Error(fmt.Sprintf("Call to SetGFF returned error: %s; %s",
			err, request.SVarName))
	}
	C.free(ptr)
}

// rpcClient contains the clients to RPC microservices
type rpcClient struct {
	isValid             bool
	name                string
	url                 string
	nwnxServiceClient   pb.NWNXServiceClient
	scorcoServiceClient pb.SCORCOServiceClient
	actionServiceClient pb.ActionServiceClient
	callActionRequest   *pb.CallActionRequest
	callActionResponse  *pb.CallActionResponse
	config              *rpcConfig
}

func (c *rpcClient) resetCall() {
	c.callActionRequest = &pb.CallActionRequest{}
	c.callActionResponse = nil
}

func (c *rpcClient) callAction() bool {
	if c.callActionRequest == nil {
		c.resetCall()
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.config.perClient.getTimeout())
	defer cancel()

	response, err := c.actionServiceClient.CallAction(ctx, c.callActionRequest)

	if err != nil {
		c.isValid = false
		log.Error(fmt.Sprintf("Error sending request: %s", err))

		return false
	}

	c.callActionRequest = nil
	c.callActionResponse = response

	return true
}

// newRpcPlugin builds and returns a new RPC plugin
func newRpcPlugin() *rpcPlugin {
	return &rpcPlugin{
		config:  &rpcConfig{},
		clients: make(map[string]rpcClient),
	}
}

var plugin *rpcPlugin // Singleton

// All exports to C library

//export NWNXCPlugin_GetID
func NWNXCPlugin_GetID(_ *C.void) *C.char {
	return C.CString(pluginID)
}

//export NWNXCPlugin_GetInfo
func NWNXCPlugin_GetInfo() *C.char {
	return C.CString("NWNX4 RPC - A better way to integrate services with NWN2")
}

//export NWNXCPlugin_GetVersion
func NWNXCPlugin_GetVersion() *C.char {
	return C.CString(pluginVersion)
}

//export NWNXCPlugin_New
func NWNXCPlugin_New(initInfo C.CPluginInitInfo) C.uint32_t {
	plugin = newRpcPlugin()

	// Setup the log file
	nwnxHomePath_ := C.GoString(initInfo.nwnx_user_path)
	logFile, err := os.OpenFile(path.Join(nwnxHomePath_, "xp_rpc.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0
	}

	// Adding the header/description to the log
	header := fmt.Sprintf(
		"%s %s \n"+
			"(c) 2021-2023 by ihatemundays (scottmunday84@gmail.com) \n", pluginName, pluginVersion)
	description := "A better way to integrate services with NWN2"

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
	var config = plugin.config
	err3 := yaml.Unmarshal(configFile, config)
	if err3 != nil {
		log.Error(err3)

		return 0
	}

	log.Info("Processing configuration file")

	// Initialize the server
	plugin.initRpcServer(config.host)

	// Build out the clients
	for name, url := range config.clients {
		log.Info(fmt.Sprintf("Adding client: %s@%s", name, url))
		plugin.addRpcClient(name, url)
	}

	log.Info("Initialized plugin")

	// Giving back address
	return C.uint32_t(uintptr(unsafe.Pointer(&plugin)))
}

//export NWNXCPlugin_Delete
func NWNXCPlugin_Delete(_ C.uint32_t) {}

//export NWNXCPlugin_GetInt
func NWNXCPlugin_GetInt(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.int {
	return plugin.getInt(sFunction, sParam1, nParam2)
}

//export NWNXCPlugin_SetInt
func NWNXCPlugin_SetInt(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	plugin.setInt(sFunction, sParam1, nParam2, nValue)
}

//export NWNXCPlugin_GetFloat
func NWNXCPlugin_GetFloat(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.float {
	return plugin.getFloat(sFunction, sParam1, nParam2)
}

//export NWNXCPlugin_SetFloat
func NWNXCPlugin_SetFloat(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	plugin.setFloat(sFunction, sParam1, nParam2, fValue)
}

//export NWNXCPlugin_GetString
func NWNXCPlugin_GetString(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, result *C.char, resultSize C.size_t) {
	response := plugin.getString(sFunction, sParam1, nParam2)
	C.strncpy_s(result, resultSize, response, C.strlen(response))
}

//export NWNXCPlugin_SetString
func NWNXCPlugin_SetString(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	plugin.setString(sFunction, sParam1, nParam2, sValue)
}

//export NWNXCPlugin_GetGFFSize
func NWNXCPlugin_GetGFFSize(_ *C.void, sVarName *C.char) C.size_t {
	return plugin.getGffSize(sVarName)
}

//export NWNXCPlugin_GetGFF
func NWNXCPlugin_GetGFF(_ *C.void, sVarName *C.char, result *C.uint8_t, resultSize C.size_t) {
	plugin.getGff(sVarName, result, resultSize)
}

//export NWNXCPlugin_SetGFF
func NWNXCPlugin_SetGFF(_ *C.void, sVarName *C.char, gffData *C.uint8_t, gffDataSize C.size_t) {
	plugin.setGff(sVarName, gffData, gffDataSize)
}

func main() {}
