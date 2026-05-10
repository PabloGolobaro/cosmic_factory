package part

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func restoreReserved(id uuid.UUID) entity.Part {
	return entity.RestorePart(id.String(), "Деталь", "", valueobject.PartTypeHull, 0, 10, 1, valueobject.PartProperties{}, time.Time{})
}

func (s *ServiceSuite) TestReleasePartsSuccess() {
	id := uuid.New()
	part := restoreReserved(id)
	filter := valueobject.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)
	s.repo.EXPECT().UpdateReservedBatch(s.ctx, mock.MatchedBy(func(parts []entity.Part) bool {
		return len(parts) == 1 && parts[0].Reserved() == 0
	})).Return(nil)

	err := s.svc.ReleaseParts(s.ctx, []string{id.String()})
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestReleasePartsInvalidUUID() {
	err := s.svc.ReleaseParts(s.ctx, []string{"bad-uuid"})
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestReleasePartsNotFound() {
	id := uuid.New()
	repoErr := errors.New("деталь не найдена")
	filter := valueobject.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return(nil, repoErr)

	err := s.svc.ReleaseParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestReleasePartsNothingToRelease() {
	id := uuid.New()
	part := restoreAvailable(id)
	filter := valueobject.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)

	err := s.svc.ReleaseParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, errs.ErrNothingToRelease)
}

func (s *ServiceSuite) TestReleasePartsUpdateError() {
	id := uuid.New()
	part := restoreReserved(id)
	dbErr := errors.New("db error")
	filter := valueobject.PartFilter{UUIDs: []string{id.String()}}

	txPassThrough(s)
	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{part}, nil)
	s.repo.EXPECT().UpdateReservedBatch(s.ctx, mock.Anything).Return(dbErr)

	err := s.svc.ReleaseParts(s.ctx, []string{id.String()})
	s.Require().ErrorIs(err, dbErr)
}
