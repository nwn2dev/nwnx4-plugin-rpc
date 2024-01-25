# NWNX4 RPC Plugin

NWNX4 RPC is a plugin for NWNX4 and Neverwinter Nights 2 (NWN2) game server through a set of remote procedure calls (RPC) to external applications. This allows developers to create custom functionality for NWN2 servers using services written in any language supported by gRPC, such as C++, C#, Python, or JavaScript.

## Features

![xp_rpc Concept Chart](docs/assets/xp_rpc-concept.svg)

- **Easy-to-use**: NWNX4 RPC provides a simple and intuitive approach to interacting with external applications, making it easy for developers to implement custom functionality in their NWN2 servers.
- **Cross-platform support**: NWNX4 RPC is designed to work with services on multiple platforms, including Windows, Linux, and macOS, allowing developers to use it in various environments. The only constant is HTTP/2.
- **Plugin integration**: NWNX4 RPC is designed to be integrated with other NWNX4 plugins seamlessly, providing a straightforward way to add custom functionality to NWN2 servers.
- **Language-agnostic**: NWNX4 RPC does not restrict the language used for writing plugins, allowing developers to use their preferred programming language.
- **Fault-tolerant**: NWNX4 has failover protections. Start and restart external applications without worry about the state of the NWN2 game server.

## Installation

1. Download the latest release of NWNX4 Plugin RPC from the [GitHub repository](https://github.com/nwn2dev/nwnx4-plugin-rpc/releases).
2. Extract the contents of the release archive into the `plugins` folder of your NWNX4 user directory.
3. Configure the plugin by modifying the `NWNX.ini` file in your NWN2 server installation directory, following the instructions provided in the plugin's documentation.
4. Modify any xp_rpc settings in xp_rpc.yml, placed in the root of your NWNX4 user directory.
5. Restart your NWNX4/NWN2 server to load the plugin.

## Usage

To create an external application, you will need to run `protoc` and build the appropriate files necessary to interact with the plugin.

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

Once you have installed your requirements, building your protobufs are easy. For example, if you use Go, the following
bash CLI command will develop your Go files.

```bash
protoc 
  --proto_path=proto proto/*.proto 
  --go_out=../some_go_project/proto 
  --go-grpc_out=../some_go_project/proto
```

Replace go_ with the language you wish to use (java_, cpp_, etc.) and the *_out value with the location you wish to place your built files.

With a created external application, configure your xp_rpc settings for the client from the xp_rpc.yml file.

```yaml
auth:
  pfxFilePath: path/to/pfx
  pfxPassword: password
log:
  logLevel: info
perClient:
  retries: 3
  delay: 5  # In seconds
  timeout: 30  # In seconds
clients:
  service: localhost:3000
```

In this example, there's an external application/client named "service" with a URL of "localhost:3000". Any requests to "service" in NWScript will be delivered to this client.

## Contributing

NWNX4 Plugin RPC is an open-source project, and contributions are welcome! If you would like to contribute to the project, you can do so by:

- Reporting bugs or issues in the [GitHub issue tracker](https://github.com/nwn2dev/nwnx4-plugin-rpc/issues).
- Submitting feature requests or suggestions in the [GitHub issue tracker](https://github.com/nwn2dev