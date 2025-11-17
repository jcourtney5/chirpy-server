# Chirpy Server

Simple twitter style clone project to practice Golang

---

### Requirements
* Golang version 1.25.3 or higher
* Postgres version 16.10 or higher

---

### How To Build And Run
* Download code from github
* Run "go install ." which will build and copy the chirpy-server binary to your go bin folder
* Make sure you have a postgres server running and know the connection string
* Create a .env file with the following values
    ```shell
    DB_URL="<your postgres connection string>"
    PLATFORM="dev, to support reset command, prod to disable it"
    JWT_SECRET="<jwt secret key, ex: openssl rand -base64 32>"
    POLKA_KEY="<key for webhook sample, ex: openssl rand -hex 32>"
    ```
* Run goose migrations to create the tables in postgres
    ```shell
    goose postgres "<postgres connection string>" up
    ```
* Run "chirpy-server" to start the server

---

### List of APIs
* GET /api/healthz
* POST /api/validate_chirp
* POST /api/login
* POST /api/refresh
* POST /api/revoke
* POST /api/users
* PUT /api/users
* POST /api/chirps
* GET /api/chirps
* GET /api/chirps/{chirpID}
* DELETE /api/chirps/{chirpID}
* POST /api/polka/webhooks
* GET /admin/metrics
* POST /admin/reset