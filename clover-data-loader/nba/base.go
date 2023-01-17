package nba

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"
)

const NBA_API_BASE_URL = "https://stats.nba.com/stats/"
const NBA_LEAGUE_ID = "00"
const WNBA_LEAGUE_ID = "10"

func NBAQuery[T any](endpoint string, params map[string]string) (*NBAResponse[T], error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression: true,
		}}
	paramStr := "params {\n"
	for k, v := range params {
		if v != "" {
			paramStr += fmt.Sprintf("    %v=%v\n", k, v)
		}
	}
	paramStr += "}"
	logrus.Infof("requesting %v %v", endpoint, paramStr)
	req, err := http.NewRequest("GET", NBA_API_BASE_URL+endpoint, nil)
	if err != nil {
		logrus.Errorf("couldn't create request: %v", err)
	}
	req.Header.Set("Host", "stats.nba.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:72.0) Gecko/20100101 Firefox/72.0")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("x-nba-stats-origin", "stats")
	req.Header.Set("x-nba-stats-token", "true")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "https://stats.nba.com/")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Cache-Control", "no-cache")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error retrieving NBA data: %v", err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("response status code: %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("couldn't create gzip reader: %v", err)
		}
		defer reader.Close()
	default:
		reader = resp.Body
	}

	var data NBAResponse[T]
	err = json.NewDecoder(reader).Decode(&data)
	if err != nil {
		logrus.Errorf("couldn't decode NBA data: %v", err)
	}
	return &data, nil
}

type NBAResponse[T any] struct {
	Resource   string            `json:"resource"`
	Parameters []string          `json:"parameters"`
	ResultSets []NBAResultSet[T] `json:"resultSets"`
}

type NBAResultSet[T any] struct {
	Name    string   `json:"name"`
	Headers []string `json:"headers"`
	RowSet  []T      `json:"rowSet"`
}

func (r *NBAResponse[T]) UnmarshalJSON(b []byte) error {
	var alias struct {
		Resource   string                 `json:"resource"`
		Parameters map[string]interface{} `json:"parameters"`
		ResultSets []NBAResultSet[T]      `json:"resultSets"`
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return err
	}
	keys := make([]string, 0, len(alias.Parameters))
	for k := range alias.Parameters {
		keys = append(keys, k)
	}
	r.Resource = alias.Resource
	r.Parameters = keys
	r.ResultSets = alias.ResultSets
	return nil
}

func (r *NBAResultSet[T]) UnmarshalJSON(b []byte) error {
	var alias struct {
		Name    string        `json:"name"`
		Headers []string      `json:"headers"`
		RowSet  []interface{} `json:"rowSet"`
	}
	if err := json.Unmarshal(b, &alias); err != nil {
		return err
	}
	r.Name = alias.Name
	r.Headers = alias.Headers
	rowSetAsJson := make([]map[string]interface{}, len(alias.RowSet))
	for rowIndex, row := range alias.RowSet {
		rowSetAsJson[rowIndex] = make(map[string]interface{})
		for headerIndex, header := range r.Headers {
			rowSetAsJson[rowIndex][header] = row.([]interface{})[headerIndex]
		}
	}
	jsonbody, err := json.Marshal(rowSetAsJson)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonbody, &r.RowSet); err != nil {
		return err
	}
	return nil
}

func ParamMap[T any](p T) map[string]string {
	params := map[string]string{}
	val := reflect.ValueOf(p)
	for i := 0; i < val.Type().NumField(); i++ {
		field := val.Type().Field(i)
		tag := field.Tag.Get("json")
		if tag != "" {
			// if val.Field(i).String() != "" {
			params[tag] = val.Field(i).String()
			// }
		}
	}
	return params
}
