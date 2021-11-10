# AVOXI DEMO APP PROMPT

Scenario (note, this is fictional, but gives a sense of the types of requests we might encounter):
Our product team has heard from several customers that we restrict users to logging in to their UI accounts from selected countries to prevent them from outsourcing their work to others.  For an initial phase one we're not going to worry about VPN connectivity, only the presented IP address.

The team has designed a solution where the customer database will hold the white listed countries and the API gateway will capture the requesting IP address, check the target customer for restrictions, and send the data elements to a new service you are going to create.  
The new service will be an HTTP-based API that receives an IP address and a white list of countries.  The API should return an indicator if the IP address is within the listed countries.  You can get a data set of IP addresses to country mappings from https://dev.maxmind.com/geoip/geoip2/geolite2/.

We do our backend development in Go (Golang) and prefer solutions in that language, but we can accept submissions in common backend languages (Java, Node.js, Python, C, C++, C#, etc).  We'll be explicitly looking at coding style, code organization, API design, and operational/maintenance aspects such as logging and error handling.

We'll also be giving bonus points for things like
- Including a Docker file for the running service
- Including a Kubernetes YAML file for running the service in an existing cluster
- Exposing the service as gRPC in addition to HTTP
- Documenting a plan for keeping the mapping data up to date.  Extra bonus points for implementing the solution.
- Other extensions to the service you think would be worthwhile.  If you do so, please include a brief description of the feature and justification for its inclusion.  Think of this as what you would have said during the design meeting to convince your team the effort was necessary.

# Web Service documentation
We have one main endpoint, located at /api/v1/checkiplocation
- Example invocation via cURL for the service running locally:
```
curl --location --request POST 'localhost:8080/api/v1/checkiplocation' \
--header 'Accept: application/json' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ=' \
--header 'Content-Type: application/json' \
--data-raw '{
    "ip_address": "71.229.168.95",
    "country_names": [
        "United States",
        "Mexico",
        "Canada"
    ]
}'
```
It will return a `302 Found` status code if the IP address successfully maps to a country provided in the `country_names` list. Otherwise it will return a `404 Not Found` status. If the request is malformed, i.e. the IP address was incomplete, or the list of countries was not included, or the JSON could not be parsed, it will return a `400 Bad Request` status.
This service is secured via Basic Authentication. The basic auth creds in the above example are `admin:password`, but can and should be changed in the .env file provided in the top-level project directory.
We assume that the IP address provided in the JSON payload will adhere to IPv4 structure, but the system can already handle IPv6 addresses as well, if the need ever arises.
We assume that the `country_names` field provided will be a list of strings. We could have asked for a comma-separated list like "Mexico,Canada,...etc" but thought that a list would be easier to parse given the existence of companies with spaces in their name, like the United States or Republic of Korea. Escaping spaces in a single string could be done, but it adds complexity when it isn't necessary.
The same consideration to spaces in a countries' name applies in our conscious decision to provide our endpoint as a `POST` request, so that we can ask for and expect a JSON object rather than building it as a `GET` request, and leaving the values in the URL itself as query parameters like `http://localhost:8080/api/v1/checkiplocation?ip=x.x.x.x&country_names=Canada,Mexico,Ireland`.

# Setup
This service will need to provide geoip credentials of the form User_ID / License_Key for us to utilize their 3rd party web service. Instructions for generating those credentials can be found at https://dev.maxmind.com/geoip/geolite2-free-geolocation-data#Access. Once generated, those credentials should be placed in the .env file with the Basic Authentication credentials. The .env file will resemble this:
```# basic auth creds
BASIC_AUTH_USERNAME=foo
BASIC_AUTH_PASSWORD=bar

# external geoIP lookup api creds
MAXMIND_USER_ID=hello
MAXIND_LICENSE_KEY=world

# environment
ENVIRONMENT=test
```
Additional environment variables that handle sensitive information should also be placed here. Check how we're using the credentials already in the env file for examples of how to use an env var once it is loaded in the workspace.
Before release of this service, make sure the `ENVIRONMENT` var is set to `prod`. Setting it to `test` at the moment logs the IP address along with whether the requests succeeded or failed.

# Running
## Local Development
Commands for running this as a standalone go application for testing purposes and local development are:
```
go build -o out/app main/main.go && ./out/app serve
```

## Docker
The assumption is this service will be run via Docker, and a Dockerfile is provided in the top-level project directory. Commands for building the Docker image and running that image are as follows:
```docker build -f Dockerfile -t avoxi-demo-app:0.0.1 .
docker run --env-file ./avoxi-demo-app.env --publish 8080:8080 avoxi-demo-app:0.0.1
```
To get the .env file into the docker workspace we provide the .env file as a arg in the `docker run` command, and then inside the Dockerfile we copy it into our final image via ```COPY --from=builder /build/avoxi-demo-app.env .```. This allows us to not expose any sensitive information to console logging. The .env file cannot be committed to VCS for obvious reasons, so check the team's sharepoint for the latest version to use in your local workspace.

# Considerations and Assumptions
In order to avoid issues keeping the mapping data up to date we are directly utilizing the geolite2 web service (documented at https://dev.maxmind.com/geoip/docs/web-services?lang=en) rather than downloading a local copy of their database. This allows us to always be on the latest data they have.
- Some pros/cons to this decision:
    - This workflow is simpler due to abstracting our data upkeep and validation needs to a 3rd party, in this case geoip.
    - Need to take into consideration changes to the 3rd party API, although they state that any future api development will be versioned according.
    - Need to be mindful of how many geoip API calls we are making from our service due to the possibility of web service throttling on their end. 
        - We could probably cut down on the breadth of outgoing geoip calls significantly by introducing a local mapping of IP addresses to the countries we pulled back from geoip, so that we can look an IP up once, cache the result, and then not have to remake the call if we receive an additional request to verify that IP address. This could either be accomplished by an in-memory cache, which would mean monitoring our memory usage, or using a separate storage solution like Redis.
    - If their service still isn't up to handling the breadth of calls we encounter then we might need to explore a different service architecture. We're currently utilizing the free service, but paying for a more robust and/or accurate version should be considered.
- If there comes a time when using the geoip web service isn't possible, we could explore automating the process of downloading their ip -> country mapping database and storing it in our environment for our use. This approach would be more complicated and involve the use of involving our own storage solution (my vote would go towards a key-value database like DynamoDB or Redis). This would solve the throttling issue but we would still be dependent on a 3rd party.

# Future Ideas
- Implement monitoring of API statistics with something like statsd and tracing with Jaeger
- Implement RPC support using gRPC or Twirp
- Examine whether a more robust security device besides Basic Authentication should be explored. JWTs might be an option, but could also just be too much hassle for what they would provide.
- Better testing coverage, specifically around the api package
- Work with team where requests will originate to start sending `X-Request-ID` header so that we can log specific request details. Currently the only useful identifying data being sent with each request is the IP address, which is probably too sensitive for logging in production code.

# Main 3rd party libraries used:
- go.uber.org/zap for logging
- github.com/spf13/viper for environment config support
- github.com/spf13/cobra for standalone CLI support
- inet.af/netaddr for IP validation/parsing support
- github.com/gin-gonic/gin for HTTP router and basic web server needs.
