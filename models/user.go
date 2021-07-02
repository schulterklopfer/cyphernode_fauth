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

package models

import (
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_fauth/cnaErrors"
)

type UserModel struct {
  gorm.Model
  Login string `json:"login" gorm:"type:varchar(30);unique_index;not null" form:"login" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\-]+$"`
  Name string `json:"name" form:"name"` // optional
  Password string `json:"password" gorm:"type:varchar(128);not null" form:"password" validate:"nonzero" sbjt:"hashPassword"`
  EmailAddress string `json:"email_address" gorm:"type:varchar(100)" form:"emailAddress" validate:"max=100,regexp=(^$|^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"`
  Roles []*RoleModel `json:"roles" gorm:"many2many:user_roles;association_autoupdate:false;gorm:association_autocreate:false" form:"roles" validate:"-"`
}

func (user *UserModel) AfterCreate( tx *gorm.DB ) (err error) {
  var allAutoAssignRoles []*RoleModel
  tx.Where( &RoleModel{ AutoAssign: true }).Find( &allAutoAssignRoles )
  for i:=0; i< len(allAutoAssignRoles); i++ {
    tx.Model(user).Association("Roles").Append(allAutoAssignRoles[i])
  }
  return
}

func (user *UserModel) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if user.ID == 0 {
    err = cnaErrors.ErrNoSuchUser
    return
  }
  return
}

func (user *UserModel) AfterDelete( tx *gorm.DB ) (err error) {
  tx.Model(user).Association("Roles").Clear()
  return
}

func (user *UserModel) BeforeSave( tx *gorm.DB ) (err error) {
  err = user.checkDuplicate(tx)
  if err != nil {
    return
  }
  err = user.checkRoles(tx)
  if err != nil {
    return
  }
  return
}

func (user *UserModel) BeforeCreate( tx *gorm.DB ) (err error) {
  err = user.checkDuplicate(tx)
  if err != nil {
    return
  }
  err = user.checkRoles(tx)
  if err != nil {
    return
  }
  return
}

func (user *UserModel) checkDuplicate( tx *gorm.DB ) error {
  var existingUsers []UserModel
  tx.Limit(1).Find( &existingUsers, "login = ? AND id != ?", user.Login, user.ID )

  if len(existingUsers) > 0 {
    return cnaErrors.ErrDuplicateUser
  }
  return nil
}

func (user *UserModel) checkRoles( tx *gorm.DB ) error {
  for i:=0; i<len(user.Roles ); i++ {
    if user.Roles[i].ID == 0 {
      return cnaErrors.ErrUserHasUnknownRole
    }
    var role RoleModel
    tx.Take( &role,  user.Roles[i].ID )
    if role.ID != user.Roles[i].ID {
      return cnaErrors.ErrUserHasUnknownRole
    }
  }
  return nil
}
