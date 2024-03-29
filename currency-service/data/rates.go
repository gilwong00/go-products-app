package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRate struct {
	log   hclog.Logger
	rates map[string]float64
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

func NewExchangeRange(l hclog.Logger) (*ExchangeRate, error) {
	er := &ExchangeRate{log: l, rates: map[string]float64{}}
	er.getRates()
	return er, nil
}

func (e *ExchangeRate) GetRate(base string, final string) (float64, error) {
	baseRate, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s", base)
	}
	finalRate, ok := e.rates[final]
	if !ok {
		return 0, fmt.Errorf("rate not found for %s", final)
	}
	return finalRate / baseRate, nil
}

// PollRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
func (e *ExchangeRate) PollRates(interval time.Duration) chan struct{} {
	channel := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// just add a random difference to the rate and return it
				// this simulates the fluctuations in currency rates
				for k, v := range e.rates {
					// change can be 10% of original value
					change := (rand.Float64() / 10)
					// is this a postive or negative change
					direction := rand.Intn(1)
					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}
					// modify the rate
					e.rates[k] = v * change
				}
				// notify updates, this will block unless there is a listener on the other end
				channel <- struct{}{}
			}
		}
	}()
	return channel
}

func (e *ExchangeRate) getRates() error {
	res, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status code 200 got %d", res.StatusCode)
	}
	defer res.Body.Close()
	cube := &Cubes{}
	xml.NewDecoder(res.Body).Decode(&cube)

	for _, c := range cube.CubeData {
		rate, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = rate
	}
	e.rates["EUR"] = 1
	return nil
}
