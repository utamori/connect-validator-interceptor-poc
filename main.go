package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s %s!", req.Msg.LastName, req.Msg.FirstName),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {

	i := connect.WithInterceptors(ValidatorInterceptor())
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter, i)
	mux.Handle(path, handler)
	log.Fatal(http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	))
}
