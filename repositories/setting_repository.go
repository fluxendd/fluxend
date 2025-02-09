package repositories

import (
	"fluxton/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
)

type SettingRepository struct {
	db *sqlx.DB
}

func NewSettingRepository(injector *do.Injector) (*SettingRepository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)

	return &SettingRepository{db: db}, nil
}

func (r *SettingRepository) List() ([]models.Setting, error) {
	query := "SELECT * FROM fluxton.settings;"

	var settings []models.Setting
	err := r.db.Select(&settings, query)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve rows: %v", err)
	}

	return settings, nil
}

func (r *SettingRepository) Update(settings []models.Setting) (bool, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return false, fmt.Errorf("could not begin transaction: %v", err)
	}
	defer tx.Rollback()

	for _, setting := range settings {
		query := "UPDATE fluxton.settings SET value = $1, default_value = $2, updated_at = NOW() WHERE name = $3;"
		_, err := tx.Exec(query, setting.Value, setting.DefaultValue, setting.Name)
		if err != nil {
			return false, fmt.Errorf("could not update setting: %v", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, fmt.Errorf("could not commit transaction: %v", err)
	}

	return true, nil
}
