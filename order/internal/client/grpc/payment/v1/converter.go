package v1

import (
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	paymentv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/proto/payment/v1"
)

var modelToProtoPaymentMethod = map[model.PaymentMethod]paymentv1.PaymentMethod{
	model.PaymentMethodCard:          paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	model.PaymentMethodSBP:           paymentv1.PaymentMethod_PAYMENT_METHOD_SBP,
	model.PaymentMethodCreditCard:    paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD,
	model.PaymentMethodInvestorMoney: paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY,
}

func paymentMethodToProto(m model.PaymentMethod) paymentv1.PaymentMethod {
	return modelToProtoPaymentMethod[m]
}
