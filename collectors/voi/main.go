package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type voi struct {
	Phone       string `json:"phone"`
	CountryCode string `json:"country_code"`
	AccessToken string `json:"access_token"`
}

func (v *voi) post(url string, body []byte, data any) error {
	headers := map[string]string{
		"x-app-version": "3.165.1",
		"x-os":          "Android",
		"x-os-version":  "28",
		"model":         "Pixel 3a XL",
		"brand":         "google",
		"manufacturer":  "Google",
		"x-app-name":    "Rider",
		"user-agent":    "okhttp/4.9.2",
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if v.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+v.AccessToken)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get token: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}

func (v *voi) getToken() (string, error) {
	type Response struct {
		Token string `json:"token"`
	}

	response := &Response{}
	body := []byte(`{"country_code":"` + v.CountryCode + `","phone_number":"` + v.Phone + `"}`)

	err := v.post("https://api.voiapp.io/v1/auth/verify/phone", body, response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (v *voi) sendOTP(token string, otp string) (string, string, error) {
	u := "https://api.voiapp.io/v2/auth/verify/code"

	type Response struct {
		VerificationStep string `json:"verificationStep"`
		AuthToken        string `json:"authToken"`
	}

	response := &Response{}
	body := []byte(`{"token":"` + token + `","otp":"` + otp + `"}`)

	err := v.post(u, body, response)
	if err != nil {
		return "", "", err
	}

	return response.AuthToken, response.VerificationStep, nil
}

func main() {
	v := &voi{
		Phone:       "15227679550",
		CountryCode: "DE",
	}

	token, err := v.getToken()
	if err != nil {
		fmt.Println("Error getting token:", err)
		return
	}

	fmt.Println("Token:", token)

	// otp := "your_otp"
	// authToken, verificationStep, err := sendOTP(token, otp)
	// if err != nil {
	// 	// Handle error
	// 	return
	// }

	// Use authToken and verificationStep
}
