package car

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/NhutNam2904/carzone/models"
	"github.com/NhutNam2904/carzone/service"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
)

type CarHandler struct {
	service service.CarServiceInterface
}

func NewCarHandler(service service.CarServiceInterface) CarHandler {
	return CarHandler{service: service}
}

func (h *CarHandler) GetCarByID(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "GetCarByID-Handler")

	defer span.End()

	//ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.GetCarById(ctx, id)
	//fmt.Print(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}

	bodyresponse, err := json.Marshal(res) // byte

	fmt.Print(string(bodyresponse))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(bodyresponse)

	if err != nil {
		log.Println("Error Writing Response: ", err)

	}

}

func (h *CarHandler) GetCarByBrand(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "GetCarByBrand-Handler")

	defer span.End()

	brand := r.URL.Query().Get("brand")
	isEngine := r.URL.Query().Get("isEngine") == "true"

	res, err := h.service.GetCarByBrand(ctx, brand, isEngine)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return
	}
	body, err := json.Marshal(res)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error: ", err)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(body)

	if err != nil {
		log.Println("Error Writing Response: ", err)
	}

	log.Printf("[INFO] Success: GetCarByBrand - Brand: %s, IsEngine: %v, ResultCount: %d", brand, isEngine, len(res))
}

func (h *CarHandler) CreateCar(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "CreateCar-Handler")

	defer span.End()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var carReq models.CarRequest

	err = json.Unmarshal(body, &carReq)

	if err != nil {
		log.Println("Error while Unmarshalling Request body  ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdCar, err := h.service.CreateCar(ctx, &carReq)

	if err != nil {
		log.Println("Error Creating Car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(createdCar)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(responseBody)

}

func (h *CarHandler) UpdateCar(w http.ResponseWriter, r *http.Request) {

	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "UpdateCar-Handler")

	defer span.End()

	params := mux.Vars(r)

	id := params["id"]

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error Reading Request Body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var carReq models.CarRequest

	err = json.Unmarshal(body, &carReq)

	if err != nil {
		log.Println("Error while Unmarshalling Request body  ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatecar, err := h.service.UpdateCar(ctx, id, &carReq)

	if err != nil {
		log.Println("Error Updating Car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(updatecar)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseBody)

}

func (h *CarHandler) DeleteCar(w http.ResponseWriter, r *http.Request) {

	tracer := otel.Tracer("CarHandler")

	ctx, span := tracer.Start(r.Context(), "Delete-Handler")

	defer span.End()

	params := mux.Vars(r)

	id := params["id"]

	cardelete, err := h.service.DeleteCar(ctx, id)

	if err != nil {
		log.Println("Error Deleting Car: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(cardelete)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseBody)

	log.Printf("[INFO] Success: Delete Car By ID: %s", id)

}
