package fetchService

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sjlleo/traceSysClient/model"
)

var (
	ErrFetchFail error = errors.New("network error")
	ErrParseFail error = errors.New("json parse error")
)

func FetchTraceList() (*model.TraceList, error) {
	// var ApiToken string = viper.Get("token").(string)
	// var ApiPrefix string = viper.Get("backcallurl").(string)
	// apiPath := ApiPrefix
	ApiToken := "cuxfDeNdaB4TZb8dDyKBD"
	ApiPrefix := "https://api.trace.ac"
	apiPath := "/api/tracelist/token/"
	url := ApiPrefix + apiPath + ApiToken
	log.Println(url)
	resp, err := http.Get(url)
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
