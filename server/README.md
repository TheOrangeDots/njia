# Flagger

A Go-based HTTP(S) redirect service for redirecting generic oAuth callback urls to tenant-specific urls 

# Why Flagger?


## Getting Started
Build using Go and then have fun with it

## Deployment

The service supports the following config options, through environment variables:
- port: the port the service listewns on. Defaults to 443 if TLS is configured, otherwise 8080
- redirectUrlTemplate (required): the url 
- certPem & keyPem (optional): the actual certificate and key strings for TLS 

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
