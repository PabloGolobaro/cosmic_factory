package converter

import (
	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/order/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/order/internal/repository/record"
)

var paymentMethodToStr = map[model.PaymentMethod]string{
	model.PaymentMethodCard:          "CARD",
	model.PaymentMethodSBP:           "SBP",
	model.PaymentMethodCreditCard:    "CREDIT_CARD",
	model.PaymentMethodInvestorMoney: "INVESTOR_MONEY",
}

var strToPaymentMethod = map[string]model.PaymentMethod{
	"CARD":           model.PaymentMethodCard,
	"SBP":            model.PaymentMethodSBP,
	"CREDIT_CARD":    model.PaymentMethodCreditCard,
	"INVESTOR_MONEY": model.PaymentMethodInvestorMoney,
}

func uuidPtrToStr(u *uuid.UUID) *string {
	if u == nil {
		return nil
	}
	s := u.String()
	return &s
}

func strPtrToUUID(s *string) *uuid.UUID {
	if s == nil {
		return nil
	}
	u := uuid.MustParse(*s)
	return &u
}

func OrderToRecord(o model.Order) record.OrderRecord {
	r := record.OrderRecord{
		OrderUUID:       o.OrderUUID.String(),
		HullUUID:        o.HullUUID.String(),
		EngineUUID:      o.EngineUUID.String(),
		ShieldUUID:      uuidPtrToStr(o.ShieldUUID),
		WeaponUUID:      uuidPtrToStr(o.WeaponUUID),
		TotalPrice:      o.TotalPrice,
		TransactionUUID: uuidPtrToStr(o.TransactionUUID),
		Status:          o.Status,
		CreatedAt:       o.CreatedAt,
	}
	if s, ok := paymentMethodToStr[o.PaymentMethod]; ok {
		r.PaymentMethod = &s
	}
	return r
}

func OrderFromRecord(r record.OrderRecord) model.Order {
	o := model.Order{
		OrderUUID:       uuid.MustParse(r.OrderUUID),
		HullUUID:        uuid.MustParse(r.HullUUID),
		EngineUUID:      uuid.MustParse(r.EngineUUID),
		ShieldUUID:      strPtrToUUID(r.ShieldUUID),
		WeaponUUID:      strPtrToUUID(r.WeaponUUID),
		TotalPrice:      r.TotalPrice,
		TransactionUUID: strPtrToUUID(r.TransactionUUID),
		Status:          r.Status,
		CreatedAt:       r.CreatedAt,
	}
	if r.PaymentMethod != nil {
		o.PaymentMethod = strToPaymentMethod[*r.PaymentMethod]
	}
	return o
}
