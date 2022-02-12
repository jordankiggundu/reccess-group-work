package db

import (
  "gorm.io/gorm"
  "gorm.io/driver/mysql"
  "fmt"
)


func Connect() (*gorm.DB, error) {
	dbUser := "root"
	dbPswd := ""
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "kinderCareRecess2022"
	
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPswd, dbHost, dbPort, dbName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
      if err != nil {
        return nil, err
      }
    //======= Migrate the schema ======
    db.AutoMigrate(&Profile{})
    db.AutoMigrate(&Pupil{})
    db.AutoMigrate(&Teacher{})
    db.AutoMigrate(&Assignment{})
    db.AutoMigrate(&AssignmentList{})
    db.AutoMigrate(&Solution{})
    db.AutoMigrate(&ActivationRequest{})
    
	return db, nil
}




