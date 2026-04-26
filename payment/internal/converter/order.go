package converter

import (
	"github.com/PabloGolobaro/cosmic_factory/payment/internal/model"
)

var paymentMethodToStr = map[model.PaymentMethod]string{
	model.PaymentMethodCard:          "CARD",
	model.PaymentMethodSBP:           "SBP",
	model.PaymentMethodCreditCard:    "CREDIT_CARD",
	model.PaymentMethodInvestorMoney: "INVESTOR_MONEY",
}

func PaymentMethodToString(m model.PaymentMethod) string {
	return paymentMethodToStr[m]
}
