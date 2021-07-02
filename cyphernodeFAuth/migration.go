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

package cyphernodeFAuth

import (
  "github.com/schulterklopfer/cyphernode_fauth/cnaErrors"
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_fauth/password"
  "github.com/schulterklopfer/cyphernode_fauth/queries"
)

const ADMIN_APP_NAME string = "Cyphernode Admin"
const ADMIN_APP_DESCRIPTION string = "Manage your cyphernode"

const ADMIN_APP_ADMIN_ROLE_NAME string = "admin"
const ADMIN_APP_ADMIN_ROLE_DESCRIPTION string = "Main admin with god mode"

const ADMIN_APP_USER_ROLE_NAME string = "user"
const ADMIN_APP_USER_ROLE_DESCRIPTION string = "Regular user"

func (cyphernodeFAuth *CyphernodeFAuth) migrate() error {

  // Create adminUser id=1, cyphernodeFAuth id=1, adminRole id=1
  adminRole := new(models.RoleModel)
  userRole := new(models.RoleModel)
  adminApp := new(models.AppModel)
  adminUser := new(models.UserModel)

  db := dataSource.GetDB()

  _ = queries.Get( adminRole, 1, true )
  _ = queries.Get( userRole, 2, true )
  _ = queries.Get( adminApp, 1, true )
  _ = queries.Get( adminUser, 1, true )

  hashedPassword, err := password.HashPassword( cyphernodeFAuth.Config.InitialAdminPassword )
  if err != nil {
    return err
  }

  tx := db.Begin()

  if adminApp.ID != 1 {
    logwrapper.Logger().Info("adding admin app")
    adminApp.ID = 1
    adminApp.Name = ADMIN_APP_NAME
    adminApp.Description = ADMIN_APP_DESCRIPTION
    adminApp.Hash = buildAdminHash( adminApp.Name, globals.CYPHERAPPS_REPO )
    adminApp.MountPoint = globals.BASE_ADMIN_MOUNTPOINT
    adminApp.Version = globals.VERSION
    adminApp.Meta = &models.Meta{
      Icon:  "data:image/png;base64,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",
      Color: "#000000",
    }
    adminApp.AccessPolicies = models.AccessPolicies{
      /* General stuff */
      {
        Patterns: []string{"favicon.ico$"},
        Roles: []string{"*"},
        Actions: []string{"options","get"},
        Effect: "allow",
      },
      /* API endpoints */
      {
        Patterns: []string{"^\\/api\\/v0\\/login$"},
        Roles: []string{"*"},
        Actions: []string{"options","post"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/users","^\\/api\\/v0\\/docker","^\\/api\\/v0\\/files"},
        Roles: []string{"admin"},
        Actions: []string{"options","get","post","patch","delete"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/apps","^\\/api\\/v0\\/status","^\\/api\\/v0\\/blocks","^\\/api\\/v0\\/users\\/me$"},
        Roles: []string{"*"},
        Actions: []string{"options","get"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/api\\/v0\\/apps"},
        Roles: []string{"admin"},
        Actions: []string{"options","post","patch"},
        Effect: "allow",
      },
      {
        Patterns: []string{"^\\/$", "^\\/_\\/"},
        Roles: []string{"*"},
        Actions: []string{"options","get"},
        Effect: "allow",
      },
    }
    if adminApp.Hash == "" {
      return cnaErrors.ErrMigrationFailed
    }
    tx.Create(adminApp)
  }

  if adminRole.ID != 1 {
    logwrapper.Logger().Info("adding admin role")
    adminRole.ID = 1
    adminRole.Name = ADMIN_APP_ADMIN_ROLE_NAME
    adminRole.Description = ADMIN_APP_ADMIN_ROLE_DESCRIPTION
    adminRole.AutoAssign = false
    adminRole.AppId = 1
    tx.Create(adminRole)
  }

  if userRole.ID != 2 {
    logwrapper.Logger().Info("adding user role")
    userRole.ID = 2
    userRole.Name = ADMIN_APP_USER_ROLE_NAME
    userRole.Description = ADMIN_APP_USER_ROLE_DESCRIPTION
    userRole.AutoAssign = true
    userRole.AppId = 1
    tx.Create(userRole)
  }

  adminAppHasAdminRole := false
  adminAppHasUserRole := false
  for i:=0; i<len(adminApp.AvailableRoles); i++ {
    if adminApp.AvailableRoles[i].ID == adminRole.ID {
      adminAppHasAdminRole = true
    }
    if adminApp.AvailableRoles[i].ID == userRole.ID {
      adminAppHasUserRole = true
    }
  }

  if !adminAppHasAdminRole {
    tx.Model(&adminApp).Association("AvailableRoles").Append(adminRole)
  }

  if !adminAppHasUserRole {
    tx.Model(&adminApp).Association("AvailableRoles").Append(userRole)
  }

  if adminUser.ID != 1 {
    logwrapper.Logger().Info("adding admin user")
    adminUser.ID = 1
    adminUser.Login = cyphernodeFAuth.Config.InitialAdminLogin
    adminUser.Password = hashedPassword
    adminUser.Name = cyphernodeFAuth.Config.InitialAdminName
    adminUser.EmailAddress = cyphernodeFAuth.Config.InitialAdminEmailAddress
    tx.Create(adminUser)
  }

  for _, role := range []*models.RoleModel{adminRole,userRole} {
    tx.Model(&adminUser).Association("Roles").Append(role)
  }

  return tx.Commit().Error

}

func buildAdminHash( label string, sourceLocation string ) string {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(label)... )
  bytes = append( bytes, []byte(sourceLocation)... )
  return helpers.TrimmedRipemd160Hash( bytes )
}

