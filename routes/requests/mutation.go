package requests

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sofa-logs-servers/models"
	"sofa-logs-servers/utils"
	"time"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type CreateForm struct {
	UserID    uint   `json:"user_id"`
	Page      string `json:"page"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

type UpdateForm struct {
	ID        string `json:"request_id"`
	UserID    uint   `json:"user_id"`
	Page      string `json:"page"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

type DeleteFrom struct {
	ID string `json:"request_id"`
}

type UpdateReq struct {
	UserID    uint   `json:"user_id"`
	Page      string `json:"page"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at"`
}

type RespMutation struct {
	Index   string `json:"_index"`
	Id      string `json:"_id"`
	Version int    `json:"_version"`
	Result  string `json:"result"`
	Shards  struct {
		Total      int `json:"total"`
		Successful int `json:"successful"`
		Failed     int `json:"failed"`
	} `json:"_shards"`
	SeqNo       int `json:"_seq_no"`
	PrimaryTerm int `json:"_primary_term"`
}

func Create(w http.ResponseWriter, r *http.Request, client *es8.Client) {

	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	form := CreateForm{}
	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	startedAt, err := time.Parse(time.RFC3339, form.StartedAt)
	if err != nil {
		utils.WriteErr(w, "BAD DATE FORMAT", http.StatusBadRequest)
		return
	}

	endedAt, err := time.Parse(time.RFC3339, form.EndedAt)
	if err != nil {
		utils.WriteErr(w, "BAD DATE FORMAT", http.StatusBadRequest)
		return
	}

	logForm := models.Log{
		UserID:    form.UserID,
		Page:      form.Page,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if logForm.Page == "" {
		utils.WriteErr(w, "EMPTY Page", http.StatusBadRequest)
		return
	}

	if logForm.StartedAt.IsZero() {
		utils.WriteErr(w, "EMPTY Start Request Time", http.StatusBadRequest)
	}

	if logForm.EndedAt.IsZero() {
		utils.WriteErr(w, "EMPTY End Request Time", http.StatusBadRequest)
	}

	f, err := jsoniter.Marshal(logForm)
	if err != nil {
		utils.WriteErr(w, "can't parse body data to json", http.StatusBadRequest)
		return
	}

	request := esapi.IndexRequest{
		Index:      "requests",
		DocumentID: uuid.New().String(),
		Body:       bytes.NewReader(f),
	}

	resp, err := request.Do(context.Background(), client)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		utils.WriteErr(w, "Could not send Request to Elastic Client", http.StatusBadRequest)
		return
	}

	fmt.Println(resp)

	if resp.IsError() {
		utils.WriteErr(w, "Error Creating the Document", http.StatusBadRequest)
		return
	}

	respDecoded := RespMutation{}
	err = jsoniter.NewDecoder(resp.Body).Decode(&respDecoded)
	if err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
	}
	utils.WriteJson(w, respDecoded)
}

func Update(w http.ResponseWriter, r *http.Request, client *es8.Client) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	form := UpdateForm{}

	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	startedAt, err := time.Parse(time.RFC3339, form.StartedAt)
	if err != nil {
		utils.WriteErr(w, "BAD DATE FORMAT", http.StatusBadRequest)
		return
	}

	endedAt, err := time.Parse(time.RFC3339, form.EndedAt)
	if err != nil {
		utils.WriteErr(w, "BAD DATE FORMAT", http.StatusBadRequest)
		return
	}

	logForm := models.Log{
		UserID: form.UserID,
		Page: form.Page,
		StartedAt: startedAt,
		EndedAt: endedAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if form.ID == "" {
		utils.WriteErr(w, "request_id IS REQUIRED", http.StatusBadRequest)
		return
	}

	if logForm.Page == "" {
		utils.WriteErr(w, "page IS REQUIRED", http.StatusBadRequest)
		return
	}

	if logForm.StartedAt.IsZero() {
		utils.WriteErr(w, "EMPTY Start Request Time", http.StatusBadRequest)
	}

	if logForm.EndedAt.IsZero() {
		utils.WriteErr(w, "EMPTY End Request Time", http.StatusBadRequest)
	}

	//encode the form to update the index

	//remove the ID field from form

	reqEncoded, err := jsoniter.Marshal(logForm)
	if err != nil {
		utils.WriteErr(w, "CAN'T PARSE BODY DATA TO JSON", http.StatusBadRequest)
	}
	request := esapi.UpdateRequest{
		Index:      "requests",
		DocumentID: form.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, reqEncoded))),
	}

	resp, err := request.Do(context.Background(), client)
	if err != nil {
		utils.WriteErr(w, "Could not send Request to Elastic Client", http.StatusBadRequest)
		return
	}

	if resp.IsError() {
		utils.WriteErr(w, "UPDATING DOCUMENT FAILED", http.StatusBadRequest)
		return
	}

	respDecoded := RespMutation{}
	err = jsoniter.NewDecoder(resp.Body).Decode(&respDecoded)
	if err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
	}
	utils.WriteJson(w, respDecoded)
}

func Delete(w http.ResponseWriter, r *http.Request, client *es8.Client) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	form := DeleteFrom{}
	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	if form.ID == "" {
		utils.WriteErr(w, "request_id IS REQUIRED", http.StatusBadRequest)
		return
	}

	request := esapi.DeleteRequest{
		Index:      "requests",
		DocumentID: form.ID,
	}

	resp, err := request.Do(context.Background(), client)
	if err != nil {
		utils.WriteErr(w, "Could not send Request to Elastic Client", http.StatusBadRequest)
		return
	}

	if resp.IsError() {
		utils.WriteErr(w, "DELETING DOCUMENT FAILED", http.StatusBadRequest)
		return
	}

	respDecoded := RespMutation{}
	err = jsoniter.NewDecoder(resp.Body).Decode(&respDecoded)
	if err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
	}
	utils.WriteJson(w, respDecoded)
}
