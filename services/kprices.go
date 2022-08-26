package services

import (
	"fmt"
	"log"
)

// PRICES
var (
	StrPricesLoadingURL     = "loading avgPrice from:"
	StrAvgPriceFailed       = "check avgPrice failed"
	strEmptyAvgPriceSymbols = "avgPrice symbols cant be empty"
	strMissingAvgPrice      = "avgPrice of one symbols is missing"
	PriceProdURL            = "https://prices.endpoints.services.klever.io/v1/prices"
)

type PriceServiceInterface interface {
	FetchAvgPrice(symbols []string) (MeanPrice, error)
}

type PriceService struct {
	url string
}

func NewPriceService(url string) PriceServiceInterface {
	if len(url) == 0 {
		url = PriceProdURL
	}

	return &PriceService{
		url: url,
	}
}

// FetchAvgPrice returns all prices of each symbols item
func (s *PriceService) FetchAvgPrice(symbols []string) (MeanPrice, error) {
	log.Println(strEmptyAvgPriceSymbols, s.url)

	meanPrice, err := AVGPrice(symbols, s.url)
	if err != nil {
		err := fmt.Errorf("%s: %s", StrAvgPriceFailed, err.Error())
		return MeanPrice{}, err
	}

	if len(meanPrice.Symbols) == 0 {
		err := fmt.Errorf(strEmptyAvgPriceSymbols)
		return MeanPrice{}, err
	}

	if len(meanPrice.Symbols) != len(symbols) {
		err := fmt.Errorf(strMissingAvgPrice)
		return MeanPrice{}, err
	}

	return meanPrice, err
}
