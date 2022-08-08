# CloudIndexer
Backend Service to link to your Cloud Drive and Search 


# Running ES and Kibana
Kibana - https://www.elastic.co/guide/en/kibana/current/docker.html
ES - https://www.elastic.co/guide/en/elasticsearch/reference/current/docker.html
docker run --name es01 --net elastic -p 9200:9200 -p 9300:9300  -e "discovery.type=single-node" 
       -e "xpack.security.enabled=false" -it docker.elastic.co/elasticsearch/elasticsearch:8.3.2

# Running ES with security disabled
```
docker run \
       -p 9200:9200 \
       -p 9300:9300 \
       -e "discovery.type=single-node" \
       -e "xpack.security.enabled=false" docker.elastic.co/elasticsearch/elasticsearch:5.6.3
```
https://stackoverflow.com/questions/47035056/how-to-disable-security-username-password-on-elasticsearch-docker-container

# Dropbox
- Docker provides Access Tokens in case you dont want to implement a OAuth2 Flow
- These tokens are short lived. They are prefixed with `sl` and valid for 4 hours.

# OAuth2 
- Get Client Key, Client Secret from registering your app on Dropbox
- First, add `http://localhost:8080/redirect` or any other URL as redirect in the OAuth2 section of the app
- Go to the authorization link in your browser something like this `https://www.dropbox.com/oauth2/authorize?client_id=MY_CLIENT_ID&redirect_uri=MY_REDIRECT_URI&response_type=code`
- The redirect URI will contain a param like `code` with the **Authorization Code**
- You can use a API endpoint to capture this redirect, or copy the Authorization code from the browser window.
- To get auth token, refresh token do either CURL/Postman or a POST Call via the API at startup.
- https://www.dropbox.com/developers/documentation/http/documentation

Sample Request
```
curl --location --request POST 'https://api.dropboxapi.com/oauth2/token' \
--header 'Authorization: Basic {Base64(CLIENT_KEY:CLIENT SECRET)}' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'code={AUTHORIZATION CODE}' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'redirect_uri=http://localhost:8080/redirect'
```

- Grab the Refresh Token and add it to Config.
- Code periodically generates access tokens. Dropbox provides short lived tokens(4 hours expiry). 
- They start with `sl`