package helpers

import "fmt"

// ResolveStatusCode resolve api error
func ResolveStatusCode(statusCode int) string {
	switch statusCode {
	case 400:
		return "HTTP 400 - Request Failed - Please check output for more information"
	case 401:
		return "HTTP 401 - Authorization Required - Provide Correct Username and API key"
	case 429:
		return "HTTP 429 - Payment Required - Amount of SSL tests exceeded on account. Check response for more information"
	default:
		returnMessage := fmt.Sprintf("HTTP %v - Unknown error from statuscake API", statusCode)
		return returnMessage
	}
}
