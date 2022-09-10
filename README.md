# NWNX4 RPC Plugin

A modern RPC plugin for NWNX4.

## Overview

xp_rpc is a plugin that allows high performance RPC communication between an NWN2 server and any set of microservices. A
microservice can be built in any supported language of both gRPC and Protocol Buffers.

## Purpose

![xp_rpc Concept Chart](docs/assets/xp_rpc-concept.svg)

Previously any plugin developed for NWNX4 would be a Windows x86 C++ DLL. This C++ DLL would have a core plugin class:
a child class of the parent class. With the advent of the new ABI interface, plugins expand the capability of supported
plugins to anything that can be developed into a C library, but still targeting the Windows x86 architecture.

Plugins were further designed with a domain-oriented approach (i.e. data, health, system). The xp_rpc plugin
potentially allows a more service-oriented approach. All domains (authentication, data, etc.) can now exist inside a
scalable, distributable, performant microservice running on practically any environment fully (or partially) decoupled
from the host itself.

## gRPC

![gRPC Concept Chart](docs/assets/grpc-concept.svg)

gRPC with protobufs allow an efficient, fast and secure approach to data contract interfaces. To read more about gRPC
and its capabilities, go [here](https://grpc.io).

## How Does It Work

The plugin uses a handful of protobuf messages and services to build a data contract between itself and the service. A
microservice can then be developed in any supported language with a thorough definition of the service implementation. A
YAML configuration should then be manually set with the list of the service names and their respective paths. With the
NWNX4 application running with the xp_rpc plugin, the application is ready to transmit requests to the microservice.

All that will be left is to send requests through the module using NWScript.

## Use

1. Include the `include/nwnx_rpc.nss` file into your module.
2. All 6 base NWNX* functions including the 2 RCO/SCO functions are available for use.

```NWScript
int RPCGetInt(string sClient, string sParam1);
void RPCSetInt(string sClient, string sParam1, int nValue);

bool RPCGetBool(string sClient, string sParam1);
void RPCSetBool(string sClient, string sParam1, bool bValue);

float RPCGetFloat(string sClient, string sParam1);
void RPCSetFloat(string sClient, string sParam1, float fValue);

string RPCGetString(string sClient, string sParam1);
void RPCSetString(string sClient, string sParam1, string sValue);

object RPCRetrieveCampaignObject(string client, string sVarName, object oObject);
int RPCStoreCampaignObject(string sClient, string sVarName, object oObject); 
```

3. Build a microservice to handle the requests.

## Configuration

Clients are connections from the RPC plugin to that of the services. To setup your xp_rpc plugin to bind to the client
on startup, you create a YAML configuration `xp_rpc.yml` at the root of your NWNX4 path. In the previous example, "
clientName" is the unique identifier to the microservice. What follows is a YAML configuration example:

```yaml
server:
  log:
    logLevel: debug
clients:
  clientName: localhost:3000
```

| Setting              | Description                                                                                          |
|----------------------|------------------------------------------------------------------------------------------------------|
| server:log:logLevel  | A string representation of the log level you wish to use (default: info)                             |
| clients[key]         | Key is the client name and how you call it through NWN 2; the value is the URL route to the service. |

## Microservices

Microservices are the real power behind this plugin. They can be developed in any architecture and many languages.

### Building your Microservice

Official support for all the gRPC programming languages can be found [here](https://grpc.io/docs/#official-support).

#### Requirements

Each language requires both the following, so you can compile the protobufs for your application:

* [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/downloads)
* gRPC (per language):
    * [C#](https://grpc.io/docs/languages/csharp/)
    * [C++](https://grpc.io/docs/languages/cpp/)
    * [Dart](https://grpc.io/docs/languages/dart/)
    * [Go](https://grpc.io/docs/languages/go/)
    * [Java](https://grpc.io/docs/languages/java/)
    * [Kotlin](https://grpc.io/docs/languages/kotlin/)
    * [Node](https://grpc.io/docs/languages/node/)
    * [Objective-C](https://grpc.io/docs/languages/objective-c/https://grpc.io/docs/languages/objective-c/)
    * [PHP](https://grpc.io/docs/languages/php/)
    * [Python](https://grpc.io/docs/languages/python/)
    * [Ruby](https://grpc.io/docs/languages/ruby/)

Some other languages are unofficially supported. They are available, but require more setup.

#### Use

Once you have installed your requirements, building your protobufs are easy. For example, if you use Go, the following
bash CLI command will develop your Go files.

```bash
protoc 
  --proto_path=proto proto/*.proto 
  --go_out=../some_go_project/proto 
  --go-grpc_out=../some_go_project/proto
```

Replace go_ with the language you wish to use (java_, cpp_, etc.) and the *_out value with the location you wish to
place your built files.

From your project folder for your microservice, include your files into your application and create a service using the
documentation mentioned above. For our example, we would build a service like so:

```go
type rpcServer struct {
    pb.UnimplementedCallServiceServer
}

func (s *rpcServer) Call(ctx context.Context, in *pb.CallRequest) (*pb.CallResponse, error) {
    ...
}
```

*pb* is an alias for the protobuf package. You can name this whatever makes sense to you.

The CallRequest contains all parameters provided with a key set through the 4 set functions listed above. This can be as
many parameters as you want and how complex you want.

The CallResponse contains all values returned with a key set through the 4 get functions listed above. This can be as
many values as you want and how complex you want.

Any of the following types are available either on the request or the response.

* Boolean
* Integer (32-bit)
* Float (32-bit)
* String (ASCII)
* GFF File (in bytes)

You might be wondering, then how do I tell it to send a request to my microservice? In other words, how do I call a
remote procedure? Easy. Set functions setup a request (if not set);
get functions send the call and return the result. Once you set variables on the client again, you reset the process.
It's quite simple, but quite powerful.

Happy coding!