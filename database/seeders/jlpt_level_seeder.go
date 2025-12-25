package seeders

import (
	"manabu-service/domain/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunJlptLevelSeeder(db *gorm.DB) {
	jlptLevels := []models.JlptLevel{
		{
			Code:        "N5",
			Name:        "JLPT N5",
			Description: "Basic level",
			LevelOrder:  5,
		},
		{
			Code:        "N4",
			Name:        "JLPT N4",
			Description: "Elementary level",
			LevelOrder:  4,
		},
		{
			Code:        "N3",
			Name:        "JLPT N3",
			Description: "Intermediate level",
			LevelOrder:  3,
		},
		{
			Code:        "N2",
			Name:        "JLPT N2",
			Description: "Pre-advanced level",
			LevelOrder:  2,
		},
		{
			Code:        "N1",
			Name:        "JLPT N1",
			Description: "Advanced level",
			LevelOrder:  1,
		},
	}

	for _, jlptLevel := range jlptLevels {
		err := db.FirstOrCreate(&jlptLevel, models.JlptLevel{Code: jlptLevel.Code}).Error
		if err != nil {
			logrus.Errorf("failed to seed jlpt level: %v", err)
			panic(err)
		}
		logrus.Infof("jlpt level %s successfully seeded", jlptLevel.Code)
	}
}
