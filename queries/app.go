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

package queries

import (
  "errors"
  "github.com/schulterklopfer/cyphernode_fauth/cnaErrors"
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "gopkg.in/validator.v2"
)

func CreateApp( app *models.AppModel ) error {
  if app.ID != 0 {
    // app must not have any ID possibly existing in DB
    return errors.New( "app ID must be 0" )
  }
  db := dataSource.GetDB()

  var existingApps []models.AppModel
  db.Limit(1).Find( &existingApps, models.AppModel{Hash: app.Hash} )

  if len(existingApps) > 0 {
    return errors.New( "app with same hash already exists" )
  }

  err := validator.Validate(app)
  if err != nil {
    return err
  }
  db.Create(app)
  return nil
}

func DeleteApp( id uint ) error {
  if id == 0 {
    return errors.New("no such app")
  }
  if id == 1 {
    return cnaErrors.ErrActionForbidden
  }
  db := dataSource.GetDB()
  var app models.AppModel
  db.Take( &app, id )
  if app.ID == 0 {
    return errors.New("no such app")
  }
  db.Unscoped().Delete( &app )
  return nil
}

func RemoveRoleFromApp(  app *models.AppModel, roleId uint ) error {
  if roleId == 1 && app.ID == 1 {
    return cnaErrors.ErrActionForbidden
  }
  //db := dataSource.GetDB()

  var role models.RoleModel

  err := Get( &role, roleId, false )

  if err != nil {
    return err
  }

  if role.ID == 0 || role.AppId != app.ID {
    return cnaErrors.ErrNoSuchRole
  }

  //db.Model(app).Association("AvailableRoles").Delete( role )
  return DeleteRole( role.ID )
}

func CreateRoleForApp( app *models.AppModel, role *models.RoleModel ) error {
  db := dataSource.GetDB()

  if role.ID != 0 {
    return cnaErrors.ErrCannotAddExistingRole
  }

  db.Model(app).Association("AvailableRoles").Append( role )
  return db.Error
}

func GetAppByHash( hash string ) (*models.AppModel, error) {
  var apps []*models.AppModel
  err := Find( &apps,  []interface{}{"hash = ?", hash}, "", 1,0,false)

  if err != nil {
    return nil, err
  }

  if len(apps) == 0 {
    return nil, nil
  }

  err = LoadRoles( apps[0] )
  if err != nil {
    return nil, err
  }

  return apps[0], nil
}

func GetAppByMountPoint( mountPoint string ) (*models.AppModel, error) {
  var apps []*models.AppModel
  err := Find( &apps,  []interface{}{"mount_point = ?", mountPoint}, "", 1,0,false)

  if err != nil {
    return nil, err
  }

  if len(apps) == 0 {
    return nil, cnaErrors.ErrNoSuchApp
  }

  err = LoadRoles( apps[0] )
  if err != nil {
    return nil, err
  }

  return apps[0], nil
}