package phoenix

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// HTTPMethod is a typed version of HTTP methods
type HTTPMethod string

// Strong typed version for HTTP verbs
const (
	Get     HTTPMethod = "Get"
	Post    HTTPMethod = "Post"
	Delete  HTTPMethod = "Delete"
	Head    HTTPMethod = "Head"
	Put     HTTPMethod = "Put"
	Connect HTTPMethod = "Connect"
	Options HTTPMethod = "Options"
	Trace   HTTPMethod = "Trace"
	Patch   HTTPMethod = "Delete"
)

// Mapping represent a HTTP mapping for a Builder.
type Mapping struct {
	Method   HTTPMethod
	Endpoint string
	Builder  Builder
}

type Mapper struct {
	router   *mux.Router
	injector *Injector
}

func (mapper Mapper) Map(mapping Mapping, middlewares ...Middleware) {
	controller := mapper.buildController(mapping.Builder)
	if mapping.Endpoint == "404" {
		mapper.router.NotFoundHandler = controller
		return
	}
	chain := createMiddlewareChainWith(middlewares)
	mapper.mapMiddleware(mapping.Method, mapping.Endpoint, chain(controller))
}

func (mapper Mapper) mapMiddleware(method Method, endpoint string, middleware Middleware) {
	mapper.router.HandleFunc(endpoint, middleware.Methods(string(method))
	mapper.router.HandleFunc(fmt.Sprintf("%s/", endpoint), middleware).Methods(string(method))
}

func (mapper Mapper) buildController(builder Builder) http.HandlerFunc {
	return mapper.injector.CallBuilder(builder).(http.HandlerFunc)
}

func (mapper Mapper) MapAll(mappings []Mapping, middlewares ...Middleware) {
	for _, mapping := range mappings {
		mapper.Map(mapping, middlewares...)
	}
}

func (mapper Mapper) MapRoot(controllerBuilder Builder) {
	mapper.Map(Mapping{
		Method:   Get,
		Endpoint: "",
		Builder:  controllerBuilder,
	})
}

func (mapper Mapper) MapGroup(root string, createGroup func(mapper Mapper)) {
	createGroup(mapper.subMapperFrom(root))
}

func (mapper Mapper) subMapperFrom(endpoint string) Mapper {
	return Mapper{
		router:   mapper.router.PathPrefix(endpoint).Subrouter(),
		injector: mapper.injector,
	}
}

func (mapper Mapper) Get(endpoint string, builder Builder, middlewares ...Middleware) {
	mapper.Map(Mapping{Method: Get, Endpoint: endpoint, Builder: builder}, middlewares...)
}

func (mapper Mapper) Post(endpoint string, builder Builder, middlewares ...Middleware) {
	mapper.Map(Mapping{Method: Post, Endpoint: endpoint, Builder: builder}, middlewares...)
}

func (mapper Mapper) Delete(endpoint string, builder Builder, middlewares ...Middleware) {
	mapper.Map(Mapping{Method: Delete, Endpoint: endpoint, Builder: builder}, middlewares...)
}

func (mapper Mapper) Put(endpoint string, builder Builder, middlewares ...Middleware) {
	mapper.Map(Mapping{Method: Put, Endpoint: endpoint, Builder: builder}, middlewares...)
}
