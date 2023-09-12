build:
`docker build https://github.com/size12/planning-poker.git`

run:
`docker run -d -ti -p 8080:8080 -e RUN_ADDRESS=:8080 -e BASE_URL=127.0.0.1 planning-poker`