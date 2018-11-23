package ipc

import "testing"

type EchoServer struct {
}

func (server *EchoServer) Handle(method, params string) *Response {
	return &Response{"OK", "ECHO: " + method + " ~ " + params}
}

func (server *EchoServer) Name() string {
	return "EchoServer"
}

func TestIpc(test *testing.T) {
	server := NewIpcServer(&EchoServer{})

	client1 := NewIpcClient(server)
	client2 := NewIpcClient(server)

	resp1, _ := client1.Call("foo", "From Client1")
	resp2, _ := client2.Call("foo", "From Client2")

	if resp1.Body != "ECHO: foo ~ From Client1" ||
		resp2.Body != "ECHO: foo ~ From Client2" {
		test.Errorf("err resp1 %s resp2 %s", resp1.Body, resp2.Body)
	}

	client1.Close()
	client2.Close()
}
