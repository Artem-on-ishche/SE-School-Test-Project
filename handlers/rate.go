package handlers

import (
	"fmt"
	"net/http"

	"gses2.app/api/rate"
)

func RateHandler(responseWriter http.ResponseWriter, _ *http.Request) {
	exchangeRate, err := rate.GetBtcToUahRate()
	if isRateWrong(exchangeRate, err) {
		sendBadRequestResponse(responseWriter, "An error has occurred")

		return
	}

	rateString := fmt.Sprintf("%v", exchangeRate)
	sendSuccessResponse(responseWriter, rateString)
}

func isRateWrong(rate float64, err error) bool {
	return err != nil || rate <= 0
}