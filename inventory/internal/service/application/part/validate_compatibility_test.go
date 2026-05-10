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

func typedPart(id uuid.UUID, pt valueobject.PartType) entity.Part {
	return entity.RestorePart(id.String(), "test", "", pt, 0, 1, 0, valueobject.PartProperties{}, time.Time{})
}

func defaultSlots(hullID, engineID uuid.UUID) valueobject.ShipSlots {
	return valueobject.ShipSlots{HullUUID: hullID.String(), EngineUUID: engineID.String()}
}

func (s *ServiceSuite) TestValidateCompatibilitySuccess() {
	hullID, engineID := uuid.New(), uuid.New()
	hull := typedPart(hullID, valueobject.PartTypeHull)
	engine := typedPart(engineID, valueobject.PartTypeEngine)
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String()}}

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{hull, engine}, nil)
	s.checker.EXPECT().Check(mock.Anything).Return(nil)

	err := s.svc.ValidateCompatibility(s.ctx, defaultSlots(hullID, engineID))
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestValidateCompatibilitySuccessAllFour() {
	hullID, engineID, shieldID, weaponID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	hull := typedPart(hullID, valueobject.PartTypeHull)
	engine := typedPart(engineID, valueobject.PartTypeEngine)
	shield := typedPart(shieldID, valueobject.PartTypeShield)
	weapon := typedPart(weaponID, valueobject.PartTypeWeapon)
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String(), shieldID.String(), weaponID.String()}}

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{hull, engine, shield, weapon}, nil)
	s.checker.EXPECT().Check(mock.Anything).Return(nil)

	slots := valueobject.ShipSlots{
		HullUUID:   hullID.String(),
		EngineUUID: engineID.String(),
		ShieldUUID: shieldID.String(),
		WeaponUUID: weaponID.String(),
	}
	err := s.svc.ValidateCompatibility(s.ctx, slots)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestValidateCompatibilityInvalidUUID() {
	slots := valueobject.ShipSlots{HullUUID: "bad-uuid", EngineUUID: uuid.New().String()}

	err := s.svc.ValidateCompatibility(s.ctx, slots)
	s.Require().ErrorIs(err, errs.ErrInvalidUUID)
}

func (s *ServiceSuite) TestValidateCompatibilityGetBatchError() {
	hullID, engineID := uuid.New(), uuid.New()
	repoErr := errors.New("db error")
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String()}}

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return(nil, repoErr)

	err := s.svc.ValidateCompatibility(s.ctx, defaultSlots(hullID, engineID))
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestValidateCompatibilityPartNotFound() {
	hullID, engineID := uuid.New(), uuid.New()
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String()}}

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{}, nil)

	err := s.svc.ValidateCompatibility(s.ctx, defaultSlots(hullID, engineID))
	s.Require().ErrorIs(err, errs.ErrPartNotFound)
}

func (s *ServiceSuite) TestValidateCompatibilityTypeMismatch() {
	hullID, engineID := uuid.New(), uuid.New()
	// Hull slot UUID maps to an engine-typed part.
	wrongPart := typedPart(hullID, valueobject.PartTypeEngine)
	engine := typedPart(engineID, valueobject.PartTypeEngine)
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String()}}

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{wrongPart, engine}, nil)

	err := s.svc.ValidateCompatibility(s.ctx, defaultSlots(hullID, engineID))
	s.Require().ErrorIs(err, errs.ErrPartTypeMismatch)
}

func (s *ServiceSuite) TestValidateCompatibilityIncompatibleParts() {
	hullID, engineID := uuid.New(), uuid.New()
	hull := typedPart(hullID, valueobject.PartTypeHull)
	engine := typedPart(engineID, valueobject.PartTypeEngine)
	filter := valueobject.PartFilter{UUIDs: []string{hullID.String(), engineID.String()}}
	checkerErr := errors.New("несовместимы")

	s.repo.EXPECT().GetBatch(s.ctx, filter).Return([]entity.Part{hull, engine}, nil)
	s.checker.EXPECT().Check(mock.Anything).Return(checkerErr)

	err := s.svc.ValidateCompatibility(s.ctx, defaultSlots(hullID, engineID))
	s.Require().ErrorIs(err, checkerErr)
}
