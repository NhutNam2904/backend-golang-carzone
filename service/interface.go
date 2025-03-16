package service

import (
	"context"

	"github.com/NhutNam2904/carzone/models"
)

type CarServiceInterface interface {
	GetCarById(ctx context.Context, id string) (*models.Car, error)
	GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error)
	DeleteCar(ctx context.Context, id string) (models.Car, error)
	UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error)
}

type EngineServiceInterface interface {
	EngineById(ctx context.Context, id string) (models.Engine, error)
	CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error)
	EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error)
	DeleteEngine(ctx context.Context, id string) (models.Engine, error)
}

//type LoginServiceInterface interface {
//GetUsernamePassword(ctx context.Context, username string) (models.Credentials, error)
//}

type UserServiceInteface interface {
	SignUp(ctx context.Context, user *models.User) (models.User, error)
}
