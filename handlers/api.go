package handlers

import (
	"encoding/json"
	"go-solid-embed/models"
	"go.uber.org/zap"
	"net/http"
)

func RootResource(w http.ResponseWriter, r *http.Request) {
	appInfo := models.App{
		Name:    "TODO Server",
		Version: "1.0.0",
	}
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(appInfo)
	if err != nil {
		zap.S().Errorw("failed to encode todos", "error", err)
		return
	}
}
