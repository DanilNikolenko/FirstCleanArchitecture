package applications

import (
	"ProjectCleanArchitecture/FirstCleanArchitecture/models"
	"context"
)

type ApplicationsRepository interface {
	GetShowedAndCancelApplications(ctx context.Context) ([]models.Application, error)
	GetRandomAliveApplication(ctx context.Context) (models.Application, error)
}
