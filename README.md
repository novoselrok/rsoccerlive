# r/soccer.live

Live highlight updates from the https://reddit.com/r/soccer subreddit. For each higlight it also fetches any possible mirrors in the comments.

It uses Go for the backend API, Vue for the frontend, and WebSockets to make it live.

# Development

* Prerequsites
  * Docker and Docker Compose
  * PostgreSQL
  * golang-migrate tool (https://github.com/golang-migrate/migrate)

* Populate .env and client/.env
  * Necessary keys are listed in .env.sample and client/.env.sample
* Create a database and save the url into `POSTGRESQL_URL`
* Run migrations: `migrate -database ${POSTGRESQL_URL} -path src/migrations up`
* Build the images: `docker-compose build`
* Run the client and api services: `docker-compose up -d`
* Watch the logs: `docker-compose logs -f`
