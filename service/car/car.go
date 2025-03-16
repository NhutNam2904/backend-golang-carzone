package car

import (
	"context"

	"github.com/NhutNam2904/carzone/models"
	"github.com/NhutNam2904/carzone/store"
	"go.opentelemetry.io/otel"
)

type CarService struct {
	store store.CarStoreInterface
}

func NewCarService(store store.CarStoreInterface) CarService {
	return CarService{
		store: store,
	}
}

func (s CarService) GetCarById(ctx context.Context, id string) (*models.Car, error) {

	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "GetCarByID-Service")

	defer span.End()

	car, err := s.store.GetCarById(ctx, id)

	//fmt.Printf("Car in layer service: %v", car)
	return car, err

}

func (s CarService) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {

	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "GetCarByBrand-Service")

	defer span.End()

	cars, err := s.store.GetCarByBrand(ctx, brand, isEngine)

	if err != nil {
		return nil, err
	}
	return cars, nil

}

func (s CarService) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {

	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "CreateCar-Service")

	defer span.End()

	err := models.ValidateCarRequest(*carReq)

	if err != nil {
		return models.Car{}, err
	}

	car, err := s.store.CreateCar(ctx, carReq)

	if err != nil {
		return models.Car{}, err
	}
	return car, err
}

func (s CarService) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "DeleteCar-Service")

	defer span.End()
	car, err := s.store.DeleteCar(ctx, id)

	if err != nil {
		return models.Car{}, err
	}

	return car, err
}

func (s CarService) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {

	tracer := otel.Tracer("CarService")

	ctx, span := tracer.Start(ctx, "UpdateCar-Service")

	defer span.End()
	car, err := s.store.UpdateCar(ctx, id, carReq)

	if err != nil {
		return models.Car{}, err
	}

	return car, err
}
