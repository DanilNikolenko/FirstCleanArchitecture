package applications

import "context"

type UseCase interface {
	GetApplication(ctx context.Context) (string, error)
	GetAdminApplications(ctx context.Context) ([]string, error)
}
