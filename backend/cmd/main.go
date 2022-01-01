package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

	"github.com/darkphnx/vehiclemanager/cmd/api"
	"github.com/darkphnx/vehiclemanager/cmd/background"
	"github.com/darkphnx/vehiclemanager/internal/authservice"
	"github.com/darkphnx/vehiclemanager/internal/models"
	"github.com/darkphnx/vehiclemanager/internal/mothistoryapi"
	"github.com/darkphnx/vehiclemanager/internal/vesapi"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	vesapiKey := flag.String("vesapi-key", "", "Vehicle Enquiry Service API Key")
	mothistoryapiKey := flag.String("mothistoryapi-key", "", "MOT History API Key")
	jwtSigningSecret := flag.String("jwt-signing-secret", "", "JWT Signing Secret")
	mongoConnectionString := flag.String("mongo-connection-string", "", "MongoDB Connection String")
	flag.Parse()

	database, err := models.InitDB(*mongoConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	vesapiClient := vesapi.NewClient(*vesapiKey, "")
	mothistoryClient := mothistoryapi.NewClient(*mothistoryapiKey, "")
	authService := authservice.NewAuthService(*jwtSigningSecret, 24, "mot.ninja")

	backgroundTasks := background.Task{
		Database:                 database,
		VehicleEnquiryServiceAPI: vesapiClient,
		MotHistoryAPI:            mothistoryClient,
	}
	go backgroundTasks.Begin()

	apiServer := api.Server{
		Database:                 database,
		VehicleEnquiryServiceAPI: vesapiClient,
		MotHistoryAPI:            mothistoryClient,
		AuthService:              authService,
	}

	mux := mux.NewRouter()

	mux.Use(api.LoggingMiddleware)

	mux.HandleFunc("/signup", apiServer.Signup).Methods("POST")
	mux.HandleFunc("/login", apiServer.Login).Methods("POST")
	mux.HandleFunc("/logout", apiServer.Logout).Methods("GET")

	apiMux := mux.PathPrefix("/api").Subrouter()
	apiMux.Use(apiServer.AuthJwtTokenMiddleware)
	apiMux.HandleFunc("/vehicles/{registration}", apiServer.VehicleShow).Methods("GET")
	apiMux.HandleFunc("/vehicles/{registration}", apiServer.VehicleDelete).Methods("DELETE")
	apiMux.HandleFunc("/vehicles", apiServer.VehicleList).Methods("GET")
	apiMux.HandleFunc("/vehicles", apiServer.VehicleCreate).Methods("POST")

	// mux.Handle("/", http.FileServer(http.Dir("./ui/build")))

	spa := spaHandler{staticPath: "ui/build", indexPath: "index.html"}
	mux.PathPrefix("/").Handler(spa)

	err = http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
