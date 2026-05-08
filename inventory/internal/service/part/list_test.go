package part

import (
	"errors"
	"time"

	"github.com/google/uuid"

	errs "github.com/PabloGolobaro/cosmic_factory/inventory/internal/errors"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/entity"
	"github.com/PabloGolobaro/cosmic_factory/inventory/internal/model/valueobject"
)

func restorePart(id uuid.UUID, name string, partType valueobject.PartType) entity.Part {
	return entity.RestorePart(id.String(), name, "", partType, 0, 0, 0, valueobject.PartProperties{}, time.Time{})
}

// --- Путь А: ids непустые → GetBatch.

func (s *ServiceSuite) TestListByIDsSuccess() {
	id1, id2 := uuid.New(), uuid.New()
	expected := []entity.Part{
		restorePart(id1, "Корпус", valueobject.PartTypeHull),
		restorePart(id2, "Двигатель", valueobject.PartTypeEngine),
	}

	s.repo.EXPECT().GetBatch(s.ctx, []uuid.UUID{id1, id2}).Return(expected, nil)

	got, err := s.svc.List(s.ctx, []string{id1.String(), id2.String()}, valueobject.PartTypeUnspecified)
	s.Require().NoError(err)
	s.Equal(expected, got)
}

func (s *ServiceSuite) TestListByIDsInvalidUUID() {
	_, err := s.svc.List(s.ctx, []string{uuid.New().String(), "bad-uuid"}, valueobject.PartTypeUnspecified)
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestListByIDsRepoError() {
	id := uuid.New()
	repoErr := errors.New("db error")

	s.repo.EXPECT().GetBatch(s.ctx, []uuid.UUID{id}).Return(nil, repoErr)

	_, err := s.svc.List(s.ctx, []string{id.String()}, valueobject.PartTypeUnspecified)
	s.Require().ErrorIs(err, repoErr)
}

// --- Путь Б: ids пустые → GetAll + filter + sort.

func (s *ServiceSuite) TestListAllSorted() {
	parts := []entity.Part{
		restorePart(uuid.New(), "Щит", valueobject.PartTypeShield),
		restorePart(uuid.New(), "Корпус", valueobject.PartTypeHull),
		restorePart(uuid.New(), "Двигатель", valueobject.PartTypeEngine),
	}

	s.repo.EXPECT().GetAll(s.ctx).Return(parts, nil)

	got, err := s.svc.List(s.ctx, nil, valueobject.PartTypeUnspecified)
	s.Require().NoError(err)
	s.Require().Len(got, 3)
	s.Equal("Двигатель", got[0].Name())
	s.Equal("Корпус", got[1].Name())
	s.Equal("Щит", got[2].Name())
}

func (s *ServiceSuite) TestListFilteredByType() {
	hullUUID := uuid.New()
	parts := []entity.Part{
		restorePart(hullUUID, "Корпус", valueobject.PartTypeHull),
		restorePart(uuid.New(), "Двигатель", valueobject.PartTypeEngine),
		restorePart(uuid.New(), "Оружие", valueobject.PartTypeWeapon),
	}

	s.repo.EXPECT().GetAll(s.ctx).Return(parts, nil)

	got, err := s.svc.List(s.ctx, nil, valueobject.PartTypeHull)
	s.Require().NoError(err)
	s.Require().Len(got, 1)
	s.Equal(hullUUID.String(), got[0].UUID())
}

func (s *ServiceSuite) TestListGetAllRepoError() {
	repoErr := errors.New("db error")

	s.repo.EXPECT().GetAll(s.ctx).Return(nil, repoErr)

	_, err := s.svc.List(s.ctx, nil, valueobject.PartTypeUnspecified)
	s.Require().ErrorIs(err, repoErr)
}
