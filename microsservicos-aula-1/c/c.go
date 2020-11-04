package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status string
}

var coupons Coupons

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	resultCouponList := makeHttpCall("http://localhost:9093")
	fmt.Println(resultCouponList)
	// coupons = resultCouponList
	valid := resultCouponList.Check(coupon)

	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}

func makeHttpCall(urlMicroservice string) Coupons {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.Get(urlMicroservice)
	if err != nil {
		// result := Result{Status: "Servidor fora do ar!"}
		result := Coupons{}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Coupons{}

	json.Unmarshal(data, &result)
	return result
}
