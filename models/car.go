package models

import (
	"errors"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Car struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Year      string    `json:"year"`
	Brand     string    `json:"brand"`
	FuelType  string    `json:"fuel_type"`
	Engine    Engine    `json:"engine"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CarRequest struct {
	Name     string  `json:"name"`
	Year     string  `json:"year"`
	Brand    string  `json:"brand"`
	FuelType string  `json:"fuel_type"`
	Engine   Engine  `json:"engine"`
	Price    float64 `json:"price"`
}

func ValidateCarRequest(carRequest CarRequest) error {

	if err := validateName(carRequest.Name); err != nil {
		return err
	}

	if err := validateBranch(carRequest.Brand); err != nil {
		return err
	}
	if err := validateYear(carRequest.Year); err != nil {
		return err
	}
	//if err := validateFueltype(carRequest.FuelType); err != nil {
	//	return err
	//}
	if err := validateEngine(carRequest.Engine); err != nil {
		return err
	}
	if err := validateCarprice(carRequest.Price); err != nil {
		return err
	}

	return nil

}

func validateName(name string) error {
	if name == "" {
		return errors.New("Name is Required")

	}
	return nil
}

func validateYear(year string) error {
	if year == "" {
		return errors.New("Year is Required")
	}
	_, err := strconv.Atoi(year)
	if err != nil {
		return errors.New("Year is must be a valid  number")
	}
	currentYear := time.Now().Year()
	yearInt, _ := strconv.Atoi(year)
	if yearInt < 1886 || yearInt > currentYear {
		return errors.New("Year must be between 1886 and current year")
	}
	return nil
}

func validateBranch(branch string) error {

	if branch == "" {
		return errors.New("Branch is Required")
	}
	return nil

}

func validateFueltype(fueltype string) error {

	fueltypes := []string{"Persol", "Diesel", "Electric", "Hybrid"}

	for _, a := range fueltypes {
		if fueltype == a {
			return nil
		}
	}

	return errors.New("FuelType in: Persol, Diesel, Electric, Hybrid")
}

func validateEngine(engine Engine) error {
	if engine.EngineID == uuid.Nil {
		return errors.New("EngineID is Required")
	}
	if engine.Displacement <= 0 {
		return errors.New("Displacement must be greater than zero")
	}

	if engine.NoOfCyclinders <= 0 {
		return errors.New("noOfCyclinders must be greater than zero")
	}

	if engine.CarRange <= 0 {
		return errors.New("carRange must be greater than zero")
	}

	return nil

}

func validateCarprice(carprice float64) error {
	if carprice <= 0 {
		return errors.New("carPrice is  unvalid")
	}
	return nil

}
