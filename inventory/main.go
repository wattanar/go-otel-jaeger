package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

func main() {

	http.HandleFunc("/inventory/update", func(w http.ResponseWriter, r *http.Request) {

		n := rand.Intn(10)

		if n < 5 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application-json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "update inventory failed.",
			})
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application-json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "update inventory successful.",
		})
	})

	http.ListenAndServe(":8081", nil)
}
