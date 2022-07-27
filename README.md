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