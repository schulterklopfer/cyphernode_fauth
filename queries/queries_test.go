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

package queries_test

import (
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_fauth/queries"
  "github.com/sirupsen/logrus"
  "math/rand"
  "os"
  "strconv"
  "testing"
  "time"
)

func TestQueries(t *testing.T) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  logwrapper.Logger().SetLevel( logrus.PanicLevel )

  dbFile := "/tmp/tests_"+strconv.Itoa(r.Intn(1000000 ))+".sqlite3"

  dataSource.Init(dbFile)

  t.Run( "Create app", createApp )
  t.Run( "Create role", createRole )
  t.Run( "Create user", createUser )
  t.Run( "Remove role from user", removeRoleFromUser )
  t.Run( "Add role to user", addRoleToUser )
  t.Run( "Get app", getApp )
  t.Run( "Find app", findApps)
  t.Run( "Get role", getRole )
  t.Run( "Find role", findRole )
  t.Run("Get user", getUser )
  t.Run("Find users", findUsers )
  t.Run( "Delete role", deleteRole )
  t.Run( "Delete user", deleteUser )
  t.Run( "Delete app", deleteApp )

  t.Run( "Hooks", func(t *testing.T) {
    app := new ( models.AppModel )
    app.Name = "App1"
    app.Description = "Description"
    _ = queries.CreateApp(app)

    user := new ( models.UserModel )
    user.Name = "User1"
    user.Login = "login1"
    user.Password = "test123"
    user.EmailAddress = "user@user.com"
    _ = queries.Create(user)

    t.Run( "Role Autoassign", roleAutoAssign )
  })

  dataSource.Close()
  os.Remove(dbFile)

}

func createApp( t *testing.T ) {

  var app *models.AppModel

  app = new(models.AppModel)
  app.Name = "First app"
  app.Description = "First app description"

  err := queries.CreateApp(app)

  err = queries.CreateApp(app)

  if err == nil {
    t.Error( "Create app with app id" )
  }

  app.ID = 0
  err = queries.CreateApp(app)

  if err == nil {
    t.Error( "Created same app twice" )
  }

  app = new(models.AppModel)
  app.Name = "Second app"
  app.Description = "Second app description"

  _ = queries.CreateApp(app)
}

func deleteApp( t *testing.T ) {
  db := dataSource.GetDB()
  var app *models.AppModel

  err := queries.DeleteApp( 0 )
  if err == nil {
    t.Error( "Deleted app with no primary key" )
  }

  app = new( models.AppModel )
  db.Take(app, 1)

  _ = queries.DeleteApp( 1 )

  app = new( models.AppModel )
  db.Take(app, 1)

  if app.ID != 0 {
    t.Error( "App was not deleted" )
  }

  app = new( models.AppModel )
  db.Take(app, 2)

  _ = queries.DeleteApp( 2 )

  app = new( models.AppModel )
  db.Take(app, 2)

  if app.ID != 0 {
    t.Error( "App was not deleted" )
  }
}

func getApp( t *testing.T ) {
  var app *models.AppModel

  app = new(  models.AppModel )
  _ = queries.Get(app, 1,true)

  if app.ID == 0 {
    t.Error( "unable to load app")
  }

  app = new(  models.AppModel )
  _ = queries.Get(app, 3,true)

  if app.ID != 0 {
    t.Error( "App id should be 0")
  }

}

func findApps( t *testing.T ) {
  var apps []*models.AppModel

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps, nil, "", -1,0, true)

  if apps == nil {
    t.Error( "unable to load all apps")
  }

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps, nil, "", 0,0,true)

  if apps == nil || len(apps) != 0 {
    t.Error( "unable to load 0 apps")
  }

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps, nil, "", 1,0,true)

  if apps == nil || len(apps) != 1 || apps[0].Name != "First app" || len(apps[0].AvailableRoles) != 2 {
    t.Error( "unable to load 1 app")
  }

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps,  nil, "", 1,1,true)

  if apps == nil || len(apps) != 1 || apps[0].Name != "Second app" {
    t.Error( "unable to load 1 app with offset 1")
  }

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps,  nil, "id desc", 1,1,true)

  if apps == nil || len(apps) != 1 || apps[0].Name != "First app" {
    t.Error( "unable to load 1 app with offset 1 order by id desc")
  }

  apps = make( []*models.AppModel, 0 )
  _ = queries.Find( &apps,  []interface{}{"name LIKE ?", "% app"}, "name", -1,0,true)

  if apps == nil || len(apps) != 2 || apps[0].Name != "First app" || apps[1].Name != "Second app" {
    t.Error( "unable to load 2 apps with order by name")
  }

}

func createRole( t *testing.T ) {

  var role *models.RoleModel

  role = new(models.RoleModel)
  role.Name = "admin"
  role.Description = "Admin of first app"
  role.AppId = 1

  err := queries.CreateRole(role)

  err = queries.CreateRole(role)

  if err == nil {
    t.Error( "Create role with role id" )
  }

  role = new(models.RoleModel)
  role.Name = "user"
  role.AutoAssign = true
  role.Description = "User of frist app"
  role.AppId = 1

  _ = queries.CreateRole(role)

  role = new(models.RoleModel)
  role.Name = "admin"
  role.Description = "Admin of second app"
  role.AppId = 2

  _ = queries.CreateRole(role)

  role = new(models.RoleModel)
  role.Name = "user"
  role.AutoAssign = true
  role.Description = "User of second app"
  role.AppId = 2

  _ = queries.CreateRole(role)

}

func deleteRole( t *testing.T ) {
  db := dataSource.GetDB()
  var role *models.RoleModel

  err := queries.DeleteRole( 0 )
  if err == nil {
    t.Error( "Deleted role with no primary key" )
  }

  role = new( models.RoleModel )
  db.Take(role, 1)

  _ = queries.DeleteRole( 1 )

  role = new( models.RoleModel )
  db.Take(role, 1)

  if role.ID != 0 {
    t.Error( "Role was not deleted" )
  }

  role = new( models.RoleModel )
  db.Take(role, 2)

  _ = queries.DeleteRole( 2 )

  role = new( models.RoleModel )
  db.Take(role, 2)

  if role.ID != 0 {
    t.Error( "Role was not deleted" )
  }

  role = new( models.RoleModel )
  db.Take(role, 3)

  _ = queries.DeleteRole( 3 )

  role = new( models.RoleModel )
  db.Take(role, 3)

  if role.ID != 0 {
    t.Error( "Role was not deleted" )
  }

  role = new( models.RoleModel )
  db.Take(role, 4)

  _ = queries.DeleteRole( 4 )

  role = new( models.RoleModel )
  db.Take(role, 4)

  if role.ID != 0 {
    t.Error( "Role was not deleted" )
  }
}

func getRole( t *testing.T ) {
  var role *models.RoleModel

  role = new(  models.RoleModel )
  _ = queries.Get(role, 1,false)

  if role.ID == 0 {
    t.Error( "unable to load role")
  }

  role = new(  models.RoleModel )
  _ = queries.Get(role, 5,false)

  if role.ID != 0 {
    t.Error( "Role id should be 0")
  }
}

func findRole( t *testing.T ) {
  var roles []*models.RoleModel

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles, nil, "", -1,0, false)

  if roles == nil {
    t.Error( "unable to load all roles")
  }

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles, nil, "", 0,0,false)

  if roles == nil || len(roles) != 0 {
    t.Error( "unable to load 0 roles")
  }

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles, nil, "", 1,0,false)

  if roles == nil || len(roles) != 1 || roles[0].Name != "admin" {
    t.Error( "unable to load 1 role")
  }

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles,  nil, "", 1,1,false)

  if roles == nil || len(roles) != 1 || roles[0].Name != "user" {
    t.Error( "unable to load 1 role with offset 1")
  }

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles,  nil, "id desc", 1,1,false)

  if roles == nil || len(roles) != 1 || roles[0].Name != "admin" {
    t.Error( "unable to load 1 role with offset 1 order by id desc")
  }

  roles = make( []*models.RoleModel, 0 )
  _ = queries.Find( &roles,  []interface{}{"name LIKE ?", "user"}, "id", 2,0,false)

  if roles == nil || len(roles) != 2 || roles[0].Name != "user" || roles[1].Name != "user" {
    t.Error( "unable to load 2 roles with order by name")
  }

}

func createUser( t *testing.T ) {

  db := dataSource.GetDB()
  role := new(models.RoleModel)

  db.Take( &role, 1)

  roles := make( []*models.RoleModel, 1 )
  roles[0]=role

  var user *models.UserModel

  user = new(models.UserModel)
  user.Login = "login1"
  user.Password = "password1"
  user.Name ="name1"
  user.EmailAddress = "email1@email.rocks"
  user.Roles = roles

  err := queries.Create( user )

  err = queries.Create( user )

  if err == nil {
    t.Error( "Create user with user id" )
  }

  user.ID = 0
  err = queries.Create( user )

  if err == nil {
    t.Error( "Created same user twice" )
  }

  user = new(models.UserModel)
  user.Login = "login2"
  user.Password = "password2"
  user.Name ="name2"
  user.EmailAddress = "email2@email.rocks"
  user.Roles = roles

  err = queries.Create( user )

  if err != nil {
    t.Error( "Failed to create second user" )
  }
}

func removeRoleFromUser( t *testing.T ) {
  // remove role id 4
  var user *models.UserModel
  user = new(  models.UserModel )
  _ = queries.Get( user, 1,true)

  queries.RemoveRoleFromUser( user, 4 )
}

func addRoleToUser( t *testing.T ) {
  // add role id 4
  var user *models.UserModel
  user = new(  models.UserModel )
  _ = queries.Get( user, 1,true)

  queries.AddRoleToUser( user, 4 )

}

func deleteUser( t *testing.T) {

  db := dataSource.GetDB()
  var user *models.UserModel

  err := queries.DeleteUser( 0 )
  if err == nil {
    t.Error( "Deleted user with no primary key" )
  }

  user = new( models.UserModel )
  db.Take( user, 1)

  _ = queries.DeleteUser( 1 )

  user = new( models.UserModel )
  db.Take( user, 1)

  if user.ID != 0 {
    t.Error( "User was not deleted" )
  }

  user = new( models.UserModel )
  db.Take( user, 2)

  _ = queries.DeleteUser( 2 )

  user = new( models.UserModel )
  db.Take( user, 2)

  if user.ID != 0 {
    t.Error( "User was not deleted" )
  }

}

func getUser( t *testing.T ) {
  var user *models.UserModel

  user = new(  models.UserModel )
  _ = queries.Get( user, 1,true)

  if user.ID == 0 {
    t.Error( "unable to load user")
  }

  user = new(  models.UserModel )
  _ = queries.Get( user, 3,true)

  if user.ID != 0 {
    t.Error( "User id should be 0")
  }

}

func findUsers( t *testing.T ) {
  var users []*models.UserModel

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users, nil, "", -1,0, true)

  if users == nil {
    t.Error( "unable to load all users")
  }

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users, nil, "", 0,0,true)

  if users == nil || len(users) != 0 {
    t.Error( "unable to load 0 users")
  }

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users, nil, "", 1,0,true)

  if users == nil || len(users) != 1 || users[0].Name != "name1" || len(users[0].Roles) != 3 {
    t.Error( "unable to load 1 user")
  }

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users,  nil, "", 1,1,true)

  if users == nil || len(users) != 1 || users[0].Name != "name2" {
    t.Error( "unable to load 1 user with offset 1")
  }

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users,  nil, "id desc", 1,1,true)

  if users == nil || len(users) != 1 || users[0].Name != "name1" {
    t.Error( "unable to load 1 user with offset 1 order by id desc")
  }

  users = make( []*models.UserModel, 0 )
  _ = queries.Find( &users,  []interface{}{"name LIKE ?", "name%"}, "name", -1,0,true)

  if users == nil || len(users) != 2 || users[0].Name != "name1" || users[1].Name != "name2" {
    t.Error( "unable to load 2 users with order by name")
  }

}

func roleAutoAssign( t *testing.T ) {
  var user *models.UserModel

  user = new (models.UserModel)
  _ = queries.Get( user, 3, true )

  var role models.RoleModel

  role.Name = "autoassign"
  role.AutoAssign = true
  role.AppId = 3

  _ = queries.CreateRole(&role)
  _ = queries.Get( user, 3, true )

  if len(user.Roles) != 1 {
    t.Error( "Autoassign failed: new role" )
  }

  user = new (models.UserModel)
  user.Login = "login"
  user.Password = "test123"
  user.EmailAddress = "user@user.com"

  _ = queries.Create( user )

  if len(user.Roles) != 1 {
    t.Error( "Autoassign failed: new user" )
  }
}

