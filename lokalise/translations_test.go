package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/17media/go-lokalise-api/lokalise"
	"github.com/17media/go-lokalise-api/model"
	"github.com/stretchr/testify/assert"
)

func TestClient_Translations_List(t *testing.T) {
	type input struct {
		projectID string
		options   lokalise.TranslationsOptions
	}
	type output struct {
		expectedOutgoingRequest outgoingRequest
		response                model.TranslationsResponse
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
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/20008339586cded200e0d8.29879849/translations",
					method: http.MethodGet,
				},
				response: model.TranslationsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
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
				err: nil,
			},
		},
		{
			name: "with options succesful json response",
			input: input{
				projectID: "20008339586cded200e0d8.29879849",
				options: lokalise.TranslationsOptions{
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
				}`,
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/20008339586cded200e0d8.29879849/translations?disable_references=1&limit=1&page=2",
					method: http.MethodGet,
				},
				response: model.TranslationsResponse{
					Paged: model.Paged{
						TotalCount: 1,
						PageCount:  2,
						Limit:      3,
						Page:       4,
					},
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
				body:       notFoundResponseBody("translation not found"),
			},
			output: output{
				expectedOutgoingRequest: outgoingRequest{
					path:   "/projects/12345/translations",
					method: http.MethodGet,
				},
				response: model.TranslationsResponse{
					Paged: model.Paged{
						TotalCount: -1,
						PageCount:  -1,
						Limit:      -1,
						Page:       -1,
					},
				},
				err: &model.Error{
					Code:    404,
					Message: "translation not found",
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			assert := assert.New(t)
			client, fixture, close := setupListClient(t, tc.serverResponse)
			defer close()

			resp, err := client.Translations.List(context.Background(), tc.input.projectID, tc.input.options)

			if tc.output.err != nil {
				assert.EqualError(err, tc.output.err.Error(), "output error not as expected")
			} else {
				assert.NoError(err, "output error not expected")
			}
			tc.output.expectedOutgoingRequest.Assert(t, fixture)
			assert.Equal(tc.output.response.Translations, resp.Translations, "response not as expected")
		})
	}
}

func TestClient_Translations_Retrieve(t *testing.T) {
	inputProjectID := "1"
	inputTranslationID := int64(2)
	mockedServerResponseBody := `{
    "project_id": "3002780358964f9bab5a92.87762498",
    "translation": {
			"translation_id": 344412,
			"key_id": 553662,
			"language_iso": "en_US",
			"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
			"modified_by": 420,
			"modified_by_email": "user@mycompany.com",
			"translation": "Hello, world!",
			"is_fuzzy": true,
			"is_reviewed": true,
			"words": 2
    }
	}`
	expectedOutgoingRequest := outgoingRequest{
		method: http.MethodGet,
		path:   "/projects/1/translations/2",
	}
	expectedResult := model.TranslationResponse{
		ProjectID: "3002780358964f9bab5a92.87762498",
		Translation: model.Translation{
			TranslationID:   344412,
			KeyID:           553662,
			LanguageISO:     "en_US",
			ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
			ModifiedBy:      420,
			ModifiedByEmail: "user@mycompany.com",
			Translation:     "Hello, world!",
			IsFuzzy:         true,
			IsReviewed:      true,
			Words:           2,
		},
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Translations.Retrieve(context.Background(), inputProjectID, inputTranslationID)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Translation, resp.Translation, "response translation not as expected")
}

func TestClient_Translations_Update(t *testing.T) {
	inputProjectID := "1"
	inputTranslationID := int64(2)
	inputTranslation := "updated translation"
	inputIsFuzzy := true
	inputIsReviewed := true
	mockedServerResponseBody := `{
    "project_id": "3002780358964f9bab5a92.87762498",
    "translation": {
			"translation_id": 344412,
			"key_id": 553662,
			"language_iso": "en_US",
			"modified_at": "2018-12-31 12:00:00 (Etc/UTC)",
			"modified_by": 420,
			"modified_by_email": "user@mycompany.com",
			"translation": "Hello, world!",
			"is_fuzzy": true,
			"is_reviewed": true,
			"words": 2
    }
	}`
	expectedOutgoingRequest := outgoingRequest{
		method: http.MethodPut,
		path:   "/projects/1/translations/2",
		body:   `{"is_fuzzy":true,"is_reviewed":true,"translation":"updated translation"}`,
	}
	expectedResult := model.TranslationResponse{
		ProjectID: "3002780358964f9bab5a92.87762498",
		Translation: model.Translation{
			TranslationID:   344412,
			KeyID:           553662,
			LanguageISO:     "en_US",
			ModifiedAt:      "2018-12-31 12:00:00 (Etc/UTC)",
			ModifiedBy:      420,
			ModifiedByEmail: "user@mycompany.com",
			Translation:     "Hello, world!",
			IsFuzzy:         true,
			IsReviewed:      true,
			Words:           2,
		},
	}
	assert := assert.New(t)
	client, fixture, close := setupClient(t, mockedServerResponseBody)
	defer close()

	resp, err := client.Translations.Update(context.Background(), inputProjectID, inputTranslationID, inputTranslation, inputIsFuzzy, inputIsReviewed)

	assert.NoError(err, "output error not expected")
	expectedOutgoingRequest.Assert(t, fixture)
	assert.Equal(expectedResult.ProjectID, resp.ProjectID, "response project id not as expected")
	assert.Equal(expectedResult.Translation, resp.Translation, "response translation not as expected")
}
