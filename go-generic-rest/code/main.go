package main

import (
	"fmt"
	"golang-crud-rest-api/controllers"
	"golang-crud-rest-api/entities"
	"golang-crud-rest-api/repos"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	LoadAppConfig()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	RegisterProductRoutes(router)
	RegisterBrandRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterProductRoutes(router *mux.Router) {
	var muxBase = "/api/products"
	router.HandleFunc(muxBase, controllers.GetProducts).Methods("GET")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controllers.GetProductById).Methods("GET")
	router.HandleFunc(muxBase, controllers.CreateProduct).Methods("POST")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controllers.UpdateProduct).Methods("PUT")
	router.HandleFunc(fmt.Sprintf("%s/{id}", muxBase), controllers.DeleteProduct).Methods("DELETE")
}

func RegisterBrandRoutes(router *mux.Router) {
	var brandRepo repos.GenericRepo[entities.Brand] = repos.NewBrandRepo()
	NewGenericRouter[entities.Brand, *repos.BrandRepo]("/api/brands", router, &brandRepo)
}
