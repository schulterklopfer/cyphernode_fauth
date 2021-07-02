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
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "gopkg.in/validator.v2"
)

func Create(model interface{} ) error {
  db := dataSource.GetDB()
  err := validator.Validate( model )
  if err != nil {
    return err
  }
  return db.Create( model ).Error
}


func Get( model interface{}, id uint, recursive bool ) error {
  db := dataSource.GetDB()
  db.Take(model, id)
  if recursive {
    err := LoadRoles(model)
    if err != nil {
      return err
    }
  }
  return nil
}

func Update( model interface{} ) error {
  db := dataSource.GetDB()

  err := validator.Validate( model )

  if err != nil {
    return err
  }
  return db.Save( model ).Error
}

func Find( out interface{}, where []interface{}, order string, limit int, offset int, recursive bool ) error {

  /*
     where == nil -> no where
     order == "" -> no order
     limit == -1 -> no limit
     offset == 0 -> no offset
  */

  db := dataSource.GetDB()

  if len(where) > 0 {
    db = db.Where( where[0].(string), where[1:] )
  }

  if order != "" {
    db = db.Order( order )
  }

  if limit != -1 {
    db = db.Limit( limit )
  }

  if offset > 0 {
    db = db.Offset( offset )
  }

  db.Find( out )

  if recursive {
    switch out.(type) {
    case *[]*models.UserModel:
      users := *out.(*[]*models.UserModel)
      for i:=0; i<len(users); i++ {
        _ = LoadRoles(users[i])
      }
    case *[]*models.AppModel:
      apps := *out.(*[]*models.AppModel)
      for i:=0; i<len(apps); i++ {
        _ = LoadRoles(apps[i])
      }
    }
  }

  return db.Error

}

func LoadRoles( in interface{} ) error {
  db := dataSource.GetDB()
  var roles []*models.RoleModel
  switch in.(type) {
  case *models.UserModel:
    if in.(*models.UserModel).ID > 0 {
      db.Model(in).Association("Roles").Find(&roles)
      in.(*models.UserModel).Roles = roles
    }
  case *models.AppModel:
    if in.(*models.AppModel).ID > 0 {
      db.Model(in).Association("AvailableRoles").Find(&roles)
      in.(*models.AppModel).AvailableRoles = roles
    }
  }
  return db.Error
}
