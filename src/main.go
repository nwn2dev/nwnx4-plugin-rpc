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
	pbCore "nwnx4.org/src/proto"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

const pluginName string = "NWNX RPC Plugin" // Plugin name passed to hook
const pluginVersion string = "0.3.0"        // Plugin version passed to hook
const pluginID string = "RPC"               // Plugin ID used for identification in the list

type config struct {
	server  *serverConfig
	clients map[string]string
}

type serverConfig struct {
	log *serverLogConfig
}

type serverLogConfig struct {
	logLevel string
}

type rpcPlugin struct {
	rpcClients map[string]rpcClient
}

// initRpcServer initializes the RPC server
// Runs on an rpcPlugin and accepts a ServerConfig
func (p *rpcPlugin) initRpcServer(serverConfig *serverConfig) {
	if serverConfig == nil {
		log.Info("No server configuration; skipping")

		return
	}

	if serverConfig.log != nil {
		// Set the log level based on what was passed if it matches a level
		for _, logLevel := range log.AllLevels {
			if strings.EqualFold(logLevel.String(), serverConfig.log.logLevel) {
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

	p.rpcClients[name] = rpcClient{
		name:               name,
		url:                url,
		eventServiceClient: pbCore.NewEventServiceClient(conn),
		sendRequest:        nil,
		sendResponse:       nil,
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

	if client.sendResponse == nil && !client.send() {
		return 0
	}

	value, exists := client.sendResponse.Data[sParam1_]
	if !exists {
		log.Warn(fmt.Sprintf("Value not declared in response: %s", sParam1_))

		return 0
	}

	switch nParam2_ {
	case 0:
		return C.int(value.GetNValue())
	case 1:
		if value.GetBValue() {
			return 1
		}

		return 0
	}

	return C.int(value.GetNValue())
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

	if client.sendRequest == nil {
		client.resetSend()
	}

	switch nParam2_ {
	case 0:
		client.sendRequest.Params[sParam1_] = &pbCore.Value{Value: &pbCore.Value_NValue{NValue: nValue_}}
		break
	case 1:
		client.sendRequest.Params[sParam1_] = &pbCore.Value{Value: &pbCore.Value_BValue{BValue: !(nValue_ == 0)}}
		break
	}
}

// getFloat the body of the NWNXGetFloat() call
func (p *rpcPlugin) getFloat(sFunction, sParam1 *C.char, _ C.int) C.float {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	// nParam2_ := int32(nParam2)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return 0.0
	}

	if client.sendResponse == nil && !client.send() {
		return 0
	}

	value, exists := client.sendResponse.Data[sParam1_]
	if !exists {
		log.Warn(fmt.Sprintf("Value not declared in response: %s", sParam1_))

		return 0
	}

	return C.float(value.GetFValue())
}

func (p *rpcPlugin) setFloat(sFunction, sParam1 *C.char, _ C.int, fValue C.float) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	// nParam2_ := int32(nParam2)
	fValue_ := float32(fValue)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	if client.sendRequest == nil {
		client.resetSend()
	}

	client.sendRequest.Params[sParam1_] = &pbCore.Value{Value: &pbCore.Value_FValue{FValue: fValue_}}
}

// getString the body of the NWNXGetString() call
func (p *rpcPlugin) getString(sFunction, sParam1 *C.char, _ C.int) *C.char {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	// nParam2_ := int32(nParam2)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return C.CString("")
	}

	if client.sendResponse == nil && !client.send() {
		return C.CString("")
	}

	value, exists := client.sendResponse.Data[sParam1_]
	if !exists {
		log.Warn(fmt.Sprintf("Value not declared in response: %s", sParam1_))

		return C.CString("")
	}

	return C.CString(value.GetSValue())
}

// setString the body of the NWNXSetString() call
func (p *rpcPlugin) setString(sFunction, sParam1 *C.char, _ C.int, sValue *C.char) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	// nParam2_ := int32(nParam2)
	sValue_ := C.GoString(sValue)
	client, ok := p.getRpcClient(sFunction_)
	if !ok {
		return
	}

	if client.sendRequest == nil {
		client.resetSend()
	}

	client.sendRequest.Params[sParam1_] = &pbCore.Value{Value: &pbCore.Value_SValue{SValue: sValue_}}
}

func (p *rpcPlugin) getGffSize(sVarName *C.char) C.size_t {
	sVarName_ := C.GoString(sVarName)
	varNameSplits := strings.SplitN(sVarName_, "///", 2)
	var clientKey string
	if len(varNameSplits) == 2 {
		clientKey = varNameSplits[0]
		sVarName_ = varNameSplits[1]
	} else {
		return 0
	}
	client, ok := p.getRpcClient(clientKey)
	if !ok {
		return 0
	}

	if client.sendResponse == nil && !client.send() {
		return 0
	}

	value, exists := client.sendResponse.Data[sVarName_]
	if !exists {
		log.Warn(fmt.Sprintf("Value not declared in response: %s", sVarName_))

		return 0
	}

	return C.size_t(len(value.GetGffValue()))
}

func (p *rpcPlugin) getGff(sVarName *C.char, result *C.uint8_t, resultSize C.size_t) {
	sVarName_ := C.GoString(sVarName)
	varNameSplits := strings.SplitN(sVarName_, "///", 2)
	var clientKey string
	if len(varNameSplits) == 2 {
		clientKey = varNameSplits[0]
		sVarName_ = varNameSplits[1]
	} else {
		return
	}
	client, ok := p.getRpcClient(clientKey)
	if !ok {
		return
	}

	if client.sendResponse == nil && !client.send() {
		return
	}

	value, exists := client.sendResponse.Data[sVarName_]
	if !exists {
		log.Warn(fmt.Sprintf("Value not declared in response: %s", sVarName_))

		return
	}

	gff := value.GetGffValue()
	C.memcpy(unsafe.Pointer(result), unsafe.Pointer(&gff[0]), resultSize)
}

func (p *rpcPlugin) setGff(sVarName *C.char, gffData *C.uint8_t, _ C.size_t) {
	sVarName_ := C.GoString(sVarName)
	varNameSplits := strings.SplitN(sVarName_, "///", 2)
	var clientKey string
	if len(varNameSplits) == 2 {
		clientKey = varNameSplits[0]
		sVarName_ = varNameSplits[1]
	} else {
		return
	}
	gffData_ := *(*[]byte)(unsafe.Pointer(gffData))
	client, ok := p.getRpcClient(clientKey)
	if !ok {
		return
	}

	if client.sendRequest == nil {
		client.resetSend()
	}

	client.sendRequest.Params[sVarName_] = &pbCore.Value{Value: &pbCore.Value_GffValue{GffValue: gffData_}}
}

// rpcClient contains the clients to RPC microservices
type rpcClient struct {
	name               string
	url                string
	eventServiceClient pbCore.EventServiceClient
	sendRequest        *pbCore.SendRequest
	sendResponse       *pbCore.SendResponse
}

func (c *rpcClient) resetSend() {
	c.sendRequest = &pbCore.SendRequest{}
	c.sendResponse = nil
}

func (c *rpcClient) send() bool {
	if c.sendRequest == nil {
		c.resetSend()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	response, err := c.eventServiceClient.Send(ctx, c.sendRequest)

	if err != nil {
		log.Error(fmt.Sprintf("Error sending request: %s", err))

		return false
	}

	c.sendRequest = nil
	c.sendResponse = response

	return true
}

// newRpcPlugin builds and returns a new RPC plugin
func newRpcPlugin() *rpcPlugin {
	return &rpcPlugin{
		rpcClients: make(map[string]rpcClient),
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
	return C.CString("NWNX4 RPC - A better way to build service-oriented applications in NWN2")
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
	logFile, err := os.OpenFile(path.Join(nwnxHomePath_, "src.log"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0
	}

	// Adding the header/description to the log
	header := fmt.Sprintf(
		"%s %s \n"+
			"(c) 2021-2022 by ihatemundays (scottmunday84@gmail.com) \n", pluginName, pluginVersion)
	description := "A better way to build service-oriented applications in NWN2"

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(logFile)
	log.Info(header)
	log.Info(description)

	// Get YAML file with services
	configFile, err2 := ioutil.ReadFile(path.Join(nwnxHomePath_, "src.yml"))
	if err2 != nil {
		log.Error(err2)

		return 0
	}
	config := config{}
	err3 := yaml.Unmarshal(configFile, &config)
	if err3 != nil {
		log.Error(err3)

		return 0
	}

	log.Info("Processing configuration file")

	// Initialize the server
	plugin.initRpcServer(config.server)

	// Build out the clients
	for name, url := range config.clients {
		log.Info(fmt.Sprintf("Adding client %s@%s", name, url))
		plugin.addRpcClient(name, url)
	}

	log.Info("Initialized plugin")

	// Giving back address
	return C.uint32_t(reflect.ValueOf(plugin).Pointer())
}

//export NWNXCPlugin_Delete
func NWNXCPlugin_Delete(_ uintptr) {}

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
