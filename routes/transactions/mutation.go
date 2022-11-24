package transactions

import (
	"bytes"
	"context"
	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sofa-logs-servers/models"
	"sofa-logs-servers/utils"
	"time"
)

type CreateForm struct {
	Amount uint   `json:"amount"`
	Date   string `json:"date"`
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
		r.Body.Close()
	}()

	form := CreateForm{}
	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, form.Date)
	if err != nil {
		utils.WriteErr(w, "BAD DATE FORMAT", http.StatusBadRequest)
		return
	}

	transactionForm := models.Transaction{
		Amount:    form.Amount,
		Date:      date,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if transactionForm.Amount == 0 {
		utils.WriteErr(w, "EMPTY Page", http.StatusBadRequest)
		return
	}

	if transactionForm.Date.IsZero() {
		utils.WriteErr(w, "EMPTY Start Request Time", http.StatusBadRequest)
	}

	f, err := jsoniter.Marshal(transactionForm)
	if err != nil {
		utils.WriteErr(w, "can't parse body data to json", http.StatusBadRequest)
		return
	}

	request := esapi.IndexRequest{
		Index:      "transactions",
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
