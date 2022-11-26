package transactions

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	zinc "github.com/zinclabs/sdk-go-zincsearch"
	"net/http"
	"sofa-logs-servers/infra/zincsearch"
	"sofa-logs-servers/models"
	"sofa-logs-servers/utils"
	"time"
)

type FindByIdForm struct {
	TransactionID string `json:"transaction_id"`
}

type FindAllForm struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
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
			Index     string             `json:"_index"`
			Id        string             `json:"_id"`
			Score     float64            `json:"_score"`
			Timestamp time.Time          `json:"@timestamp"`
			Source    models.Transaction `json:"_source"`
		}
	}
}

func FindAll(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {
	fmt.Println(r.Body)
	defer func() {
		err := r.Body.Close()
		if err != nil {
			utils.WriteErr(w, "ERROR IN REQUEST", http.StatusInternalServerError)
			return
		}
	}()

	form := FindAllForm{}

	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "WRONG BODY FORMAT", http.StatusBadRequest)
		return
	}
	query := *zinc.NewMetaZincQuery() // V1ZincQuery | Query
	subQuery := *zinc.NewMetaQuery()
	subQuery.SetMatchAll(map[string]interface {
	}{})

	query.SetQuery(subQuery)
	query.SetSort([]string{"+date"})

	_, res, err := zincClient.Client.Search.Search(zincClient.Ctx, "transactions").Query(query).Execute()

	if err != nil {
		utils.WriteErr(w, "Error searching the Document", http.StatusBadRequest)
		return
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	resDecoded := RespFind{}
	if err := json.NewDecoder(res.Body).Decode(&resDecoded); err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE FROM ELASTIC", http.StatusInternalServerError)
		return
	}

	var returnedArray []models.Transaction
	fmt.Println(form.StartDate, form.EndDate)
	for _, hit := range resDecoded.Hits.Hits {
		if hit.Source.Date.After(form.StartDate) && hit.Source.Date.Before(form.EndDate) {
			returnedArray = append(returnedArray, hit.Source)
		}
	}

	utils.WriteJson(w, returnedArray)
}

func FindById(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			utils.WriteErr(w, "ERROR IN REQUEST", http.StatusInternalServerError)
			return
		}
	}()

	form := FindByIdForm{}
	err := jsoniter.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		utils.WriteErr(w, "BAD BODY FORMAT", http.StatusBadRequest)
		return
	}

	if form.TransactionID == "" {
		utils.WriteErr(w, "TRANSACTIONS ID IS REQUIRED !!!", http.StatusBadRequest)
		return
	}

	query := *zinc.NewMetaZincQuery() // V1ZincQuery | Query
	metaQuery := *zinc.NewMetaTermQuery()
	metaQuery.SetValue(form.TransactionID)
	subQuery := *zinc.NewMetaQuery()
	subQuery.SetTerm(map[string]zinc.MetaTermQuery{
		"_id": metaQuery,
	})

	query.SetQuery(subQuery)

	_, res, err := zincClient.Client.Search.Search(zincClient.Ctx, "transactions").Query(query).Execute()
	fmt.Println(res)
	if err != nil {
		utils.WriteErr(w, "Error searching the Document", http.StatusBadRequest)
		return
	}

	defer func() {
		err := res.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	resDecoded := RespFind{}

	if err := json.NewDecoder(res.Body).Decode(&resDecoded); err != nil {
		utils.WriteErr(w, "ERROR PARSING RESPONSE", http.StatusInternalServerError)
		return
	}

	utils.WriteJson(w, resDecoded.Hits.Hits)

}
