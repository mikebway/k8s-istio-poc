# Running Istio on Minikube locally

## Objective

Demonstrate how Istio can be deployed on top of Minikube to:

1. Manage path routing for a trivial two service web application
2. Demonstrate the configuration of the Istio ingress gateway to utilize an [Envoy external authorization filter](https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/ext_authz_filter)
   to evaluate requests and modify both request and response headers. 

Login and logout is crudely illustrated by the creation and removal of a session cookie containing the username.
This is clearly not a secure pattern; do not use it for anything more than demonstration purposes.

**NOTE:** While it would be possible to perform much if not all of the Kubernetes and Istio configuration with just one
or two manifest YAML files, this project takes an incremental approach of building up the cluster one Lego brick at a 
time with smaller, single purpose YAML files and single `kubectl`/`istioctl` commands rather than scripts. 

If it looks like this project was assembled brick-by-brick, that's because it's true; learning how
to configure Kubernetes and Istio was achieved one small step at a time.

## Installation and configuration

1. [Minikube installation and basic configuration](Install.md)
2. [Build and deploy the authtest service](svc-authtest)
3. [Build and deploy the login service](svc-login.md)
4. [Configure Istio ingress](istio.md)
5. [Configure Istio service mesh](mesh.md)
6. [Configure Istio visualization](visualize.md)<p>
   **Optional:**
7. [Build and deploy the ingress authorization filter](svc-extauth.md)
8. [Configure the ingress authorization policy](authz-policy.md)

## Testing without the authorization filter

After completing installation and configuration through step 6 you should be able to:

1. Point a browser to http://localhost and see a simple text response that looks something like this, echoing the `/`
   path of your request and a count 1 times that the [authtest](../authtest) service has responded to a request.
   ```text
   Path:		"/"
   Count:		1

   HEADERS (25)
   =======

   Accept: [text/html,application/xhtml+xml,application/xml;q=0.9,image...
   Accept-Encoding: [gzip, deflate, br, zstd]
   Accept-Language: [en-US,en;q=0.9]
   Cache-Control: [max-age=0]
   Sec-Ch-Ua: ["Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome...
   Sec-Ch-Ua-Mobile: [?0]
   Sec-Ch-Ua-Platform: ["macOS"]
   Sec-Fetch-Dest: [document]
   Sec-Fetch-Mode: [navigate]
   Sec-Fetch-Site: [none]
   Sec-Fetch-User: [?1]
   Upgrade-Insecure-Requests: [1]
   User-Agent: [Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWe...
   X-B3-Sampled: [0]
   X-B3-Spanid: [70b28d3aafbee86f]
   X-B3-Traceid: [a0305281afeff95b70b28d3aafbee86f]
   X-Envoy-Attempt-Count: [1]
   X-Envoy-Decorator-Operation: [authtest.authtest.svc.cluster.local:80...
   X-Envoy-Internal: [true]
   X-Envoy-Peer-Metadata: [ChoKCkNMVVNURVJfSUQSDBoKS3ViZXJuZXRlcwodCgxJ...
   X-Envoy-Peer-Metadata-Id: [router~10.244.0.16~istio-ingressgateway-8...
   X-Forwarded-For: [10.244.0.1]
   X-Forwarded-Proto: [http]
   X-Request-Id: [90ee77dd-0395-42be-8148-bcc766f0f246]
   ```
2. Repeating with different URL paths (other than `/login` and `/logout`) will show the `Path:` value changing the
   to match and the `Count:` value increasing, i.e. the count of times that the [authtest](../authtest) service
   has responded to a request.
   
3. Going to http://localhost/login will prompt you to add a `user=` query parameter:
   ```text
   To login, add a user=username query parameter to this ULR path
   ```
   
4. Going to http://localhost/login?user=micky-mouse will create a session cookie containing that name and
   redirect to http://localhost/dashboard. You should now see the session cookie value in the [authtest](../authtest)
   response.
   ```text
   Path:		"/dashboard"
   Count:		5
   
   HEADERS (24)
   =======
   
   Accept: [text/html,application/xhtml+xml,application/xml;q=0.9,image...
   Accept-Encoding: [gzip, deflate, br, zstd]
   Accept-Language: [en-US,en;q=0.9]
   Cookie: [session=micky-mouse]
   Sec-Ch-Ua: ["Chromium";v="128", "Not;A=Brand";v="24", "Google Chrome...
   ... etc ...
   ```
5. Going to http://localhost/logout will reset the session cookie and, in effect, "log you out," and display: 
   ```text
   user micky-mouse has been logged out
   ```
6. Returning to any non `/login` or `/logout` path will again display the [authtest](../authtest)
   response but with the session cookie removed or at least emptied (the behavior depends on the browser type used).

## Testing with the authorization filter

If the installation and configuration is completed through stages 7 and 8, then two additional headers will be found in 
the http://localhost/whatever echo text: This shall contain a JWT as a bearer token. If the JWT is pasted into 

* `X-Extauth-Was-Her` containing the text `blah, blah, blah`
* `X-Extauth-Authorization` containing a JWT bearer token

Pasting the JWT text from the `X-Extauth-Authorization` into the JWT debugger form at https://jwt.io, the payload will 
be seen to contain the username copied from the session cookie. The payload portion of the JWT will look something like 
this after visiting http://localhost/login?user=micky-mouse:

```json
{
  "exp": 1726777195,
  "iat": 1726777165,
  "nbf": 1726777165,
  "sub": "micky-mouse"
}
```

## Visualizing the service mesh

Assuming that you have started the Kiali dashboard web app as described under [Configure Istio visualization](visualize.md),
you should now be able to see a picture of the service mesh at http://localhost:20001/kiali/console/graph/namespaces.
If nothing shows up, make sure that you **Select all** check in the Namespace dropdown.

![Kiali service mesh visualization](kiali.png)

## Suspending Minikube

You can stop Minikube and Colima to get your memory and CPU back at any time by killing any `minikube dashboard` and
`minikube tunnel` shell commands that you have running with Ctrl-C and then execute the following:

```shell
minikube stop

colima stop
```

The next time you start Minikube, the services that you deployed will be brought back to life; you do not
have to go through the installation and deployment steps again. For example:

```shell
colima start -c6 -m12 

minikube start
```
