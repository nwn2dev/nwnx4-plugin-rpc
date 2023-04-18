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
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

const pluginID string = "RPC"               // Plugin ID used for identification in the list
const pluginName string = "NWNX RPC Plugin" // Plugin name passed to hook
const pluginDescription string = "A better way to integrate services with NWN2"
const pluginVersion string = "0.3.1" // Plugin version passed to hook
const pluginContact string = "(c) 2021-2023 by ihatemundays (scottmunday84@gmail.com)"

const logFilename string = "xp_rpc.log"
const configFilename string = "xp_rpc.yml"

const rpcPluginSeparator string = "!"

var plugin *rpcPlugin // Singleton

// All exports to C library

//export NWNXCPlugin_GetID
func NWNXCPlugin_GetID(_ *C.void) *C.char {
	return C.CString(pluginID)
}

//export NWNXCPlugin_GetVersion
func NWNXCPlugin_GetVersion() *C.char {
	return C.CString(pluginVersion)
}

//export NWNXCPlugin_GetInfo
func NWNXCPlugin_GetInfo() *C.char {
	info := fmt.Sprintf("%s v%s - %s", pluginName, pluginVersion, pluginDescription)

	return C.CString(info)
}

//export NWNXCPlugin_New
func NWNXCPlugin_New(initInfo C.CPluginInitInfo) C.uint32_t {
	// Set up the log file
	nwnxHomePath_ := C.GoString(initInfo.nwnx_user_path)
	logFile, err := os.OpenFile(path.Join(nwnxHomePath_, logFilename), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0
	}

	// Adding the header/description to the log
	header := fmt.Sprintf("%s v%s \n %s \n", pluginName, pluginVersion, pluginContact)
	description := pluginDescription

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(logFile)
	log.Info(header)
	log.Info(description)
	log.SetLevel(log.TraceLevel)

	plugin = newRpcPlugin()

	// Get YAML file with services
	configFilepath := path.Join(nwnxHomePath_, configFilename)
	configFile, err := ioutil.ReadFile(configFilepath)
	if err != nil {
		log.Error(err)

		return 0
	}

	err = yaml.Unmarshal(configFile, &plugin.config)
	if err != nil {
		log.Error(err)

		return 0
	}

	// Set up the RPC plugin from the configuration
	plugin.init()

	// Giving back address; still maintained by the plugin itself
	return C.uint32_t(reflect.ValueOf(plugin).Pointer())
}

//export NWNXCPlugin_Delete
func NWNXCPlugin_Delete(_ *C.void) {}

//export NWNXCPlugin_GetInt
func NWNXCPlugin_GetInt(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.int {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1":   sParam1_,
		"nParam2":   nParam2_,
	}).Debug("NWNXGetInt()")

	return C.int(plugin.getInt(sFunction_, sParam1_, nParam2_))
}

//export NWNXCPlugin_SetInt
func NWNXCPlugin_SetInt(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, nValue C.int) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	nValue_ := int32(nValue)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1":   sParam1_,
		"nParam2":   nParam2_,
		"nValue":    nValue_,
	}).Debug("NWNXSetInt()")

	plugin.setInt(sFunction_, sParam1_, nParam2_, nValue_)
}

//export NWNXCPlugin_GetFloat
func NWNXCPlugin_GetFloat(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int) C.float {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1":   sParam1_,
		"nParam2":   nParam2_,
	}).Debug("NWNXGetFloat()")

	return C.float(plugin.getFloat(sFunction_, sParam1_, nParam2_))
}

//export NWNXCPlugin_SetFloat
func NWNXCPlugin_SetFloat(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, fValue C.float) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	fValue_ := float32(fValue)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1_":  sParam1_,
		"nParam2_":  nParam2_,
		"fValue":    fValue_,
	}).Debug("NWNXSetFloat()")

	plugin.setFloat(sFunction_, sParam1_, nParam2_, fValue_)
}

//export NWNXCPlugin_GetString
func NWNXCPlugin_GetString(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, result *C.char, resultSize C.size_t) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1":   sParam1_,
		"nParam2":   nParam2_,
	}).Debug("NWNXGetString()")

	response := C.CString(plugin.getString(sFunction_, sParam1_, nParam2_))

	// Get the pointer to the memory
	responseSize := C.strlen(response)
	responsePtr := unsafe.Pointer(response)
	defer C.free(responsePtr)

	// Copy the response over to the result
	C.strncpy_s(result, resultSize, response, responseSize)
}

//export NWNXCPlugin_SetString
func NWNXCPlugin_SetString(_ *C.void, sFunction, sParam1 *C.char, nParam2 C.int, sValue *C.char) {
	sFunction_ := C.GoString(sFunction)
	sParam1_ := C.GoString(sParam1)
	nParam2_ := int32(nParam2)
	sValue_ := C.GoString(sValue)
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sParam1":   sParam1_,
		"nParam2":   nParam2_,
		"sValue":    sValue_,
	}).Debug("NWNXSetString()")

	plugin.setString(sFunction_, sParam1_, nParam2_, sValue_)
}

//export NWNXCPlugin_GetGFFSize
func NWNXCPlugin_GetGFFSize(_ *C.void, sVarName *C.char) C.size_t {
	sFunction_, sVarName_, sParam2_, nParam2_ := "", C.GoString(sVarName), "", rpcCallActionParam2Default
	splits := strings.SplitN(C.GoString(sVarName), rpcPluginSeparator, 3)
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else if len(splits) == 3 {
		sFunction_, sParam2_, sVarName_ = splits[0], splits[1], splits[2]
	}
	if v, err := strconv.Atoi(sParam2_); err == nil {
		nParam2_ = int32(v)
	}
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sVarName":  sVarName_,
		"nParam2":   nParam2_,
	}).Debug("SCORCOGetGFFSize()")

	return C.size_t(plugin.getGffSize(sFunction_, sVarName_, nParam2_))
}

//export NWNXCPlugin_GetGFF
func NWNXCPlugin_GetGFF(_ *C.void, sVarName *C.char, result *C.uint8_t, resultSize C.size_t) {
	sFunction_, sVarName_, sParam2_, nParam2_ := "", C.GoString(sVarName), "", rpcCallActionParam2Default
	splits := strings.SplitN(C.GoString(sVarName), rpcPluginSeparator, 3)
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else if len(splits) == 3 {
		sFunction_, sParam2_, sVarName_ = splits[0], splits[1], splits[2]
	}
	if v, err := strconv.Atoi(sParam2_); err == nil {
		nParam2_ = int32(v)
	}
	log.WithFields(log.Fields{
		"sFunction": sFunction_,
		"sVarName":  sVarName_,
		"nParam2":   nParam2_,
	}).Debug("SCORCOGetGFF()")

	response := plugin.getGff(sFunction_, sVarName_, nParam2_)
	if response == nil {
		log.Error("GFF response is empty")

		return
	}

	// Get the pointer to the memory
	responseSize := C.size_t(len(response))
	responsePtr := unsafe.Pointer(&response[0])
	if resultSize < responseSize {
		log.Errorf("%d response size is too large for the %d result size", uint32(responseSize), uint32(resultSize))

		return
	}

	// Copy the response over to the result
	resultPtr := unsafe.Pointer(result)
	C.memcpy(resultPtr, responsePtr, responseSize)
}

//export NWNXCPlugin_SetGFF
func NWNXCPlugin_SetGFF(_ *C.void, sVarName *C.char, gffData *C.uint8_t, gffDataSize C.size_t) {
	sFunction_, sVarName_, sParam2_, nParam2_ := "", C.GoString(sVarName), "", rpcCallActionParam2Default
	gffDataSize_ := uint32(gffDataSize)
	gffData_ := C.GoBytes(unsafe.Pointer(gffData), C.int(gffDataSize))
	splits := strings.SplitN(C.GoString(sVarName), rpcPluginSeparator, 3)
	if len(splits) == 2 {
		sFunction_, sVarName_ = splits[0], splits[1]
	} else if len(splits) == 3 {
		sFunction_, sParam2_, sVarName_ = splits[0], splits[1], splits[2]
	}
	if v, err := strconv.Atoi(sParam2_); err == nil {
		nParam2_ = int32(v)
	}
	log.WithFields(log.Fields{
		"sFunction":   sFunction_,
		"sVarName":    sVarName_,
		"nParam2":     nParam2_,
		"gffData":     gffData_,
		"gffDataSize": gffDataSize_,
	}).Debug("NWNXSetGFF()")

	plugin.setGff(sFunction_, sVarName_, nParam2_, gffData_, gffDataSize_)
}

func main() {}
