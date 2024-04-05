package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/adrianowsh/go-multithreading/dto"
)

const url_viacep = "https://viacep.com.br/ws"
const url_brasilapi = "https://brasilapi.com.br/api/cep/v1"
const time_to_request = 1 * time.Second

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx, cancel := context.WithTimeout(context.Background(), time_to_request)
	defer cancel()

	go func() {
		resultViaCep, err := getResultViaCep(ctx, "06401015")
		if err != nil {
			log.Fatalln("timeout viacep", err)
		}
		println("viacep => ", resultViaCep)
		wg.Done()
	}()

	go func() {
		resultBrasilApi, err := getResultBrasilApi(ctx, "06401015'")
		if err != nil {
			log.Fatalln("timeout brasilapi", err)
		}
		println("brasil api =>", resultBrasilApi)
		wg.Done()
	}()

	wg.Wait()
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
	resultViaCep, _ := json.Marshal(resp)

	return string(resultViaCep), nil
}

func getResultBrasilApi(ctx context.Context, cep string) (string, error) {
	resp, err := getBrasilApi(ctx, cep)
	if err != nil {
		return "", err
	}
	resultBrasilApi, _ := json.Marshal(resp)

	return string(resultBrasilApi), nil
}
