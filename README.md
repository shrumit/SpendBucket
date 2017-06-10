# Spend Bucket

*A full-stack RESTful single-page application to track shared expenses in a group of people.*

## Server

Written with the intention of leveraging only the standard Go libraries, without 3rd-party frameworks or routers.

#### Compiling

Ensure that the project directory is included in the GOPATH. Run `go install` in `server`.

## Site

Built with Vue.js, Vuex and superagent; bundled using Webpack. Uses Material Design UI components from muse-ui and vue-material.

#### Compiling

Run `npm run build` in site directory to generate a production build in the `dist` directory.

## Database

Schema is written for MySQL. Dummy data is encased in a MySQL Event that resets the dummy account to default every 10 minutes.

The choice of data store is completely decoupled from the server application and can be swapped for anything implementing the DataController interface in `controller.go`. The current MySQL controller can be found in `controllerImpl.go`.

## License

GNU AGPL v3