package converter

import (
	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	orderv1 "github.com/PabloGolobaro/cosmic_factory/shared/pkg/openapi/order/v1"
)

var modelToDto = map[model.PaymentMethod]orderv1.PaymentMethod{
	model.PaymentMethodCard:          orderv1.PaymentMethodCARD,
	model.PaymentMethodSBP:           orderv1.PaymentMethodSBP,
	model.PaymentMethodCreditCard:    orderv1.PaymentMethodCREDITCARD,
	model.PaymentMethodInvestorMoney: orderv1.PaymentMethodINVESTORMONEY,
}

var strToModel = map[string]model.PaymentMethod{
	"CARD":           model.PaymentMethodCard,
	"SBP":            model.PaymentMethodSBP,
	"CREDIT_CARD":    model.PaymentMethodCreditCard,
	"INVESTOR_MONEY": model.PaymentMethodInvestorMoney,
}

func modelPaymentMethodToDto(m model.PaymentMethod) (orderv1.PaymentMethod, bool) {
	v, ok := modelToDto[m]
	return v, ok
}

func PaymentMethodFromString(s string) (model.PaymentMethod, bool) {
	v, ok := strToModel[s]
	return v, ok
}

func OrderFromCreateRequest(req *orderv1.CreateOrderRequest) model.Order {
	order := model.Order{
		UserUUID:   req.UserUUID,
		HullUUID:   req.HullUUID,
		EngineUUID: req.EngineUUID,
	}
	if v, ok := req.ShieldUUID.Get(); ok {
		order.ShieldUUID = &v
	}
	if v, ok := req.WeaponUUID.Get(); ok {
		order.WeaponUUID = &v
	}
	return order
}

func OrderToDto(order model.Order) *orderv1.OrderDto {
	dto := &orderv1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		HullUUID:   order.HullUUID,
		EngineUUID: order.EngineUUID,
		TotalPrice: order.TotalPrice,
		Status:     orderv1.OrderStatus(order.Status),
		CreatedAt:  order.CreatedAt,
	}
	if order.ShieldUUID != nil {
		dto.ShieldUUID = orderv1.NewOptNilUUID(*order.ShieldUUID)
	}
	if order.WeaponUUID != nil {
		dto.WeaponUUID = orderv1.NewOptNilUUID(*order.WeaponUUID)
	}
	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderv1.NewOptNilUUID(*order.TransactionUUID)
	}
	if pm, ok := modelPaymentMethodToDto(order.PaymentMethod); ok {
		dto.PaymentMethod = orderv1.NewOptNilPaymentMethod(pm)
	}
	return dto
}
