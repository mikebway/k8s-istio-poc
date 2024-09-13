# Running Istio on Minikube locally

```diff
- UNDER CONSTRUCTION
-
- This project is incomplete and may never be completed!!   
-
- Currently incomplete:
-    * Istio is not yet configured
-    * Only one service is built and deployed 
```

## Objective

Demonstrate how Istio can be deployed on top of Minikube to manage path routing for a trivial two
service web application.

Login and logout is crudely illustrated by the creation and removal of a session cookie containing the username.
This is clearly not a secure pattern; do not use it for anything more than demonstration purposes.

## Installation and configuration

1. [Minikube installation and basic configuration](Install.md)
2. [Build and deploy the authtest service](svc-authtest)
3. [Build and deploy the login service](svc-login.md)
4. [Configuring Istio ingress](istio.md)

## Testing

After completing installation and configurationyou should be able to:

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
   response but with the session cookie removed.