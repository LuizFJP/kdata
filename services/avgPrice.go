package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// MeanPrice struct
type MeanPrice struct {
	Symbols []struct {
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	} `json:"symbols"`
}

// AVGPrice get coin average price based on various exchanges prices
func AVGPrice(symbol []string, url string) (MeanPrice, error) {
	param := map[string]interface{}{
		"names": symbol,
	}

	reqBody, err := json.Marshal(param)
	if err != nil {
		return MeanPrice{}, err
	}

	var resp *http.Response
	attempts := 3 // TODO - set up the attempts request variable at env/config variable

loopAttempts:
	for attempts > 0 {
		resp, err = http.Post(url, "application/json", bytes.NewBuffer(reqBody))
		switch {
		case err != nil:
			log.Println(err)
		case resp.StatusCode != 200:
			msg := GetHttpStatusMsgByCode(resp.StatusCode)
			err = fmt.Errorf("code: %v, msg: %v", resp.StatusCode, msg)
			if err != nil {
				log.Println(err)
			}
		default:
			break loopAttempts
		}
		attempts--
		time.Sleep(500 * time.Millisecond) // TODO - set up the time sleep variable at env/config variable
		continue
	}

	// reach attemps should return with last error
	if attempts == 0 {
		return MeanPrice{}, err
	}

	var meanPrice MeanPrice
	errDecode := json.NewDecoder(resp.Body).Decode(&meanPrice)
	if errDecode != nil {
		return MeanPrice{}, errDecode
	}

	for _, data := range meanPrice.Symbols {
		if data.Price == 0 {
			return MeanPrice{}, fmt.Errorf("price miss symbol failed", data.Name)
		}
	}
	return meanPrice, nil
}

func GetHttpStatusMsgByCode(code int) string {
	var msg string
	switch code {
	case http.StatusBadRequest:
		msg = http.StatusText(http.StatusBadRequest)
	case http.StatusForbidden:
		msg = http.StatusText(http.StatusForbidden)
	case http.StatusNotFound:
		msg = http.StatusText(http.StatusNotFound)
	default:
		msg = "not mapped error"
	}
	return msg
}
