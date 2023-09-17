package router

import (
	"github.com/Go-Lang/go-PostgreSQL/middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/stocks/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/stocks", middleware.GetAllStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/stocks/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/deletestocks/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")
	return router

}
