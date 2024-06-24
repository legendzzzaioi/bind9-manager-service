// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	config "bind9-manager-service/internal/handler/config"
	record "bind9-manager-service/internal/handler/record"
	zone "bind9-manager-service/internal/handler/zone"
	"bind9-manager-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/config",
				Handler: config.GetConfigHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/config",
				Handler: config.UpdateConfigHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/records",
				Handler: record.GetRecordsHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/records",
				Handler: record.CreateRecordHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/records",
				Handler: record.UpdateRecordHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/records",
				Handler: record.DeleteRecordHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/zones",
				Handler: zone.GetZonesHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/zones",
				Handler: zone.CreateZoneHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/zones",
				Handler: zone.UpdateZoneHandler(serverCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/zones",
				Handler: zone.DeleteZoneHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1"),
	)
}
