package main

import (
	"216/internal/orchestrator/Controllers"
	"216/internal/orchestrator/Database"
	"216/internal/orchestrator/Middlewares"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/v1/login", Middlewares.SetMiddlewareJSON(Controllers.Login)).Methods("POST")
	router.HandleFunc("/api/v1/register", Middlewares.SetMiddlewareJSON(Controllers.CreateUser)).Methods("POST")

	router.HandleFunc("/api/v1/add/evaluation_arithmetic_expressions", Middlewares.SetMiddlewareAuthentication(Controllers.AddArithmeticExpressions)).Methods("POST")
	router.HandleFunc("/api/v1/arithmetic_expressions", Middlewares.SetMiddlewareAuthentication(Controllers.ArithmeticExpressionsList)).Methods("GET")
	router.HandleFunc("/api/v1/arithmetic_expression/{id}", Middlewares.SetMiddlewareAuthentication(Controllers.ArithmeticExpression)).Methods("GET")
	router.HandleFunc("/api/v1/arithmetic_operations", Middlewares.SetMiddlewareAuthentication(Controllers.AvailableArithmeticOperations)).Methods("GET")
	router.HandleFunc("/api/v1/task", Controllers.TaskForExecution).Methods("GET")
	router.HandleFunc("/api/v1/task/{id}", Controllers.TaskForExecutionPart).Methods("GET")

	router.HandleFunc("/api/v1/result/{id}", Middlewares.SetMiddlewareAuthentication(Controllers.SetResult)).Methods("PUT")
	router.HandleFunc("/api/v1/arithmetic_operations/{id}", Middlewares.SetMiddlewareAuthentication(Controllers.SetLeadTimeToArithmeticOperation)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/v1/available_calculators", Middlewares.SetMiddlewareAuthentication(Controllers.AvailableComputingResource)).Methods("GET")
	router.HandleFunc("/api/v1/profile", Middlewares.SetMiddlewareAuthentication(Controllers.Profile)).Methods("GET")

}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, token, x-request-id, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}
	//Redis.InitRedis(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD"))
	Database.Connect(os.Getenv("DB_CONNECTION_STRING"))
	Database.Migrate()
	Database.Seeder()
	//Rabbit.Connect(os.Getenv("RABBIT_STRING"))
	//exp, err := json.Marshal(Entities.ArithmeticExpressions{})
	//if err != nil {
	//	log.Println(err)
	//}
	//Rabbit.Send("Calculation", string(exp))

	router := mux.NewRouter().StrictSlash(true)
	RegisterRoutes(router)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), enableCORS(router)))
}
