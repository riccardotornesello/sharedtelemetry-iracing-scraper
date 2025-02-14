package logic

import (
	"gorm.io/gorm"
	"riccardotornesello.it/sharedtelemetry/iracing/cars_models"
)

func GetCarBrands(db *gorm.DB) (map[string]cars_models.Brand, error) {
	var brands []*cars_models.Brand

	err := db.Find(&brands).Error
	if err != nil {
		return nil, err
	}

	brandsMap := make(map[string]cars_models.Brand)
	for _, brand := range brands {
		brandsMap[brand.Name] = *brand
	}

	return brandsMap, nil
}

func GetCarModelsById(db *gorm.DB, ids []int) (map[int]cars_models.Car, error) {
	var models []*cars_models.Car

	err := db.Where("id IN ?", ids).Find(&models).Error
	if err != nil {
		return nil, err
	}

	modelsMap := make(map[int]cars_models.Car)
	for _, model := range models {
		modelsMap[*model.ID] = *model
	}

	return modelsMap, nil
}
