package store

import (
	"database/sql"

	"github.com/A-pen-app/kickstart/models"
	"github.com/lib/pq"
)

func parseError(err error) error {
	if err == sql.ErrNoRows {
		return models.ErrorNotFound
	}
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	switch pqErr.Code.Name() {
	case "string_data_right_truncation":
		return models.ErrorWrongParams
	case "not_null_violation":
		return models.ErrorWrongParams
	case "foreign_key_violation":
		return models.ErrorWrongParams
	case "unique_violation":
		return models.ErrorDuplicateEntry
	case "integrity_constraint_violation":
		return models.ErrorDuplicateEntry
	}
	return err
}
