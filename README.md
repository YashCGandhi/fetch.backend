# fetch.backend
fetch backend task

Make sure docker daemon is running

1. cd fetch.backend
2. docker build . -t fetch.backend:latest
3. docker run -p 9090:9090 fetchtest:latest

Then test the api on localhost:9090
