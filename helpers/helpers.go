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

package helpers

import (
  "crypto/rand"
  "encoding/base64"
  "encoding/json"
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/password"
  "golang.org/x/crypto/ripemd160"
  "io"
  "os"
  "reflect"
  "regexp"
  "strings"
  "time"
)

func SliceIndex(limit int, predicate func(i int) bool) int {
  for i := 0; i < limit; i++ {
    if predicate(i) {
      return i
    }
  }
  return -1
}

func SetInterval(someFunc func(), milliseconds int, async bool) chan bool {

  // How often to fire the passed in function
  // in milliseconds
  interval := time.Duration(milliseconds) * time.Millisecond

  // Setup the ticket and the channel to signal
  // the ending of the interval
  ticker := time.NewTicker(interval)
  clear := make(chan bool)

  // Put the selection in a go routine
  // so that the for loop is none blocking
  go func() {
    for {

      select {
      case <-ticker.C:
        if async {
          // This won't block
          go someFunc()
        } else {
          // This will block
          someFunc()
        }
      case <-clear:
        ticker.Stop()
        return
      }

    }
  }()

  // We return the channel so we can pass in
  // a value to it to clear the interval
  return clear

}

func EndpointIsPublic( endpoint string ) bool {
  for i:=0; i<len( globals.ENDPOINTS_PUBLIC_PATTERNS); i++ {
    pattern := globals.ENDPOINTS_PUBLIC_PATTERNS[i]
    matches, err := regexp.MatchString( pattern, endpoint )
    if matches && err == nil {
      return true
    }
  }
  return false
}

func RandomString(length int, encodeToString func([]byte) string ) string {
  randomBytes := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
    return ""
  }
  return strings.TrimRight( encodeToString( randomBytes), "=" )
}

func AbsoluteURL( path string ) string {
  return AbsoluteURLFromHostEnvKey( globals.BASE_URL_EXTERNAL_ENV_KEY, path )
}

func AbsoluteURLFromHostEnvKey( hostEnvKEy string, path string ) string {
  return AbsoluteURLFromHost( GetenvOrDefault( hostEnvKEy ), path )
}

func AbsoluteURLFromHost( host string, path string ) string {
  for strings.HasSuffix( host,"/") {
    // remove last character
    host = host[:len(host)-1]
  }

  for strings.HasPrefix( path,"/") {
    // remove last character
    path = path[1:len(path)]
  }

  return host+"/"+path
}


func SetByJsonTag( obj interface{}, values *map[string]interface{} ) {

  // evaluate sbjt tag actions like hashing passwords
  structType := reflect.TypeOf(obj).Elem()
  //mutableObject := reflect.ValueOf(obj).Elem()
  for jsonFieldName, jsonFieldValue := range *values {
    for i := 0; i < structType.NumField(); i++ {
      field := structType.Field(i)
      jsonTag, hasJsonTag := field.Tag.Lookup("json")
      sbjtTag, hasSbjtTag := field.Tag.Lookup("sbjt")

      if hasSbjtTag && hasJsonTag && jsonTag == jsonFieldName {
        switch sbjtTag {
        case "hashPassword":
          if reflect.TypeOf(jsonFieldValue).Kind() == reflect.String {
            hashedPassword, _ := password.HashPassword(jsonFieldValue.(string))
            (*values)[jsonFieldName] = hashedPassword
          }
          break
        }
      }
    }
  }

  jsonStringBytes, _ := json.Marshal( values )
  _ = json.Unmarshal( jsonStringBytes, obj )

}

func GetenvOrDefault( key string ) string {
  value := os.Getenv( key )
  if value == "" {
    defaultValue, _ := globals.DEFAULTS[key]
    return defaultValue
  }
  return value
}

func TrimmedRipemd160Hash( bytes []byte ) string {
  hasher := ripemd160.New()
  hasher.Write(bytes)
  return strings.Trim(base64.URLEncoding.EncodeToString(hasher.Sum(nil)), "=" )
}

func TokenFromBearerAuthHeader(authHeader string) string {
  if authHeader == "" {
    return ""
  }

  parts := strings.Split(authHeader, "Bearer")
  if len(parts) != 2 {
    return ""
  }

  token := strings.TrimSpace(parts[1])
  if len(token) < 1 {
    return ""
  }

  return token
}