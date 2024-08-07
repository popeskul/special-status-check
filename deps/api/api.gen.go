//go:build go1.22

// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.1-0.20240623082115-ae781fb10484 DO NOT EDIT.
package api

import (
	"context"
	"fmt"
	"net/http"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns 500 Internal Server Error
	// (GET /internal)
	GetInternal(w http.ResponseWriter, r *http.Request)
	// Returns random status
	// (POST /random)
	PostRandom(w http.ResponseWriter, r *http.Request)
	// Returns 200 OK
	// (GET /success)
	GetSuccess(w http.ResponseWriter, r *http.Request)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// GetInternal operation middleware
func (siw *ServerInterfaceWrapper) GetInternal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetInternal(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// PostRandom operation middleware
func (siw *ServerInterfaceWrapper) PostRandom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostRandom(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// GetSuccess operation middleware
func (siw *ServerInterfaceWrapper) GetSuccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetSuccess(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       *http.ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m *http.ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m *http.ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/internal", wrapper.GetInternal)
	m.HandleFunc("POST "+options.BaseURL+"/random", wrapper.PostRandom)
	m.HandleFunc("GET "+options.BaseURL+"/success", wrapper.GetSuccess)

	return m
}

type GetInternalRequestObject struct {
}

type GetInternalResponseObject interface {
	VisitGetInternalResponse(w http.ResponseWriter) error
}

type GetInternal500Response struct {
}

func (response GetInternal500Response) VisitGetInternalResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

type PostRandomRequestObject struct {
}

type PostRandomResponseObject interface {
	VisitPostRandomResponse(w http.ResponseWriter) error
}

type PostRandom200Response struct {
}

func (response PostRandom200Response) VisitPostRandomResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PostRandom500Response struct {
}

func (response PostRandom500Response) VisitPostRandomResponse(w http.ResponseWriter) error {
	w.WriteHeader(500)
	return nil
}

type GetSuccessRequestObject struct {
}

type GetSuccessResponseObject interface {
	VisitGetSuccessResponse(w http.ResponseWriter) error
}

type GetSuccess200Response struct {
}

func (response GetSuccess200Response) VisitGetSuccessResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Returns 500 Internal Server Error
	// (GET /internal)
	GetInternal(ctx context.Context, request GetInternalRequestObject) (GetInternalResponseObject, error)
	// Returns random status
	// (POST /random)
	PostRandom(ctx context.Context, request PostRandomRequestObject) (PostRandomResponseObject, error)
	// Returns 200 OK
	// (GET /success)
	GetSuccess(ctx context.Context, request GetSuccessRequestObject) (GetSuccessResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// GetInternal operation middleware
func (sh *strictHandler) GetInternal(w http.ResponseWriter, r *http.Request) {
	var request GetInternalRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetInternal(ctx, request.(GetInternalRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetInternal")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetInternalResponseObject); ok {
		if err := validResponse.VisitGetInternalResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostRandom operation middleware
func (sh *strictHandler) PostRandom(w http.ResponseWriter, r *http.Request) {
	var request PostRandomRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostRandom(ctx, request.(PostRandomRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostRandom")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostRandomResponseObject); ok {
		if err := validResponse.VisitPostRandomResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetSuccess operation middleware
func (sh *strictHandler) GetSuccess(w http.ResponseWriter, r *http.Request) {
	var request GetSuccessRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetSuccess(ctx, request.(GetSuccessRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetSuccess")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetSuccessResponseObject); ok {
		if err := validResponse.VisitGetSuccessResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
