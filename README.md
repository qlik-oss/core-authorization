# [EXPERIMENTAL] Core Authorization

This repo consists of two parts:
1. Three Qlik Associative Engines set up with different access rules. 
2. A set of test cases, both positive an negative, that check the access to the engines for different JWTs.

## Running

** Note: You need Go installed and the repo located correctly on the GO path <gopath>/src/github.com/qlik-oss/core-authorization


Start the engines using docker-compose

```bash
docker-compose up -d
```

Run the test cases using

```bash
go test -v ./access_tests
```


