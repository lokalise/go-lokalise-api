package lokalise_test

import "fmt"

func notFoundResponseBody(message string) string {
	return fmt.Sprintf(`{
		"error": {
			"code": 404,
			"message": "%s"
		}
	}`, message)
}
