# pull official base image
FROM node:alpine

# set working directory
WORKDIR /app

# add `/app/node_modules/.bin` to $PATH
ENV PATH /app/node_modules/.bin:$PATH
RUN npm install -g npm@8.18.0
# install app dependencies
COPY ./package.json ./
COPY ./package-lock.json ./
RUN npm install react-scripts@3.4.1 -g --silent
RUN npm install -g serve --silent
RUN npm install --legacy-peer-deps --silent

# add app (e.g. public/index.html) to /app
COPY ./ ./

# start app
RUN npm run build

CMD ["serve", "-s", "build"]
