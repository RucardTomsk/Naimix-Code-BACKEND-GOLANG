package migration

import (
	"errors"
	"fmt"
	"github.com/RucardTomsk/Naimix-Code-BACKEND-GOLANG/internal/domain/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
)

func Migrate(
	db *gorm.DB,
	adminID uuid.UUID,
	adminUserName string,
	adminEmail string,
	adminPassword string,
) error {
	if err := db.AutoMigrate(
		&entity.Session{},
		&entity.User{},
		&entity.Company{},
		&entity.Vacancy{},
		&entity.Candidate{},
	); err != nil {
		//relationship doesn't exist
		if !strings.Contains(err.Error(), "42P07") {
			fmt.Println(err)
			return err
		}
	}

	if err := adminMigration(db, adminID, adminUserName, adminEmail, adminPassword); err != nil {
		return err
	}

	return nil
}

func adminMigration(db *gorm.DB,
	adminID uuid.UUID,
	adminUserName string,
	adminEmail string,
	adminPassword string) error {

	tx := db.First(&entity.User{}, adminID)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return tx.Error
		}
	}

	if tx.RowsAffected == 0 {
		user := &entity.User{
			Email:    adminEmail,
			Name:     adminUserName,
			Password: adminPassword,
		}
		user.ID = adminID

		if err := db.Create(user).Error; err != nil {
			return err
		}
	}

	return nil
}
