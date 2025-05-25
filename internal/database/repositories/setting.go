package repositories

import (
	"fluxend/internal/domain/setting"
	"fluxend/pkg"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/samber/do"
	"strings"
)

type SettingRepository struct {
	db *sqlx.DB
}

func NewSettingRepository(injector *do.Injector) (setting.Repository, error) {
	db := do.MustInvoke[*sqlx.DB](injector)
	return &SettingRepository{db: db}, nil
}

func (r *SettingRepository) List() ([]setting.Setting, error) {
	query := "SELECT * FROM fluxend.settings;"

	var settings []setting.Setting
	err := r.db.Select(&settings, query)
	if err != nil {
		return nil, pkg.FormatError(err, "select", pkg.GetMethodName())
	}

	return settings, nil
}

func (r *SettingRepository) Get(name string) (setting.Setting, error) {
	query := "SELECT * FROM fluxend.settings WHERE name = $1;"
	var settingItem setting.Setting
	err := r.db.Get(&settingItem, query, name)
	if err != nil {
		return setting.Setting{}, pkg.FormatError(err, "select", pkg.GetMethodName())
	}

	return settingItem, nil
}

func (r *SettingRepository) CreateMany(settings []setting.Setting) (bool, error) {
	var valuePlaceholders []string
	var args []interface{}

	for i, currentSetting := range settings {
		placeholder := fmt.Sprintf("($%d, $%d, $%d)", 3*i+1, 3*i+2, 3*i+3)
		valuePlaceholders = append(valuePlaceholders, placeholder)
		args = append(args, currentSetting.Name, currentSetting.Value, currentSetting.DefaultValue)
	}

	query := fmt.Sprintf("INSERT INTO fluxend.settings (name, value, default_value) VALUES %s;",
		strings.Join(valuePlaceholders, ", "))

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return false, fmt.Errorf("could not create settings: %v", err)
	}

	return true, nil
}

func (r *SettingRepository) Update(settings []setting.Setting) (bool, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return false, pkg.FormatError(err, "transactionBegin", pkg.GetMethodName())
	}
	defer tx.Rollback()

	for _, currentSetting := range settings {
		query := "UPDATE fluxend.settings SET value = $1, default_value = $2, updated_at = NOW() WHERE name = $3;"
		_, err := tx.Exec(query, currentSetting.Value, currentSetting.DefaultValue, currentSetting.Name)
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
