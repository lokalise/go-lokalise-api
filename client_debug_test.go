package lokalise

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

// bufferLogger captures resty debug output.
type bufferLogger struct{ buf bytes.Buffer }

func (l *bufferLogger) Errorf(format string, v ...interface{}) { fmt.Fprintf(&l.buf, format, v...) }
func (l *bufferLogger) Warnf(format string, v ...interface{})  { fmt.Fprintf(&l.buf, format, v...) }
func (l *bufferLogger) Debugf(format string, v ...interface{}) { fmt.Fprintf(&l.buf, format, v...) }

func TestWithDebug_RedactsApiToken(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	if err := WithDebug(true)(client); err != nil {
		t.Fatal(err)
	}
	logger := &bufferLogger{}
	client.httpClient.SetLogger(logger)

	mux.HandleFunc("/projects/"+testProjectID, func(w http.ResponseWriter, r *http.Request) {
		// Redaction only applies to the log: the request must still carry the real token.
		testHeader(t, r, apiTokenHeader, testApiToken)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"project_id": "`+testProjectID+`"}`)
	})

	if _, err := client.Projects().Retrieve(testProjectID); err != nil {
		t.Fatalf("Projects.Retrieve returned error: %v", err)
	}

	out := logger.buf.String()
	if !strings.Contains(out, "~~~ REQUEST ~~~") {
		t.Fatalf("expected debug output, got %q", out)
	}
	if strings.Contains(out, testApiToken) {
		t.Errorf("debug output leaks the API token: %q", out)
	}
	if !strings.Contains(out, apiTokenHeader+": "+redactedValue) {
		t.Errorf("expected %s to be logged as %s, got %q", apiTokenHeader, redactedValue, out)
	}
}
