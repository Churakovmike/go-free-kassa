package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var (
	baseUrl = "https://www.free-kassa.ru/api.php?type=json"
	ordersUrl = "https://www.free-kassa.ru/export.php?type=json"
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

type Balance struct {
	Type string      `json:"type"`
	Desc string      `json:"desc"`
	Data BalanceData `json:"data"`
}

type BalanceData struct {
	Balance string `json:"balance"`
}

func newClient(cfg *Config) (*Client, error) {
	return &Client{config: cfg}, nil
}

func generateApiSignature(merchantId string, secondSecret string) string {
	hash := merchantId + secondSecret
	hasher := md5.New()
	hasher.Write([]byte(hash))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (client *Client) getBalance() (*http.Response, error) {
	m := map[string]string{
		"merchant_id": client.config.MerchantId,
		"action":      "get_balance",
		"s":           generateApiSignature(client.config.MerchantId, client.config.SecondSecret),
	}
	form := buildForm(m)
	req, _ := createRequest("POST", baseUrl, form)
	httpClient := http.Client{}
	return httpClient.Do(req)
}

func (client *Client) getOrder(orderId string, intid string) (*http.Response, error) {
	m := map[string]string{
		"merchant_id": client.config.MerchantId,
		"s":           generateApiSignature(client.config.MerchantId, client.config.SecondSecret),
		"action":      "check_order_status",
		"order_id":    orderId,
		"intid":       intid,
	}
	form := buildForm(m)
	req, _ := createRequest("POST", baseUrl, form)
	httpClient := http.Client{}
	return httpClient.Do(req)
}

func (client *Client) exportOrders(dateFrom string, dateTo string, limit string, offset string, status string) (*http.Response, error) {
	m := map[string]string{
		"merchant_id": client.config.MerchantId,
		"s":           generateApiSignature(client.config.MerchantId, client.config.SecondSecret),
		"action": "get_orders",
		"date_from" : dateFrom,
		"date_to" : dateTo,
		"status": status,
		"limit" : limit,
		"offset": offset,
	}
	form := buildForm(m)
	req, _ := createRequest("POST", ordersUrl, form)
	httpClient := http.Client{}
	return httpClient.Do(req)
}

func buildForm(m map[string]string) string {
	form := url.Values{}
	for k, v := range m {
		form.Add(k, url.QueryEscape(v))
	}

	return form.Encode()
}

func createRequest(method string, url string, query string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(query))
	if err != nil {
		fmt.Println("Create request error", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, err
}
