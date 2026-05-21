package handlers

import (
	"avocato-db/src/core/ledger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type BlockResponse struct {
	Index     uint64      `json:"index"`
	UUID      string      `json:"uuid"`
	Hash      string      `json:"hash"`
	PrevHash  string      `json:"prev_hash"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

// NewGetBlockHandler godoc
// @Summary      Get block by UUID
// @Description  Retrieves a specific block from the ledger using its unique UUID. Indices and hashes are not allowed for external lookups.
// @Tags         ledger
// @Produce      json
// @Param        id   path      string  true  "Block UUID"
// @Success      200  {object}  BlockResponse
// @Failure      400  {string}  string "Invalid UUID format"
// @Failure      404  {string}  string "Block not found"
// @Router       /v1/block/{id} [get]
func NewGetBlockHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Path: /v1/block/{id}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) < 4 {
			http.Error(w, "Invalid block ID", http.StatusBadRequest)
			return
		}
		id := parts[3]

		// STRICT VALIDATION: Only allow UUIDs
		if _, err := uuid.Parse(id); err != nil {
			http.Error(w, "Invalid UUID format. Access by Index or Hash is forbidden for security reasons.", http.StatusBadRequest)
			return
		}

		block, err := service.GetBlock(r.Context(), id)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to retrieve block: %v", err), http.StatusInternalServerError)
			return
		}

		if block == nil {
			http.Error(w, "Block not found", http.StatusNotFound)
			return
		}

		resp := BlockResponse{
			Index:     block.Index,
			UUID:      block.UUID,
			Hash:      fmt.Sprintf("%x", block.Hash),
			PrevHash:  fmt.Sprintf("%x", block.PrevHash),
			Timestamp: block.Timestamp,
			Payload:   json.RawMessage(block.Payload),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
