package part

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	partUUID := uuid.New()
	expected := model.Part{
		UUID:     partUUID,
		Name:     "Корпус",
		PartType: model.PartTypeHull,
	}

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

	s.repo.EXPECT().Get(s.ctx, partUUID).Return(model.Part{}, repoErr)

	_, err := s.svc.Get(s.ctx, partUUID.String())
	s.Require().ErrorIs(err, repoErr)
}
