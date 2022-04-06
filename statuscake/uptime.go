package statuscake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"statuscakectl/helpers"
	"strconv"
	"strings"
)

// UptimeCheckResult response from api
type UptimeCheckResult struct {
	TestID      int     `json:"TestID"`
	Paused      bool    `json:"Paused"`
	TestType    string  `json:"TestType"`
	WebsiteName string  `json:"WebsiteName"`
	WebsiteURL  string  `json:"WebsiteURL"`
	ContactID   int     `json:"ContactID"`
	Status      string  `json:"Status"`
	Uptime      float64 `json:"Uptime"`
}

// UptimeCheckDetailed response from api
type UptimeCheckDetailed struct {
	Method        string `json:"Method"`
	TestID        int    `json:"TestID"`
	TestType      string `json:"TestType"`
	Paused        bool   `json:"Paused"`
	WebsiteName   string `json:"WebsiteName"`
	URI           string `json:"URI"`
	ContactGroup  string `json:"ContactGroup"`
	ContactID     int    `json:"ContactID"`
	ContactGroups []struct {
		ID    int    `json:"ID"`
		Name  string `json:"Name"`
		Email string `json:"Email"`
	} `json:"ContactGroups"`
	Status           string   `json:"Status"`
	Tags             []string `json:"Tags"`
	Uptime           int      `json:"Uptime"`
	CheckRate        int      `json:"CheckRate"`
	Timeout          int      `json:"Timeout"`
	LogoImage        string   `json:"LogoImage"`
	Confirmation     string   `json:"Confirmation"`
	FinalEndpoint    string   `json:"FinalEndpoint"`
	WebsiteHost      string   `json:"WebsiteHost"`
	NodeLocations    []string `json:"NodeLocations"`
	FindString       string   `json:"FindString"`
	DoNotFind        bool     `json:"DoNotFind"`
	LastTested       string   `json:"LastTested"`
	NextLocation     string   `json:"NextLocation"`
	Processing       bool     `json:"Processing"`
	ProcessingState  string   `json:"ProcessingState"`
	ProcessingOn     string   `json:"ProcessingOn"`
	DownTimes        string   `json:"DownTimes"`
	TriggerRate      string   `json:"TriggerRate"`
	Sensitive        bool     `json:"Sensitive"`
	EnableSSLWarning bool     `json:"EnableSSLWarning"`
	FollowRedirect   bool     `json:"FollowRedirect"`
	Dnsip            string   `json:"DNSIP"`
	DNSServer        string   `json:"DNSServer"`
	CustomHeader     string   `json:"CustomHeader"`
	PostRaw          string   `json:"PostRaw"`
	UseJar           int      `json:"UseJar"`
	StatusCodes      []string `json:"StatusCodes"`
}

// ListUptime listing all uptime checks in the account
func ListUptime(api, user, key string) []UptimeCheckResult {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", api+"/API/Tests/", nil)
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

	var uptimeCheckSlice []UptimeCheckResult
	e := json.Unmarshal(responseBody, &uptimeCheckSlice)
	if e != nil {
		log.Println(e)
		log.Println("Error: Failed to parse response body")
	}
	// debug
	// log.Println(string(string(responseBody)))

	return uptimeCheckSlice
}

// CreateUptimeCheck create an uptime check
func CreateUptimeCheck(domain string, checkrate, timeout, confirmation, virus, donotfind, realbrowser, trigger, sslalert, follow int, contacts, testType, findstring, api, user, key string) bool {
	target, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Please make sure to enter a valid domain (e.g https://www.domain.com)")
		return false
	}
	if target.Scheme == "" {
		fmt.Printf("Please add url scheme http/https to your domain %v\n", domain)
		return false
	}

	p := url.Values{}
	p.Add("WebsiteName", domain)
	p.Add("WebsiteURL", domain)
	p.Add("CheckRate", strconv.Itoa(checkrate))
	p.Add("Timeout", strconv.Itoa(timeout))
	p.Add("Confirmation", strconv.Itoa(confirmation))
	p.Add("Virus", strconv.Itoa(virus))
	p.Add("RealBrowser", strconv.Itoa(realbrowser))
	p.Add("TriggerRate", strconv.Itoa(trigger))
	if testType == "HTTP" {
		p.Add("EnableSSLAlert", strconv.Itoa(sslalert))
		p.Add("FollowRedirect", strconv.Itoa(follow))
	}
	p.Add("ContactGroup", contacts)
	p.Add("TestType", testType)
	if len(findstring) > 0 {
		p.Add("FindString", findstring)
		p.Add("DoNotFind", strconv.Itoa(donotfind))
	}
	payload := strings.NewReader(p.Encode())

	client := &http.Client{}
	request, _ := http.NewRequest("PUT", api+"/API/Tests/Update", payload)
	request.Header.Add("Username", user)
	request.Header.Add("API", key)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	// debug
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(responseBody))

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return false
	}
	return true
}

// DeleteUptimeCheck delete uptime check by id
func DeleteUptimeCheck(id int, api, user, key string) bool {
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", api+"/API/Tests/Details/", nil)
	// headers
	request.Header.Add("Username", user)
	request.Header.Add("API", key)
	// params
	q := request.URL.Query()
	q.Add("TestID", strconv.Itoa(id))
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	// debug
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(responseBody))

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return false
	}
	return true
}

// GetDetailedTestData gets detailed uptime data for test
func GetDetailedTestData(id int, api, user, key string) UptimeCheckDetailed {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", api+"/API/Tests/Details/", nil)
	// headers
	request.Header.Add("Username", user)
	request.Header.Add("API", key)
	// params
	q := request.URL.Query()
	q.Add("TestID", strconv.Itoa(id))
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return UptimeCheckDetailed{}
	}
	defer resp.Body.Close()

	// debug
	// responseBody, _ := ioutil.ReadAll(resp.Body)
	// log.Println(string(responseBody))

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return UptimeCheckDetailed{}
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	var uptimeCheckData UptimeCheckDetailed
	e := json.Unmarshal(responseBody, &uptimeCheckData)
	if e != nil {
		log.Println(e)
		log.Println("Error: Failed to parse response body")
	}

	return uptimeCheckData
}
