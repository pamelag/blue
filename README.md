# blue
CRUD APIs for Articles and their Tags

## setup
The project is dockerized and has a docker-compose file that points to Dockerfiles for the API and DB. The Dockerfile for database is located in db directory and the Dockerfile for the API is in the root. The table structure for article and tag tables is in the pg.sql file of db directory.

1. Go version 1.11 and upwards is required to build this project. Please check your go version using the command ```go version```
2. I have used go1.11 to create the project and distributed go.mod which contains all the dependencies and can run a project outside the $GOPATH.
3. Create a directory structure of the form ```ffx/github.com/pamelag```. I have used github code layout.
4. Next cd to the directory pamelag and run the command ```git clone https://github.com/pamelag/blue.git```
5. Now the directory structure would be ```ffx/github.com/pamelag/blue```
6. To run the project docker must be installed and the project needs to be built using ```docker-compose build``` 
7. Next the command ```docker-compose up``` should be used to start the containers

## package structure
1. **content** is the domain package which contains the main entities Article, Tag, the factory methods and Repository interfaces
2. **article** package contains the Article services for AddArticle and GetArticle. The article package depends on the content package.
3. **tag** package has implementation of the GetTag service. The tag package has a dependency on the content package.
4. **server** package has all the handlers for routing, decoding requests and encoding responses. The handlers in the server package invoke the service functions in the article and tag packages.
5. **postgres** package has all the database statements and queries. It has the implementations of the ArticleRepository and TagRepository interfaces defined in the content package
6. **main** package has the main.go file which initializes all the repositories and services and handles the wiring and dependency injection of repositories into services and the services to the http server. The config.go reads all configurations for the host and port and authentication details.
7 **db** package has the pg.sql file which contains the DDL for the article and tag tables
