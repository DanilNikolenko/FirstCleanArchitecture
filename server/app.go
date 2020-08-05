package server

import (
	"FirstCleanArchitecture/applications"
	appHTTP "FirstCleanArchitecture/applications/delivery/http"
	"FirstCleanArchitecture/applications/repository/storage"
	"FirstCleanArchitecture/applications/usecase"

	"github.com/labstack/gommon/log"

	"fmt"
	"net/http"
)

type App struct {
	httpServer *http.Server

	ApplicationsUC applications.UseCase
}

func NewApp() *App {
	app := storage.NewApplicationsRepository()
	return &App{
		ApplicationsUC: usecase.NewApplicationsUseCase(app),
	}
}

func (a *App) Run(port string) error {

	mux := http.NewServeMux()

	appHTTP.RegisterHTTPEndpoints(mux, a.ApplicationsUC)

	a.httpServer = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	fmt.Println("starting server at " + port)

	if err := a.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	return nil
}
