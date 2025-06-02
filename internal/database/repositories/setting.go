package repositories

import (
	"fluxend/internal/domain/setting"
	"fluxend/internal/domain/shared"
	"fmt"
	"github.com/samber/do"
	"strings"
)

type SettingRepository struct {
	db shared.DB
}

func NewSettingRepository(injector *do.Injector) (setting.Repository, error) {
	db := do.MustInvoke[shared.DB](injector)
	return &SettingRepository{db: db}, nil
}

func (r *SettingRepository) List() ([]setting.Setting, error) {
	query := "SELECT * FROM fluxend.settings;"

	var settings []setting.Setting
	return settings, r.db.SelectList(&settings, query)
}

func (r *SettingRepository) Get(name string) (setting.Setting, error) {
	query := "SELECT * FROM fluxend.settings WHERE name = $1;"

	var settingItem setting.Setting
	return settingItem, r.db.Get(&settingItem, query, name)
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

	_, err := r.db.ExecWithRowsAffected(query, args...)
	return err == nil, err
}

func (r *SettingRepository) Update(settings []setting.Setting) (bool, error) {
	err := r.db.WithTransaction(func(tx shared.Tx) error {
		for _, currentSetting := range settings {
			query := "UPDATE fluxend.settings SET value = $1, default_value = $2, updated_at = NOW() WHERE name = $3;"
			if _, err := tx.Exec(query, currentSetting.Value, currentSetting.DefaultValue, currentSetting.Name); err != nil {
				return fmt.Errorf("could not update setting: %v", err)
			}
		}
		return nil
	})

	return err == nil, err
}
