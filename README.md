build:
`docker build https://github.com/size12/planning-poker.git -t planning-poker`

run:
`docker run -d -ti -p 8080:8080 -e RUN_ADDRESS=:8080 -e planning-poker`

example:
`http://195.2.81.167:8080/rooms/create`
