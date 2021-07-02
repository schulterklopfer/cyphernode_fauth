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

package dataSource_test

import (
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/sirupsen/logrus"
  "testing"
)

func TestDataSource(t *testing.T) {
  logwrapper.Logger().SetLevel( logrus.PanicLevel )
  dbDsn := "host=localhost port=5432 user=cnadmin password=cnadmin dbname=cnadmin sslmode=disable"
  dataSource.Init(dbDsn)

  t.Run("testCreateApp", testCreateApp )
  t.Run("testLoadApp", testLoadApp )
  t.Run("testLoadRole", testLoadRole )
  t.Run("testCreateUser", testCreateUser )
  t.Run("testLoadUser", testLoadUser )

  dataSource.Close()

}

func testCreateApp(t *testing.T) {

  app1 := new(models.AppModel)
  app1.Name = "app1"
  app1.Description = "description"

  role1 := new(models.RoleModel)
  role1.Name = "role1"
  role1.Description = "description"

  role2 := new(models.RoleModel)
  role2.Name = "role2"
  role2.Description = "description"

  roles1 := [2]*models.RoleModel{role1, role2}

  app1.AvailableRoles = roles1[:]

  db := dataSource.GetDB()
  db.Create(app1)

  if app1.ID == 0 || role1.ID == 0 || role2.ID == 0 {
    t.Error("Failed to insert app")
  }

}

func testLoadApp(t *testing.T) {
  var app models.AppModel
  db := dataSource.GetDB()
  db.First(&app, 1)

  var roles []*models.RoleModel

  db.Model(&app).Association("AvailableRoles").Find(&roles)

  if app.Name != "app1" || app.ID != 1 || roles[0].ID != 1 || roles[1].ID != 2 {
    t.Error("Failed to load app")
  }
}

func testLoadRole(t *testing.T) {
  var role models.RoleModel
  db := dataSource.GetDB()
  db.First(&role, 1)

  var app models.AppModel

  db.First(&app, role.AppId)

  if app.Name != "app1" || app.ID != 1 {
    t.Error("Failed to load role")
  }
}

func testCreateUser(t *testing.T) {
  user := new(models.UserModel)
  user.Login = "login"
  user.Name = "Test user"
  user.EmailAddress = "user@email.com"

  var app models.AppModel

  db := dataSource.GetDB()
  db.First(&app, 1)

  var roles []*models.RoleModel

  db.Model(&app).Association("AvailableRoles").Find(&roles)

  user.Roles = roles[0:1]

  db.Create(user)

  if user.ID == 0 {
    t.Error("Failed to insert user")
  }
}

func testLoadUser(t *testing.T) {
  var user models.UserModel

  db := dataSource.GetDB()
  db.First(&user, 1)

  var roles []*models.RoleModel

  db.Model(&user).Association("Roles").Find(&roles)

  if user.Login != "login" || user.ID != 1 || roles[0].ID != 1 {
    t.Error("Failed to load user")
  }
}
