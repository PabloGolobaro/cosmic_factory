package part

import (
	"errors"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model"
)

// --- Путь А: ids непустые → GetBatch.

func (s *ServiceSuite) TestListByIDsSuccess() {
	id1, id2 := uuid.New(), uuid.New()
	expected := []model.Part{
		{UUID: id1, Name: "Корпус", PartType: model.PartTypeHull},
		{UUID: id2, Name: "Двигатель", PartType: model.PartTypeEngine},
	}

	s.repo.EXPECT().GetBatch(s.ctx, []uuid.UUID{id1, id2}).Return(expected, nil)

	got, err := s.svc.List(s.ctx, []string{id1.String(), id2.String()}, model.PartTypeUnspecified)
	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *ServiceSuite) TestListByIDsInvalidUUID() {
	_, err := s.svc.List(s.ctx, []string{uuid.New().String(), "bad-uuid"}, model.PartTypeUnspecified)
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestListByIDsRepoError() {
	id := uuid.New()
	repoErr := errors.New("db error")

	s.repo.EXPECT().GetBatch(s.ctx, []uuid.UUID{id}).Return(nil, repoErr)

	_, err := s.svc.List(s.ctx, []string{id.String()}, model.PartTypeUnspecified)
	s.Require().ErrorIs(err, repoErr)
}

// --- Путь Б: ids пустые → GetAll + filter + sort.

func (s *ServiceSuite) TestListAllSorted() {
	parts := []model.Part{
		{UUID: uuid.New(), Name: "Щит", PartType: model.PartTypeShield},
		{UUID: uuid.New(), Name: "Корпус", PartType: model.PartTypeHull},
		{UUID: uuid.New(), Name: "Двигатель", PartType: model.PartTypeEngine},
	}

	s.repo.EXPECT().GetAll(s.ctx).Return(parts, nil)

	got, err := s.svc.List(s.ctx, nil, model.PartTypeUnspecified)
	s.Require().NoError(err)
	s.Require().Len(got, 3)
	s.Equal("Двигатель", got[0].Name)
	s.Equal("Корпус", got[1].Name)
	s.Equal("Щит", got[2].Name)
}

func (s *ServiceSuite) TestListFilteredByType() {
	hullUUID := uuid.New()
	parts := []model.Part{
		{UUID: hullUUID, Name: "Корпус", PartType: model.PartTypeHull},
		{UUID: uuid.New(), Name: "Двигатель", PartType: model.PartTypeEngine},
		{UUID: uuid.New(), Name: "Оружие", PartType: model.PartTypeWeapon},
	}

	s.repo.EXPECT().GetAll(s.ctx).Return(parts, nil)

	got, err := s.svc.List(s.ctx, nil, model.PartTypeHull)
	s.Require().NoError(err)
	s.Require().Len(got, 1)
	s.Equal(hullUUID, got[0].UUID)
}

func (s *ServiceSuite) TestListGetAllRepoError() {
	repoErr := errors.New("db error")

	s.repo.EXPECT().GetAll(s.ctx).Return(nil, repoErr)

	_, err := s.svc.List(s.ctx, nil, model.PartTypeUnspecified)
	s.Require().ErrorIs(err, repoErr)
}
