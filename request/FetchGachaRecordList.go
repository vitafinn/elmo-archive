package request

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type GachaResponseBody struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type GachaRecordListData struct {
	RecordList []struct {
		PoolId         int64 `json:"pool_id"`
		ItemId         int64 `json:"item"`
		GachaTimestamp int64 `json:"time"`
	} `json:"list"`
	Next string `json:"next"`
}

func FetchGachaRecordList(gachaUrl, accessToken, next string, poolType int64) (data GachaRecordListData, err error) {
	// PROTOCOL ENFORCEMENT: Upgrade HTTP to HTTPS to prevent 301 redirects
	// from dropping the POST body and causing 404 errors.
	if strings.HasPrefix(gachaUrl, "http://") {
		gachaUrl = strings.Replace(gachaUrl, "http://", "https://", 1)
	}

	values := url.Values{}
	if next != "" {
		values.Set("next", next)
	}
	values.Set("type_id", strconv.FormatInt(poolType, 10))

	req, err := http.NewRequest("POST", gachaUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return GachaRecordListData{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", accessToken)

	client, err := NewHttpClient()
	if err != nil {
		return GachaRecordListData{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return GachaRecordListData{}, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return GachaRecordListData{}, err
	}

	if len(bodyBytes) > 0 && bodyBytes[0] != '{' {
		return GachaRecordListData{}, errors.Errorf("Server rejected request (expected JSON): %s", string(bodyBytes))
	}

	var body GachaResponseBody
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return GachaRecordListData{}, err
	}

	if body.Code == 0 {
		err = json.Unmarshal(body.Data, &data)
		if err != nil {
			return GachaRecordListData{}, err
		}

		return data, nil
	} else {
		return GachaRecordListData{}, errors.Errorf("%s(Code:%d)", body.Message, body.Code)
	}
}
