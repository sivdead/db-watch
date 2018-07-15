# db-watch
### intro
  db-watch is a simple go util to watch database connection status  and kill the long time running connections.
### installation
  `go get github.com/sivdead/db-watch`
  since this is a little util, you need build it by youself if you want:)
### usage
  #### 1. you need to create a config  json file like below:
  ```json
    [
      {
        "name": "product_test", //name of the datasource
        "addr": "127.0.0.1:3306", //url
        "username": "root", //username
        "password": "123456", //password
        "schema": "schema", //database you want to watch
        "command": ["Query"], //which command you want to watch("query","sleep",...etc)
        "timeout": 3, //timeout time by seconds that you want to log or kill
        "action": "kill",//what kind of action you want to take,it has two option "log" and "kill"
        "cron": "@every 5s"//the scheduler time, you can use cron pattern
      },
      {
        "name": "customer_test",
        "addr": "127.0.0.1:3307",
        "username": "root",
        "password": "Vvvv",
        "schema": "schema1",
        "command": ["Query"],
        "timeout": 3,
        "action": "kill",
        "cron": "@every 5s"
      }
    ]
```
#### 2.start the binary use -c or -config-file to specify the config file
