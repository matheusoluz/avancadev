package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

var coupons Coupons

func main() {
	coupons = Coupons{
		Coupon: []Coupon{
			{Code: "abc"},
			{Code: "123"},
			{Code: "qwe"},
		},
	}

	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	jsonResult, err := json.Marshal(coupons)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}
