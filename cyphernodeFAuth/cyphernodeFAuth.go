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
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_fauth/appList"
  "github.com/schulterklopfer/cyphernode_fauth/dataSource"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "golang.org/x/sync/errgroup"
)

type Config struct {
  DatabaseDsn string
  InitialAdminLogin string
  InitialAdminPassword string
  InitialAdminName string
  InitialAdminEmailAddress string
  DisableAuth bool
}

type CyphernodeFAuth struct {
  Config         *Config
  engineInternal *gin.Engine
  engineExternal *gin.Engine
  engineAuth     *gin.Engine
  routerGroups   map[string]*gin.RouterGroup
}

var instance *CyphernodeFAuth

func NewCyphernodeFAuth(config *Config) *CyphernodeFAuth {
  instance = new(CyphernodeFAuth)
  instance.Config = config
  return instance
}

func Get() *CyphernodeFAuth {
  return instance
}

func (cyphernodeFAuth *CyphernodeFAuth) Init() error {

  err := dataSource.Init(cyphernodeFAuth.Config.DatabaseDsn)
  if err != nil {
    logwrapper.Logger().Error("Failed to connect to database" )
    return err
  }

  cyphernodeFAuth.routerGroups = make(map[string]*gin.RouterGroup)
  err = cyphernodeFAuth.migrate()
  if err != nil {
    logwrapper.Logger().Error("Failed to init database" )
    return err
  }

  cyphernodeFAuth.engineAuth = gin.New()
  cyphernodeFAuth.initAuthHandlers()

  err = appList.Init( helpers.GetenvOrDefault( globals.CYPHERAPPS_INSTALL_DIR_ENV_KEY ) )
  if err != nil {
    logwrapper.Logger().Error("Failed to init applist" )
    return err
  }

  return nil
}

func (cyphernodeFAuth *CyphernodeFAuth) Engine() *gin.Engine {
  return cyphernodeFAuth.engineExternal
}

func (cyphernodeFAuth *CyphernodeFAuth) Start() {

  var g errgroup.Group

  g.Go(func() error {
    return  cyphernodeFAuth.engineAuth.Run(":3032")
  })

  if err := g.Wait(); err != nil {
    logwrapper.Logger().Fatal(err)
  }

}
