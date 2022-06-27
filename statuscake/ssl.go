package statuscake

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"statuscakectl/helpers"
	"strings"
)

// SSLCheckResult result from api
type SSLCheckResult struct {
	ID            string        `json:"id"`
	Checkrate     int           `json:"checkrate"`
	Paused        bool          `json:"paused"`
	Domain        string        `json:"domain"`
	IssuerCn      string        `json:"issuer_cn"`
	CertScore     string        `json:"cert_score"`
	CipherScore   string        `json:"cipher_score"`
	CertStatus    string        `json:"cert_status"`
	Cipher        string        `json:"cipher"`
	ValidFromUtc  string        `json:"valid_from_utc"`
	ValidUntilUtc string        `json:"valid_until_utc"`
	MixedContent  []interface{} `json:"mixed_content"`
	Flags         struct {
		IsExtended bool `json:"is_extended"`
		HasPfs     bool `json:"has_pfs"`
		IsBroken   bool `json:"is_broken"`
		IsExpired  bool `json:"is_expired"`
		IsMissing  bool `json:"is_missing"`
		IsRevoked  bool `json:"is_revoked"`
		HasMixed   bool `json:"has_mixed"`
	} `json:"flags"`
	ContactGroups  []interface{} `json:"contact_groups"`
	AlertAt        string        `json:"alert_at"`
	LastReminder   int           `json:"last_reminder"`
	AlertReminder  bool          `json:"alert_reminder"`
	AlertExpiry    bool          `json:"alert_expiry"`
	AlertBroken    bool          `json:"alert_broken"`
	AlertMixed     bool          `json:"alert_mixed"`
	LastUpdatedUtc string        `json:"last_updated_utc"`
}

// ListSSLChecks listing all ssl checks in the account
func ListSSLChecks(api, user, key string) []SSLCheckResult {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", api+"/API/SSL/", nil)
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
	var sslCheckSlice []SSLCheckResult

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return sslCheckSlice
	}

	e := json.Unmarshal(responseBody, &sslCheckSlice)
	if e != nil {
		log.Println("Failed to parse response body")
	}
	// debug
	// log.Println(string(responseBody))

	return sslCheckSlice
}

// CreateSSLStatuscakeCheck to create a ssl certificate check in statuscake
func CreateSSLStatuscakeCheck(domain string, checkrate int, contacts string, api string, user string, key string) bool {
	target, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Please make sure to enter a valid domain (e.g domain.com)")
		return false
	}
	target.Scheme = "https"
	domain = fmt.Sprint(target)

	p := url.Values{}
	p.Add("domain", domain)
	p.Add("checkrate", fmt.Sprintf("%v", checkrate))
	p.Add("contact_groups", contacts)
	p.Add("alert_at", "1,7,30")
	p.Add("alert_broken", "true")
	p.Add("alert_expiry", "true")
	p.Add("alert_reminder", "true")
	p.Add("alert_mixed", "true")

	payload := strings.NewReader(p.Encode())

	client := &http.Client{}
	request, _ := http.NewRequest("PUT", api+"/API/SSL/Update", payload)
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

// DeleteSSLCheck delete ssl check by id
func DeleteSSLCheck(id, api, user, key string) bool {
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", api+"/API/SSL/Update", nil)
	// headers
	request.Header.Add("Username", user)
	request.Header.Add("API", key)
	// params
	q := request.URL.Query()
	q.Add("id", id)
	request.URL.RawQuery = q.Encode()

	resp, err := client.Do(request)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	// debug
	// responseBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		message := helpers.ResolveStatusCode(resp.StatusCode)
		log.Println(message)
		return false
	}
	return true
}
