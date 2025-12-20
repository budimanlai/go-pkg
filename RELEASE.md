v1.0.6
- Add package Storage
- Add Auth Middleware:
  * Key Provider: Base Key Provider, DB Api Key Provider
  * Basic HTTP Auth (Authorization: Basic username:password)
  * JWT Auth (Authorization: Bearer xxxxxxx)
  * Header Auth (x-api-key)
  * Query String Auth (?access-token=xxxx)