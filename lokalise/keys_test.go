package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Keys_List(t *testing.T) {
	type input struct {
		projectID string
		options   lokalise.ListKeysOptions
	}
	type output struct {
		expectedOutgoingRequest outgoingRequest
		response                model.KeysResponse
		err                     error
	}
	tt := []struct {
		name           string
		input          input
		serverResponse serverListResponse
		output         output
	}{
		{
			name: "no pagination succesful json response",
			input: input{
				projectID: "20008339586cded200e0d8.29879849",
			},
			serverResponse: serverListResponse{
				statusCode: http.StatusOK,
				body: `{
					"project_id": "20008339586cded200e0d8.29879849",
					"keys":[{
						"translations": [{
							"translation_id": 344412,
							"key_id": 553662,
							"language_iso": "en_US",
							"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
							"modified_by": 420,
							"modified_by_email": "user@mycompany.com",
							"translation": "Hello, world!",
							"is_fuzzy": true,
							"is_reviewed": false,
							"words": 2
						}]
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/20008339586cded200e0d8.29879849/keys",
					method: http.MethodGet,
				},
				response: model.KeysResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Keys: []model.Key{
						model.Key{
							Translations: []model.Translation{
								{
									TranslationID:   344412,
									KeyID:           553662,
									LanguageISO:     "en_US",
									ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
									ModifiedBy:      420,
									ModifiedByEmail: "user@mycompany.com",
									Translation:     "Hello, world!",
									IsFuzzy:         true,
									IsReviewed:      false,
									Words:           2,
								},
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "with options succesful json response",
			input: input{
				projectID: "20008339586cded200e0d8.29879849",
				options: lokalise.ListKeysOptions{
					DisableReferences: true,
					PageOptions: lokalise.PageOptions{
						Limit: 1,
						Page:  2,
					},
				},
			},
			serverResponse: serverListResponse{
				statusCode:            http.StatusOK,
				pagedTotalCountHeader: "1",
				pagedPageCountHeader:  "2",
				pagedLimitHeader:      "3",
				pagedPageHeader:       "4",
				body: `{
					"project_id": "20008339586cded200e0d8.29879849",
					"keys":[{
						"translations": [{
							"translation_id": 344412,
							"key_id": 553662,
							"language_iso": "en_US",
							"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
							"modified_by": 420,
							"modified_by_email": "user@mycompany.com",
							"translation": "Hello, world!",
							"is_fuzzy": true,
							"is_reviewed": false,
							"words": 2
						}]
					}]
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/20008339586cded200e0d8.29879849/keys?disable_references=1&limit=1&page=2",
					method: http.MethodGet,
				},
				response: model.KeysResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
					Keys: []model.Key{
						model.Key{
							Translations: []model.Translation{
								{
									TranslationID:   344412,
									KeyID:           553662,
									LanguageISO:     "en_US",
									ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
									ModifiedBy:      420,
									ModifiedByEmail: "user@mycompany.com",
									Translation:     "Hello, world!",
									IsFuzzy:         true,
									IsReviewed:      false,
									Words:           2,
								},
							},
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "404 error response",
			input: input{
				projectID: "12345",
			},
			serverResponse: serverListResponse{
				statusCode: http.StatusNotFound,
				body:       notFoundResponseBody("keys not found"),
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/12345/keys",
					method: http.MethodGet,
				},
				response: model.KeysResponse{
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
				},
				err: &model.Error{
					Code:    404,
					Message: "keys not found",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			client, fixture, close := setupListClient(t, tc.serverResponse)
			defer close()

			resp, err := client.Keys.List(context.Background(), tc.input.projectID, tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			tc.output.expectedOutgoingRequest.Assert(t, fixture)
			assert.Equal(tc.output.response.Keys, resp.Keys, "response not as expected")
		})
	}
}

func TestClient_Keys_Create(t *testing.T) {
	inputProjectID := "1"
	inputKeys := []model.Key{
		model.Key{
			KeyName:   "testName",
			Platforms: []string{"ios", "android"},
			Translations: []model.Translation{
				{
					LanguageISO: "en_US",
					Translation: "Hello, world!",
				},
			},
		},
	}
	mockedServerResponseBody := `{
    "project_id": "3002780358964f9bab5a92.87762498",
    "keys":[{
		"translations": [{
			"translation_id": 344412,
			"key_id": 553662,
			"language_iso": "en_US",
			"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
			"modified_by": 420,
			"modified_by_email": "user@mycompany.com",
			"translation": "Hello, world!",
			"is_fuzzy": true,
			"is_reviewed": false,
			"words": 2
		}]
	}],
	"errors":[{
		"message":"test msg",
		"code":400,
		"key":{
			"key_name":"name"
		}
	}]
	}`
	expectedOutgoingRequest := outgoingRequest{
		method: http.MethodPost,
		path:   "/projects/1/keys",
		body:   `{"keys":[{"key_name":"testName","filenames":{},"platforms":["ios","android"],"translations":[{"language_iso":"en_US","translation":"Hello, world!"}]}]}`,
	}
	expectedResult := model.KeysResponse{
		ProjectID: "3002780358964f9bab5a92.87762498",
		Keys: []model.Key{
			model.Key{
				Translations: []model.Translation{
					{
						TranslationID:   344412,
						KeyID:           553662,
						LanguageISO:     "en_US",
						ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
						ModifiedBy:      420,
						ModifiedByEmail: "user@mycompany.com",
						Translation:     "Hello, world!",
						IsFuzzy:         true,
						IsReviewed:      false,
						Words:           2,
					},
				},
			},
		},
		Errors: []model.ErrorKeys{
			model.ErrorKeys{
				model.Error{
					Message: "test msg",
					Code:    400,
				},
				model.ErrorKey{
					KeyName: "name",
				},
			},
		},
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Keys.Create(context.Background(), inputProjectID, inputKeys)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Keys, resp.Keys, "response keys not as expected")
}
