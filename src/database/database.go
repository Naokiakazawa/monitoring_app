package database

import (	
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type Record struct {
    gorm.Model
    GroupID	string
    Tx_Hash	string
    From	string
    To		string
}

func gormConnect() *gorm.DB {
    DBMS := "mysql"
    USER := "root"
    PASSWORD := "password"
    PROTOCOL := "tcp(mysql)"
    DBNAME := "test_db"
    OPTION := "?charset=utf8&parseTime=True&loc=Local"

    CONNECTION := USER + ":" + PASSWORD + "@" + PROTOCOL + "/" + DBNAME + OPTION

    db,err := gorm.Open(DBMS, CONNECTION)
    if err != nil {
        panic("Connection Failed!")
    }
    return db
}

func CreateLogs(footprint *Footprint) {
    db := gormConnect()
    defer db.Close()

    if !db.NewRecord(footprint) {
        panic("Could not create new record")
    }
    db.Create(footprint)
}