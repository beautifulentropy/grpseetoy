# gRPSee

A gRPC server and client with verbose logging enabled, intended as a toy for
testing gRPC service keepalive and client reconnect behavior.

## Usage

### Server

```shell
$ go run cmd/server/main.go
```

### Client

```shell
$ go run cmd/client/main.go
```

## HTTP/2 Debugging

To see HTTP/2 frame information, set the `GODEBUG` environment variable like so:

```shell
$ export GODEBUG=http2debug=2
```

## Example Output

Here you can watch server and client negotiate a connection, then the server
close the connection after 30 seconds of inactivity (+/- 10% jitter), followed
by the client succesfully reconnecting.

### Server

```shell
$ go run ./cmd/server/main.go
2023/04/26 12:34:09 INFO: [core] [Server #1] Server created
2023/04/26 12:34:09 INFO: [core] [Server #1 ListenSocket #2] ListenSocket created
2023/04/26 12:34:18 INFO: [core] CPU time info is unavailable on non-linux environments.
2023/04/26 12:34:47 INFO: [transport] transport: closing: EOF
2023/04/26 12:34:47 INFO: [transport] transport: loopyWriter exiting with error: transport closed by client
```

### Client

```shell
$ go run ./cmd/client/main.go
2023/04/26 12:21:14 INFO: [core] [Channel #1] Channel created
2023/04/26 12:21:14 INFO: [core] [Channel #1] original dial target is: "127.0.0.1:1337"
2023/04/26 12:21:14 INFO: [core] [Channel #1] dial target "127.0.0.1:1337" parse failed: parse "127.0.0.1:1337": first path segment in URL cannot contain colon
2023/04/26 12:21:14 INFO: [core] [Channel #1] fallback to scheme "passthrough"
2023/04/26 12:21:14 INFO: [core] [Channel #1] parsed dial target is: {Scheme:passthrough Authority: URL:{Scheme:passthrough Opaque: User: Host: Path:/127.0.0.1:1337 RawPath: OmitHost:false ForceQuery:false RawQuery: Fragment: RawFragment:}}
2023/04/26 12:21:14 INFO: [core] [Channel #1] Channel authority set to "127.0.0.1:1337"
2023/04/26 12:21:14 INFO: [core] [Channel #1] Resolver state updated: {
  "Addresses": [
    {
      "Addr": "127.0.0.1:1337",
      "ServerName": "",
      "Attributes": null,
      "BalancerAttributes": null,
      "Type": 0,
      "Metadata": null
    }
  ],
  "ServiceConfig": null,
  "Attributes": null
} (resolver returned new addresses)
2023/04/26 12:21:14 INFO: [core] [Channel #1] Channel switches to new LB policy "pick_first"
2023/04/26 12:21:14 INFO: [core] [Channel #1 SubChannel #2] Subchannel created
2023/04/26 12:21:14 INFO: [core] [Channel #1] Channel Connectivity change to CONNECTING
2023/04/26 12:21:14 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to CONNECTING
2023/04/26 12:21:14 INFO: [core] [Channel #1 SubChannel #2] Subchannel picks a new address "127.0.0.1:1337" to connect
2023/04/26 12:21:14 INFO: [core] pickfirstBalancer: UpdateSubConnState: 0x1400000e858, {CONNECTING <nil>}
2023/04/26 12:21:14 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to READY
2023/04/26 12:21:14 INFO: [core] pickfirstBalancer: UpdateSubConnState: 0x1400000e858, {READY <nil>}
2023/04/26 12:21:14 INFO: [core] [Channel #1] Channel Connectivity change to READY
2023/04/26 12:21:44 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to IDLE
2023/04/26 12:21:44 INFO: [core] blockingPicker: the picked transport is not ready, loop back to repick
2023/04/26 12:21:44 INFO: [transport] transport: closing: connection error: desc = "received goaway and there are no active streams"
2023/04/26 12:21:44 INFO: [transport] transport: loopyWriter exiting with error: transport closed by client
2023/04/26 12:21:44 INFO: [core] pickfirstBalancer: UpdateSubConnState: 0x1400000e858, {IDLE <nil>}
2023/04/26 12:21:44 INFO: [core] [Channel #1] Channel Connectivity change to IDLE
2023/04/26 12:21:44 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to CONNECTING
2023/04/26 12:21:44 INFO: [core] [Channel #1 SubChannel #2] Subchannel picks a new address "127.0.0.1:1337" to connect
2023/04/26 12:21:44 INFO: [core] pickfirstBalancer: UpdateSubConnState: 0x1400000e858, {CONNECTING <nil>}
2023/04/26 12:21:44 INFO: [core] [Channel #1] Channel Connectivity change to CONNECTING
2023/04/26 12:21:44 INFO: [core] [Channel #1 SubChannel #2] Subchannel Connectivity change to READY
2023/04/26 12:21:44 INFO: [core] pickfirstBalancer: UpdateSubConnState: 0x1400000e858, {READY <nil>}
2023/04/26 12:21:44 INFO: [core] [Channel #1] Channel Connectivity change to READY
```

### HTTP/2 Frames

```shell
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read DATA flags=END_STREAM stream=183 len=12 data="\x00\x00\x00\x00\a\n\x05world"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote WINDOW_UPDATE len=4 (conn) incr=12
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote HEADERS flags=END_HEADERS stream=183 len=2
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote DATA stream=183 len=18 data="\x00\x00\x00\x00\r\n\vhello world"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote HEADERS flags=END_STREAM|END_HEADERS stream=183 len=2
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read WINDOW_UPDATE len=4 (conn) incr=18
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read HEADERS flags=END_HEADERS stream=185 len=7
2023/04/25 19:07:47 http2: decoded hpack field header field ":method" = "POST"
2023/04/25 19:07:47 http2: decoded hpack field header field ":scheme" = "http"
2023/04/25 19:07:47 http2: decoded hpack field header field ":path" = "/helloworld.Greeter/SayHello"
2023/04/25 19:07:47 http2: decoded hpack field header field ":authority" = "localhost:1337"
2023/04/25 19:07:47 http2: decoded hpack field header field "content-type" = "application/grpc"
2023/04/25 19:07:47 http2: decoded hpack field header field "user-agent" = "grpc-go/1.54.0"
2023/04/25 19:07:47 http2: decoded hpack field header field "te" = "trailers"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read DATA flags=END_STREAM stream=185 len=12 data="\x00\x00\x00\x00\a\n\x05world"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote WINDOW_UPDATE len=4 (conn) incr=12
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote HEADERS flags=END_HEADERS stream=185 len=2
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote DATA stream=185 len=18 data="\x00\x00\x00\x00\r\n\vhello world"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote HEADERS flags=END_STREAM|END_HEADERS stream=185 len=2
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read WINDOW_UPDATE len=4 (conn) incr=18
2023/04/25 19:07:47 http2: Framer 0x140000f4000: read PING len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
2023/04/25 19:07:47 http2: Framer 0x140000f4000: wrote PING flags=ACK len=8 ping="\x02\x04\x10\x10\t\x0e\a\a"
```
