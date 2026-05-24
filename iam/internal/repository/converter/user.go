package converter

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/PabloGolobaro/cosmic_factory/iam/internal/model"
	"github.com/PabloGolobaro/cosmic_factory/iam/internal/repository/record"
)

func UserToRecord(u model.User) record.UserRecord {
	return record.UserRecord{
		UUID:         u.UUID.String(),
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}

func UserFromRecord(r record.UserRecord) (model.User, error) {
	id, err := uuid.Parse(r.UUID)
	if err != nil {
		return model.User{}, fmt.Errorf("некорректный UUID записи: %w", err)
	}

	return model.User{
		UUID:         id,
		Login:        r.Login,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}, nil
}
