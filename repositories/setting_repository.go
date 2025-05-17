package repositories

import (
	"fluxton/models"
	"fluxton/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"strings"
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
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}

	return settings, nil
}

func (r *SettingRepository) CreateMany(settings []models.Setting) (bool, error) {
	var valuePlaceholders []string
	var args []interface{}

	for i, setting := range settings {
		placeholder := fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3)
		valuePlaceholders = append(valuePlaceholders, placeholder)
		args = append(args, setting.Name, setting.Value, setting.DefaultValue)
	}

	query := fmt.Sprintf("INSERT INTO fluxton.settings (name, value, default_value) VALUES %s;",
		strings.Join(valuePlaceholders, ", "))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("could not create settings: %v", err)
	}

	return true, nil
}

func (r *SettingRepository) Update(settings []models.Setting) (bool, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return false, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
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
		return false, pkg.FormatError(err, "transactionCommit", pkg.GetMethodName())
	}

	return true, nil
}
