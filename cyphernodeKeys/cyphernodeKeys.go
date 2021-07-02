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

package cyphernodeKeys

import (
  "bufio"
  "bytes"
  "crypto/hmac"
  "crypto/sha256"
  "encoding/base64"
  "encoding/hex"
  "fmt"
  "github.com/pkg/errors"
  "github.com/schulterklopfer/cyphernode_fauth/helpers"
  "github.com/schulterklopfer/cyphernode_fauth/logwrapper"
  "os"
  "strings"
  "sync"
  "time"
)

type CyphernodeKeys struct {
  KeysConfigFilePath    string
  ActionsConfigFilePath string
  LastKeysUpdate        time.Time
  LastActionsUpdate     time.Time

  // label -> key
  keys                      map[string]string

  // label -> groups
  groups                    map[string][]string

  // action -> group
  actions                   map[string]string

  lastKeysConfigFileInfo    os.FileInfo
  lastActionsConfigFileInfo os.FileInfo

  loadKeysMutex    sync.Mutex
  loadActionsMutex sync.Mutex
}


var instance *CyphernodeKeys
var once sync.Once

func initOnce( keysConfigFilePath string, actionsConfigFilePath string ) error {
  var initOnceErr error
  once.Do(func() {
    keysConfigFile, err := os.Open( keysConfigFilePath )
    if err != nil {
      initOnceErr = err
      return
    }
    actionsConfigFile, err := os.Open( actionsConfigFilePath )
    if err != nil {
      initOnceErr = err
      return
    }
    instance = &CyphernodeKeys{
      KeysConfigFilePath: keysConfigFilePath,
      ActionsConfigFilePath: actionsConfigFilePath,
    }
    err = instance.parseKeysConfigFile(keysConfigFile)
    if err != nil {
      initOnceErr = err
      return
    }
    err = instance.parseActionsConfigFile(actionsConfigFile)
    if err != nil {
      initOnceErr = err
      return
    }
    helpers.SetInterval(instance.checkConfigFilesChange, 1000, false)
  })
  return initOnceErr
}

func Init( keysConfigFilePath string, actionsConfigFilePath string ) error {
  if instance == nil {
    err := initOnce( keysConfigFilePath, actionsConfigFilePath )
    if err != nil {
      return err
    }
  }
  return nil
}

func Instance() *CyphernodeKeys {
  return instance
}


/* legacy: parse strange key file format
kapi_id="001";kapi_key="a27f9e73fdde6a5005879c273c9aea5e8d917eec77bbdfd73272c0af9b4c6b7a";kapi_groups="watcher";eval ugroups_${kapi_id}=${kapi_groups};eval ukey_${kapi_id}=${kapi_key}
*/

func (cyphernodeKeys *CyphernodeKeys) parseKeysConfigFile(file *os.File) error {
  cyphernodeKeys.loadKeysMutex.Lock()
  defer cyphernodeKeys.loadKeysMutex.Unlock()
  cyphernodeKeys.keys = make(map[string]string)
  cyphernodeKeys.groups = make(map[string][]string)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := []byte(scanner.Text())
    fieldsKV :=bytes.Split( bytes.Trim(line, " "), []byte(";") )

    if len(fieldsKV) < 3 {
      // Something like an empty line
      continue
    }

    // only first 3 kv pairs are relevant
    var keyLabel string
    var keyHex string
    var keyGroups []string
    for fkv := 0; fkv<3; fkv++ {
      kv := bytes.Split( bytes.Trim(fieldsKV[fkv], " "), []byte("=") )

      switch string(kv[0]) {
      case "kapi_id":
        keyLabel = string(bytes.Trim(kv[1],"\""))
        break
      case "kapi_key":
        keyHex = string(bytes.Trim(kv[1],"\""))
        break
      case "kapi_groups":
        keyGroups = strings.Split(string(bytes.Trim(kv[1],"\"")),",")
        for i:=0; i< len(keyGroups); i++ {
          keyGroups[i] = strings.Trim( keyGroups[i], " ")
        }
        break
      }

    }
    if keyLabel != "" {
      if keyHex != "" {
        (*cyphernodeKeys).keys[keyLabel] = keyHex
      }
      if len(keyGroups) > 0 {
        (*cyphernodeKeys).groups[keyLabel] = keyGroups
      }
    }
  }
  return scanner.Err()
}

func (cyphernodeKeys *CyphernodeKeys) parseActionsConfigFile(file *os.File) error {
  cyphernodeKeys.loadActionsMutex.Lock()
  defer cyphernodeKeys.loadActionsMutex.Unlock()
  cyphernodeKeys.actions = make(map[string]string)
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    line := scanner.Text()

    if strings.HasPrefix(line, "#") {
      continue
    }

    kv := strings.Split( strings.Trim(line, " "), "=")

    if len( kv ) != 2 {
      continue
    }

    cyphernodeKeys.actions[strings.TrimPrefix( kv[0],"action_" )]=kv[1]

  }
  return scanner.Err()
}

func (cyphernodeKeys *CyphernodeKeys) BearerFromKey( keyLabel string ) (string, error) {
  cyphernodeKeys.loadKeysMutex.Lock()
  defer cyphernodeKeys.loadKeysMutex.Unlock()
  if keyHex, ok := (*cyphernodeKeys).keys[keyLabel]; ok {
    header := "{\"alg\":\"HS256\",\"typ\":\"JWT\"}"
    payload := fmt.Sprintf("{\"id\":\"%s\",\"exp\":%d}", keyLabel, time.Now().Unix()+10 )

    h64 := base64.StdEncoding.EncodeToString( []byte(header) )
    p64 := base64.StdEncoding.EncodeToString( []byte(payload) )
    toSign := h64+"."+p64
    h := hmac.New( sha256.New, []byte(keyHex) )
    h.Write([]byte(toSign))
    sha := hex.EncodeToString(h.Sum(nil))
    return "Bearer "+toSign+"."+sha, nil
  }
  return "", errors.New("No such key with label "+keyLabel )
}

// TOOO: we should handle all keys as bytes from hex string... this is strange
func (cyphernodeKeys *CyphernodeKeys) KeyForLabel( keyLabel string ) string {
  if key, exists := cyphernodeKeys.keys[keyLabel]; exists {
    return key
  }
  return ""
}

func (cyphernodeKeys *CyphernodeKeys) CheckSignature( keyLabel string, signed string, expected string ) bool {

  if keyHex, exists := cyphernodeKeys.keys[keyLabel]; exists {
    h := hmac.New( sha256.New, []byte(keyHex) )
    h.Write([]byte(signed))
    return hex.EncodeToString(h.Sum(nil)) == expected
  }
  return false
}

func (cyphernodeKeys *CyphernodeKeys) ActionAllowed( keyLabel string, action string ) bool {
  cyphernodeKeys.loadKeysMutex.Lock()
  defer cyphernodeKeys.loadKeysMutex.Unlock()
  cyphernodeKeys.loadActionsMutex.Lock()
  defer cyphernodeKeys.loadActionsMutex.Unlock()

  if group, exists0 := cyphernodeKeys.actions[action]; exists0 {
    // we found a group for this action
    if groups, exists1 := cyphernodeKeys.groups[keyLabel]; exists1 {
      // we found groups for the key label
      if helpers.SliceIndex( len(groups), func(i int) bool {
        return groups[i] == group
      }) != -1 {
        // group of action is in groups of key. all is good
        return true
      }
    }
  }
  return false
}

func (cyphernodeKeys *CyphernodeKeys) checkConfigFilesChange() {
  cyphernodeKeys.checkKeysConfigFileChange()
  cyphernodeKeys.checkActionsConfigFileChange()
}

func (cyphernodeKeys *CyphernodeKeys) checkKeysConfigFileChange() {
  fileInfo, err := os.Stat( cyphernodeKeys.KeysConfigFilePath )
  if err != nil {
    logwrapper.Logger().Error( err.Error() )
    return
  }
  if cyphernodeKeys.lastKeysConfigFileInfo != nil && (
      cyphernodeKeys.lastKeysConfigFileInfo.Size() != fileInfo.Size() ||
          cyphernodeKeys.lastKeysConfigFileInfo.ModTime().Before( fileInfo.ModTime() ) ) {
    file, err := os.Open( cyphernodeKeys.KeysConfigFilePath )
    if err != nil {
      logwrapper.Logger().Error( err.Error() )
    }
    err = cyphernodeKeys.parseKeysConfigFile( file )
    if err != nil {
      logwrapper.Logger().Error( err.Error() )
    }
    cyphernodeKeys.LastKeysUpdate = time.Now()
    cyphernodeKeys.lastKeysConfigFileInfo = fileInfo
  }
}

func (cyphernodeKeys *CyphernodeKeys) checkActionsConfigFileChange() {
  fileInfo, err := os.Stat( cyphernodeKeys.ActionsConfigFilePath )
  if err != nil {
    logwrapper.Logger().Error( err.Error() )
    return
  }
  if cyphernodeKeys.lastActionsConfigFileInfo != nil && (
      cyphernodeKeys.lastActionsConfigFileInfo.Size() != fileInfo.Size() ||
          cyphernodeKeys.lastActionsConfigFileInfo.ModTime().Before( fileInfo.ModTime() ) ) {
    file, err := os.Open( cyphernodeKeys.ActionsConfigFilePath )
    if err != nil {
      logwrapper.Logger().Error( err.Error() )
    }
    err = cyphernodeKeys.parseKeysConfigFile( file )
    if err != nil {
      logwrapper.Logger().Error( err.Error() )
    }
    cyphernodeKeys.LastActionsUpdate = time.Now()
    cyphernodeKeys.lastActionsConfigFileInfo = fileInfo
  }
}