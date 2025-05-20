# Running the application

1. Start Docker Desktop and open a terminal in the main project folder: token-transfer-api.
2. Type the following command to build and start containers:

```
docker compose -f 'docker-compose.yml' up -d --build
``` 

3. Wait until Docker finishes setup. 
4. In order to open GraphQL Playground open `http://localhost:8080` in your browser. Now you can run your mutation in the GraphQL Playground code editor, for example:

```
mutation {
  transfer(
    from_address: "0x0000000000000000000000000000000000000000", 
    to_address: "0x1111111111111111111111111111111111111111", 
    amount: 1000
  )
}
```

Output given for this example:

```
{
  "data": {
    "transfer": 999000
  }
}
```

5. In order to run tests, use the following command in the opened terminal (in the same directory):

```
docker-compose exec app go test -v ./tests
```