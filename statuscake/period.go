package statuscake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"statuscakectl/helpers"
)

// Period object response from api
type Period struct {
	Status     string `json:"Status"`
	StatusID   string `json:"StatusID"`
	Start      string `json:"Start"`
	End        string `json:"End"`
	StartUnix  int    `json:"Start_Unix"`
	EndUnix    int    `json:"End_Unix"`
	Additional string `json:"Additional"`
	Period     string `json:"Period"`
}

type PeriodErrorMessage struct {
	ErrNo  int    `json:"ErrNo"`
	Access bool   `json:"Access"`
	Client string `json:"Client"`
	TestID string `json:"TestID"`
	Error  string `json:"Error"`
}

// ListPeriods listing all period objects for test
func ListPeriods(api, user, key string, testID int) []Period {
	client := &http.Client{}
	urlString := fmt.Sprintf("%s/API/Tests/Periods/?TestID=%v", api, testID)
	request, _ := http.NewRequest("GET", urlString, nil)
	request.Header.Add("Username", user)
	request.Header.Add("API", key)

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var Periods []Period

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return Periods
	}

	e := json.Unmarshal(responseBody, &Periods)
	if e != nil {

		// debug
		// log.Println(string(responseBody))

		var errorResponse PeriodErrorMessage
		e = json.Unmarshal(responseBody, &errorResponse)
		if e != nil {
			log.Println("Failed to parse response body")
			return nil
		}

		log.Printf("Failed to get periods: %v", errorResponse.Error)

		return nil

	}

	return Periods
}
