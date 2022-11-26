package requests

import (
	"net/http"
	"sofa-logs-servers/infra/zincsearch"
	"sofa-logs-servers/utils"
	"time"

	jsoniter "github.com/json-iterator/go"
)

type CreateForm struct {
	UserID    uint      `json:"user_id"`
	Page      string    `json:"page"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

type UpdateForm struct {
	ID        string    `json:"request_id"`
	UserID    uint      `json:"user_id"`
	Page      string    `json:"page"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

type DeleteFrom struct {
	ID string `json:"request_id"`
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

	if form.Page == "" {
		utils.WriteErr(w, "EMPTY Page", http.StatusBadRequest)
		return
	}

	if form.StartedAt.IsZero() {
		utils.WriteErr(w, "EMPTY Start Request Time", http.StatusBadRequest)
	}

	if form.EndedAt.IsZero() {
		utils.WriteErr(w, "EMPTY End Request Time", http.StatusBadRequest)
	}

	document := map[string]interface{}{
		"UserID":    form.UserID,
		"Page":      form.Page,
		"StartedAt": form.StartedAt,
		"EndedAt":   form.EndedAt,
	} // map[string]interface{} | Document

	_, res, err := zincClient.Client.Document.Index(zincClient.Ctx, "requests").Document(document).Execute()

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

	if form.ID == "" {
		utils.WriteErr(w, "request_id IS REQUIRED", http.StatusBadRequest)
		return
	}

	if form.Page == "" {
		utils.WriteErr(w, "page IS REQUIRED", http.StatusBadRequest)
		return
	}

	if form.StartedAt.IsZero() {
		utils.WriteErr(w, "EMPTY Start Request Time", http.StatusBadRequest)
	}

	if form.EndedAt.IsZero() {
		utils.WriteErr(w, "EMPTY End Request Time", http.StatusBadRequest)
	}

	document := map[string]interface{}{
		"UserID":    form.UserID,
		"Page":      form.Page,
		"StartedAt": form.StartedAt,
		"EndedAt":   form.EndedAt,
	} // map[string]interface{} | Document

	_, res, err := zincClient.Client.Document.Update(zincClient.Ctx, "requests", form.ID).Document(document).Execute()

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

	_, res, err := zincClient.Client.Document.Delete(zincClient.Ctx, "requests", form.ID).Execute()
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
