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

package forwardAuth

import (
  "encoding/hex"
  "errors"
  "fmt"
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_fauth/queries"
  "net/http"
)

func ForwardAppAuth( c *gin.Context ) {

  // get symmetrically signed token
  tokenString := helpers.TokenFromBearerAuthHeader( c.Request.Header.Get("authorization") )

  if tokenString == "" {
    c.Status(http.StatusUnauthorized)
    return
  }


  _, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Don't forget to validate the alg is what you expect:
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
      return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
    }

    appIdFloat, exists := token.Claims.(jwt.MapClaims)["id"]

    if !exists {
      return nil, errors.New("No app id in claims")
    }

    appId := uint(appIdFloat.(float64))
    var app models.AppModel
    err := queries.Get( &app, appId,false )

    if err != nil {
      return nil, err
    }

    return hex.DecodeString( app.Secret )
  })

  if err != nil {
    c.Status(http.StatusUnauthorized)
    return
  }



  c.Next()
}
