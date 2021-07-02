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
  "github.com/dgrijalva/jwt-go"
  "github.com/gin-gonic/gin"
  "github.com/schulterklopfer/cyphernode_fauth/cyphernodeKeys"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "net/http"
  "strings"
)

func ForwardGatekeeperAuth(c *gin.Context) {

  //secret := []byte("my_secret_key")

  uriInAp := c.Request.Header.Get("x-forwarded-uri")

  if uriInAp == "/" || uriInAp == "" {
    c.Status(http.StatusUnauthorized)
    return
  }

  action := strings.Split( strings.TrimPrefix(uriInAp,"/"), "/" )[0]

  if action == "" {
    c.Status(http.StatusUnauthorized)
    return
  }

  tokenString := helpers.TokenFromBearerAuthHeader( c.Request.Header.Get("authorization") )
  token, _ := jwt.Parse(tokenString, nil)

  claims, ok := token.Claims.(jwt.MapClaims)

  if !ok {
    c.Status(http.StatusUnauthorized)
    return
  }

  if value, exists := claims["id"]; exists {
    keyLabel, ok := value.(string)
    if !ok {
      c.Status(http.StatusUnauthorized)
      return
    }
    tokenParts := strings.Split( token.Raw, "." )
    // custom legacy jwt signing stuff...
    // TODO: use standard jwt tokens
    if !cyphernodeKeys.Instance().CheckSignature(keyLabel, strings.Join( tokenParts[0:2], "." ), tokenParts[2] ) {
      c.Status(http.StatusUnauthorized)
      return
    }

    if cyphernodeKeys.Instance().ActionAllowed( keyLabel, action ) {
      c.Status(http.StatusOK)
      return
    }
  }

  c.Status(http.StatusUnauthorized)

}
