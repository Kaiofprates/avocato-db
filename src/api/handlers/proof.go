package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

type ProofResponse struct {
	Index uint64   `json:"index"`
	Root  string   `json:"root"`
	Proof []string `json:"proof"`
}

func NewProofHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// In a full MMR implementation, we would extract the index from the URL path 
		// and calculate the Merkle path. For this POC, we return a stub.
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		resp := ProofResponse{
			Index: 0, // Should be parsed from parts[3]
			Root:  "stub_root",
			Proof: []string{"stub_proof_1", "stub_proof_2"},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
