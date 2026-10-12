package logic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

const webshopApiUrl = "https://gf2-us-webshop-api.sunborngame.com/shop/product_list"

func FetchTrueUidFromWebshop(accessToken string) (string, error) {
	reqBody := []byte(`{}`)
	req, err := http.NewRequest("POST", webshopApiUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", errors.WithStack(err)
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Safari/537.36")
	req.Header.Set("Origin", "https://gf2-webshop-us.sunborngame.com")
	req.Header.Set("Referer", "https://gf2-webshop-us.sunborngame.com/")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.WithStack(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("webshop API HTTP status: %d", resp.StatusCode)
	}
	return parseWebshopUID(body)
}

func parseWebshopUID(body []byte) (string, error) {
	var parsedResponse map[string]interface{}
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&parsedResponse); err != nil {
		return "", errors.Wrap(err, "failed to parse webshop JSON response")
	}

	var rawUid interface{}
	if dataMap, ok := parsedResponse["data"].(map[string]interface{}); ok {
		if uid, exists := dataMap["uid"]; exists {
			rawUid = uid
		}
	}

	if rawUid == nil {
		if uid, exists := parsedResponse["uid"]; exists {
			rawUid = uid
		}
	}

	if rawUid == nil {
		return "", errors.New("webshop API response did not contain a uid field")
	}

	uidStr := fmt.Sprintf("%v", rawUid)
	uid, parseErr := strconv.ParseUint(uidStr, 10, 64)
	if parseErr != nil || uid == 0 {
		return "", errors.New("webshop API returned empty or zero UID")
	}

	return uidStr, nil
}
