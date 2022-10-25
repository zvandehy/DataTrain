# CLOVER ANALYTICS
NBA Player Sportsbook Analyzer

This tool aims to allow anyone to review available Player Props on DFS sites, create their own prediction models, and analyze the performance of these models. Sportsbooks have all of the edge: they have dedicated line makers, computing power, favorable betting lines, knowledge of the amount of money being placed on each bet, and they only show customers a limited set of relevant data to their player props. This tool attempts tp take some of that edge back by creating an easy way for users to get all of the data that they want about a player and their past performances. Users can create custom prediction models by simply selecting what data they believe is relevant to a prediction, and weights associated with the importannce of that data. The tool uses this custom model to automatically collect, analyze, and calculate a prediction. Then the user can review the prediction and how that model performs compared to actual results and sportsbook lines.

# React
Front End of Web App

http://www.clover-analytics.com

Requires React and local libraries
- ```cd react/nba-analyzer```
- ```npm i -g npm```
- ```npm i -g react react-scripts typescript ajv``` 
- ```npm i```
- ```npm start```

# NBA_GraphQL
Golang implementation of a GraphQL backend to aggregate data from DB or other API sources.

https://clover-backend.fly.dev/nba

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

# Workflow Process
1. On JIRA, locate the ticket number that you are working on. E.g. “NBA-35”
2. Open a terminal to the DataTrain/ repository
3. Make sure you are on the main branch and that the branch is up to date
    1. ```git switch main```
    2. ```git pull```
4. Create a new branch for this ticket
    1. ```git checkout -b NBA-35```
5. Make changes to the code for the ticket
6. Stage these changes by either:
    1. ```git add <filename>``` or ```git add .```
    2. Clicking the “+” button next to the changed files on the source control tab
7. Commit these changes by either:
    1. ```git commit -m “my commit message”```
    2. Clicking the checkmark in the source control tab and typing in a commit message
8. Push changes to the repository
    1. ```git push```
    2. NOTE: the first time you run ```git push``` on this branch, it may ask you to set the upstream branch because it does not yet exist on the GitHub repository. For example: ```git push --set-upstream origin NBA-35```
9. Repeat steps 5-8 as needed
10. Go to https://github.com/zvandehy/DataTrain
    1. Also, notice that on JIRA your branch is automatically linked to your GitHub branch. This means that when looking at anyone’s ticket in PR, you can go directly from their ticket to the their branch / PR.
11. There should be a notification that your branch was just created and give you the option to open a pull request
    1. If this option is not showing, you can initiate a Pull Request from the “Branches” page https://github.com/zvandehy/DataTrain/branches
12. Scroll down to compare your changes with the master branch & verify that everything is right
13. Add any comments/description for context
14. Create Pull Request
15. Get a teammate to review and approve the code
16. Merge Pull Request
