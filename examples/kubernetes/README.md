```console
$ kubectl get configmaps,serviceaccounts,services,statefulsets,deployments,pods -l app=sanaa -n default
NAME             DATA      AGE
cm/sanaa-redis   1         34m

NAME              SECRETS   AGE
sa/sanaa-redis    1         34m
sa/sanaa-server   1         12m
sa/sanaa-worker   1         11m

NAME               CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
svc/sanaa-redis    100.65.130.128   <none>        6379/TCP   34m
svc/sanaa-server   100.71.225.132   <none>        80/TCP     12m

NAME                       DESIRED   CURRENT   AGE
statefulsets/sanaa-redis   1         1         34m

NAME                  DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
deploy/sanaa-server   1         1         1            1           12m
deploy/sanaa-worker   1         1         1            1           11m

NAME                               READY     STATUS    RESTARTS   AGE
po/sanaa-redis-0                   1/1       Running   0          34m
po/sanaa-server-68757d6fb-nv7kh    1/1       Running   0          12m
po/sanaa-worker-6b6c479dc9-sx4hq   1/1       Running   0          11m
```
