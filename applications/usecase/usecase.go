package usecase

import (
	"FirstCleanArchitecture/applications"
	"context"
	"strconv"

	"github.com/labstack/gommon/log"
)

type ApplicationsUseCase struct {
	AppsStorage applications.ApplicationsRepository
}

func NewApplicationsUseCase(AR applications.ApplicationsRepository) *ApplicationsUseCase {
	return &ApplicationsUseCase{
		AppsStorage: AR,
	}
}

func (r *ApplicationsUseCase) GetApplication(ctx context.Context) (string, error) {
	//get application
	app, err := r.AppsStorage.GetRandomAliveApplication(ctx)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return app.Name, nil
}

func (r *ApplicationsUseCase) GetAdminApplications(ctx context.Context) ([]string, error) {
	// get applications
	Active, Cancel, err := r.AppsStorage.GetShowedAndCancelApplications(ctx)
	if err != nil {
		log.Error(err)
		return []string{}, err
	}

	// crete slice rez
	rez := make([]string, 0)

	for _, app := range Active {
		rez = append(rez, "Active("+app.Name+"-"+strconv.Itoa(app.Count)+")")
	}

	for _, app := range Cancel {
		rez = append(rez, "Cancel("+app.Name+"-"+strconv.Itoa(app.Count)+")")
	}

	return rez, nil
}
