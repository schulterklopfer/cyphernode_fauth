/*
 * MIT License
 *
 * Copyright (c) 2021 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package dataSource

import (
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  //_ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/schulterklopfer/cyphernode_fauth/cnaErrors"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/models"
)

var db *gorm.DB

func GetDB() *gorm.DB {
  return db
}

func Init( dsn string ) error {
  if db != nil {
    return nil
  }
  var err error
  logwrapper.Logger().Info( "Opening database "+dsn)

  db, err = gorm.Open(postgres.New(postgres.Config{
    DSN: dsn,
    PreferSimpleProtocol: true, // disables implicit prepared statement usage
  }), &gorm.Config{})

  //db, err = gorm.Open("sqlite3", dsn )
  if err != nil {
    logwrapper.Logger().Panic("failed to connect to database "+err.Error() )
    return err
  }
  err = AutoMigrate()
  if err != nil {
    return err
  }
  return nil
}

func Close() {
  if db == nil {
    return
  }

  sqlDB, err := db.DB()
  if err != nil {
    logwrapper.Logger().Panic("failed to close " + err.Error())
  }
  defer sqlDB.Close()
  db = nil
}

func AutoMigrate() error {
  if db == nil {
    return cnaErrors.ErrDatabaseNotInitialised
  }
  logwrapper.Logger().Info( "Migrating database")
  return db.AutoMigrate(
    &models.UserModel{},
    &models.AppModel{},
    &models.RoleModel{} )
}