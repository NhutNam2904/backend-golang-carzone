package car

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/NhutNam2904/carzone/models"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type Store struct {
	db          *sql.DB
	redisClient *redis.Client
}

func New(db *sql.DB, redisClient *redis.Client) *Store {
	return &Store{db: db,
		redisClient: redisClient}
}

func (s Store) GetCarById(ctx context.Context, id string) (*models.Car, error) {
	tracer := otel.Tracer("CarStore")

	ctx, span := tracer.Start(ctx, "GetCarByID-Store")

	defer span.End()

	var car models.Car

	query := `SELECT c.id, c.name,c.year,c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range FROM car c JOIN 
	engine e ON c.engine_id = e.id WHERE c.id = $1`

	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&car.ID,
		&car.Name,
		&car.Year,
		&car.Brand,
		&car.FuelType,
		&car.Engine.EngineID,
		&car.Price,
		&car.CreatedAt,
		&car.UpdatedAt,
		&car.Engine.EngineID,
		&car.Engine.Displacement,
		&car.Engine.NoOfCyclinders,
		&car.Engine.CarRange)

	if err != nil {
		if err == sql.ErrNoRows {
			return &car, nil
		}
		return &car, err
	}
	//fmt.Print(&car)
	return &car, nil

}

func (s Store) GetCarByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {

	tracer := otel.Tracer("CarStore")

	ctx, span := tracer.Start(ctx, "GetCarByBrand-Store")

	defer span.End()
	var cars []models.Car
	var query string

	key := fmt.Sprintf("Brand:%s", brand)
	cachedData, err := s.redisClient.Get(ctx, key).Result()

	if err == nil {

		if err := json.Unmarshal([]byte(cachedData), &cars); err == nil {
			log.Println("Cache hit")
			return cars, nil
		}

	}

	log.Println("Cache miss, querying database")

	if isEngine {
		query = `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.engine_id, c.price, c.created_at, c.updated_at, e.id, e.displacement, e.no_of_cylinders, e.car_range 
				FROM car c 
				JOIN engine e ON c.engine_id = e.id 
				WHERE c.brand = $1`
	} else {
		query = `SELECT c.id, c.name, c.year, c.brand, c.fuel_type, c.price, c.created_at, c.updated_at 
				FROM car c 
				WHERE c.brand = $1`
	}

	rows, err := s.db.QueryContext(ctx, query, brand)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		if isEngine {
			var engine models.Engine
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Engine.EngineID,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
				&engine.EngineID,
				&engine.Displacement,
				&engine.NoOfCyclinders,
				&engine.CarRange,
			)
			if err != nil {
				return nil, err
			}

			car.Engine = engine
		} else {
			err := rows.Scan(
				&car.ID,
				&car.Name,
				&car.Year,
				&car.Brand,
				&car.FuelType,
				&car.Price,
				&car.CreatedAt,
				&car.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
		}
		cars = append(cars, car)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	jsonData, err := json.Marshal(cars)

	if err != nil {
		log.Println("Cannot marshar json to cache data to redis")
	}

	err = s.redisClient.Set(ctx, key, jsonData, 60).Err()

	if err != nil {
		log.Println("Failed to set data to redis: ", err)
	}

	return cars, nil
}

func (s Store) CreateCar(ctx context.Context, carReq *models.CarRequest) (models.Car, error) {

	tracer := otel.Tracer("CarStore")

	ctx, span := tracer.Start(ctx, "CreateCar-Store")

	defer span.End()

	var createdCar models.Car

	var engineID uuid.UUID

	err := s.db.QueryRowContext(ctx, "SELECT id FROM engine WHERE id=$1", carReq.Engine.EngineID).Scan(&engineID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return createdCar, errors.New("Engine ID does not exists in the engine table")
		}
		return createdCar, err

	}

	carID := uuid.New()
	createdAt := time.Now()

	updatedAt := createdAt

	newCar := models.Car{
		ID:        carID,
		Name:      carReq.Name,
		Year:      carReq.Year,
		Brand:     carReq.Brand,
		FuelType:  carReq.FuelType,
		Engine:    carReq.Engine,
		Price:     carReq.Price,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	tx, err := s.db.BeginTx(ctx, nil)

	if err != nil {
		return createdCar, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `INSERT INTO car(id, name, year,brand,fuel_type, engine_id, price, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	          RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		newCar.ID,
		newCar.Name,
		newCar.Year,
		newCar.Brand,
		newCar.FuelType,
		newCar.Engine.EngineID,
		newCar.Price,
		newCar.CreatedAt,
		newCar.UpdatedAt,
	).Scan(&createdCar.ID,
		&createdCar.Name,
		&createdCar.Year,
		&createdCar.Brand,
		&createdCar.FuelType,
		&createdCar.Engine.EngineID,
		&createdCar.Price,
		&createdCar.CreatedAt,
		&createdCar.UpdatedAt,
	)
	if err != nil {
		return createdCar, err
	}

	return createdCar, nil

}

func (s Store) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	tracer := otel.Tracer("CarStore")

	ctx, span := tracer.Start(ctx, "DeleteCar-Store")

	defer span.End()

	var deleteCar models.Car

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return deleteCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = s.db.QueryRowContext(ctx,
		"SELECT id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at FROM car WHERE id = $1", id).
		Scan(
			&deleteCar.ID,
			&deleteCar.Name,
			&deleteCar.Year,
			&deleteCar.Brand,
			&deleteCar.FuelType,
			&deleteCar.Engine.EngineID, // Truyền vào trường Engine.EngineID
			&deleteCar.Price,
			&deleteCar.CreatedAt,
			&deleteCar.UpdatedAt,
		)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Car{}, errors.New("Car not found")
		}
		return models.Car{}, err
	}

	result, err := tx.ExecContext(ctx, "DELETE FROM car WHERE id =$1", id)

	if err != nil {
		return models.Car{}, err
	}
	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return models.Car{}, err
	}

	if rowsAffected == 0 {
		return models.Car{}, errors.New("No rows were affected")
	}

	return deleteCar, nil

}

func (s Store) UpdateCar(ctx context.Context, id string, carReq *models.CarRequest) (models.Car, error) {
	tracer := otel.Tracer("CarStore")

	ctx, span := tracer.Start(ctx, "UpdateCar-Store")

	defer span.End()

	var updatedCar models.Car

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return updatedCar, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `
	WITH updated_car AS (
    UPDATE car
    SET name = $2, 
        year = $3, 
        brand = $4, 
        fuel_type = $5, 
        engine_id = $6, 
        price = $7, 
        updated_at = $8
    WHERE id = $1
RETURNING id, name, year, brand, fuel_type, engine_id, price, created_at, updated_at
)
SELECT 
    updated_car.id, 
	updated_car.name, 
    updated_car.year, 
    updated_car.brand, 
    updated_car.fuel_type, 
    updated_car.engine_id, 
    updated_car.price, 
    updated_car.created_at, 
    updated_car.updated_at,
    engine.displacement, 
    engine.no_of_cylinders, 
    engine.car_range
FROM 
    updated_car
JOIN 
    engine ON updated_car.engine_id = engine.id;
 `

	err = tx.QueryRowContext(ctx, query,
		id,
		carReq.Name,
		carReq.Year,
		carReq.Brand,
		carReq.FuelType,
		carReq.Engine.EngineID,
		carReq.Price,
		time.Now(),
	).Scan(&updatedCar.ID,
		&updatedCar.Name,
		&updatedCar.Year,
		&updatedCar.Brand,
		&updatedCar.FuelType,
		&updatedCar.Engine.EngineID,
		&updatedCar.Price,
		&updatedCar.CreatedAt,
		&updatedCar.UpdatedAt,
		&updatedCar.Engine.Displacement,
		&updatedCar.Engine.NoOfCyclinders,
		&updatedCar.Engine.CarRange,
	)
	if err != nil {
		return updatedCar, err
	}

	return updatedCar, nil

}
