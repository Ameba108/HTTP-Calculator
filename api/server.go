package server

import (
	"calc/pkg/calculator"
	"encoding/json"
	"log"

	"net/http"
)

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionRespose struct {
	Error  string  `json:"error"`
	Result float64 `json:"result"`
}

func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	req := &ExpressionRequest{}

	if err = dec.Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := &ExpressionRespose{}
	resp.Result, err = calculator.Calc(req.Expression)
	if err != nil {
		if err.Error() == "деление на ноль" {
			resp.Error = "деление на ноль"
			resp.Result = 0
			w.WriteHeader(http.StatusBadRequest)
		} else if err.Error() == "неизвестный символ" {
			resp.Error = "неизвестный символ"
			resp.Result = 0
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		resp.Error = ""
	}
	w.Header().Set("Content-Type", "application/json")
	if resp.Error != "" {
		w.WriteHeader(http.StatusInternalServerError)
	}

	enc := json.NewEncoder(w)
	if err := enc.Encode(resp); err != nil {
		log.Fatal(err)
	}
}

func HandleRequest() {
	http.HandleFunc("/", ExpressionHandler)
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
