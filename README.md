# Njia

Njia --- Swahili for Route

A Go-based HTTP(S) redirect service for redirecting generic oAuth callback urls to tenant-specific urls 

# Why Njia?
In oAuth flows, the User-Agent is redirected by the OP back to the RP. However, the OP will only redirect back to preregistered redirectUrls for the RP's oAuth Client.

When deploying a SaaS webapp on a url that is specifc per tenant, like https://tenant1.example.com or https://example.com/tenant1, each tenant-specific url would have to be added as a redirectUrl, which is cumbersome and doesn't scale well

Hence Njia was created: by encoding the bit of the actual redirect url that is dynamic for tenants in the state parameter (which should be used!), Njia will read it, merge it into a preconfigured redirectUrlTemplate and redirect the User-Agent to the tenant-specific url 

What does the flow look like?:
- Njia is started (by the RP) with `redirectUrlTemplate` set to for example `https://*.example.com` and exposed to the internet on `https://redirect.example.com`
- RP registers an oAuth Client with the OP with a generid redirect url: https://redirect.example.com
- RP starts an oAuth flow in the browser (User-Agent) for a user of tenant1, by redirecting to the oauth endpoint of the OP, sending `&state=tenant1:theActualRandomStateValue` as state parameter  
- at some point in the oAuth flow, the OP redirect the User-Agent back to the RP on `\https://redirect.example.com?someOauthParams=xxx&state=tenant1:theActualRandomStateValue
- Njia receives the incoming request and:
  - builds the new, tenant-specific redirectUrl, starting with the redirectUrlTemplate value (https://*.example.com)
  - reads the state parameter from the requestUrl of the incoming request and splits it by colon (`:`)
  - replaces the astrix in the new redirectUrl with the first part of the state value (`tenant1`)
  - appends the path of the current requestUrl (if any) to the new redirectUrl
  - copies over all query parameters of the incoming request to the new redirect url, except the state parameter
  - appends a new state parameter with just the second value of the incoming state parameter (`theActualRandomStateValue`)
  - redirects the User-Agent to the resulting url: https://tenant1.example.com?someOauthParams=xxx&state=theActualRandomStateValue

## Getting Started
Build using Go and then have fun with it

## Deployment

The service supports the following config options, through environment variables:
- port: the port the service listens on. Defaults to 443 if TLS is configured, otherwise 8080
- redirectUrlTemplate (required): the template for the url that Njia needs to redirect incoming GETs to. The template must contain a astrix, which will be replaced by the first part of the state parameter after splitting the state parameter by colon (`:`) 
- certPem & keyPem (optional): the actual certificate and key strings for TLS 

### Docker

Build your Docker container

`cd server`  
`docker build -t njia .`

and run it with 

`docker run --name njia --env redirectUrlTemplate='https://*.example.com' -p 8080:8080 njia`

> the server starts on port `8080` by default but a different port can be specified with `--env port=9999`

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
