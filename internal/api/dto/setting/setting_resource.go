package setting

import (
	"fluxton/internal/domain/setting"
)

type SettingResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Value        string `json:"value"`
	DefaultValue string `json:"defaultValue"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

func SettingResource(setting *setting.Setting) SettingResponse {
	return SettingResponse{
		ID:           setting.ID,
		Name:         setting.Name,
		Value:        setting.Value,
		DefaultValue: setting.DefaultValue,
		CreatedAt:    setting.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    setting.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func SettingResourceCollection(settings []setting.Setting) []SettingResponse {
	resourcesettings := make([]SettingResponse, len(settings))
	for i, setting := range settings {
		resourcesettings[i] = SettingResource(&setting)
	}

	return resourcesettings
}
