# WINK.gg

Wink is a video game streaming website similar to Twitch built with Golang and VanillaJS.

# Docker

First, clone the repository from github:

`https://github.com/comebacknader/wink.git`

Then change into the directory: 

`cd wink`

Wink is designed to be run with Docker Compose. Run the following command:

`docker-compose up`

Then go to localhost:8080 and the application should load up.

To run the application without docker compose, first create a network:

`docker network create wink-net`

Build the web application from the /wink directory with the following command: 

`docker build -t comebacknader/winkgg .`

Change into the /wink-psql directory: 

`cd wink-psql`

Build the database image with the following command: 

`docker build -t comebacknader/wink-psql .`

To run the postgres container, use the following command:

`docker container run -d  --network wink-net -v wink-psql:/var/lib/postgresql/data --name wink-db-container comebacknader/wink-psql`

To run the golang container, use the following command:

`docker container run -d -p 8080:8080 --name wink-container-name --network wink-net -v "$(pwd)"/:/go/src/github.com/comebacknader/wink comebacknader/winkgg`

To watch the CSS files with SASS use the following:

`docker exec -it wink-container-name bash`

Then issue the following command:

`sass --watch assets/css/main.scss assets/css/main.css` 

