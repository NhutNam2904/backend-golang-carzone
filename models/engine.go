package models

import (
	"errors"

	"github.com/google/uuid"
)

type Engine struct {
	EngineID       uuid.UUID `json:"engine_id"`
	Displacement   int64     `json:"displacement"`
	NoOfCyclinders int64     `json:"noOfCyclinders"`
	CarRange       int64     `json:"carRange"`
}

type EngineRequest struct {
	Displacement   int64 `json:"displacement"`
	NoOfCyclinders int64 `json:"noOfCyclinders"`
	CarRange       int64 `json:"carRange"`
}

func ValidateEngineRequest(enginerequest EngineRequest) error {

	if err := validateDisplacement(enginerequest.Displacement); err != nil {
		return err
	}

	if err := validateCarRange(enginerequest.CarRange); err != nil {
		return err
	}

	if err := validateNoOfCyclinders(enginerequest.NoOfCyclinders); err != nil {
		return err
	}
	return nil

}

func validateDisplacement(displacement int64) error {

	if displacement <= 0 {
		return errors.New("Displacement must be greater than zero")
	}
	return nil
}

func validateNoOfCyclinders(noOfCyclinder int64) error {
	if noOfCyclinder <= 0 {
		return errors.New("noOfCyclinder is must be greater than 0")
	}
	return nil
}

func validateCarRange(carRange int64) error {
	if carRange <= 0 {
		return errors.New("CaRange is mus be than 0")
	}
	return nil
}
