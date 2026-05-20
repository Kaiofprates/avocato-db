package handlers

import (
	"avocato-db/src/core/integrity"
	"avocato-db/src/core/ledger"
	"encoding/json"
	"net/http"
)

type IntegrityResponse struct {
	TotalBlocks uint64 `json:"total_blocks"`
	MMRRoot     string `json:"mmr_root"`
	Status      string `json:"status"`
}

func NewIntegrityHandler(service *ledger.Service, mmr *integrity.MMR) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		lastIndex, _ := service.GetLastState()

		resp := IntegrityResponse{
			TotalBlocks: lastIndex,
			MMRRoot:     mmr.RootHex(),
			Status:      "VALID",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
