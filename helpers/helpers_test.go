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

package helpers_test

import (
  "github.com/schulterklopfer/cyphernode_fauth/globals"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "os"
  "testing"
)

type testStruct struct {
  Aint int `json:"aint"`
  Bint32 int32 `json:"bint32"`
  Cint64 int64 `json:"cint64"`
  Dstring string `json:"dstring"`
  Ebool bool `json:"ebool"`
  Ffloat32 float32 `json:"ffloat32"`
  Gfloat64 float64 `json:"gfloat64"`
}

func TestSetByJsonTag(t *testing.T) {

  target := testStruct{ 1,1,1,"foo", false, 1.0, 1.0 }

  newValues := map[string]interface{}{
    "aint": 2,
    "bint32": int32(3),
    "cint64": int64(4),
    "dstring": "bar",
    "ebool": true,
    "ffloat32": float32(2.0),
    "gfloat64": float64(3.0),
  }

  helpers.SetByJsonTag(  &target, &newValues )

  if target.Aint != 2 ||
      target.Bint32 != 3 ||
      target.Cint64 != 4 ||
      target.Dstring != "bar" ||
      target.Ebool != true ||
      target.Ffloat32 != 2.0 ||
      target.Gfloat64 != 3.0 {
    t.Error( "Set value failed")
  }

}

func TestAbsoluteURL( t *testing.T ) {

  _ = os.Setenv( globals.BASE_URL_EXTERNAL_ENV_KEY, "http://www.foo.com")

  a := helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  _ = os.Setenv( globals.BASE_URL_EXTERNAL_ENV_KEY, "http://www.foo.com/")

  a = helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  _ = os.Setenv( globals.BASE_URL_EXTERNAL_ENV_KEY, "http://www.foo.com//")

  a = helpers.AbsoluteURL( "bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "/bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

  a = helpers.AbsoluteURL( "//bar" )

  if a != "http://www.foo.com/bar" {
    t.Errorf( "%s should be %s", a, "http://www.foo.com/bar" )
  }

}