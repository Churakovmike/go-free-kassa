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
)

type Client struct {
	config     *Config
	httpClient *http.Client
}

type Balance struct {
	Type string `json:"type"`
	Desc string `json:"desc"`
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
		"merchant_id" : client.config.MerchantId,
		"action": "get_balance",
		"s" : generateApiSignature(client.config.MerchantId, client.config.SecondSecret),
	}
	form := buildForm(m)
	println(form)
	method := "POST"
	req, err := http.NewRequest(method, baseUrl, strings.NewReader(form))
	if err != nil {
		fmt.Println("Ошибка при попытке получения баланса ", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	hc := http.Client{}
	return hc.Do(req)
}

func buildForm(m map[string]string) string {
	form := url.Values{}
	for k, v := range m {
		form.Add(k, url.QueryEscape(v))
	}

	return form.Encode()
}
