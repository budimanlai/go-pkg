v3.0.0
- Migrate to Fiber v3
- Module path changed to github.com/budimanlai/go-pkg/v3
- v1.x.x branch preserved for legacy Fiber v2 support
- Breaking changes from Fiber v3:
  * Handler signature: func(c fiber.Ctx) error (Ctx is now interface)
  * c.BodyParser() → c.Bind().Body()
  * keyauth: KeyLookup field replaced by Extractor (extractors.FromHeader / extractors.FromQuery)
  * basicauth: ContextUsername and ContextPassword fields removed
  * basicauth: Authorizer now receives (user, pass string, c fiber.Ctx) bool
- Fix: JWTAuth default ContextKey corrected to "user"

v1.0.7
- tambah package base. Pindahan dari go-core/base

v1.0.6
- Add package Storage
- Add Auth Middleware:
  * Key Provider: Base Key Provider, DB Api Key Provider
  * Basic HTTP Auth (Authorization: Basic username:password)
  * JWT Auth (Authorization: Bearer xxxxxxx)
  * Header Auth (x-api-key)
  * Query String Auth (?access-token=xxxx)