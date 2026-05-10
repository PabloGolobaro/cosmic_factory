package part

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func restoreAvailable(id uuid.UUID) entity.Part {
	return entity.RestorePart(id.String(), "Деталь", "", valueobject.PartTypeHull, 0, 10, 0, valueobject.PartProperties{}, time.Time{})
}

func restoreExhausted(id uuid.UUID) entity.Part {
	return entity.RestorePart(id.String(), "Деталь", "", valueobject.PartTypeHull, 0, 0, 0, valueobject.PartProperties{}, time.Time{})
}

func txPassThrough(s *ServiceSuite) {
	s.txManager.EXPECT().Do(s.ctx, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
}

func (s *ServiceSuite) TestReservePartsSuccess() {
	id := uuid.New()
	part := restoreAvailable(id)
	filter := model.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)
	s.repo.EXPECT().UpdateReservedBatch(s.ctx, mock.MatchedBy(func(parts []entity.Part) bool {
		return len(parts) == 1 && parts[0].Reserved() == 1
	})).Return(nil)

	err := s.svc.ReserveParts(s.ctx, []string{id.String()})
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestReservePartsInvalidUUID() {
	err := s.svc.ReserveParts(s.ctx, []string{"bad-uuid"})
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestReservePartsNotFound() {
	id := uuid.New()
	repoErr := errors.New("деталь не найдена")
	filter := model.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return(nil, repoErr)

	err := s.svc.ReserveParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestReservePartsOutOfStock() {
	id := uuid.New()
	part := restoreExhausted(id)
	filter := model.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)

	err := s.svc.ReserveParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, errs.ErrOutOfStock)
}

func (s *ServiceSuite) TestReservePartsUpdateError() {
	id := uuid.New()
	part := restoreAvailable(id)
	dbErr := errors.New("db error")
	filter := model.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)
	s.repo.EXPECT().UpdateReservedBatch(s.ctx, mock.Anything).Return(dbErr)

	err := s.svc.ReserveParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, dbErr)
}
