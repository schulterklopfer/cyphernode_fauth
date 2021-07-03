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
  "github.com/schulterklopfer/cyphernode_fauth/globals"
)

type RoleModel struct {
  gorm.Model
  Name string `json:"name" gorm:"type:varchar(30)" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_-]+$"`
  Description string `json:"description" gorm:"type:varchar(255)"`
  AutoAssign bool `json:"autoAssign" gorm:"default false"`
  AppId uint `json:"appId"`
  Users []*UserModel `json:"users" gorm:"many2many:user_roles;"`
}

func ( role *RoleModel ) AfterDelete( tx *gorm.DB ) {
  role.removeFromAllUsers( tx )
}

func ( role *RoleModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if role.ID == 0 {
    err = globals.ErrNoSuchRole
    return
  }
  return
}

func ( role *RoleModel ) AfterSave( tx *gorm.DB ) {
  role.AfterUpdate( tx )
}

func ( role *RoleModel ) AfterUpdate( tx *gorm.DB ) {
  if role.AutoAssign {
    role.addToAllUsers( tx )
  } else {
    role.removeFromAllUsers( tx )

    // reassign to admin user
    var adminUser UserModel
    tx.First( &adminUser, 1 )
    if adminUser.ID == 1 {
      tx.Model(adminUser).Association("Roles").Append( role )
    }
  }
}

func ( role *RoleModel) AfterCreate( tx *gorm.DB )  {
  // all roles are given to the admin user
  var adminUser UserModel
  tx.First( &adminUser, 1 )
  if adminUser.ID == 1 {
    tx.Model(adminUser).Association("Roles").Append( role )
  }

  if !role.AutoAssign {
    return
  }
  role.addToAllUsers( tx )
}

func ( role *RoleModel) addToAllUsers( tx *gorm.DB ) {
  var allUsers []*UserModel
  tx.Find( &allUsers )
  for i:=0; i< len(allUsers); i++ {
    tx.Model(allUsers[i]).Association("Roles").Append( role )
  }
}

func ( role *RoleModel) removeFromAllUsers( tx *gorm.DB ) {
  tx.Model(role).Association("Users" ).Clear()
}


