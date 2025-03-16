package engine

import (
	"context"

	"github.com/NhutNam2904/carzone/models"
	"github.com/NhutNam2904/carzone/store"
	"go.opentelemetry.io/otel"
)

type EngineService struct {
	store store.EngineStoreInterface
}

func NewEngineService(store store.EngineStoreInterface) EngineService {
	return EngineService{store: store}
}

func (s EngineService) EngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "EngineByID-Service")

	defer span.End()
	engine, err := s.store.EngineById(ctx, id)
	if err != nil {
		return models.Engine{}, err
	}

	return engine, nil
}

func (s EngineService) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "CreateEngine-Service")

	defer span.End()
	engine, err := s.store.CreateEngine(ctx, engineReq)
	if err != nil {
		return models.Engine{}, err
	}

	return engine, nil
}

func (s EngineService) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "EngineUpdate-Service")

	defer span.End()

	engine, err := s.store.EngineUpdate(ctx, id, engineReq)

	if err != nil {
		return models.Engine{}, err
	}

	return engine, nil
}

func (s EngineService) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineService")

	ctx, span := tracer.Start(ctx, "DeleteEngine-Service")

	defer span.End()
	engine, err := s.store.DeleteEngine(ctx, id)

	if err != nil {
		return models.Engine{}, err
	}
	return engine, nil
}
