package part

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func (s *ServiceSuite) TestGetSuccess() {
	partUUID := uuid.New()
	expected := entity.RestorePart(partUUID.String(), "Корпус", "", valueobject.PartTypeHull, 0, 0, 0, valueobject.PartProperties{}, time.Time{})

	s.repo.EXPECT().Get(s.ctx, partUUID).Return(expected, nil)

	got, err := s.svc.Get(s.ctx, partUUID.String())
	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *ServiceSuite) TestGetInvalidUUID() {
	_, err := s.svc.Get(s.ctx, "not-a-uuid")
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestGetRepoError() {
	partUUID := uuid.New()
	repoErr := errors.New("деталь не найдена")

	s.repo.EXPECT().Get(s.ctx, partUUID).Return(entity.Part{}, repoErr)

	_, err := s.svc.Get(s.ctx, partUUID.String())
	s.Require().ErrorIs(err, repoErr)
}
