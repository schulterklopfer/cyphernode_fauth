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

package main

import (
  "github.com/schulterklopfer/cyphernode_fauth/cyphernodeFAuth"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "github.com/sirupsen/logrus"
  "log"
  "net/http"
  "os"
  _ "net/http/pprof"
)


func main() {

  go func() {
    log.Println( http.ListenAndServe("localhost:6060", nil))
  }()

  logwrapper.Logger().SetLevel(logrus.TraceLevel)

  app := cyphernodeFAuth.NewCyphernodeFAuth( &cyphernodeFAuth.Config{
      DatabaseDsn: helpers.GetenvOrDefault(globals.CNA_ADMIN_DATABASE_DSN_ENV_KEY ),
      InitialAdminLogin: helpers.GetenvOrDefault( globals.CNA_ADMIN_LOGIN_ENV_KEY ),
      InitialAdminPassword: helpers.GetenvOrDefault(globals.CNA_ADMIN_PASSWORD_ENV_KEY ),
      InitialAdminName: helpers.GetenvOrDefault(globals.CNA_ADMIN_NAME_ENV_KEY ),
      InitialAdminEmailAddress: helpers.GetenvOrDefault(globals.CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY ),
    },
  )
  err := app.Init()
  if err != nil {
    println("Error in application init: ", err.Error() )
    os.Exit(1)
  }
  app.Start()
}