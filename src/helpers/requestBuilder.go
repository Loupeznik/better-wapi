package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/loupeznik/better-wapi/src/models"
)

func BuildRequest(baseUrl string, model *models.Request) *http.Request {
	requestBody, err := json.Marshal(model)

	if err != nil {
		panic("Could not parse request body")
	}

	payload := strings.NewReader(fmt.Sprintf("request=%s", requestBody))

	request, err := http.NewRequest("POST", baseUrl, payload)

	if err != nil {
		panic(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return request
}
