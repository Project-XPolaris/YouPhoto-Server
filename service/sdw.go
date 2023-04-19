package service

import "github.com/projectxpolaris/youphoto/database"

func SaveSDWConfig(name string, config string, userId uint) (*database.SdwConfig, error) {
	saveConfig := &database.SdwConfig{
		Name:   name,
		Config: config,
		UserId: userId,
	}
	var existConfig database.SdwConfig
	err := database.Instance.Where("user_id = ? and name = ?", userId, name).First(&existConfig).Error
	if err == nil {
		saveConfig.ID = existConfig.ID
	}
	err = database.Instance.Save(saveConfig).Error
	if err != nil {
		return nil, err
	}
	return saveConfig, nil
}

func DeleteSDWConfig(id uint, userId uint) error {
	var sdwConfig database.SdwConfig
	err := database.Instance.Where("user_id = ?", userId).First(&sdwConfig, id).Error
	if err != nil {
		return err
	}
	return database.Instance.Delete(&database.SdwConfig{}, id).Error
}

func GetSDWConfigList(userId uint) ([]*database.SdwConfig, error) {
	var result []*database.SdwConfig
	err := database.Instance.Where("user_id = ?", userId).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
