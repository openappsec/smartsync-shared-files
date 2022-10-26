package rest

import (
	"context"
	"net/http"

	"openappsec.io/smartsync-shared-files/internal/app/utils"

	"github.com/go-chi/chi"
	"openappsec.io/ctxutils"
	healthhandlers "openappsec.io/health/http/rest"
	"openappsec.io/httputils/middleware"
	"openappsec.io/log"
)

const (
	// keys in the error-responses.json file which contains the proper error responses bodies
	timeoutErrorBodyKey     = "timeout-error"
	noAgentIDErrorBodyKey   = "no-agent-id-error"
	noProfileIDErrorBodyKey = "no-profile-id-error"
	noTenantIDErrorBodyKey  = "no-tenant-id-error"
)

// newRouter returns a router including method, path, name, and handler
// the router object will parse the request and pass it to the proper function
func (a *Adapter) newRouter(ctx context.Context) *chi.Mux {
	router := chi.NewRouter()
	log.AddContextField(ctxutils.ContextKeyCallingService)

	// load error bodies from configs/error-responses.json
	errorBodyTimeout := utils.CreateErrorBody(ctx, timeoutErrorBodyKey)
	errorBodyAgentID := utils.CreateErrorBody(ctx, noAgentIDErrorBodyKey)
	errorBodyProfileID := utils.CreateErrorBody(ctx, noProfileIDErrorBodyKey)
	errorTenantBodyID := utils.CreateErrorBody(ctx, noTenantIDErrorBodyKey)

	// set server timeout on requests
	// in this project it is set to 15 seconds
	// look at the NewHTTPAdapter function in server.go which loads
	// the timeout as an environment variable
	router.Use(middleware.Timeout(a.wait, errorBodyTimeout))

	router.Group(func(router chi.Router) {

		// k8s automatically does a "health check" for us upon deploying a service to the cluster
		// this check's purpose is to avoid deploying a service which failed to initialize it's code
		// you can look at the /deployments/k8s/helm-chart/templates/deployment.yaml file how we configured it
		// the official documentation: https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
		router.Route("/health", func(r chi.Router) {
			r.Get("/live", healthhandlers.LivenessHandler(a.healthSvc).ServeHTTP)
			r.Get("/ready", healthhandlers.ReadinessHandler(a.healthSvc).ServeHTTP)
		})
	})

	defaultErrorBody := utils.CreateErrorBody(ctx, "default-error")
	router.Group(func(router chi.Router) {

		router.Route("/api", func(r chi.Router) {
			// Logs "new incoming request" upon receiving the request
			// Logs the request duration after returning a response
			r.Use(middleware.Logging(defaultErrorBody))
			r.Use(middleware.Tracing)

			// create middlewares that will parse the headers (in this case x-tenant-id, x-profile-id and x-agent-id)
			// and save them to the context. To extract them (usually done in the handler) - use ExtractString function from ctxutils package
			// you can remove all/some of them and create your own middlewares
			r.Use(middleware.TenantID(errorTenantBodyID)) // returns error if the request doesn't include "x-tenant-id" as header
			r.Use(func(next http.Handler) http.Handler {
				return middleware.HeaderToContext(next, "X-Agent-Id", ctxutils.ContextKeyAgentID, false,
					errorBodyAgentID)
			})
			r.Use(func(next http.Handler) http.Handler {
				return middleware.HeaderToContext(next, "X-Profile-Id", ctxutils.ContextKeyProfileID, false,
					errorBodyProfileID)
			})
			r.Use(middleware.CorrelationID(defaultErrorBody))  // search for header "x-trace-id"
			r.Use(middleware.CallingService(defaultErrorBody)) // search for optional header "X-Calling-Service"

			r.Get("/", a.GetFilesList)
			r.Get("/*", a.GetFile)
			r.Put("/*", a.PutFile)
		})
	})

	return router
}
