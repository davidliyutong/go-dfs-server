# Login design


## Login API

```shell
jwt=$(curl -X POST  localhost:27903/auth/login -H 'Content-Type: application/json' -d '{"accessKey":"12345678","secretKey":"xxxxxxxx"}' | jq -r ".token")
curl -XGET -H "Content-Type: application/json" -H "Authorization: Bearer ${jwt}"  http://127.0.0.1:27903/heartbeat/
curl -XGET http://127.0.0.1:27903/ping/
```
