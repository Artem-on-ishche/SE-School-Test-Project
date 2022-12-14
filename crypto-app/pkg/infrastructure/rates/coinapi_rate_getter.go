package rates

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/resty.v0"

	"gses2.app/api/pkg/application"
	"gses2.app/api/pkg/config"
	"gses2.app/api/pkg/domain/models"
	"gses2.app/api/pkg/domain/services"
)

const coinAPIRequestFormatString = "https://rest.coinapi.io/v1/exchangerate/%s/%s"

type receivedCoinAPIResponse struct {
	Time string  `json:"time"`
	Rate float64 `json:"rate"`
}

type CoinAPIClientFactory struct {
	Cacher CacherRateService
	Logger services.Logger
}

func (factory CoinAPIClientFactory) CreateRateService() ExchangeRateServiceChain {
	return &exchangeRateService{
		concreteRateClient: coinAPIClient{},
		cacher:             factory.Cacher,
		logger:             factory.Logger,
	}
}

type coinAPIClient struct{}

func (c coinAPIClient) name() string {
	return "Coinbase"
}

func (c coinAPIClient) getAPIRequestURLForGivenCurrencies(pair models.CurrencyPair) string {
	return fmt.Sprintf(coinAPIRequestFormatString, pair.Base.Name, pair.Quote.Name)
}

func (c coinAPIClient) getAPIRequest() *resty.Request {
	return resty.R().SetHeader("X-CoinAPI-Key", config.CoinAPIKeyValue)
}

func (c coinAPIClient) parseResponseBody(responseBody []byte) (*parsedResponse, error) {
	var result receivedCoinAPIResponse

	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, application.ErrAPIResponseUnmarshallError
	}

	timestamp, err := time.Parse(timeLayout, result.Time)
	if err != nil {
		return nil, application.ErrAPIResponseUnmarshallError
	}

	return &parsedResponse{
		price: result.Rate,
		time:  timestamp,
	}, nil
}
