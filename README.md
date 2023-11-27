# Cyclone
Cyclone is a Gin API that connects Bruharmy's Powershell module to a RESTful API. Ingested by Kamino.

## Known Issues
### "Same-Origin-Policy" Errors on Firefox and Safari
If you are receiving a problem where Kamino is getting its Axios requests blocked, you may have to accept the self-signed certificates for the API endpoit, specifically.

Roughly, you'll have to do the following steps:

Go to > https://<cyclone_url>:8080 > Add exception
