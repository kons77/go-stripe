package cards

import (
	// "github.com/stripe/stripe-go/v72"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentintent"
	"github.com/stripe/stripe-go/v81/paymentmethod"
)

// card type holds necessary information to talk to Stripe
type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

// Charge - meaningful name for CreatePaymentIntent
func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

// CreatePaymentIntent gets payment intent
func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret

	// create a payment intent
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	// params.AddMetadata("key", "value") - just for reference

	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

// GetPaymentMethod gets the payment method by payment intend id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret

	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

// RetrievePaymentIntent gets an existing payment intent by id
func (c *Card) RetrievePaymentIntent(id string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}

	return pi, nil
}

/*  JUST USE pi.LatestCharge.ID
// RetrieveCharge an existing charge by paiment intent id
func (c *Card) RetrieveCharge(id string) (*stripe.Charge, error) {
	Stripe has deprecated direct access to Charges in PaymentIntents as part of their newer API versions.
	Instead, you should use payment_intent.
	LatestCharge to get the associated charge ID and then retrieve the charge separately if needed.
	stripe.Key = c.Secret

	pi, err := paymentintent.Get(id, nil)
	if err != nil {
		return nil, err
	}

	// Retrieve the charge using LatestCharge
	if pi.LatestCharge == nil {
		return nil, fmt.Errorf("no charge found for this payment intent")
	}

	// Retrieve the charge details
	chargeID := pi.LatestCharge.ID

	ch, err := charge.Get(chargeID, nil)
	if err != nil {
		return nil, err
	}

	return ch, nil
}
*/

func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card was expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CSV code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "Amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "Amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	default:
		msg = "Your card was declined"
	}
	return msg
}
