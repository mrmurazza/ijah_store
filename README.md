# Quickstart Guide to Set up Dev Workspace

1. Git clone https://github.com/mrmurazza/ijah_store.git
2. Make sure you have Go installed and set up Go workspace
3. Install sqlite3 driver library for database driver and Gin-Gonic for web framework using `go get` like this:

   `go get github.com/mattn/go-sqlite3`
   
    `go get -u github.com/gin-gonic/gin`
    
4. Go to your project folder using `cd` and run your project using:

    `$ go run main.go`
    
5. For deployment, run this command to build binary file

    `$ go build main.go`