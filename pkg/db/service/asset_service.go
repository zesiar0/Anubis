package service

import (
	"Anubis/pkg/db"
	"Anubis/pkg/logger"
	"Anubis/pkg/model"
)

type AssetService struct {
}

func (a *AssetService) GetAssetsByUser(username string) []model.Asset {
	var assets []model.Asset

	if err := db.DB.Where("username = ?", username).Find(&assets); err != nil {
		logger.Errorf("Get user's assets err: %v", err)
	}

	return assets
}
