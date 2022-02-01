# DataTrain
NBA player props analyzer
Current implementations are more of a prototype

# React
Front End of Web App

Requires React and local libraries
- ```cd react/nba-analyzer```
- ```npm i```
- ```npm start```

# NBA_GraphQL
Golang implementation of a GraphQL backend to aggregate data from DB or other API sources

Requires golang to be installed
- ```cd nba_graphql```
- ```go run server.go```


- use `go run runner/runner.go` to use the gql generator after updating the graphql schema

# Python API
Python (Notebook) with scripts for extracting data from NBA.com/stats and storing it in MongoDB

Requires pymongo-srv, along with other basic libraries

# DEMO
To run the program locally make sure that
- data is up to date by running the python api
- the graphql backend is running
- the react app is running
