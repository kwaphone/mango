package mango

import (
	"bytes"
	"http"
	"http/httptest"
	"log"
	"testing"
	"runtime"
)

var loggerBuffer = &bytes.Buffer{}

func loggerTestServer(env Env) (Status, Headers, Body) {
	env.Logger().Println("Never gonna give you up")
	return 200, Headers{}, Body("Hello World!")
}

func init() {
	runtime.GOMAXPROCS(4)
}

func TestLogger(t *testing.T) {
	// Compile the stack
	loggerStack := new(Stack)
	custom_logger := log.New(loggerBuffer, "prefixed:", 0)
	loggerStack.Middleware(Logger(custom_logger))
	loggerApp := loggerStack.Compile(loggerTestServer)

	// Request against it
	request, err := http.NewRequest("GET", "http://localhost:3000/", nil)
	status, _, _ := loggerApp(Env{"mango.request": &Request{request}})

	if err != nil {
		t.Error(err)
	}

	if status != 200 {
		t.Error("Expected status to equal 200, got:", status)
	}

	expected := "prefixed:Never gonna give you up\n"
	if loggerBuffer.String() != expected {
		t.Error("Expected logger to print: \"", expected, "\" got: \"", loggerBuffer.String(), "\"")
	}
}

func BenchmarkLogger(b *testing.B) {
	b.StopTimer()

	stack := new(Stack)
	custom_logger := log.New(loggerBuffer, "prefixed:", 0)
	stack.Middleware(Logger(custom_logger))
	testServer := httptest.NewServer(stack.HandlerFunc(loggerTestServer))
	defer testServer.Close()
	address := testServer.URL

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		http.Get(address)
	}
	b.StopTimer()
}
