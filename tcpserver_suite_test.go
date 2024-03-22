package tcpserver_test

import (
	"context"
	"testing"

	"github.com/jtarchie/tcpserver"
	"github.com/jtarchie/tcpserver/handlers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/phayes/freeport"
)

func TestTcp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TCP Suite")
}

func startServer(handler tcpserver.Handler) (int, *tcpserver.Server) {
	port, err := freeport.GetFreePort()
	Expect(err).NotTo(HaveOccurred())

	server, err := tcpserver.NewServer(context.TODO(), uint(port), 1)
	Expect(err).NotTo(HaveOccurred())

	go func() {
		defer GinkgoRecover()

		err := server.Listen(context.TODO(), handler)
		Expect(err).NotTo(HaveOccurred())
	}()

	return port, server
}

var _ = Describe("When starting a TCP Server", func() {
	It("accepts a connection", func() {
		port, server := startServer(&handlers.Echo{})
		defer server.Close()

		response, err := tcpserver.Write(port, "echo\r\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Equal("echo"))

		response, err = tcpserver.Write(port, "hello\r\n")
		Expect(err).NotTo(HaveOccurred())
		Expect(response).To(Equal("hello"))
	})

	When("the handler errors on the client", func() {
		It("server continues accepting connections", func() {
			port, server := startServer(&handlers.Error{})
			defer server.Close()

			_, err := tcpserver.Write(port, "echo\r\n")
			Expect(err).NotTo(HaveOccurred())

			_, err = tcpserver.Write(port, "error\r\n")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("could not read line"))

			_, err = tcpserver.Write(port, "echo\r\n")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
