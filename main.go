package main

import (
	"context"
	"database/sql"

	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NhutNam2904/carzone/driver"
	"github.com/gorilla/mux"

	carHandler "github.com/NhutNam2904/carzone/handler/car"
	engineHandler "github.com/NhutNam2904/carzone/handler/engine"

	//loginHandler "github.com/NhutNam2904/carzone/handler/login"

	//middleware "github.com/NhutNam2904/carzone/middleware"
	carService "github.com/NhutNam2904/carzone/service/car"
	engineService "github.com/NhutNam2904/carzone/service/engine"
	carStore "github.com/NhutNam2904/carzone/store/car"
	engineStore "github.com/NhutNam2904/carzone/store/engine"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	traceProvider, err := startTracing()

	if err != nil {
		log.Fatal("Failed to Start Tracing : %v", err)
	}

	defer func() {

		if err := traceProvider.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to Shut Down Tracing: %v ", err)
		}
	}()
	driver.InitRedis()

	rd := driver.GetRedis()
	driver.StartUpDB()
	defer driver.CloseDB()

	db := driver.GetDBCarManageMent()

	//loginStore := loginStore.New(db, rd)
	//loginService := loginService.NewLoginService(loginStore)

	carStore := carStore.New(db, rd)
	carService := carService.NewCarService(carStore)

	engineStore := engineStore.New(db)
	engineService := engineService.NewEngineService(engineStore)

	carHandler := carHandler.NewCarHandler(carService)
	engineHandler := engineHandler.NewEngineHandler(engineService)
	//loginHandler := loginHandler.NewLoginHandler(loginService)

	router := mux.NewRouter()

	router.Use(otelmux.Middleware("CarZone"))

	schemaFile := "store/schema.sql"

	if err := excuteSchemaFile(db, schemaFile); err != nil {
		log.Fatal("Error while executing the schema file: ", err)
	}

	//router.HandleFunc("/login", loginHandler.LoginHandlerUsernamePassowrd).Methods("POST")

	//router := router.PathPrefix("/").Subrouter()
	//router.Use(middleware.AuthMiddleware)

	router.HandleFunc("/cars/{id}", carHandler.GetCarByID).Methods("GET")
	router.HandleFunc("/cars", carHandler.GetCarByBrand).Methods("GET")
	router.HandleFunc("/cars", carHandler.CreateCar).Methods("POST")
	router.HandleFunc("/cars/{id}", carHandler.UpdateCar).Methods("PUT")
	router.HandleFunc("/cars/{id}", carHandler.DeleteCar).Methods("DELETE")

	router.HandleFunc("/engine/{id}", engineHandler.GetEngineByID).Methods("GET")
	router.HandleFunc("/engine", engineHandler.CreateEngine).Methods("POST")
	router.HandleFunc("/engine/{id}", engineHandler.EngineUpdate).Methods("PUT")
	router.HandleFunc("/engine/{id}", engineHandler.DeleteEngine).Methods("DELETE")

	//

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server is running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}

func excuteSchemaFile(db *sql.DB, fileName string) error {
	sqlFile, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlFile))

	if err != nil {

		return err

	}

	return nil
}

func startTracing() (*trace.TracerProvider, error) {
	header := map[string]string{
		"Content-Type": "application/json",
	}

	// Tạo OTLP exporter với Jaeger làm backend
	exporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("http://192.168.249.100:4318"), // Đảm bảo endpoint đúng
			otlptracehttp.WithHeaders(header),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Error creating exporter: %w", err)
	}

	// Tạo TracerProvider
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(
			exporter,
			trace.WithMaxExportBatchSize(trace.DefaultMaxExportBatchSize),
			trace.WithBatchTimeout(trace.DefaultScheduleDelay),
		),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("CarZone"),
			),
		),
	)

	log.Println("Successfully Started Tracing")

	return tracerProvider, nil

}
