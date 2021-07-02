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

package appList

import (
  "github.com/SatoshiPortal/cam/storage"
  camUtils "github.com/SatoshiPortal/cam/utils"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_fauth/queries"
  "os"
  "sync"
  "time"
)

var appList *AppList

type AppList struct {
  LastFileInfo os.FileInfo
  LastUpate time.Time
  InstalledApps *storage.InstalledAppsIndex
  mutex sync.Mutex
}

func Get() *AppList {
  return appList
}

func Init( whitelistFilePath string ) error {
  if appList != nil {
    return nil
  }
  appList = &AppList{ InstalledApps: &storage.InstalledAppsIndex{} }
  err := appList.load()
  if err != nil {
    return err
  }
  helpers.SetInterval(appList.checkForChange, 1000, false)
  return nil
}

func (appList *AppList) checkForChange() {
  fileInfo, err := os.Stat( camUtils.GetInstalledAppsIndexFilePath() )
  if err != nil {
    return
  }
  if appList.LastFileInfo != nil && (
      appList.LastFileInfo.Size() != fileInfo.Size() ||
      appList.LastFileInfo.ModTime().Before( fileInfo.ModTime() ) ) {
    err := appList.load()
    if err != nil {
      logwrapper.Logger().Errorf( "Failed to load whitelist: %s", err.Error() )
    }
    appList.LastUpate = time.Now()
  }
  appList.LastFileInfo = fileInfo
}

/*
func (appList *AppList) ContainsClientID( clientID string ) bool {
  appList.mutex.Lock()
  defer appList.mutex.Unlock()
  if len(appList.entries) == 0 {
    return false
  }

  if helpers.SliceIndex( len(appList.entries), func(i int) bool {
    return appList.entries[i].Hash == clientID
  } ) >= 0 {
    return true
  }
  return false
}
*/

func (appList *AppList) load() error {

  appList.mutex.Lock()
  defer appList.mutex.Unlock()
  err := appList.InstalledApps.Load()
  if err != nil {
    return err
  }

  err = appList.syncToDb()
  if err != nil {
    return err
  }

  appList.LastUpate = time.Now()
  return nil
}

func (appList *AppList) syncToDb() error {

  // 1) go through apps in applist and see if they exist in the db
  // if not, create them
  logwrapper.Logger().Debug("Syncing to database")
  for _, app := range appList.InstalledApps.Apps {

    appFromDb, err := queries.GetAppByHash( app.GetHash() )
    if err != nil {
      return err
    }

    if appFromDb != nil {
      // app exists, dont insert, but check available roles

      // 1) check all roles for app in applist and see if they exist in the db
      // if not create them

      logwrapper.Logger().Debug("found app in database: "+app.Name )
      logwrapper.Logger().Debug("checking roles" )

      // update access policies and other properties
      appFromDb.MountPoint = app.MountPoint
      appFromDb.Name = app.Name
      appFromDb.Secret = app.Secret
      appFromDb.Version = app.Candidates[0].Version.Raw
      appFromDb.Meta = &models.Meta{app.Meta.Icon, app.Meta.Color }
      appFromDb.AccessPolicies = app.Candidates[0].AccessPolicies

      err := queries.Update( appFromDb )
      if err != nil {
        return err
      }

      for _, role := range app.Candidates[0].AvailableRoles {
        found := false
        for _, roleFromDb := range appFromDb.AvailableRoles {
          if roleFromDb.Name == role.Name {
            found = true
            break
          }
        }
        if !found {
          // needs to be created in db
          err := queries.CreateRoleForApp( appFromDb, &models.RoleModel{
            Name:        role.Name,
            Description: role.Description,
            AutoAssign:  role.AutoAssign,
          })
          logwrapper.Logger().Debug("creating new role in database: "+role.Name )

          if err != nil {
            return err
          }
        }
      }

      // 2) check all roles for app in db and see if they exist in the app from the applist
      // if not delete them

      // reload from db in case roles changed
      err = queries.LoadRoles(appFromDb)
      if err != nil {
        return err
      }

      for _, roleFromDb := range appFromDb.AvailableRoles {
        found := false
        for _, role := range app.Candidates[0].AvailableRoles {
          if roleFromDb.Name == role.Name {
            found = true
            break
          }
        }
        if !found {
          // needs to be created in db
          err := queries.RemoveRoleFromApp( appFromDb, roleFromDb.ID )
          logwrapper.Logger().Debug("removing role from database: "+roleFromDb.Name )

          if err != nil {
            return err
          }

          err = queries.DeleteRole( roleFromDb.ID )
          if err != nil {
            return err
          }
        }
      }
      continue
    }

    // app does not exist, create it in db
    appFromDb = &models.AppModel{
      Hash:           app.GetHash(),
      Secret:         app.Secret,
      MountPoint:     app.MountPoint,
      Name:           app.Name,
      Meta:           &models.Meta{app.Meta.Icon, app.Meta.Color },
      Version:        app.Candidates[0].Version.Raw,
      AccessPolicies: app.Candidates[0].AccessPolicies,
    }

    err = queries.CreateApp( appFromDb )
    logwrapper.Logger().Debug("creating app in database: "+appFromDb.Name )

    if err != nil {
      return err
    }
    for _, role := range app.Candidates[0].AvailableRoles {
      err = queries.CreateRoleForApp( appFromDb, &models.RoleModel{
        Name:        role.Name,
        Description: role.Description,
        AutoAssign:  role.AutoAssign,
      })
      if err != nil {
        return err
      }
    }
  }

  // 2) go through apps in database and see if they exist in the applist
  // if not, delete them

  var appsFromDb []*models.AppModel
  // exclude app id == 1, cause its the admin app
  err := queries.Find( &appsFromDb, []interface{}{"id != ?", 1 }, "", -1,0,true)
  if err != nil {
    return err
  }

  for _, appFromDb := range appsFromDb {
    found := false
    for _, app := range appList.InstalledApps.Apps {
      if app.GetHash() == appFromDb.Hash {
        found = true
        break
      }
    }

    if !found {
      // delete
      err := queries.DeleteApp( appFromDb.ID )
      logwrapper.Logger().Debug("removing app from database: "+appFromDb.Name )

      if err != nil {
        return err
      }
    }
  }

  return nil
}
