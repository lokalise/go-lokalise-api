package lokalise

import (
	"fmt"
)

const (
	pathPaymentCards = "payment_cards"
)

type PaymentCardService struct {
	BaseService
}

type PaymentCard struct {
	WithCreationTime
	CardId int64  `json:"card_id"`
	Last4  string `json:"last4"`
	Brand  string `json:"brand"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service request/response objects
// _____________________________________________________________________________________________________________________

type CreatePaymentCard struct {
	Number   string `json:"number"`
	CVC      string `json:"cvc"`
	ExpMonth int64  `json:"exp_month"`
	ExpYear  int64  `json:"exp_year"`
}

type PaymentCardsResponse struct {
	Paged
	WithUserID
	Cards []PaymentCard `json:"payment_cards"`
}

type PaymentCardResponse struct {
	WithUserID
	Card PaymentCard `json:"payment_card"`
}

type DeletePaymentCardResponse struct {
	CardID  int64 `json:"card_id"`
	Deleted bool  `json:"card_deleted"`
}

// ‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾‾
// Service methods
// _____________________________________________________________________________________________________________________

func (c *PaymentCardService) Create(card CreatePaymentCard) (r PaymentCard, err error) {
	resp, err := c.post(c.Ctx(), pathPaymentCards, &r, card)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *PaymentCardService) List() (r PaymentCardsResponse, err error) {
	resp, err := c.getList(c.Ctx(), pathPaymentCards, &r, c.PageOpts())

	if err != nil {
		return
	}
	applyPaged(resp, &r.Paged)
	return r, apiError(resp)
}

func (c *PaymentCardService) Retrieve(cardID int64) (r PaymentCardResponse, err error) {
	resp, err := c.get(c.Ctx(), fmt.Sprintf("%s/%d", pathPaymentCards, cardID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}

func (c *PaymentCardService) Delete(cardID int64) (r DeletePaymentCardResponse, err error) {
	resp, err := c.delete(c.Ctx(), fmt.Sprintf("%s/%d", pathPaymentCards, cardID), &r)

	if err != nil {
		return
	}
	return r, apiError(resp)
}
