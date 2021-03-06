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

package globals

import "errors"

const VERSION = "v0.1.0"

/** env keys **/
const BASE_URL_EXTERNAL_ENV_KEY string = "BASE_URL_EXTERNAL"
const BASE_URL_INTERNAL_ENV_KEY string = "BASE_URL_INTERNAL"
const OIDC_SSO_COOKIE_DOMAIN_ENV_KEY string = "OIDC_SSO_COOKIE_DOMAIN"
const CNA_ADMIN_LOGIN_ENV_KEY string = "CNA_ADMIN_LOGIN"
const CNA_ADMIN_PASSWORD_ENV_KEY string = "CNA_ADMIN_PASSWORD"
const CNA_ADMIN_NAME_ENV_KEY string = "CNA_ADMIN_NAME"
const CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY string = "CNA_ADMIN_EMAIL_ADDRESS"
const CNA_ADMIN_DATABASE_DSN_ENV_KEY string = "CNA_ADMIN_DATABASE_DSN"
const CNA_ADMIN_CONFIG7Z_FILE_ENV_KEY string = "CNA_ADMIN_CONFIG7Z_FILE"
const CNA_ADMIN_CLIENT7Z_FILE_ENV_KEY string = "CNA_ADMIN_CLIENT7Z_FILE"

const CNA_SESSION_COOKIE_NAME_ENV_KEY = "CNA_SESSION_COOKIE_NAME"
const CNA_COOKIE_SECRET_ENV_KEY = "CNA_COOKIE_SECRET"
const CNA_STATIC_FILE_DIR_ENV_KEY = "CNA_STATIC_FILE_DIR"
const CYPHERAPPS_INSTALL_DIR_ENV_KEY = "CYPHERAPPS_INSTALL_DIR"
const GATEKEEPER_HOST_ENV_KEY = "GATEKEEPER_HOST"
const GATEKEEPER_PORT_ENV_KEY = "GATEKEEPER_PORT"
const KEYS_FILE_ENV_KEY = "CYPHERNODE_KEYS_FILE"
const ACTIONS_FILE_ENV_KEY = "CYPHERNODE_ACTIONS_FILE"
const CERT_FILE_ENV_KEY = "CYPHERNODE_CERT_FILE"


const BASE_ADMIN_MOUNTPOINT string = "admin"

/** urls and endpoints **/
const FORWARD_AUTH_ENDPOINTS_AUTH = "/public"
const PROXY_GATEKEEPER_ENDPOINTS_AUTH = "/gatekeeper"

const UNAUTHORIZED_REDIRECT_URL string = "/admin"

const CYPHERAPPS_REPO string = "git://github.com/SatoshiPortal/cypherapps.git"


/** useful vars **/
var ENDPOINTS_PUBLIC_PATTERNS = [...]string{".*/+favicon.ico$"}

/** defaults **/

var DEFAULTS = map[string]string{
  BASE_URL_EXTERNAL_ENV_KEY:       "http://www.cna.localhost:3030",
  BASE_URL_INTERNAL_ENV_KEY:       "http://www.cna.localhost:3031",
  OIDC_SSO_COOKIE_DOMAIN_ENV_KEY:  "www.cna.localhost",
  CNA_COOKIE_SECRET_ENV_KEY:       "thisIsTheDefaultSecret",
  CNA_ADMIN_LOGIN_ENV_KEY:         "admin",
  CNA_ADMIN_PASSWORD_ENV_KEY:      "admin",
  CNA_ADMIN_NAME_ENV_KEY:          "admin",
  CNA_ADMIN_EMAIL_ADDRESS_ENV_KEY: "admin@admin.com",
  CNA_ADMIN_DATABASE_DSN_ENV_KEY:  "host=db port=5432 user=cnadmin password=cnadmin dbname=cnadmin sslmode=disable",
  CNA_ADMIN_CONFIG7Z_FILE_ENV_KEY: "/data/config.7z",
  CNA_ADMIN_CLIENT7Z_FILE_ENV_KEY: "/data/client.7z",
  CNA_STATIC_FILE_DIR_ENV_KEY:     "/ui",
  CYPHERAPPS_INSTALL_DIR_ENV_KEY:  "/apps",
  KEYS_FILE_ENV_KEY:               "/keys.properties",
  ACTIONS_FILE_ENV_KEY:            "/api.properties",
  CERT_FILE_ENV_KEY:               "/cert.pem",
  GATEKEEPER_HOST_ENV_KEY:         "gatekeeper",
  GATEKEEPER_PORT_ENV_KEY:         "2009",
  CNA_SESSION_COOKIE_NAME_ENV_KEY: "io.cyphernode.session",
}


var ErrDuplicateUser = errors.New("user already exists")
var ErrUserHasUnknownRole = errors.New("user has unknown role")
var ErrNoSuchUser = errors.New( "no such user" )
var ErrNoSuchRole = errors.New( "no such role" )
var ErrCannotAddExistingRole = errors.New( "cannot add existing role to app" )
var ErrUserAlreadyHasRole = errors.New( "user already has role" )
var ErrNoSuchApp = errors.New( "no such app" )
var ErrMigrationFailed = errors.New( "migration failed" )
var ErrDatabaseNotInitialised = errors.New( "database not initialised")
var ErrActionForbidden = errors.New( "action forbidden" )