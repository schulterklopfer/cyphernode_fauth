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
  "fmt"
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/models"
  "github.com/schulterklopfer/cyphernode_fauth/queries"
  "net/http"
  "strings"
)

func ForwardUserAuth(c *gin.Context) {

  //secret := []byte("my_secret_key")

  prefix := c.Request.Header.Get("x-forwarded-prefix")
  forwardedHost := c.Request.Header.Get("x-forwarded-host")
  forwardedProto := c.Request.Header.Get("x-forwarded-proto")

  if prefix == "" {
    c.Redirect( http.StatusTemporaryRedirect, forwardedProto+"://"+forwardedHost+globals.UNAUTHORIZED_REDIRECT_URL )
    return
  }

  // x-forwarded-prefix header idetentifies the app we want to auth against
  mountPoint := prefix[1:]

  app, err := queries.GetAppByMountPoint( mountPoint )

  if err != nil {
    c.Header("X-Status-Reason", err.Error() )
    c.Status(http.StatusInternalServerError)
    return
  }

  uriInAp := c.Request.Header.Get("x-forwarded-uri")
  method := c.Request.Header.Get("x-forwarded-method")

  // access not granted. See if we have a valid token
  // and check access again
  tokenString := helpers.TokenFromBearerAuthHeader( c.Request.Header.Get("authorization") )

  if tokenString == "" {
    // lets see if there is a cookie where we can get the auth from
    if sessionCookie, err := c.Request.Cookie(helpers.GetenvOrDefault(globals.CNA_SESSION_COOKIE_NAME_ENV_KEY)); err == nil {
      tokenString = sessionCookie.Value
    }
  }

  // Parse takes the token string and a function for looking up the key. The latter is especially
  // useful if you use multiple keys for your application.  The standard is to use 'kid' in the
  // head of the token to identify which key to use, but the parsed token (head and claims) is provided
  // to the callback, providing flexibility.
  var token *jwt.Token
  if tokenString != "" {
    token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
      // Don't forget to validate the alg is what you expect:
      if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
      }
      return []byte(helpers.GetenvOrDefault(globals.CNA_COOKIE_SECRET_ENV_KEY)), nil
    })

    if err == nil {
      // set correct headers for cypherapp down the line
      parts := strings.Split( token.Raw, "." )
      c.Header("X-Auth-User-Claims", parts[1] )
    }

  }

  // check for public access
  accessGranted := false
  for _, accessPolicy := range app.AccessPolicies {
    accessGranted = accessPolicy.Check( method, uriInAp, nil )
    if accessGranted {
      break
    }
  }

  if accessGranted {
    c.Status(http.StatusOK)
    return
  }

  if token != nil {
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

      subject, exists := claims["id"]

      if !exists {
        c.Header("X-Status-Reason", "no subject claims" )
        c.Redirect( http.StatusTemporaryRedirect, forwardedProto+"://"+forwardedHost+globals.UNAUTHORIZED_REDIRECT_URL )
        return
      }

      userId := uint(subject.(float64))
      var user models.UserModel
      err := queries.Get( &user, userId,true )

      if err != nil || user.ID == 0 {
        c.Redirect( http.StatusTemporaryRedirect, forwardedProto+"://"+forwardedHost+globals.UNAUTHORIZED_REDIRECT_URL )
        return
      }

      var roleNames []string;
      for _, role := range user.Roles {
        if role.AppId != app.ID {
          continue
        }
        roleNames = append( roleNames, role.Name )
      }

      accessGranted := false
      for _, accessPolicy := range app.AccessPolicies {
        accessGranted = accessPolicy.Check( method, uriInAp, roleNames )
        if accessGranted {
          break
        }
      }

      if accessGranted {
        c.Status(http.StatusOK)
        return
      }
    }
  }

  c.Redirect( http.StatusTemporaryRedirect, forwardedProto+"://"+forwardedHost+globals.UNAUTHORIZED_REDIRECT_URL )

}
