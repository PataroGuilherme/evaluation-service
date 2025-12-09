package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// healthHandler retorna um status básico de saúde do serviço.
func (a *App) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// evaluationHandler é o handler público chamado em /evaluate.
// Ele chama EvaluateUserFlag, que é o método correto definido em evaluator.go.
func (a *App) evaluationHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	flagName := r.URL.Query().Get("flag_name")

	if userID == "" || flagName == "" {
		http.Error(w, "Parâmetros obrigatórios: user_id, flag_name", http.StatusBadRequest)
		return
	}

	log.Printf("Recebida solicitação de avaliação: user_id=%s, flag=%s", userID, flagName)

	// CHAMADA CORRETA:
	result, err := a.EvaluateUserFlag(userID, flagName)
	if err != nil {
		log.Printf("Erro ao avaliar flag '%s': %v", flagName, err)
		http.Error(w, `{"error":"Erro interno ao avaliar a flag"}`, http.StatusInternalServerError)
		return
	}

        if a.SqsSvc != nil && a.SqsQueueURL != "" {
            go a.sendEvaluationEvent(userID, flagName, result)
        }

	resp := map[string]interface{}{
		"user_id":   userID,
		"flag_name": flagName,
		"result":    result,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
