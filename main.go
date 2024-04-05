package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/adrianowsh/go-multithreading/dto"
)

const url_viacep = "https://viacep.com.br/ws"
const url_brasilapi = "https://brasilapi.com.br/api/cep/v1"
const time_to_request = 1 * time.Second

func main() {
	chViaCep := make(chan string)
	chBrasilApi := make(chan string)

	ctx, cancel := context.WithTimeout(context.Background(), time_to_request)
	defer cancel()

	go func(ch chan<- string) {
		resultViaCep, err := getResultViaCep(ctx, "06401015")
		if err != nil {
			log.Fatalln("timeout viacep", err)
		}
		ch <- fmt.Sprintf("brasil via cep => %s", resultViaCep)

	}(chViaCep)

	go func(ch chan<- string) {
		resultBrasilApi, err := getResultBrasilApi(ctx, "06401015'")
		if err != nil {
			log.Fatalln("timeout brasilapi", err)
		}
		ch <- fmt.Sprintf("brasil api => %s", resultBrasilApi)

	}(chBrasilApi)

	select {
	case respViaCep := <-chViaCep:
		println("viacep api -> ", respViaCep)
	case respBrasilApi := <-chBrasilApi:
		println("brasil api -> ", respBrasilApi)
	case <-ctx.Done():
		println("context timeout")
	}
}

func getViaCepAPi(ctx context.Context, cep string) (*dto.ViaCep, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s/json", url_viacep, cep), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result dto.ViaCep
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getBrasilApi(ctx context.Context, cep string) (*dto.BrasilApi, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", url_brasilapi, cep), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result dto.BrasilApi
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func getResultViaCep(ctx context.Context, cep string) (string, error) {
	resp, err := getViaCepAPi(ctx, cep)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.Marshal(resp)
	return string(jsonResult), nil
}

func getResultBrasilApi(ctx context.Context, cep string) (string, error) {
	resp, err := getBrasilApi(ctx, cep)
	if err != nil {
		return "", err
	}
	jsonResult, _ := json.Marshal(resp)
	return string(jsonResult), nil
}
