# connect-go unary interceptor sample for protoc-gen-grpc

# How to test

run connect server

`go run example`

call incorrect request using curl, etc.

Here is the schema with validation [proto/greet/v1/greet.proto](proto/greet/v1/greet.proto)

```
curl \
    --header "Content-Type: application/json" \
    --data '{"first_name": "toolongfirst_name","last_name":"s"}' \
    http://localhost:8080/greet.v1.GreetService/Greet
```

The following response will be returned

```json
// Response
{
  "code": "invalid_argument",
  "message": "Client specified an invalid argument. Check error details for more information.",
  "details": [
    {
      "type": "google.rpc.BadRequest",
      "value": "ClAKFkdyZWV0UmVxdW...",
      "debug": {
        "@type": "type.googleapis.com/google.rpc.BadRequest",
        "fieldViolations": [
          {
            "field": "GreetRequest.FirstName",
            "description": "value length must be between 2 and 10 runes, inclusive"
          },
          {
            "field": "GreetRequest.FirstName",
            "description": "value does not have prefix \"foo\""
          },
          {
            "field": "GreetRequest.LastName",
            "description": "value length must be between 2 and 10 runes, inclusive"
          }
        ]
      }
    }
  ]
}
```
