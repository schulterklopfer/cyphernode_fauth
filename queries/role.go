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
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "gopkg.in/validator.v2"
)

func CreateRole( role *models.RoleModel ) error {

  if role.ID != 0 {
    // role must not have any ID possibly existing in DB
    return errors.New( "role ID must be 0" )
  }

  db := dataSource.GetDB()
  err := validator.Validate(role)
  if err != nil {
    return err
  }
  db.Create(role)
  return nil
}

func DeleteRole( id uint ) error {
  if id == 0 {
    return errors.New("no such role")
  }
  if id == 1 {
    return globals.ErrActionForbidden
  }
  db := dataSource.GetDB()
  var role models.RoleModel
  db.Take( &role, id )
  if role.ID == 0 {
    return errors.New("no such role")
  }
  db.Unscoped().Delete( &role)
  role.ID = 0
  return nil
}

func UsersForRole( users *[]*models.UserModel, role *models.RoleModel ) error {
  if role == nil {
    return errors.New("no such role")
  }
  db := dataSource.GetDB()
  db.Model(role).Association("Users").Find( users )
  return nil
}

func AllRoles( allRoles *[]models.RoleModel ) error {
  db := dataSource.GetDB()
  return db.Find( allRoles ).Error
}