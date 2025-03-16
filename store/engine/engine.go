package engine

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/NhutNam2904/carzone/models"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
)

type EngineStore struct {
	//dba,
	//dbb
	db *sql.DB
}

func New(db *sql.DB) EngineStore {
	return EngineStore{db: db}
}

func (e EngineStore) EngineById(ctx context.Context, id string) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "EngineByID-Store")

	defer span.End()
	var get_engine_byid models.Engine

	err := e.db.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id =$1", id).Scan(
		&get_engine_byid.EngineID,
		&get_engine_byid.Displacement,
		&get_engine_byid.NoOfCyclinders,
		&get_engine_byid.CarRange)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Engine{}, errors.New("Engine ID not found")
		}
		return models.Engine{}, err
	}

	return get_engine_byid, nil

}
func (e EngineStore) CreateEngine(ctx context.Context, engineReq *models.EngineRequest) (models.Engine, error) {
	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "CreateEngine-Store")

	defer span.End()

	tx, err := e.db.BeginTx(ctx, nil)

	if err != nil {
		return models.Engine{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	engineID := uuid.New()

	query := `INSERT INTO engine(id, displacement, no_of_cylinders,car_range) VALUES ($1, $2, $3, $4)
	          RETURNING id, displacement, no_of_cylinders ,car_range`

	_, err = tx.Exec(query,
		engineID,
		engineReq.Displacement,
		engineReq.NoOfCyclinders,
		engineReq.CarRange,
	)

	if err != nil {
		return models.Engine{}, err
	}

	engine_created := models.Engine{
		EngineID:       engineID,
		Displacement:   engineReq.Displacement,
		NoOfCyclinders: engineReq.NoOfCyclinders,
		CarRange:       engineReq.CarRange,
	}

	return engine_created, nil

}

func (e EngineStore) EngineUpdate(ctx context.Context, id string, engineReq *models.EngineRequest) (models.Engine, error) {

	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "UpdateEngine-Store")

	defer span.End()
	// Kiểm tra xem engine có tồn tại không
	var existingID uuid.UUID
	err := e.db.QueryRowContext(ctx, "SELECT id FROM engine WHERE id = $1", id).Scan(&existingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Engine{}, errors.New("Engine ID not found")
		}
		return models.Engine{}, err
	}

	// Bắt đầu transaction
	tx, err := e.db.BeginTx(ctx, nil)
	if err != nil {
		return models.Engine{}, err
	}

	// Deferred function quản lý transaction
	var txErr error
	defer func() {
		if txErr != nil {
			tx.Rollback()
		} else {
			txErr = tx.Commit()
		}
	}()

	// Câu lệnh SQL cập nhật
	query := `
        UPDATE engine 
        SET displacement = $1, no_of_cylinders = $2, car_range = $3, updated_at = $4
        WHERE id = $5
    `

	// Thực thi truy vấn
	results, txErr := tx.ExecContext(ctx, query,
		engineReq.Displacement,
		engineReq.NoOfCyclinders,
		engineReq.CarRange,
		time.Now(),
		id,
	)
	if txErr != nil {
		return models.Engine{}, txErr
	}

	// Kiểm tra số hàng bị ảnh hưởng
	rowsAffected, txErr := results.RowsAffected()
	if txErr != nil {
		return models.Engine{}, txErr
	}
	if rowsAffected == 0 {
		return models.Engine{}, errors.New("No engine updated")
	}

	// Tạo đối tượng trả về
	engine := models.Engine{
		EngineID:       existingID,
		Displacement:   engineReq.Displacement,
		NoOfCyclinders: engineReq.NoOfCyclinders,
		CarRange:       engineReq.CarRange,
	}

	return engine, nil
}

func (e EngineStore) DeleteEngine(ctx context.Context, id string) (models.Engine, error) {

	tracer := otel.Tracer("EngineStore")

	ctx, span := tracer.Start(ctx, "DeleteEngine-Store")

	defer span.End()
	var engine_deleted_byid models.Engine

	err := e.db.QueryRowContext(ctx, "SELECT id, displacement, no_of_cylinders, car_range FROM engine WHERE id =$1", id).Scan(&engine_deleted_byid.EngineID,
		&engine_deleted_byid.Displacement,
		&engine_deleted_byid.NoOfCyclinders,
		&engine_deleted_byid.CarRange)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Engine{}, errors.New("Engine ID not found")
		}
		return models.Engine{}, err
	}

	tx, err := e.db.BeginTx(ctx, nil)

	if err != nil {
		return models.Engine{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	query := `DELETE FROM engine WHERE id = $1`

	results, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		tx.Rollback()
		return models.Engine{}, err
	}

	rowefftected, err := results.RowsAffected()

	if err != nil {
		return models.Engine{}, err
	}

	if rowefftected == 0 {
		return models.Engine{}, errors.New("No engine deleted")
	}

	return engine_deleted_byid, nil

}
