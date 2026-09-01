package query

import (
	"context"
)

type GetApplication struct {
	ApplicationUUID string
}

type GetApplicationHandler interface {
	Handle(ctx context.Context, query GetApplication) (Application, error)
}

type ApplicationReadModel interface {
	GetApplication(ctx context.Context, uuid string) (Application, error)
}

type getApplicationHandler struct {
	 readModel ApplicationReadModel
}

func NewGetApplicationHandler(readModel ApplicationReadModel) GetApplicationHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return getApplicationHandler{readModel: readModel}
}

func (h getApplicationHandler) Handle(ctx context.Context, query GetApplication) (Application, error) {
	return h.readModel.GetApplication(ctx, query.ApplicationUUID)
}
