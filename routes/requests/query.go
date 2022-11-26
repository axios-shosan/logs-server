package requests

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	zinc "github.com/zinclabs/sdk-go-zincsearch"
	"net/http"
	"sofa-logs-servers/infra/zincsearch"
	"sofa-logs-servers/models"
	"sofa-logs-servers/utils"
)

type FindForm struct {
	RequestID string `json:"request_id"`
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

func FindById(w http.ResponseWriter, r *http.Request, zincClient zincsearch.ZincClient) {
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

	if form.RequestID == "" {
		utils.WriteErr(w, "REQUEST ID EMPTY !!!", http.StatusBadRequest)
		return
	}

	query := *zinc.NewMetaZincQuery() // V1ZincQuery | Query
	metaQuery := *zinc.NewMetaTermQuery()
	metaQuery.SetValue(form.RequestID)
	subQuery := *zinc.NewMetaQuery()
	subQuery.SetTerm(map[string]zinc.MetaTermQuery{
		"_id": metaQuery,
	})

	query.SetQuery(subQuery)

	_, res, err := zincClient.Client.Search.Search(zincClient.Ctx, "requests").Query(query).Execute()
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

func FindAll(w http.ResponseWriter, _ *http.Request, zincClient zincsearch.ZincClient) {

	query := *zinc.NewMetaZincQuery() // V1ZincQuery | Query
	subQuery := *zinc.NewMetaQuery()
	subQuery.SetMatchAll(map[string]interface {
	}{})
	query.SetQuery(subQuery)

	_, res, err := zincClient.Client.Search.Search(zincClient.Ctx, "requests").Query(query).Execute()

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

	utils.WriteJson(w, resDecoded.Hits.Hits)
}
