# Simply Do

A ToDo List Web Application for anyone!
Visit https://simply-do.herokuapp.com to see it in action.

Simply Do uses `Vue` for the frontend, `Go` for the backend, `Firebase` for data storage, and `Heroku` for deployment.

### Developing

1.  Start the server locally:

    ```
    go run main.go server
    ```

2.  Visit `localhost:8080`, the webpage is served by files from `public/`.

> Note: changes to frontend won't require a server restart, but changes to the backend will need a server restart.

Make commands:

```shell
$ make
build-docker                   Builds the docker image
build-linux-binary             Builds the linux binary
clean                          Cleans up the built binaries
deploy                         Deploys to Heroku. Requires to be logged in on Heroku Registry
help                           List targets & descriptions
run                            Starts the server locally on port 8080 or $PORT if set
```
