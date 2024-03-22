# TCP Server Abstraction Library

## Overview

This project is a simple TCP Server Abstraction Library created to provide a
straightforward and customizable way to create TCP servers in Golang. It
abstracts away the complex parts of managing TCP servers and provides a clear
interface that software developers can use easily. Moreover, it efficiently
limits the number of goroutines in use, ensuring optimum performance.

## Features

- Concurrent management of multiple TCP connections,
- Customizable TCP connection handlers,
- Limiting the number of goroutines,
- Easy to use `client` and `server` types for inbound and outbound messages,
- Error handlers on the client and server side.

## Installing the Library

To install this library you can use the well-known `go get` to download it to
your GOPATH:

```bash
go get github.com/jtarchie/tcpserver
```

## Usage

The library consists of several functionalities. Here are some starting points:

1. Creating a new server:

```go
ctx := context.TODO()
server, err := tcp.NewServer(ctx, 8080, 1) // listen on port 8080, with 1 limited goroutine
if err != nil {
    log.Fatalf("error starting server: %v", err)
}
```

2. Handling incoming connections with a custom handler:

```go
type CustomHandler struct{}
func (*CustomHandler) OnConnection(_ context.Context, rw io.ReadWriter) error {
    _, err := rw.Write([]byte("Hello"))
    return err
}
server.Listen(ctx, &CustomHandler{})
```

3. Closing connections:

```go
err := server.Close()
if err != nil {
    log.Print("error closing server: $v", err)
}
```

4. Client operations:

This was build to have an easy way to test clients.

```go
response, err := tcp.Write(8080, "message")
if err != nil {
	log.Fatalf("error writing to server: %v", err)
}
fmt.Printf("response from server: %s\n", response)
```

## Contributing

If you wish to contribute to this project, create an issue to discuss.
