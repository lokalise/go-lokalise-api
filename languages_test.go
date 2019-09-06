package lokalise_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/lokalise/go-lokalise-api"
	"github.com/stretchr/testify/assert"
)

func TestClient_Languages_SystemLanguages(t *testing.T) {
	t.Run("no pagination", func(t *testing.T) {
		assert := assert.New(t)
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusOK,
			body: ` {
				    "languages": [
					{
					    "lang_id": 640,
					    "lang_iso": "en",
					    "lang_name": "English",
					    "is_rtl": false,
					    "plural_forms": [
						"one", "other"
					    ]
					},
					{
					    "lang_id": 597,
					    "lang_iso": "ru",
					    "lang_name": "Russian",
					    "is_rtl": true,
					    "plural_forms": [
						"one", "few", "many", "other"
					    ]
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/system/languages",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListLanguagesResponse{
			Paged: lokalise.Paged{
				TotalCount: -1,
				PageCount:  -1,
				Limit:      -1,
				Page:       -1,
			},
			Languages: []lokalise.Language{
				{
					LangID:      640,
					LangISO:     "en",
					LangName:    "English",
					IsRTL:       false,
					PluralForms: []string{"one", "other"},
				},
				{
					LangID:      597,
					LangISO:     "ru",
					LangName:    "Russian",
					IsRTL:       true,
					PluralForms: []string{"one", "few", "many", "other"},
				},
			},
		}

		resp, err := client.Languages.SystemLanguages(context.Background(), lokalise.LanguagesOptions{})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("with pagination", func(t *testing.T) {
		assert := assert.New(t)
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode:            http.StatusOK,
			pagedTotalCountHeader: "2",
			pagedPageCountHeader:  "2",
			pagedLimitHeader:      "1",
			pagedPageHeader:       "2",
			body: ` {
				    "languages": [
					{
					    "lang_id": 597,
					    "lang_iso": "ru",
					    "lang_name": "Russian",
					    "is_rtl": true,
					    "plural_forms": [
						"one", "few", "many", "other"
					    ]
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/system/languages?limit=1&page=2",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListLanguagesResponse{
			Paged: lokalise.Paged{
				TotalCount: 2,
				PageCount:  2,
				Limit:      1,
				Page:       2,
			},
			Languages: []lokalise.Language{{
				LangID:      597,
				LangISO:     "ru",
				LangName:    "Russian",
				IsRTL:       true,
				PluralForms: []string{"one", "few", "many", "other"},
			},
			},
		}

		resp, err := client.Languages.SystemLanguages(context.Background(), lokalise.LanguagesOptions{
			PageOptions: lokalise.PageOptions{
				Limit: 1,
				Page:  2,
			},
		})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("http error", func(t *testing.T) {
		assert := assert.New(t)
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusNotFound,
			body:       notFoundResponseBody("Languages not found"),
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/system/languages",
			method: http.MethodGet,
		}
		expectedError := &lokalise.Error{Code: 404, Message: "Languages not found"}

		_, err := client.Languages.SystemLanguages(context.Background(), lokalise.LanguagesOptions{})
		assert.EqualError(err, expectedError.Error())
		expectedRequest.Assert(t, fixture)
	})
}

func TestClient_Languages_List(t *testing.T) {
	t.Run("no pagination", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusOK,
			body: ` {
				    "project_id": "20008339586cded200e0d8.29879849",
				    "languages": [
					{
					    "lang_id": 640,
					    "lang_iso": "en",
					    "lang_name": "English",
					    "is_rtl": false,
					    "plural_forms": [
						"one", "other"
					    ]
					},
					{
					    "lang_id": 597,
					    "lang_iso": "ru",
					    "lang_name": "Russian",
					    "is_rtl": true,
					    "plural_forms": [
						"one", "few", "many", "other"
					    ]
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/languages",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListLanguagesResponse{
			ProjectID: projectID,
			Paged: lokalise.Paged{
				TotalCount: -1,
				PageCount:  -1,
				Limit:      -1,
				Page:       -1,
			},
			Languages: []lokalise.Language{
				{
					LangID:      640,
					LangISO:     "en",
					LangName:    "English",
					IsRTL:       false,
					PluralForms: []string{"one", "other"},
				},
				{
					LangID:      597,
					LangISO:     "ru",
					LangName:    "Russian",
					IsRTL:       true,
					PluralForms: []string{"one", "few", "many", "other"},
				},
			},
		}

		resp, err := client.Languages.List(context.Background(), projectID, lokalise.LanguagesOptions{})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("with pagination", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode:            http.StatusOK,
			pagedTotalCountHeader: "2",
			pagedPageCountHeader:  "2",
			pagedLimitHeader:      "1",
			pagedPageHeader:       "2",
			body: ` {
				    "project_id": "20008339586cded200e0d8.29879849",
				    "languages": [
					{
					    "lang_id": 597,
					    "lang_iso": "ru",
					    "lang_name": "Russian",
					    "is_rtl": true,
					    "plural_forms": [
						"one", "few", "many", "other"
					    ]
					}
				    ]
				}`,
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/languages?limit=1&page=2",
			method: http.MethodGet,
		}
		expectedResponse := lokalise.ListLanguagesResponse{
			Paged: lokalise.Paged{
				TotalCount: 2,
				PageCount:  2,
				Limit:      1,
				Page:       2,
			},
			ProjectID: projectID,
			Languages: []lokalise.Language{{
				LangID:      597,
				LangISO:     "ru",
				LangName:    "Russian",
				IsRTL:       true,
				PluralForms: []string{"one", "few", "many", "other"},
			}},
		}

		resp, err := client.Languages.List(context.Background(), projectID, lokalise.LanguagesOptions{
			PageOptions: lokalise.PageOptions{
				Limit: 1,
				Page:  2,
			},
		})
		assert.NoError(err)
		expectedRequest.Assert(t, fixture)
		assert.Equal(expectedResponse, resp)
	})
	t.Run("http error", func(t *testing.T) {
		assert := assert.New(t)
		const projectID = "20008339586cded200e0d8.29879849"
		client, fixture, stop := setupListClient(t, serverListResponse{
			statusCode: http.StatusNotFound,
			body:       notFoundResponseBody("Languages not found"),
		})
		defer stop()
		expectedRequest := &outgoingRequest{
			path:   "/projects/20008339586cded200e0d8.29879849/languages",
			method: http.MethodGet,
		}
		expectedError := &lokalise.Error{Code: 404, Message: "Languages not found"}

		_, err := client.Languages.List(context.Background(), projectID, lokalise.LanguagesOptions{})
		assert.EqualError(err, expectedError.Error())
		expectedRequest.Assert(t, fixture)
	})
}

func TestClient_Languages_Create(t *testing.T) {
	assert := assert.New(t)
	const projectID = "3002780358964f9bab5a92.87762498"
	client, fixture, stop := setupClient(t, `
		{
		    "project_id": "3002780358964f9bab5a92.87762498",
		    "languages": [
			{
			    "lang_id": 640,
			    "lang_iso": "en",
			    "lang_name": "English",
			    "is_rtl": false,
			    "plural_forms": [
				"one", "other"
			    ]
			},
			{
			    "lang_id": 597,
			    "lang_iso": "ru",
			    "lang_name": "Russian",
			    "is_rtl": false,
			    "plural_forms": [
				"one", "few", "many", "other"
			    ]
			}
		    ]
		}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92.87762498/languages",
		method: http.MethodPost,
		body:   `{"languages":[{"lang_iso":"en"},{"lang_iso":"ru"}]}`,
	}
	expectedResponse := lokalise.CreateLanguageResponse{
		ProjectID: projectID,
		Languages: []lokalise.Language{
			{
				LangID:      640,
				LangISO:     "en",
				LangName:    "English",
				IsRTL:       false,
				PluralForms: []string{"one", "other"},
			},
			{
				LangID:      597,
				LangISO:     "ru",
				LangName:    "Russian",
				IsRTL:       false,
				PluralForms: []string{"one", "few", "many", "other"},
			},
		},
	}

	resp, err := client.Languages.Create(context.Background(), projectID, []lokalise.CustomLanguage{
		{LangISO: "en"}, {LangISO: "ru"},
	})
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}

func TestClient_Languages_Retrieve(t *testing.T) {
	assert := assert.New(t)
	const projectID = "3002780358964f9bab5a92.87762498"
	client, fixture, stop := setupClient(t, `
		{
		    "project_id": "3002780358964f9bab5a92.87762498",
		    "language": {
			    "lang_id": 640,
			    "lang_iso": "en",
			    "lang_name": "English",
			    "is_rtl": false,
			    "plural_forms": [
				"one", "other"
			    ]
		    }
		}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92.87762498/languages/42",
		method: http.MethodGet,
	}
	expectedResponse := lokalise.RetrieveLanguageResponse{
		ProjectID: projectID,
		Language: lokalise.Language{
			LangID:      640,
			LangISO:     "en",
			LangName:    "English",
			IsRTL:       false,
			PluralForms: []string{"one", "other"},
		},
	}

	resp, err := client.Languages.Retrieve(context.Background(), projectID, 42)
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}

func TestClient_Languages_Update(t *testing.T) {
	assert := assert.New(t)
	const projectID = "3002780358964f9bab5a92.87762498"
	client, fixture, stop := setupClient(t, `
		{
		    "project_id": "3002780358964f9bab5a92.87762498",
		    "language": {
			"lang_id": 640,
			"lang_iso": "en-US",
			"lang_name": "English",
			"is_rtl": false,
			"plural_forms": [
			    "one", "zero", "few", "other"
			]
		    }
		}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92.87762498/languages/42",
		method: http.MethodPut,
		body:   `{"lang_iso":"en-US","plural_forms":["one","zero","few","other"]}`,
	}
	expectedResponse := lokalise.UpdateLanguageResponse{
		ProjectID: projectID,
		Language: lokalise.Language{
			LangID:      640,
			LangISO:     "en-US",
			LangName:    "English",
			IsRTL:       false,
			PluralForms: []string{"one", "zero", "few", "other"},
		},
	}

	resp, err := client.Languages.Update(context.Background(), projectID, 42, lokalise.Language{
		LangISO:     "en-US",
		PluralForms: []string{"one", "zero", "few", "other"},
	})
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}

func TestClient_Languages_Delete(t *testing.T) {
	assert := assert.New(t)
	const projectID = "3002780358964f9bab5a92.87762498"
	client, fixture, stop := setupClient(t, `
		{
		    "project_id": "3002780358964f9bab5a92.87762498",
		    "language_deleted": true
		}`,
	)
	defer stop()
	expectedRequest := &outgoingRequest{
		path:   "/projects/3002780358964f9bab5a92.87762498/languages/42",
		method: http.MethodDelete,
	}
	expectedResponse := lokalise.DeleteLanguageResponse{
		ProjectID:       projectID,
		LanguageDeleted: true,
	}

	resp, err := client.Languages.Delete(context.Background(), projectID, 42)
	assert.NoError(err)
	expectedRequest.Assert(t, fixture)
	assert.Equal(expectedResponse, resp)
}
