# pull official base image
FROM node:alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH

# install app dependencies
COPY react/nba-analyzer/package.json ./
COPY react/nba-analyzer/package-lock.json ./
RUN npm install --silent
RUN npm install react-scripts@3.4.1 -g --silent
RUN npm i -g serve --silent
RUN npm install --save chart.js react-chartjs-2

# add app
COPY react/nba-analyzer/ ./

# start app
RUN npm run build

CMD ["serve", "-s", "build"]
