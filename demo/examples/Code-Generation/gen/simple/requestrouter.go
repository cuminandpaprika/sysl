// Code generated by sysl DO NOT EDIT.
package simple

import (
	"context"
	"net/http"

	"github.com/anz-bank/sysl-go/common"
	"github.com/anz-bank/sysl-go/core"
	"github.com/anz-bank/sysl-go/handlerinitialiser"
	"github.com/anz-bank/sysl-go/validator"
	"github.com/go-chi/chi"
)

// GenCallback callbacks used by the generated code
type GenCallback interface {
	AddMiddleware(ctx context.Context, r chi.Router)
	BasePath() string
	Config() validator.Validator
	HandleError(ctx context.Context, w http.ResponseWriter, kind common.Kind, message string, cause error)
	DownstreamTimeoutContext(ctx context.Context) (context.Context, context.CancelFunc)
}

// Router interface for simple
type Router interface {
	Route(router *chi.Mux)
}

// ServiceRouter for simple API
type ServiceRouter struct {
	gc               GenCallback
	svcHandler       *ServiceHandler
	basePathFromSpec string
}

// swaggerFile is a struct to store the swagger file system
type swaggerFile struct {
	file http.FileSystem
}

// swagger will receive the embedded swagger file if it is generated by the resource application
var swagger = swaggerFile{}

// NewServiceRouter creates a new service router for simple
func NewServiceRouter(gc GenCallback, svcHandler *ServiceHandler) handlerinitialiser.HandlerInitialiser {
	return &ServiceRouter{gc, svcHandler, ""}
}

// WireRoutes ...
//nolint:funlen
func (s *ServiceRouter) WireRoutes(ctx context.Context, r chi.Router) {
	r.Route(core.SelectBasePath(s.basePathFromSpec, s.gc.BasePath()), func(r chi.Router) {
		s.gc.AddMiddleware(ctx, r)
		r.Get("/", s.svcHandler.GetHandler)
		r.Get("/foobar", s.svcHandler.GetFoobarListHandler)
	})
}

// Config ...
func (s *ServiceRouter) Config() interface{} {
	return s.gc.Config()
}

// Name ...
func (s *ServiceRouter) Name() string {
	return "simple"
}
