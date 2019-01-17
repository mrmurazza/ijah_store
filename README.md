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


# Notes Regarding the Assessment:

For this "Toko Ijah" Inventory App, based on the provided spreadsheet, here are several points that 
I could point out as a base assumption on developing this app:

1. There are 4 entities: Item, Restock Order, Restock Reception, Purchase Order
2. Item is created on when a restock of a new item is ordered.
3. When a restock order is created, it will not immediately add an item's stock unless someone has received 
it (by receive, it means someone has to input a Restock Reception)
4. A restock order has Invoice Id and it has to be unique for each order. While having Invoice Id is not mandatory, 
pending reception is forbidden when invoice id is missing. It means that when invoice id is missing, 
all of the item has to be received by the time the order is created or is input.
5. A purchase order has Order Id, just as before, it is not mandatory. However, in one order id, 
there could consist of more than 2 items. That is why the POST API of purchase order creation could 
accept a list of item purchase details. 
6. When a purchase order is posted and passed the validation, it will be created and immediately reduced 
an item's stock
