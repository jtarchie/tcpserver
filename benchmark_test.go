package tcpserver_test

import (
	"context"
	"io"
	"log"
	"runtime"
	"testing"

	"github.com/jtarchie/tcpserver"
	"github.com/jtarchie/tcpserver/handlers"
	"github.com/phayes/freeport"
)

func BenchmarkSerial(b *testing.B) {
	port, err := freeport.GetFreePort()
	if err != nil {
		b.Fatalf("no free port: %v", err)
	}

	server, err := tcpserver.NewServer(context.TODO(), uint(port), 1)
	if err != nil {
		b.Fatalf("initialized server: %v", err)
	}

	//nolint: errcheck
	go server.Listen(context.TODO(), &handlers.Echo{})
	defer server.Close()

	client, err := tcpserver.NewClient(port)
	if err != nil {
		b.Fatalf("initialized server: %v", err)
	}
	defer client.Close()

	log.SetOutput(io.Discard)

	b.ResetTimer() // Start timing now.

	for i := 0; i < b.N; i++ {
		err = client.WriteString("Hello, World\n")
		if err != nil {
			b.Fatalf("write failed: %v", err)
		}

		_, _ = client.ReadlineString()
	}
}

func BenchmarkParallel(b *testing.B) {
	port, err := freeport.GetFreePort()
	if err != nil {
		b.Fatalf("no free port: %v", err)
	}

	server, err := tcpserver.NewServer(context.TODO(), uint(port), uint(runtime.GOMAXPROCS(0)))
	if err != nil {
		b.Fatalf("initialized server: %v", err)
	}

	//nolint: errcheck
	go server.Listen(context.TODO(), &handlers.Echo{})
	defer server.Close()

	log.SetOutput(io.Discard)

	b.ResetTimer() // Start timing now.
	b.RunParallel(func(pb *testing.PB) {
		client, err := tcpserver.NewClient(port)
		if err != nil {
			b.Fatalf("initialized server: %v", err)
		}
		defer client.Close()

		for pb.Next() {
			err = client.WriteString("Hello, World\n")
			if err != nil {
				b.Fatalf("write failed: %v", err)
			}

			_, _ = client.ReadlineString()
		}
	})
}
