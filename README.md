# sidecar-security

Security sidecar forwards authenticated http/REST request to origin. The sidecar is a proxy forwarding all requests to origin on port 80.

### Deploy sample serivce POD (sidecare + origin service app) in Kubernetes 

This example deploys the sidecare with origin as https://httpbin.org. 

```
kubectl apply -f deploy/deployment.yaml
```

### Test the Sidecar

Deploy the test curl client in cluster

```
kubectl run curl --image=radial/busyboxplus:curl -i --tty 
```

Note the cluster IP of deployed Sidecar service

```
kubectl get svc origin-httpbin

```

Jump in the curl pod to test 

```
kubectl exec curl-69c656fd45-wtfx6 -it -- sh
```

Test sidecare with IP and valid credentials

```
curl -v -H "user: foo" http://10.102.69.227:8000/anything
> GET /anything HTTP/1.1
> User-Agent: curl/7.35.0
> Host: 10.102.69.227:8000
> Accept: */*
> user: foo
> 
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 422
< Content-Type: application/json
< Date: Sun, 01 Mar 2020 07:43:25 GMT
< Server: gunicorn/19.9.0
< 
{
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept": "*/*", 
    "Accept-Encoding": "gzip", 
    "Host": "10.102.69.227:8000", 
    "User": "foo", 
    "User-Agent": "curl/7.35.0", 
    "X-Forwarded-Host": "10.102.69.227:8000", 
    "X-Origin-Host": "localhost:80"
  }, 
  "json": null, 
  "method": "GET", 
  "origin": "172.17.0.6", 
  "url": "http://10.102.69.227:8000/anything"
}

curl -v -H "user: foo" http://origin-httpbin:8000/anything
> GET /anything HTTP/1.1
> User-Agent: curl/7.35.0
> Host: origin-httpbin:8000
> Accept: */*
> user: foo
> 
< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Origin: *
< Content-Length: 425
< Content-Type: application/json
< Date: Sun, 01 Mar 2020 07:43:40 GMT
< Server: gunicorn/19.9.0
< 
{
  "args": {}, 
  "data": "", 
  "files": {}, 
  "form": {}, 
  "headers": {
    "Accept": "*/*", 
    "Accept-Encoding": "gzip", 
    "Host": "origin-httpbin:8000", 
    "User": "foo", 
    "User-Agent": "curl/7.35.0", 
    "X-Forwarded-Host": "origin-httpbin:8000", 
    "X-Origin-Host": "localhost:80"
  }, 
  "json": null, 
  "method": "GET", 
  "origin": "172.17.0.6", 
  "url": "http://origin-httpbin:8000/anything"
}

```

Test with Invalid creds

```

curl -v -H "user: foobar" http://10.102.69.227:8000/anything
> GET /anything HTTP/1.1
> User-Agent: curl/7.35.0
> Host: 10.102.69.227:8000
> Accept: */*
> user: foobar
> 
< HTTP/1.1 403 Forbidden
< Date: Sun, 01 Mar 2020 07:48:20 GMT
< Content-Length: 0


curl -v -H "user: foobar" http://origin-httpbin:8000/anything
> GET /anything HTTP/1.1
> User-Agent: curl/7.35.0
> Host: origin-httpbin:8000
> Accept: */*
> user: foobar
> 
< HTTP/1.1 403 Forbidden
< Date: Sun, 01 Mar 2020 07:45:01 GMT
< Content-Length: 0
< 

```

### Building from source 

1. Clone repo
2. change TAG_HUB to your container repp. 
3. run ```make dpush```
