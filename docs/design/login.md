# Login design


## Login API

```shell
jwt=$(curl -X POST  localhost:27903/auth/login -H 'Content-Type: application/json' -d '{"accessKey":"bccaf38c966e02fc","secretKey":"n51js4034N921iEWb8t3J938zQ0S7AjM"}' | jq -r ".token")
curl -XGET -H "Content-Type: application/json" -H "Authorization: Bearer ${jwt}"  http://127.0.0.1:27903/v1/info
curl -XGET -H "Content-Type: application/json" -H "Authorization: Bearer ${jwt}"  http://127.0.0.1:27903/info
curl -XPOST -H "Content-Type: application/json" -H "Authorization: Bearer ${jwt}"  http://127.0.0.1:27903/auth/refresh
curl -XGET http://127.0.0.1:27903/ping
```

## Login CLI

```shell
go-dfs-client login dfs://127.0.0.1
go-dfs-client login dfss://127.0.0.1
go-dfs-client login dfs://127.0.0.1:27904
go-dfs-client login dfs://::27904

go-dfs-client login dfs://127.0.0.1 --accessKey=12345678 --secretKey=xxxxxxxx
go-dfs-client login dfs://127.0.0.1:27903 --accessKey=12345678 --secretKey=xxxxxxxx
```
