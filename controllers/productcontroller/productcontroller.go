package productcontroller

import (
	"go-jwt-api/helpers"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	data := []map[string]interface{}{
		{
			"id":           1,
			"nama_product": "Kemeja",
			"price":        2000,
		},
		{
			"id":           2,
			"nama_product": "Baju Kaos",
			"price":        2000,
		},
		{
			"id":           3,
			"nama_product": "Sepatu",
			"price":        2000,
		},
	}

	helpers.ResponseJSON(w, http.StatusOK, data)
}
