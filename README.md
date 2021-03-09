# gomicrofront
A project of micro front-end with Go language

## Usage
- Clone the repository
- Build with ```go build```

### CLI
```
./microfrontend init
```

Creates the folders:
- plugins 
    - Where components will be created
- public
    - Folder that will be exposed via file server
- server
    - Containing the main server file


```
./microfrontend build
```
Builds every plugin and the server

```
./microfrontend new-plugin <name>
```
Scaffolds a new plugin in plugins folder

```
./microfrontend serve
```
Starts the server in port 8080

### Server

After running ```microfrontend init```, the file created at ```./server/entry.go``` will contain 2 functions: ```Router``` and ```Entrypoint```.

### Router function
The router function receives a pointer for ServeMux, with this you can set backend routes

### Entrypoint function
The entrypoint function is executed at every request at first, if an error is returned, the request is stoped and the error is retrived to the client.


### Plugins
Every plugin can be understood as a component for the front-end and as a route for the back-end

After creating a plugin with ```microfrontend new-plugin <name>```, its folder will contain: ```template.html```, ```soul.go``` and ```body.js```.

```body.js``` contains the script for component declaration

```soul.go``` contains the route retrives the script given its request, can be used as a middleware on component retrive.

> :warning: **The project still at its very beguining, do not attempt to use in production unless you have reviewed the code**







