# Spend Bucket

## About

A full-stack RESTful single-page application to track shared expenses in a group of people.

### Server

Written with the intention of using only the standard Go libraries, without using 3rd-party frameworks or routers. Consequently, the only 3rd-party dependencies are `go-sql-driver/mysql` and `dgrijalva/jwt-go`.

**Compiling**

Ensure that the project directory is included in the GOPATH. Then run `go install` from the `server` directory.

### Site

Built with Vue.js and Vuex, and bundled using Webpack. The design (non-strictly) adheres to Material Design guidelines by using components from muse-ui and vue-material.

**Compiling**

Run `npm run build` in site directory to generate a production build. Files in `dist` may then be served using a web server such as Nginx.


### Database

The schema is written for MySQL and can be loaded from `schema.sql` in the  `sql` directory. Dummy data is provided in `dummy.sql`, encased in a MySQL Event that resets the dummy account to default every 10 minutes.

The server application is decoupled from the choice of data store. Simply implement the DataController interface found in `controller.go` . `controllerImpl.go` contains the current MySQL-targeted implementation.

## License

GNU Affero GPLv3 licensed. See LICENSE for details.