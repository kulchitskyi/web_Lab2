package main

import (
	"os"
	"log"
	"net/http"

	"deu/internal/places"
	"deu/internal/repository"
	"deu/internal/users"
	"deu/pkg/router"
	"deu/pkg/db"
)

func main() {

	port := os.Getenv("SERVER_PORT")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is not set.")
	}
	gormDB := db.InitDB(dsn)

	userRepo := repository.NewPostgresUserRepository(gormDB)
	placeRepo := repository.NewPostgresPlaceRepository(gormDB)
	userPlaceRepo := repository.NewPostgresUserPlaceRepository(gormDB)

	userService := users.NewUserService(userRepo, userPlaceRepo, placeRepo)
	placeService := places.NewPlaceService(placeRepo)

	userHandler := &users.Handler{Service: userService}
	placeHandler := &places.Handler{Service: placeService} 

	r := router.NewRouter(router.Config{
		UserHandler: userHandler,
		PlaceHandler: placeHandler,
	})
	

	log.Println("Listening on port :" + port)

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if req.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        r.ServeHTTP(w, req)
    })

	log.Fatal(http.ListenAndServe(":" + port, handler))
}