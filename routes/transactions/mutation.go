package transactions

import (
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sofa-logs-servers/infra/zincsearch"
	"sofa-logs-servers/utils"
	"time"
)

type CreateForm struct {
	Amount uint      `json:"amount"`
	Date   time.Time `json:"date"`
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

type UpdateForm struct {
	ID        string    `json:"transaction_id"`
	Amount    uint      `json:"amount"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

type DeleteFrom struct {
	ID string `json:"transaction_id"`
}

type CreateRes struct {
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

func Create(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {

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

	if form.Date.IsZero() {
		utils.WriteErr(w, "EMPTY DATE", http.StatusBadRequest)
		return
	}

	if form.Amount == 0 {
		utils.WriteErr(w, "EMPTY AMOUNT", http.StatusBadRequest)
	}

	document := map[string]interface{}{
		"Amount":    form.Amount,
		"Date":      form.Date,
		"CreatedAt": time.Now(),
	} // map[string]interface{} | Document

	_, res, err := zincClient.Client.Document.Index(zincClient.Ctx, "transactions").Document(document).Execute()

	if err != nil {
		utils.WriteErr(w, "Could not send Request to zinc-search Client", http.StatusBadRequest)
		return
	}

	if res.StatusCode != 200 {
		utils.WriteErr(w, "BAD RESPONSE FROM ZINC WHILE CREATING DOCUMENT", http.StatusBadRequest)
		return
	}

	respDecoded := CreateRes{}
	err = jsoniter.NewDecoder(res.Body).Decode(&respDecoded)
	if err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
	}
	utils.WriteJson(w, respDecoded)
}

func Update(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {
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

	//

	if form.ID == "" {
		utils.WriteErr(w, "transaction_id IS REQUIRED", http.StatusBadRequest)
		return
	}

	if form.Date.IsZero() {
		utils.WriteErr(w, "DATE IS REQUIRED", http.StatusBadRequest)
		return
	}

	if form.Amount == 0 {
		utils.WriteErr(w, "AMOUNT IS REQUIRED", http.StatusBadRequest)
		return
	}

	if form.CreatedAt.IsZero() {
		utils.WriteErr(w, "CREATED_AT IS REQUIRED", http.StatusBadRequest)
		return
	}

	document := map[string]interface{}{
		"Date":       form.Date,
		"Amount":     form.Amount,
		"created_at": form.CreatedAt,
	} // map[string]interface{} | Document

	_, res, err := zincClient.Client.Document.Update(zincClient.Ctx, "transactions", form.ID).Document(document).Execute()

	if err != nil {
		utils.WriteErr(w, "Could not send Request to zinc-search Client", http.StatusBadRequest)
		return
	}

	if res.StatusCode != 200 {
		utils.WriteErr(w, "BAD RESPONSE FROM ZINC WHILE UPDATING DOCUMENT", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, "Log Updated")

}

func Delete(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {
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

	_, res, err := zincClient.Client.Document.Delete(zincClient.Ctx, "transactions", form.ID).Execute()
	if err != nil {
		utils.WriteErr(w, "Error deleting the Document", http.StatusBadRequest)
		return
	}

	if res.StatusCode != 200 {
		utils.WriteErr(w, "BAD RESPONSE FROM ZINC WHILE DELETING DOCUMENT", http.StatusBadRequest)
		return
	}

	utils.WriteJson(w, "Document Deleted")
}
