package handlers

import (
	"avocato-db/src/core/ledger"
	"encoding/json"
	"fmt"
	"net/http"
)

type AppendRequest struct {
	Payload interface{} `json:"payload"`
}

type AppendResponse struct {
	Status string `json:"status"`
	Data   struct {
		Index     uint64 `json:"index"`
		Hash      string `json:"hash"`
		PrevHash  string `json:"prev_hash"`
		Timestamp int64  `json:"timestamp"`
	} `json:"data"`
}

// NewAppendHandler godoc
// @Summary      Append data to the ledger
// @Description  Creates a new immutable block with the provided payload
// @Tags         ledger
// @Accept       json
// @Produce      json
// @Param        payload  body      AppendRequest  true  "Data payload"
// @Success      201      {object}  AppendResponse
// @Failure      400      {string}  string "Invalid request body"
// @Failure      500      {string}  string "Internal server error"
// @Router       /v1/append [post]
func NewAppendHandler(service *ledger.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AppendRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		block, err := service.Append(r.Context(), req.Payload)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to append block: %v", err), http.StatusInternalServerError)
			return
		}

		resp := AppendResponse{
			Status: "success",
		}
		resp.Data.Index = block.Index
		resp.Data.Hash = fmt.Sprintf("%x", block.Hash)
		resp.Data.PrevHash = fmt.Sprintf("%x", block.PrevHash)
		resp.Data.Timestamp = block.Timestamp

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
