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
  "database/sql/driver"
  "encoding/json"
  "errors"
  "fmt"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/jinzhu/gorm"
  "github.com/schulterklopfer/cyphernode_fauth/cnaErrors"
)

type AccessPolicies []*storage.AccessPolicy

// this is used to create jsonb Data which then
// will ne saved to the db by gorm
func (aps AccessPolicies) Value() (driver.Value, error) {
  jsonValue, err := json.Marshal(aps)
  if err != nil {
    return nil, err
  }
  return string(jsonValue), nil
}

// convert jsonb Data from database back into
// a struct
func (aps *AccessPolicies) Scan(value interface{}) error {
  jsonValue, ok := value.([]byte)
  if !ok {
    return errors.New(fmt.Sprint("Failed to unmarshal access policies:", value))
  }
  err := json.Unmarshal(jsonValue, aps)
  return err
}


type Meta struct {
  Icon  string `json:"icon,omitempty"`
  Color string `json:"color,omitempty"`
}


// this is used to create jsonb Data which then
// will ne saved to the db by gorm
func (meta Meta) Value() (driver.Value, error) {
  jsonValue, err := json.Marshal(meta)
  if err != nil {
    return nil, err
  }
  return string(jsonValue), nil
}

// convert jsonb Data from database back into
// a struct
func (meta *Meta) Scan(value interface{}) error {
  jsonValue, ok := value.([]byte)
  if !ok {
    return errors.New(fmt.Sprint("Failed to unmarshal access policies:", value))
  }
  err := json.Unmarshal(jsonValue, meta)
  return err
}


type AppModel struct {
  gorm.Model
  Hash           string         `json:"hash" gorm:"type:varchar(32);unique_index;not null"`
  Secret         string         `json:"-" gorm:"type:varchar(32);unique_index;not null"`
  MountPoint     string         `json:"mountPoint" gorm:"type:varchar(32);unique_index;not null"`
  Name           string         `json:"name" gorm:"type:varchar(30);not null" validate:"min=3,max=30,regexp=^[a-zA-Z0-9_\\- ]+$"`
  Description    string         `json:"description" gorm:"type:varchar(255)"`
  Version        string         `json:"version" gorm:"type:varchar(255)"`
  AvailableRoles []*RoleModel   `json:"availableRoles" gorm:"foreignkey:AppId;preload"`
  AccessPolicies AccessPolicies `json:"accessPolicies,omitempty" gorm:"type:jsonb;default:'null'"`
  Meta           *Meta          `json:"meta,omitempty" gorm:"type:jsonb;default:'null'"`
}

func ( app *AppModel ) AfterDelete( tx *gorm.DB ) {
  var roles []RoleModel
  tx.Model(app).Association("AvailableRoles" ).Find(&roles)
  for i:=0; i< len(roles); i++ {
    tx.Delete( roles[i] )
    // Why do I have to call this manually?
    roles[i].AfterDelete( tx )
  }
}

func ( app *AppModel ) BeforeDelete( tx *gorm.DB ) (err error) {
  // very important. if no check, will delete all users if ID == 0
  if app.ID == 0 {
    err = cnaErrors.ErrNoSuchApp
    return
  }
  return
}