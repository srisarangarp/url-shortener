# prerequisites
Foundation DB version 6.2.30 needs to be configured locally, to make this work!!

# URL - shortener
This application can be used to shorten the given URL

## Installation and Running

you call pull the image from [docker](https://hub.docker.com/repository/docker/sreesa7144/url-shortener) and can run like specified:
```bash
docker pull sreesa7144/url-shortener:latest
docker run -d -p 8080:8080 sreesa7144/url-shortener:latest
#You can run the curl command by passing the url as form data like this
curl -X POST -F 'url=<Enter url to shorten here>' localhost:8080/url
```

## Output
The output will look like this
```json
{"shorten_URL":"https://infc.com/<shortened-url>"}
```
infc.com is the assumed domain name here!!
