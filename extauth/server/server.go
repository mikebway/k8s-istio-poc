package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoyt "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"

	"github.com/mikebway/envoy-auth-poc/extauth/userjwt"
)

// AuthorizationServer is the struct definition that acts as the Envoy external authorization service class definition.
//
// When we have values that we need to share across invocations to an instance of the service but not more
// widely, they can be defined here.
type AuthorizationServer struct{}

// denied is a shorthand function for reject a request as unauthorized in some fashion
func denied(code int32, body string) *auth.CheckResponse {
	return &auth.CheckResponse{
		Status: &status.Status{Code: code},
		HttpResponse: &auth.CheckResponse_DeniedResponse{
			DeniedResponse: &auth.DeniedHttpResponse{
				Status: &envoyt.HttpStatus{
					Code: envoyt.StatusCode(code),
				},
				Body: body,
			},
		},
	}
}

// allowed is a shorthand function for generating a success response to Envoy, let this request go through
func allowed(authHeader *core.HeaderValue) *auth.CheckResponse {

	// Start a slice of header with the header that we always set
	headers := []*core.HeaderValueOption{
		{
			Header: &core.HeaderValue{
				Key:   "x-extauth-was-here",
				Value: "blah, blah, blah",
			},
		},
	}

	// If we are to add an authorization header, add it
	if authHeader != nil {
		headers = append(headers, &core.HeaderValueOption{
			Header: authHeader,
		})
	}

	// Construct and return the response with our header(s)
	return &auth.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &auth.CheckResponse_OkResponse{
			OkResponse: &auth.OkHttpResponse{
				Headers: headers,
			},
		},
	}
}

// Check implements Envoy Authorization service. Proto file:
// https://github.com/envoyproxy/envoy/blob/main/api/envoy/service/auth/v3/external_auth.proto
func (a *AuthorizationServer) Check(_ context.Context, req *auth.CheckRequest) (*auth.CheckResponse, error) {

	// Log the target URL of the request that we are verifying authorization for
	var msgBuilder strings.Builder
	httpRequest := req.Attributes.Request.Http
	msgBuilder.WriteString("######### START REQUEST #########\n")
	msgBuilder.WriteString("Request: " + httpRequest.Protocol + " " + httpRequest.Method + " " + httpRequest.Scheme + "://" + httpRequest.Host + httpRequest.Path + "\n")

	// Log the incoming headers and context extensions
	msgBuilder.WriteString(sprintIncomingHeaders(req))

	// Unpack and log the cookies
	cookies := readCookies(httpRequest.Headers["cookie"])
	msgBuilder.WriteString(sprintCookies(cookies))

	// Build an authorization header if there is a session cookie
	authHeader, err := authorize(cookies["session"])
	if err != nil {
		// Log the error but carry on anyways; the authHeader will have been set with a blank value
		msgBuilder.WriteString(err.Error() + "\n")
	}

	// Assemble our response here
	var response *auth.CheckResponse
	response = allowed(authHeader)

	// Send our built message to the logs
	msgBuilder.WriteString("######### END REQUEST #########\n\n")
	log.Print(msgBuilder.String())

	// Default to allowing the request (for now)
	return response, nil
}

// authorize examines the request to see if there is a session cookie and, if there is, returns a populated
// "x-extauth-authorization" header. Otherwise, returns an empty "x-extauth-authorization" header.
//
// The header value will always be returned, even in the unexpected situate where an error is also returned.
func authorize(sessionCookie *RequestCookie) (*core.HeaderValue, error) {

	// Either way, we return a header. The only question is what its value might be
	authHeader := &core.HeaderValue{
		Key: "x-extauth-authorization",
	}

	// If we have a session cookie then we have a key to look up the session JWT to set as
	// a bearer token in the header
	if sessionCookie != nil && sessionCookie.Value != "" {
		jwt, err := userjwt.CreateJWT(sessionCookie.Value)
		if err != nil {
			return authHeader, fmt.Errorf("failed to create JWT: %w", err)
		}
		authHeader.Value = "Bearer " + jwt
	}

	// Return whatever authorization we could scrape together
	return authHeader, nil
}

// sprintIncomingHeaders renders both the regular request headers and context extensions that the filter
// received to a string for logging.
func sprintIncomingHeaders(req *auth.CheckRequest) string {

	// Gather the incoming headers as a formatted JSON structure in a string builder
	var msgBuilder strings.Builder
	msgBuilder.WriteString("\n=== Headers ===\n")
	jsonBytes, err := json.MarshalIndent(req.Attributes.Request.Http.Headers, "", "  ")
	if err == nil {
		msgBuilder.WriteString(string(jsonBytes))
	} else {
		msgBuilder.WriteString("failed to marshal headers: " + err.Error())
	}
	msgBuilder.WriteString("\n===============\n")

	// Add the headers / parameters that were sent just for us and will not be passed through to the
	// "upstream" services deeper inside our private network
	jsonBytes, err = json.MarshalIndent(req.Attributes.ContextExtensions, "", "  ")
	msgBuilder.WriteString("\n=== Context Extensions ===\n")
	if err == nil {
		msgBuilder.WriteString(string(jsonBytes))
	} else {
		msgBuilder.WriteString("failed to marshal context extensions: " + err.Error())
	}
	msgBuilder.WriteString("\n==========================\n")

	// All done, return the string we have been building
	return msgBuilder.String()
}

// sprintCookies renders the cookies that the filter received to a string for logging
func sprintCookies(cookies map[string]*RequestCookie) string {

	// Gather the incoming headers as a formatted JSON structure in a string builder
	var msgBuilder strings.Builder
	msgBuilder.WriteString("\n=== Cookies ===\n")
	jsonBytes, err := json.MarshalIndent(cookies, "", "  ")
	if err == nil {
		msgBuilder.WriteString(string(jsonBytes))
	} else {
		msgBuilder.WriteString("failed to marshal cookies: " + err.Error())
	}
	msgBuilder.WriteString("\n===============\n")

	// All done, return the string we have been building
	return msgBuilder.String()
}
