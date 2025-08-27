# Changelog

## [0.2.0](https://github.com/hyprmcp/mcp-gateway/compare/0.1.2...0.2.0) (2025-08-27)


### Features

* add authorization handler to inject openid scope if missing, return scope in DCR response ([#32](https://github.com/hyprmcp/mcp-gateway/issues/32)) ([e2978e9](https://github.com/hyprmcp/mcp-gateway/commit/e2978e912bc0c841be15318ad807af8def8e2068))


### Bug Fixes

* jwk refresh in background context and reduced min/max interval ([#33](https://github.com/hyprmcp/mcp-gateway/issues/33)) ([c92e2f2](https://github.com/hyprmcp/mcp-gateway/commit/c92e2f250584cd759e61786081f931f42e1cd450))

## [0.1.2](https://github.com/hyprmcp/mcp-gateway/compare/0.1.1...0.1.2) (2025-08-26)


### Bug Fixes

* add request path to metadata and protected resource path ([#31](https://github.com/hyprmcp/mcp-gateway/issues/31)) ([cf4b2c0](https://github.com/hyprmcp/mcp-gateway/commit/cf4b2c04d7913c6ac66bfd4211b5983e33f3324c))


### Docs

* improve github link in example docs ([8fe53f6](https://github.com/hyprmcp/mcp-gateway/commit/8fe53f6b71072b36fd3d26e0974fee193bba8ab9))

## [0.1.1](https://github.com/hyprmcp/mcp-gateway/compare/0.1.0...0.1.1) (2025-08-26)


### Other

* add auth proxy listener for advanced use cases ([#28](https://github.com/hyprmcp/mcp-gateway/issues/28)) ([4d7aa06](https://github.com/hyprmcp/mcp-gateway/commit/4d7aa06d50ee0f9f2e7d227cd913c3ab79e4b484))


### Docs

* change to https github url ([d68df0a](https://github.com/hyprmcp/mcp-gateway/commit/d68df0a7ff6b4d43da13884bf5263f9ee033112d))
* explicitly expose port 9000 for the gateway demo ([a985680](https://github.com/hyprmcp/mcp-gateway/commit/a98568038da502e9352b8e54098c7b33a9abda00))
* increase waitlist button size ([30b1d8a](https://github.com/hyprmcp/mcp-gateway/commit/30b1d8ad03facd53be21b8fdf254e9a91f80bf07))

## [0.1.0](https://github.com/hyprmcp/mcp-gateway/compare/0.1.0-alpha.6...0.1.0) (2025-08-25)


### Bug Fixes

* **deps:** update module github.com/lestrrat-go/httprc/v3 to v3.0.1 ([#5](https://github.com/hyprmcp/mcp-gateway/issues/5)) ([ca2f8d4](https://github.com/hyprmcp/mcp-gateway/commit/ca2f8d47b7faec572029b86e76f27b7674e63f77))
* **deps:** update module github.com/lestrrat-go/jwx/v3 to v3.0.10 ([#6](https://github.com/hyprmcp/mcp-gateway/issues/6)) ([91115cb](https://github.com/hyprmcp/mcp-gateway/commit/91115cb5c4ded8539b081b4530d850cff96e465c))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v0.3.0 ([#23](https://github.com/hyprmcp/mcp-gateway/issues/23)) ([53ac569](https://github.com/hyprmcp/mcp-gateway/commit/53ac5693166321d7ac75fed84d7b7dfb1e0cfd3b))
* **deps:** update module google.golang.org/grpc to v1.75.0 ([#24](https://github.com/hyprmcp/mcp-gateway/issues/24)) ([a4b29c6](https://github.com/hyprmcp/mcp-gateway/commit/a4b29c6969f0a398f93ddcd8b9ba9377ad691e7c))


### Other

* add license ([#11](https://github.com/hyprmcp/mcp-gateway/issues/11)) ([32eac1f](https://github.com/hyprmcp/mcp-gateway/commit/32eac1f321cf9c9005f26f349d3620ef1299c872))
* add release-please ([#14](https://github.com/hyprmcp/mcp-gateway/issues/14)) ([afa876b](https://github.com/hyprmcp/mcp-gateway/commit/afa876b6458bf08ae0bd5ac30caf827cd12f3a36))
* Configure Renovate ([#4](https://github.com/hyprmcp/mcp-gateway/issues/4)) ([e449bf5](https://github.com/hyprmcp/mcp-gateway/commit/e449bf5575cc9de5afb07b1d3fa095b2ca28b12a))
* **deps:** update actions/checkout action to v4.3.0 ([#20](https://github.com/hyprmcp/mcp-gateway/issues/20)) ([ed159de](https://github.com/hyprmcp/mcp-gateway/commit/ed159dec779e164a3bde1104cc059ec6f6033282))
* **deps:** update actions/checkout action to v5 ([#25](https://github.com/hyprmcp/mcp-gateway/issues/25)) ([4934a48](https://github.com/hyprmcp/mcp-gateway/commit/4934a48eb4787add1b44ff6f837ebc97414ded54))
* **deps:** update actions/download-artifact action to v5 ([#26](https://github.com/hyprmcp/mcp-gateway/issues/26)) ([6e5c3e4](https://github.com/hyprmcp/mcp-gateway/commit/6e5c3e4e409d5f47556647c792bbb2593ab13853))
* **deps:** update dependency go to v1.25.0 ([#8](https://github.com/hyprmcp/mcp-gateway/issues/8)) ([e1bf63f](https://github.com/hyprmcp/mcp-gateway/commit/e1bf63f5a17850f837b9a87690ecc55020f4a1f3))
* **deps:** update docker/login-action action to v3.5.0 ([#21](https://github.com/hyprmcp/mcp-gateway/issues/21)) ([2be62b3](https://github.com/hyprmcp/mcp-gateway/commit/2be62b345c2543f774404e89d949f7f18bb62cd2))
* **deps:** update docker/login-action action to v3.5.0 ([#9](https://github.com/hyprmcp/mcp-gateway/issues/9)) ([b5ff46e](https://github.com/hyprmcp/mcp-gateway/commit/b5ff46ea5b81ef02b8241cd082699c61777a4838))
* **deps:** update docker/metadata-action action to v5.8.0 ([#15](https://github.com/hyprmcp/mcp-gateway/issues/15)) ([b71aaff](https://github.com/hyprmcp/mcp-gateway/commit/b71aaff61d8106ad09f66e17ae692ef2644d0e89))
* **deps:** update golang docker tag to v1.25 ([#16](https://github.com/hyprmcp/mcp-gateway/issues/16)) ([6dec004](https://github.com/hyprmcp/mcp-gateway/commit/6dec0041667b697a2b62d07b476789626bca57cf))
* **deps:** update googleapis/release-please-action action to v4.3.0 ([#22](https://github.com/hyprmcp/mcp-gateway/issues/22)) ([0127bef](https://github.com/hyprmcp/mcp-gateway/commit/0127bef9fe92ee6a1fe88735b69e80d01432a76e))
* **deps:** update sigstore/cosign-installer action to v3.9.2 ([#19](https://github.com/hyprmcp/mcp-gateway/issues/19)) ([336d427](https://github.com/hyprmcp/mcp-gateway/commit/336d427df8a25ac60e51cd808afbd3db5d9822f9))
* rename to github.com/hyprmcp/mcp-gateway ([#12](https://github.com/hyprmcp/mcp-gateway/issues/12)) ([6a4cc1f](https://github.com/hyprmcp/mcp-gateway/commit/6a4cc1f30537e9d3bab4d981865f99aa34f1ce21))
* upgarde mcp-who-am-i to multiarch version ([cc77504](https://github.com/hyprmcp/mcp-gateway/commit/cc77504b02de1cb27d38f1b1d6a96ad374941ed4))


### Docs

* add who-am-i example ([#3](https://github.com/hyprmcp/mcp-gateway/issues/3)) ([14a7d32](https://github.com/hyprmcp/mcp-gateway/commit/14a7d3245a7549985aadb485da964dc945fd75fe))
* improve kbd buttons ([d227ad6](https://github.com/hyprmcp/mcp-gateway/commit/d227ad60ac1fdf42125ce78019c811af9235a988))
* remove obsolete commands ([#13](https://github.com/hyprmcp/mcp-gateway/issues/13)) ([489630b](https://github.com/hyprmcp/mcp-gateway/commit/489630b21da4b98b4e15f3739f220df1858bb233))


### CI

* release 0.1.0 ([c195c44](https://github.com/hyprmcp/mcp-gateway/commit/c195c44d6d7c4fa7742621955f1c6e711e04c120))
