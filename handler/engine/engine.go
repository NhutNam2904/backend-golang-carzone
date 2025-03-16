package engine

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

type EngineHandler struct {
	service service.EngineServiceInterface
}

func NewEngineHandler(service service.EngineServiceInterface) *EngineHandler {
	return &EngineHandler{
		service: service,
	}
}

func (e *EngineHandler) GetEngineByID(w http.ResponseWriter, r *http.Request) {

	tracer := otel.Tracer("EnginerHandler")

	ctx, span := tracer.Start(r.Context(), "GetEngineByID-Handler")

	defer span.End()

	params := mux.Vars(r)

	id := params["id"]

	getenginebyid, err := e.service.EngineById(ctx, id)

	if err != nil {
		log.Println("Error Get Engine by ID: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		errorMessage := fmt.Sprintf("Error Get Engine by ID: %s", err)

		_, _ = w.Write([]byte(errorMessage))

		return

	}

	responseBody, err := json.Marshal(getenginebyid)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseBody)

}

func (e *EngineHandler) CreateEngine(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EnginerHandler")

	ctx, span := tracer.Start(r.Context(), "CreateEngine-Handler")

	defer span.End()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error when read all body: ", err)
		w.WriteHeader(http.StatusInternalServerError)

		errorMessage := fmt.Sprintf("Error when read all body:%s", err)
		_, _ = w.Write([]byte(errorMessage))

		return
	}

	var engine models.EngineRequest

	err = json.Unmarshal(body, &engine)

	if err != nil {
		log.Println("Error while Unmarshalling Request body  ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	createdengine, err := e.service.CreateEngine(ctx, &engine)

	if err != nil {
		log.Println("Error Creating Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(createdengine)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(responseBody)

}

func (e *EngineHandler) EngineUpdate(w http.ResponseWriter, r *http.Request) {
	tracer := otel.Tracer("EnginerHandler")

	ctx, span := tracer.Start(r.Context(), "UpdateEngine-Handler")

	defer span.End()
	params := mux.Vars(r)

	id := params["id"]

	body, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("Error Reading Request Body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var engine models.EngineRequest

	err = json.Unmarshal(body, &engine)

	if err != nil {
		log.Println("Error while Unmarshalling Request body  ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	updatengine, err := e.service.EngineUpdate(ctx, id, &engine)

	if err != nil {
		log.Println("Error Updating Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(updatengine)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseBody)

}

func (e *EngineHandler) DeleteEngine(w http.ResponseWriter, r *http.Request) {

	tracer := otel.Tracer("EnginerHandler")

	ctx, span := tracer.Start(r.Context(), "DeleteEngine-Handler")

	defer span.End()
	params := mux.Vars(r)

	id := params["id"]

	enginedelete, err := e.service.DeleteEngine(ctx, id)

	if err != nil {
		log.Println("Error Deleting Engine: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	responseBody, err := json.Marshal(enginedelete)

	if err != nil {
		log.Println("Error while marshalling: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(responseBody)

}
