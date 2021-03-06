package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/riclava/dds/cluster/friends"

	restful "github.com/emicklei/go-restful"
	"github.com/riclava/dds/api/controller"
	"github.com/riclava/dds/api/models"
	"github.com/riclava/dds/cluster/config"
	"github.com/riclava/dds/cluster/tasks"
)

const (
	// RequestLogString is a template for request log message.
	RequestLogString = "[%s] Incoming %s %s %s request from: %s"

	// ResponseLogString is a template for response log message.
	ResponseLogString = "[%s] Outcoming response to %s with %d status code"
)

// APIHandler is a representation of API handler
type APIHandler struct {
	Config  *config.Config
	Friends *friends.Friends
}

// CreateAPIHandler create an API handler for restful API
func CreateAPIHandler(cfg *config.Config, frands *friends.Friends) (http.Handler, error) {

	container := restful.NewContainer()
	apiHandler := APIHandler{
		Config:  cfg,
		Friends: frands,
	}

	webService := new(restful.WebService)
	webService.Filter(logRequestAndResponse)

	webService.Path("/api/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	container.Add(webService)

	webService.Route(
		webService.GET("/").To(apiHandler.handleGetRoot).Writes(models.Response{}))
	webService.Route(
		webService.GET("/task").To(apiHandler.handleTaskList).Writes(tasks.HTTPTasks{}))
	webService.Route(
		webService.POST("/task").To(apiHandler.handleTaskAdd).Writes(models.Response{}))
	webService.Route(
		webService.POST("/friend").To(apiHandler.handleFriendPost).Writes(models.Response{}))
	webService.Route(
		webService.DELETE("/friend").To(apiHandler.handleFriendDelete).Writes(models.Response{}))

	return container, nil
}

// logRequestAndReponse is a web-service filter function used for request and response logging.
func logRequestAndResponse(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	log.Printf(formatRequestLog(request))
	chain.ProcessFilter(request, response)
	log.Printf(formatResponseLog(response, request))
}

// formatRequestLog formats request log string.
func formatRequestLog(request *restful.Request) string {
	uri := ""
	if request.Request.URL != nil {
		uri = request.Request.URL.RequestURI()
	}

	return fmt.Sprintf(RequestLogString, time.Now().Format(time.RFC3339), request.Request.Proto,
		request.Request.Method, uri, request.Request.RemoteAddr)
}

// formatResponseLog formats response log string.
func formatResponseLog(response *restful.Response, request *restful.Request) string {
	return fmt.Sprintf(ResponseLogString, time.Now().Format(time.RFC3339),
		request.Request.RemoteAddr, response.StatusCode())
}

func (apiHandler *APIHandler) handleGetRoot(request *restful.Request, response *restful.Response) {
	controller.Index(request, response)
}

func (apiHandler *APIHandler) handleTaskList(request *restful.Request, response *restful.Response) {
	controller.ListTask(request, response)
}

func (apiHandler *APIHandler) handleTaskAdd(request *restful.Request, response *restful.Response) {
	controller.AddTask(request, response, apiHandler.Config)
}

func (apiHandler *APIHandler) handleFriendPost(request *restful.Request, response *restful.Response) {
	controller.AddFriend(request, response, apiHandler.Friends)
}

func (apiHandler *APIHandler) handleFriendDelete(request *restful.Request, response *restful.Response) {
	controller.DeleteFriend(request, response, apiHandler.Friends)
}
