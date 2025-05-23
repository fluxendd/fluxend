package mapper

import (
	settingDto "fluxton/internal/api/dto/setting"
	settingDomain "fluxton/internal/domain/setting"
)

func ToSettingResource(setting *settingDomain.Setting) settingDto.Response {
	return settingDto.Response{
		ID:           setting.ID,
		Name:         setting.Name,
		Value:        setting.Value,
		DefaultValue: setting.DefaultValue,
		CreatedAt:    setting.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    setting.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ToSettingResourceCollection(settings []settingDomain.Setting) []settingDto.Response {
	settingResource := make([]settingDto.Response, len(settings))
	for i, currentSetting := range settings {
		settingResource[i] = ToSettingResource(&currentSetting)
	}

	return settingResource
}
