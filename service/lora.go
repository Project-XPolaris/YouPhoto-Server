package service

import "github.com/projectxpolaris/youphoto/database"

func SaveLoraConfig(name string, config string, userId uint) (*database.LoraConfig, error) {
	saveConfig := &database.LoraConfig{
		Name:   name,
		Config: config,
		UserId: userId,
	}
	var existConfig database.LoraConfig
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

func DeleteLoraConfig(id uint, userId uint) error {
	var loraConfig database.LoraConfig
	err := database.Instance.Where("user_id = ?", userId).First(&loraConfig, id).Error
	if err != nil {
		return err
	}
	return database.Instance.Delete(&database.LoraConfig{}, id).Error
}

func GetLoraConfigList(userId uint) ([]*database.LoraConfig, error) {
	var result []*database.LoraConfig
	err := database.Instance.Where("user_id = ?", userId).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
