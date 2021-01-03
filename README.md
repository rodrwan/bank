# bank

Sample repo for CQRS and Event sourcing.

This simulate a bank.

# Service list

| Name | Description | Port |
| --- | --- | --- |
| Accounts | Service to handle bank accounts. | 8010-8011 |
| Charges | Service to handle charges made to an account. | |
| Deposits | Service to handle the incomming money. | |
| Denormalizer | Service to hold a collection of views, this is intended to read data. | - |
| Movements | Service to record every transaction made by an account. | |
| Server | GraphQL server to handle user request to the system. | |
| Session | Service to handle user session, basically a JWT wrapper. | 8090 |
| Users | Service to handle users information. | 8020-8021 | 
| Withdrawals | Service to handle withdrawals from accounts. | |


To run:

```sh
$ make dev
```