package http

import (
	"ProjectCleanArchitecture/FirstCleanArchitecture/applications"
	"net/http"
)

func RegisterHTTPEndpoints(router *http.ServeMux, uc applications.UseCase) {
	h := NewHandler(uc)

	router.HandleFunc("/request", h.GetRequest)
	router.HandleFunc("/admin/request", h.GetAdminRequest)
}
