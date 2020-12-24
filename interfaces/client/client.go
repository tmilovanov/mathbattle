package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"mathbattle/models/mathbattle"
)

func sendReq(method string, endpoint string, object interface{}) (*http.Response, error) {
	if object != nil {
		log.Printf("%s %s object: %v", method, endpoint, object)
	} else {
		log.Printf("%s %s", method, endpoint)
	}

	var body io.Reader = nil
	if object != nil {
		jsonStr, err := json.Marshal(object)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonStr)
	}

	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	if object != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func PostJsonRecieveJson(endpoint string, send interface{}, recieve interface{}) error {
	resp, err := sendReq("POST", endpoint, send)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(recieve)
	if err != nil {
		return fmt.Errorf("Failed to decode response, error: %v", err)
	}

	return nil
}

func PostJsonRecieveNone(endpoint string, object interface{}) error {
	resp, err := sendReq("POST", endpoint, object)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	return nil
}

func PostNoneRecieveNone(endpoint string) error {
	resp, err := sendReq("POST", endpoint, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	return nil
}

func SendGetNoneRecieveJson(endpoint string, object interface{}) error {
	resp, err := sendReq("GET", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(object)
	if err != nil {
		return fmt.Errorf("Failed to decode response, error: %v", err)
	}

	return nil
}

func SendGetJsonRecieveJson(endpoint string, send interface{}, recieve interface{}) error {
	resp, err := sendReq("GET", endpoint, send)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(recieve)
	if err != nil {
		return fmt.Errorf("Failed to decode response, error: %v", err)
	}

	return nil
}

func PutJsonRecieveNone(endpoint string, object interface{}) error {
	resp, err := sendReq("PUT", endpoint, object)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	return nil
}

func DeleteRecieveNone(endpoint string) error {
	resp, err := sendReq("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return mathbattle.ErrNotFound
		}
		return fmt.Errorf("Unexpected HTTP status: %d", resp.StatusCode)
	}

	return nil
}
