package fetchService

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sjlleo/traceSysClient/model"
	"github.com/spf13/viper"
)

var (
	ErrFetchFail error = errors.New("network error")
	ErrParseFail error = errors.New("json parse error")
)

func FetchTraceList() (*model.TraceList, error) {
	var ApiToken string = viper.Get("token").(string)
	var ApiPrefix string = viper.Get("backcallurl").(string)
	apiPath := ApiPrefix + "/api/tracelist/token/" + ApiToken
	// ApiToken := "cuxfDeNdaB4TZb8dDyKBD"
	// ApiPrefix := "https://api.trace.ac"
	log.Println(apiPath)
	resp, err := http.Get(apiPath)
	if err != nil {
		fmt.Println(err)
		return nil, ErrFetchFail
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var res model.TraceList
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, ErrParseFail
	}
	return &res, nil
}
