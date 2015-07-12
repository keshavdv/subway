# About

This is a tunnelling application that is heavily based on [ngrok](https://github.com/inconshreveable/ngrok) that tunnels data over a WebSocket instead of standard TCP ports

### Building

To build the server and client binaries:

    make server
    make client

The binaries can be found in the bin/ directory.

### Protocol

##### Establishing tunnels:

1. Client connects to server via a WebSocket (auth included during creation)
2. Client creates a stream (control stream) and sends TunnelRequest to server to create tunnels
3. Server responds with TunnelReply to acknowledge (providing external connection URL)


#####  Tunneling:

1. Server sends ProxyRequest to client over control stream
2. Client creates new stream (the proxy stream) and responds with ProxyStart message
3. Server replies with information about client with a ProxyClientInfo message
4. Data is proxied byte-for-byte over the proxy stream