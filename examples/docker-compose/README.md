```console
$ docker-compose up -d
```

```console
$ docker-compose ps
         Name                       Command               State           Ports
---------------------------------------------------------------------------------------
dockercompose_redis_1    docker-entrypoint.sh redis ...   Up      6379/tcp
dockercompose_server_1   sanaa server --binding-por ...   Up      0.0.0.0:32771->80/tcp
dockercompose_worker_1   sanaa worker --concurrency ...   Up
```
