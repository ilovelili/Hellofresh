# HelloFresh Senior Backend Developer Test

## Technology
* go 1.7
* postgres 9.5
* mongodb 3.2.3

## Dependencies
* [gorilla/mux - URL router and dispatcher](https://github.com/gorilla/mux)
* [mgo - MongoDB driver](https://gopkg.in/mgo.v2)
* [pq - PostgreSQL driver](https://github.com/lib/pq)
* [ginkgo - BDD Testing Framework](https://github.com/onsi/ginkgo)
* [gomega - matcher/assertion library](https://github.com/onsi/gomega)

## Run
* `docker-compose build && docker-compose up`
* access http://localhost:8080

## Endpoints
| Name   | Method      | URL                            | Auth Needed   |
| ---    | ---         | ---                            | ---           |
| List   | `GET`       | `/recipes`                     | No            |
| Create | `POST`      | `/recipes`                     | Yes           |
| Get    | `GET`       | `/recipes/{id}`                | No            |
| Update | `PUT/PATCH` | `/recipes/{id}`                | Yes           |
| Delete | `DELETE`    | `/recipes/{id}`                | Yes           |
| Rate   | `PUT/PATCH` | `/recipes/{id}/rate/{rate}`    | Yes           |
| Search | `GET`       | `/recipes/search/{search}`     | No            |

## Database
Data Persisted to both Postgres and MongoDB (Redis is not implemented). The default Database is MongoDB. To switch database, you can:
* comment out the mongodb container and uncomment the postgres container and switch the link as well in docker-compose.yml
* update config.json under src/hellofresh folder (Or you can rename config.json.postgresexample in the same folder to config.json directly)

## DataTable
1. recipe
    * ID - Bson ObjectId(mongodb) or SERIAL(postgres)
    * Name - string
    * Prep - Date
    * Difficulty - int
    * Vegetarian - bool
2. reciperate
    * ID - Bson ObjectId(mongodb) or SERIAL(postgres)
    * RecipeID - string
    * Rate - int
    * User - string
    * Modified - Date

## Auth
Basic Auth is used to protect create, update, delete operations. The username and password is hellofresh/hellofresh

## Test
go test has been merged into Dockerfile so test will automatically run after `docker-compose up`. If you want to run test manually, move to src/hellofresh (where hellofresh_suite_test.go is) and run `ginkgo -v`

## Known issue
For some reason after mongodb container started up, the golang mongodb driver can't access 127.0.0.1:27017 programtically (but I can access 127.0.0.1:27017 by mongo shell directly). So I tried to bind mongod ip to a staitic ip in docker-compose.yml and to access the static ip instead of localhost, but still failed...
To solve this issue, I have to `docker inspect` the mongo container and check the container ip and hardcode it in config.json.. which is very ugly.. But right now I don't know what kind of issue it is and how to solve it. Maybe it is a docker + windows issue, or my vpn setting issue (I am using vpn since I am now in China), or docker mongo image issue or mongo driver issue. So, as a conclusion, if you failed to access the endpints, try to change the server ip to mongodb container ip in config.json.

## Contact
min ju <route666@live.cn>