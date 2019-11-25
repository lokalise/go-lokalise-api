package lokalise

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestPaymentCardService_Create(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/payment_cards",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "POST")
			testHeader(t, r, apiTokenHeader, testApiToken)
			data := `{
				"number": "4242424242424242",
				"cvc": "123",
				"exp_month": 5,
				"exp_year": 2021
			}`

			req := new(bytes.Buffer)
			_ = json.Compact(req, []byte(data))

			testBody(t, r, req.String())

			_, _ = fmt.Fprint(w, `{
				"card_id": 1234,
				"last4": "4242",
				"brand": "Visa",
				"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
				"created_at_timestamp": 1546257600
			}`)
		},
	)

	r, err := client.PaymentCards().Create(CreatePaymentCard{
		Number:   "4242424242424242",
		CVC:      "123",
		ExpMonth: 5,
		ExpYear:  2021,
	})
	if err != nil {
		t.Errorf("PaymentCards.Create returned error: %v", err)
	}

	want := PaymentCard{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		CardId: 1234,
		Last4:  "4242",
		Brand:  "Visa",
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("PaymentCards.Create returned %+v, want %+v", r, want)
	}
}

func TestPaymentCardService_Delete(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/payment_cards/%d", 22193), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		testMethod(t, r, "DELETE")
		testHeader(t, r, apiTokenHeader, testApiToken)
		body := `{
			"card_id": 22193,
			"card_deleted": true
		}`
		_, _ = fmt.Fprint(w, string(body))
	})

	r, err := client.PaymentCards().Delete(22193)
	if err != nil {
		t.Errorf("PaymentCards.Delete returned error: %v", err)
	}

	want := DeletePaymentCardResponse{
		CardID:  22193,
		Deleted: true,
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("PaymentCards.Delete returned %+v, want %+v", r, want)
	}
}

func TestPaymentCardService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		"/payment_cards",
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"user_id": 420,
				"payment_cards": [
					{
						"card_id": 22192,
						"last4": "5678",
						"brand": "MasterCard",
						"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
						"created_at_timestamp": 1546257600
					}
				]
			}`)
		})

	r, err := client.PaymentCards().List()
	if err != nil {
		t.Errorf("PaymentCards.List returned error: %v", err)
	}

	want := PaymentCardsResponse{
		Paged: Paged{
			TotalCount: -1,
			PageCount:  -1,
			Page:       -1,
			Limit:      -1,
		},
		WithUserID: WithUserID{
			UserID: 420,
		},
		Cards: []PaymentCard{
			{
				CardId: 22192,
				Last4:  "5678",
				Brand:  "MasterCard",
				WithCreationTime: WithCreationTime{
					CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
					CreatedAtTs: 1546257600,
				},
			},
		},
	}

	if !reflect.DeepEqual(r, want) {
		t.Errorf("PaymentCards.List returned %+v, want %+v", r, want)
	}
}

func TestPaymentCardService_Retrieve(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc(
		fmt.Sprintf("/payment_cards/%d", 22193),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			testMethod(t, r, "GET")
			testHeader(t, r, apiTokenHeader, testApiToken)

			_, _ = fmt.Fprint(w, `{
				"user_id": 420,
				"payment_card": {
					"card_id": 22193,
					"last4": "1234",
					"brand": "Visa",
					"created_at": "2018-12-31 12:00:00 (Etc/UTC)",
					"created_at_timestamp": 1546257600
				}
			}`)
		})

	r, err := client.PaymentCards().Retrieve(22193)
	if err != nil {
		t.Errorf("PaymentCards.Retrieve returned error: %v", err)
	}

	want := PaymentCard{
		WithCreationTime: WithCreationTime{
			CreatedAt:   "2018-12-31 12:00:00 (Etc/UTC)",
			CreatedAtTs: 1546257600,
		},
		CardId: 22193,
		Last4:  "1234",
		Brand:  "Visa",
	}

	if !reflect.DeepEqual(r.Card, want) {
		t.Errorf("PaymentCards.Retrieve returned %+v, want %+v", r.Card, want)
	}
}
