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
    
    
# Project Structure 
```
.
├── src/app
│   ├── item                         > package contains Item-related code 
│   ├── purchase                     > package contains purchase-related code
│   ├── request                      > package contains request struct code
│   ├── response                     > package contains response struct code
│   ├── restock                      > package contains restock-related code
│   ├── util                         
│   └── controllers.go               > file contains input/output preparation
├── README.md
├── ijah_store.db                   > SQLite DB file
├── inventory_report_17_01_19.csv   > Example output of inventory reports in CSV
├── main.go                         > Main Go file
└── sales_report_17_01_19.csv       > Example output of sales reports in CSV 

```

# Project Component

Basically, this project is separated into several components: 
1. Controller : responsible on preparing the inputs and serving the outputs and the one to call the business logic code.
2. Handler : responsible on handling the main core of the business logic. (ex: `item_handler.go, restock_handler.go, & purchase.handler)`
3. Model : responsible on containing the struct and queries. (ex: `item.go, restock_resception.go, etc`)

