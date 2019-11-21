package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCommentService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/comments", testProjectID, 12345),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"comments": [
					{
						"comment": "This is a test."
					}
				]
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"comments": [
				{
					"comment_id": 44444,
					"key_id": 12345,
					"comment": "This is a test.",
					"added_by": 420,
					"added_by_email": "commenter@mycompany.com",
					"added_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"added_at_timestamp": 1546257600
				}
			]
		}`)
		})

	r, err := client.Comments().Create(testProjectID, 12345, []NewComment{
		{Comment: "This is a test."},
	})
	if err != nil {
		t.Errorf("Comments.Create returned error: %v", err)
	}

	want := []Comment{
		{
			CommentID:    44444,
			KeyID:        12345,
			Comment:      "This is a test.",
			AddedBy:      420,
			AddedByEmail: "commenter@mycompany.com",
			AddedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
			AddedAtTs:    1546257600,
		},
	}

	if !reflect.DeepEqual(r.Comments, want) {
		t.Errorf("Comments.Create returned %+v, want %+v", r.Comments, want)
	}
}

func TestCommentService_ListProject(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/comments", testProjectID),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"comments": [
				{
					"comment_id": 44444
				}
			]
		}`)
		})

	r, err := client.Comments().ListProject(testProjectID)
	if err != nil {
		t.Errorf("Comments.ListProject returned error: %v", err)
	}

	want := []Comment{
		{
			CommentID: 44444,
		},
	}

	if !reflect.DeepEqual(r.Comments, want) {
		t.Errorf("Comments.ListProject returned %+v, want %+v", r.Comments, want)
	}
}

func TestCommentService_ListByKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/comments", testProjectID, 12345),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
			"project_id": "`+testProjectID+`",
			"comments": [
				{
					"comment_id": 44444
				}
			]
		}`)
		})

	r, err := client.Comments().ListByKey(testProjectID, 12345)
	if err != nil {
		t.Errorf("Comments.ListByKey returned error: %v", err)
	}

	want := []Comment{
		{
			CommentID: 44444,
		},
	}

	if !reflect.DeepEqual(r.Comments, want) {
		t.Errorf("Comments.ListByKey returned %+v, want %+v", r.Comments, want)
	}
}

func TestCommentService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/comments/%d", testProjectID, 12345, 44444),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"comment":	{
					"comment_id": 44444
				}
			}`)
		})

	r, err := client.Comments().Retrieve(testProjectID, 12345, 44444)
	if err != nil {
		t.Errorf("Comments.ListByKey returned error: %v", err)
	}

	want := Comment{
		CommentID: 44444,
	}

	if !reflect.DeepEqual(r.Comment, want) {
		t.Errorf("Comments.ListByKey returned %+v, want %+v", r.Comment, want)
	}
}

func TestCommentService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/projects/%s/keys/%d/comments/%d", testProjectID, 12345, 44444),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "DELETE")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"project_id": "`+testProjectID+`",
				"comment_deleted": true
			}`)
		},
	)

	r, err := client.Comments().Delete(testProjectID, 12345, 44444)
	if err != nil {
		t.Errorf("Comments.Delete returned error: %v", err)
	}

	want := DeleteCommentResponse{
		WithProjectID: WithProjectID{
			ProjectID: testProjectID,
		},
		IsDeleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("Comments.Delete returned %+v, want %+v", r, want)
	}
}
