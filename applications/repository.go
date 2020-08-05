package applications

import (
	"FirstCleanArchitecture/models"
	"context"
)

type ApplicationsRepository interface {
	GetShowedAndCancelApplications(ctx context.Context) ([]models.Application, []models.Application, error)
	GetRandomAliveApplication(ctx context.Context) (models.Application, error)
}
