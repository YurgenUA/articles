package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/yurgenua/golang-crud-rest-api/entities"
	"github.com/yurgenua/golang-crud-rest-api/repos"

	"github.com/gorilla/mux"
)

func main() {
	LoadAppConfig()

	// Initialize the router
	router := mux.NewRouter().StrictSlash(true)

	RegisterBrandRoutes(router)

	// Start the server
	log.Println(fmt.Sprintf("Starting Server on port %s", AppConfig.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", AppConfig.Port), router))
}

func RegisterBrandRoutes(router *mux.Router) {
	var brandRepo repos.GenericRepo[entities.Brand] = repos.NewBrandRepo()
	NewGenericRouter[entities.Brand, *repos.BrandRepo]("/api/brands", router, &brandRepo)
}
