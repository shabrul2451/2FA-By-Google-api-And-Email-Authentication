package main

import (
	container "cloud.google.com/go/container/apiv1"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"log"
	"net/http"
)

type GKEReleaseNotes struct {
	Version string `json:"version"`
}

func GetGKEReleaseNotes(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectId := "default"
	zone := "us-central1"
	client, err := container.NewClusterManagerClient(ctx, option.WithUserAgent("123656749011-compute@developer.gserviceaccount.com"))
	log.Println(err)
	if err != nil {
		http.Error(w, "Error creating GKE client", http.StatusInternalServerError)
		return
	}
	defer client.Close()

	req := &containerpb.GetServerConfigRequest{
		ProjectId: projectId, // Replace with your GCP project ID
		Zone:      zone,      // Replace with your desired GCP zone
	}

	resp, err := client.GetServerConfig(ctx, req)
	if err != nil {
		http.Error(w, "Error getting GKE server config", http.StatusInternalServerError)
		return
	}

	var releaseNotes []GKEReleaseNotes

	for _, validVersion := range resp.GetValidMasterVersions() {
		gkeReleaseNotes := GKEReleaseNotes{
			Version: validVersion,
		}

		releaseNotes = append(releaseNotes, gkeReleaseNotes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(releaseNotes)
}

func main() {
	r := mux.NewRouter()

	// Define the API endpoint for getting the list of GKE release notes versions
	r.HandleFunc("/api/gke/release-notes", GetGKEReleaseNotes).Methods("GET")

	port := "8080"
	fmt.Printf("Starting server on port %s...\n", port)
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
