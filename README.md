# go-service-example
This project is another attempt to create some good example service on Go.
The main logic you can find in the file "[transactions/service](internal/transactions/services/service.go)"

# Structure
The project is trying to follow "[Standard Go Project Layout](https://github.com/golang-standards/project-layout)".
If you want to dive deeply into the topic we strongly recommend to watch this video ин [Kat Zien](https://github.com/katzien) "[How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)"

# Setup
```
make setup
```

# Run
```
docker-compose up
```

# Examples
There are 3 scripts:
- [1_create_users.sh](examples/1_create_users.sh) creates 2 users which we need for making funds trasfer and third user (for fees) we create with migration.
- [2_create_wallets.sh](examples/2_create_wallets.sh) creates 3 USD wallets.
- [3_make_transaction.sh](examples/3_make_transaction.sh) makes funds transfer from one ise to another.