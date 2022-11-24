package requests

import (
	"bytes"
	"encoding/json"
	"fmt"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sofa-logs-servers/models"
	"sofa-logs-servers/utils"
)

type FindForm struct {
	RequestID uuid.UUID `json:"request_id"`
}

type RespFind struct {
	Took     int         `json:"took"`
	TimedOut bool        `json:"timed_out"`
	Shards   interface{} `json:"_shards"`
	Hits     struct {
		Total struct {
			Value    int    `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string     `json:"_index"`
			Id     string     `json:"_id"`
			Score  float64    `json:"_score"`
			Source models.Log `json:"_source"`
		}
	}
}

func FindById(w http.ResponseWriter, r *http.Request, client *es8.Client) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	form := FindForm{}
	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	if form.RequestID.String() == "" {
		utils.WriteErr(w, "UUID EMPTY !!!", http.StatusBadRequest)
		return
	}

	var buff bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				"_id": form.RequestID.String(),
			},
		},
	}

	err = jsoniter.NewEncoder(&buff).Encode(query)
	if err != nil {
		utils.WriteErr(w, "ERROR ENCODING QUERY", http.StatusInternalServerError)
		return
	}

	resp, err := client.Search(
		client.Search.WithIndex("requests"),
		client.Search.WithBody(&buff),
		client.Search.WithTrackScores(true),
		client.Search.WithPretty(),
	)

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		utils.WriteErr(w, "ERROR Seearching For Request", http.StatusInternalServerError)

	}

	if resp.IsError() {
		utils.WriteErr(w, "Error parsing the response body", http.StatusBadRequest)
		return
	}

	resDecoded := RespFind{}

	if err := json.NewDecoder(resp.Body).Decode(&resDecoded); err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE", http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, resDecoded.Hits.Hits[0])

	fmt.Println(resp)

}

func FindAll(w http.ResponseWriter, _ *http.Request, client *es8.Client) {

	var (
		buff bytes.Buffer
	)
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	err := jsoniter.NewEncoder(&buff).Encode(query)
	if err != nil {
		utils.WriteErr(w, "ERROR ENCODING QUERY", http.StatusInternalServerError)
		return
	}

	resp, err := client.Search(
		client.Search.WithIndex("requests"),
		client.Search.WithBody(&buff),
		client.Search.WithTrackScores(true),
		client.Search.WithPretty(),
	)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		utils.WriteErr(w, "ERROR SEARCHING FOR REQUESTS", http.StatusInternalServerError)
	}

	if resp.IsError() {
		utils.WriteErr(w, "Error searching the Document", http.StatusBadRequest)
		return
	}

	resDecoded := RespFind{}
	if err := json.NewDecoder(resp.Body).Decode(&resDecoded); err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, resDecoded.Hits.Hits)
}
