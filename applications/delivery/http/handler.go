package http

import (
	"ProjectCleanArchitecture/FirstCleanArchitecture/applications"
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/common/log"
)

type Handler struct {
	useCase applications.UseCase
}

func NewHandler(useCase applications.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetRequest(w http.ResponseWriter, r *http.Request) {
	s, err := h.useCase.GetApplication(context.Background())
	if err != nil {
		log.Error(err)
		fmt.Fprintln(w, "ERROR in GetRequest")
	}
	fmt.Fprintln(w, s)
}

func (h *Handler) GetAdminRequest(w http.ResponseWriter, r *http.Request) {
	s, err := h.useCase.GetAdminApplications(context.Background())
	if err != nil {
		log.Error(err)
		fmt.Fprintln(w, "ERROR in GetAdminRequest")
	}
	fmt.Fprintln(w, s)
}
