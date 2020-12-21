# go-service-example
The main logic you can find in [transactions/service](internal/transactions/service.go)

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