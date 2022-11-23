package fetchService

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsafe"

	"github.com/sjlleo/traceSysClient/model"
	"github.com/spf13/viper"
)

func PostResult(r *model.Report) error {

	var ApiToken string = viper.Get("token").(string)
	var ApiPrefix string = viper.Get("backcallurl").(string)
	apiPath := ApiPrefix + "/api/result/add"
	// ApiToken := "cuxfDeNdaB4TZb8dDyKBD"
	// ApiPrefix := "https://api.trace.ac"
	// apiPath := "/api/result/add"
	r.Token = ApiToken
	bytesData, err := json.Marshal(r)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", apiPath, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//byte数组直接转成string，优化内存
	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println(*str)
	return nil
}
