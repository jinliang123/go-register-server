package router

import (
	"time"

	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/choerodon/go-register-server/pkg/eureka/apps"
	"github.com/choerodon/go-register-server/pkg/eureka/metrics"
	"github.com/choerodon/go-register-server/pkg/eureka/repository"
)

const (
	APIPath = "/eureka"
)

type EurekaAppsService struct {
	appRepo *repository.ApplicationRepository
}

func NewEurekaAppsService(appRepo *repository.ApplicationRepository) *EurekaAppsService {
	s := &EurekaAppsService{
		appRepo: appRepo,
	}

	return s
}

func (es *EurekaAppsService) Register() {
	glog.Info("Register eureka app APIs")

	ws := new(restful.WebService)
	ws.Path(APIPath).Produces(restful.MIME_JSON, restful.MIME_XML)

	// GET /eureka/apps
	ws.Route(ws.GET("/apps").To(es.listEurekaApps).
		Doc("Get all apps")).Produces("application/json")

	ws.Route(ws.GET("/apps/delta").To(es.listEurekaAppsDelta).
		Doc("Get all apps delta")).Produces("application/json")

	ws.Route(ws.POST("/apps/{app-name}").To(es.registerEurekaApp).
		Doc("get a user").Produces("application/json").
		Param(ws.PathParameter("app-name", "app name").DataType("string")))

	ws.Route(ws.PUT("/apps/{app-name}/{instance-id}").To(es.renew).
		Doc("renew").
		Param(ws.PathParameter("app-name", "app name").DataType("string")).
		Param(ws.PathParameter("instance-id", "instance id").DataType("string")))
	restful.Add(ws)
}

// listEurekaApps handles the request to list eureka apps.
func (es *EurekaAppsService) listEurekaApps(request *restful.Request, response *restful.Response) {
	start := time.Now()

	metrics.RequestCount.With(prometheus.Labels{"path": request.Request.RequestURI}).Inc()
	applicationResources := es.appRepo.GetApplicationResources()
	response.WriteAsJson(applicationResources)

	finish := time.Now()
	cost := finish.Sub(start).Nanoseconds()

	metrics.FetchProcessTime.Set(float64(cost))
}
func (es *EurekaAppsService) listEurekaAppsDelta(request *restful.Request, response *restful.Response) {
	metrics.RequestCount.With(prometheus.Labels{"path": request.Request.RequestURI}).Inc()
	applicationResources := &apps.ApplicationResources{
		Applications: &apps.Applications{
			VersionsDelta:   2,
			AppsHashcode:    "app_hashcode",
			ApplicationList: make([]*apps.Application, 0),
		},
	}
	response.WriteAsJson(applicationResources)
}

func (es *EurekaAppsService) registerEurekaApp(request *restful.Request, response *restful.Response) {
	metrics.RequestCount.With(prometheus.Labels{"path": request.Request.RequestURI}).Inc()
	glog.Info("Receive registry from ", request.PathParameter("app-name"))
}

func (es *EurekaAppsService) renew(request *restful.Request, response *restful.Response) {
	metrics.RequestCount.With(prometheus.Labels{"path": request.Request.RequestURI}).Inc()
	//appName := request.PathParameter("app-name")
	//instanceId := request.PathParameter("instance-id")
	//response.WriteAsJson(es.appRepo.Renew(appName, instanceId))
}
