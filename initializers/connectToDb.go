package initializers

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func ConnectToDb() {
	var err error
	//postgres://xygvbyyy:L4PN9-6XDAKiggJvzSNSp2uLMDQJrU-O@kala.db.elephantsql.com/xygvbyyy
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect DB")
	}
}
