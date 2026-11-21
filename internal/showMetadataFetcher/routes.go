package showmetadatafetcher

import (
	"encoding/json"
	"net/http"
)

func (s *Service) MetadataFetcherRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /syncMetadata/", s.syncMetadata)
}

func (s *Service) syncMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req syncMetadataRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	inserted, err := s.InsertForDirectory(ctx, req.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	res := syncMetadataResponse{Inserted: inserted}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

type syncMetadataRequest struct {
	Path string `json:"path"`
}

type syncMetadataResponse struct {
	Inserted int64 `json:"inserted"`
}
