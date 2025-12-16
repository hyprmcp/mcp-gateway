# Changelog

## [0.4.0](https://github.com/hyprmcp/mcp-gateway/compare/v0.3.1...v0.4.0) (2025-12-16)


### Features

* add TLS support for gRPC connections to Dex ([#86](https://github.com/hyprmcp/mcp-gateway/issues/86)) ([b150636](https://github.com/hyprmcp/mcp-gateway/commit/b15063651a80cc95b934bb9e2f96e0747db2d63e))


### Bug Fixes

* **deps:** update module github.com/lestrrat-go/httprc/v3 to v3.0.2 ([#88](https://github.com/hyprmcp/mcp-gateway/issues/88)) ([75c1c07](https://github.com/hyprmcp/mcp-gateway/commit/75c1c07c7e9305252b5aae7bdb1ed5984a0d62f2))
* **deps:** update module github.com/lestrrat-go/jwx/v3 to v3.0.12 ([#80](https://github.com/hyprmcp/mcp-gateway/issues/80)) ([61c98ad](https://github.com/hyprmcp/mcp-gateway/commit/61c98adb6d7517211291e15fc58542b5caa9c285))
* **deps:** update module github.com/spf13/cobra to v1.10.2 ([#87](https://github.com/hyprmcp/mcp-gateway/issues/87)) ([89f131c](https://github.com/hyprmcp/mcp-gateway/commit/89f131c32dc9c2b921c1891dc3fd9439cc285518))
* **deps:** update module google.golang.org/grpc to v1.77.0 ([#75](https://github.com/hyprmcp/mcp-gateway/issues/75)) ([59ab1cf](https://github.com/hyprmcp/mcp-gateway/commit/59ab1cf2ecc831beafcbcb26cba9704e2ebfce3e))


### Other

* **deps:** update actions/checkout action to v6 ([#85](https://github.com/hyprmcp/mcp-gateway/issues/85)) ([440d185](https://github.com/hyprmcp/mcp-gateway/commit/440d1853c240d71e10822af4eddb517160834f73))
* **deps:** update dependency go to v1.25.5 ([#76](https://github.com/hyprmcp/mcp-gateway/issues/76)) ([35db0ea](https://github.com/hyprmcp/mcp-gateway/commit/35db0ea25140a61d81e975d1d205630ca5032eac))
* **deps:** update docker/metadata-action action to v5.10.0 ([#83](https://github.com/hyprmcp/mcp-gateway/issues/83)) ([8a5ce3a](https://github.com/hyprmcp/mcp-gateway/commit/8a5ce3aa0235c7c25b76a5bbebb7380f6b664024))
* **deps:** update github artifact actions (major) ([#90](https://github.com/hyprmcp/mcp-gateway/issues/90)) ([d7c96bc](https://github.com/hyprmcp/mcp-gateway/commit/d7c96bc71ebb148778ec4cd097e1b5358b9faef7))

## [0.3.1](https://github.com/hyprmcp/mcp-gateway/compare/v0.3.0...v0.3.1) (2025-10-29)


### Other

* **deps:** update github artifact actions (major) ([#82](https://github.com/hyprmcp/mcp-gateway/issues/82)) ([c1a8df0](https://github.com/hyprmcp/mcp-gateway/commit/c1a8df0dacbd9e23a4eb30c3e30b78f225f3e8fa))
* **deps:** update googleapis/release-please-action action to v4.4.0 ([#81](https://github.com/hyprmcp/mcp-gateway/issues/81)) ([3db34e3](https://github.com/hyprmcp/mcp-gateway/commit/3db34e3944c15dc2dbd63cb60cd3a105185635ad))
* **deps:** update sigstore/cosign-installer action to v4 ([#79](https://github.com/hyprmcp/mcp-gateway/issues/79)) ([1724f81](https://github.com/hyprmcp/mcp-gateway/commit/1724f813dcd39c9a7a3c9214a387228cf5b81b6d))
* use hyprmcp dex for who-am-i demo ([#74](https://github.com/hyprmcp/mcp-gateway/issues/74)) ([7d29ed9](https://github.com/hyprmcp/mcp-gateway/commit/7d29ed9d46034e09665691fc5a34f219c88f469a))

## [0.3.0](https://github.com/hyprmcp/mcp-gateway/compare/v0.2.6...v0.3.0) (2025-10-06)


### Features

* add proxying existing protected resource metadata ([#73](https://github.com/hyprmcp/mcp-gateway/issues/73)) ([5c12538](https://github.com/hyprmcp/mcp-gateway/commit/5c12538e0ae68e38ee0f1ff690720ca6a3526aae))


### Bug Fixes

* **deps:** update module github.com/google/jsonschema-go to v0.3.0 ([#67](https://github.com/hyprmcp/mcp-gateway/issues/67)) ([45865b5](https://github.com/hyprmcp/mcp-gateway/commit/45865b55cae03ff89970a2fcb02e93b2542eb0c1))


### Other

* **deps:** update docker/login-action action to v3.6.0 ([#70](https://github.com/hyprmcp/mcp-gateway/issues/70)) ([e40be1f](https://github.com/hyprmcp/mcp-gateway/commit/e40be1fb41d71df25a742a8abe64974ba8f8a64a))
* upgrade @hyprmcp/mcp-install-instructions-generator to 0.2.0 ([#68](https://github.com/hyprmcp/mcp-gateway/issues/68)) ([ab3e98f](https://github.com/hyprmcp/mcp-gateway/commit/ab3e98f7fcbb9e56be1db5ae707cb737632ae1de))

## [0.2.6](https://github.com/hyprmcp/mcp-gateway/compare/v0.2.5...v0.2.6) (2025-09-15)


### Other

* **deps:** upgrade @hyprmcp/mcp-install-instructions-generator to 0.1.1 ([#64](https://github.com/hyprmcp/mcp-gateway/issues/64)) ([fc3c8db](https://github.com/hyprmcp/mcp-gateway/commit/fc3c8db3750d7503098f3b4fe565c06758bc4c7a))

## [0.2.5](https://github.com/hyprmcp/mcp-gateway/compare/v0.2.4...v0.2.5) (2025-09-15)


### Bug Fixes

* add showing html response also if auth is disabled ([#61](https://github.com/hyprmcp/mcp-gateway/issues/61)) ([39e5729](https://github.com/hyprmcp/mcp-gateway/commit/39e5729bf4fbf78e107dde5e933f697fef3956ea))
* **deps:** update module github.com/google/jsonschema-go to v0.2.3 ([#57](https://github.com/hyprmcp/mcp-gateway/issues/57)) ([38be463](https://github.com/hyprmcp/mcp-gateway/commit/38be4634e90e2f6ea7598ca1462be0363c86017c))
* **deps:** update module github.com/lestrrat-go/jwx/v3 to v3.0.11 ([#60](https://github.com/hyprmcp/mcp-gateway/issues/60)) ([d171699](https://github.com/hyprmcp/mcp-gateway/commit/d1716996f0dd847329872b2bdeb09995d089612d))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v0.5.0 ([#59](https://github.com/hyprmcp/mcp-gateway/issues/59)) ([204e1a2](https://github.com/hyprmcp/mcp-gateway/commit/204e1a2c697e42953c2adc8d08c6978be9743cb2))


### Other

* **deps:** update sigstore/cosign-installer action to v3.10.0 ([#58](https://github.com/hyprmcp/mcp-gateway/issues/58)) ([94afbd7](https://github.com/hyprmcp/mcp-gateway/commit/94afbd7c25b6abde7d0796a1f2dc2deaecc55572))
* final rename to hyprmcp ([#63](https://github.com/hyprmcp/mcp-gateway/issues/63)) ([02baf70](https://github.com/hyprmcp/mcp-gateway/commit/02baf70a9d3a3f0499b54ba07be7caca4ae170f7))

## [0.2.4](https://github.com/hyprmcp/mcp-gateway/compare/v0.2.3...v0.2.4) (2025-09-10)


### Bug Fixes

* **deps:** update module github.com/google/jsonschema-go to v0.2.1 ([#53](https://github.com/hyprmcp/mcp-gateway/issues/53)) ([916eeb5](https://github.com/hyprmcp/mcp-gateway/commit/916eeb5a4312177c369c3d91c5d0011aa7fba800))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v0.4.0 ([#52](https://github.com/hyprmcp/mcp-gateway/issues/52)) ([568900b](https://github.com/hyprmcp/mcp-gateway/commit/568900b8829f14cb80e501563ac9c9a5932a4b53))
* **deps:** update module google.golang.org/grpc to v1.75.1 ([#55](https://github.com/hyprmcp/mcp-gateway/issues/55)) ([6be6137](https://github.com/hyprmcp/mcp-gateway/commit/6be6137374e81b473aa77f6d4311ec0056b05b16))


### Other

* add responding with html if request has "Accept: text/html" header ([#54](https://github.com/hyprmcp/mcp-gateway/issues/54)) ([876ad98](https://github.com/hyprmcp/mcp-gateway/commit/876ad98a7fdff6ed042a575ce65eb83faa413a54))
* **deps:** update dependency go to v1.25.1 ([#49](https://github.com/hyprmcp/mcp-gateway/issues/49)) ([715b945](https://github.com/hyprmcp/mcp-gateway/commit/715b945d4ae01ffc86d4420c2e65d0cc5266ef0d))


### Docs

* add demo video to example README ([#48](https://github.com/hyprmcp/mcp-gateway/issues/48)) ([8ec1fa0](https://github.com/hyprmcp/mcp-gateway/commit/8ec1fa08a7b78a99cbea93b8670298c7107ec390))
* Update README.md ([#51](https://github.com/hyprmcp/mcp-gateway/issues/51)) ([111628b](https://github.com/hyprmcp/mcp-gateway/commit/111628b7197b9ddf71dab1cb0cac63f159b9499f))

## [0.2.3](https://github.com/hyprmcp/mcp-gateway/compare/0.2.2...v0.2.3) (2025-09-03)


### Bug Fixes

* **deps:** update github.com/google/jsonschema-go digest to 7d3a774 ([#47](https://github.com/hyprmcp/mcp-gateway/issues/47)) ([d426690](https://github.com/hyprmcp/mcp-gateway/commit/d426690a6528f8ba3837a00e6743712782afd7a3))
* **deps:** update module github.com/dexidp/dex/api/v2 to v2.4.0 ([#42](https://github.com/hyprmcp/mcp-gateway/issues/42)) ([a90d24f](https://github.com/hyprmcp/mcp-gateway/commit/a90d24f2b69a3ac74daa57ea02bbc3840a3b0d73))
* **deps:** update module github.com/modelcontextprotocol/go-sdk to v0.3.1 ([#40](https://github.com/hyprmcp/mcp-gateway/issues/40)) ([4678d91](https://github.com/hyprmcp/mcp-gateway/commit/4678d91a4a36ea63241cdf9c54eb42fa8c9d4ff3))
* **deps:** update module github.com/spf13/cobra to v1.10.1 ([#43](https://github.com/hyprmcp/mcp-gateway/issues/43)) ([5302710](https://github.com/hyprmcp/mcp-gateway/commit/5302710d5237e0ab16c13a3f6e7bef38e0cc6c9c))


### Other

* **release:** enable include "v" in tag ([#46](https://github.com/hyprmcp/mcp-gateway/issues/46)) ([f22ba35](https://github.com/hyprmcp/mcp-gateway/commit/f22ba35de8bd58ff6bbd874afcc7d6be709990d7))


### Docs

* readme update ([#45](https://github.com/hyprmcp/mcp-gateway/issues/45)) ([df116a2](https://github.com/hyprmcp/mcp-gateway/commit/df116a21de12d7892423b4e5141bff90fa9e0839))
* rename docker-compose config ([#41](https://github.com/hyprmcp/mcp-gateway/issues/41)) ([218d546](https://github.com/hyprmcp/mcp-gateway/commit/218d546ade7b756c79c9bd65f6c851d6977794d6))

## [0.2.2](https://github.com/hyprmcp/mcp-gateway/compare/0.2.1...0.2.2) (2025-08-29)


### Bug Fixes

* omit empty ClientSecret from DCR response ([#38](https://github.com/hyprmcp/mcp-gateway/issues/38)) ([0e3f7d7](https://github.com/hyprmcp/mcp-gateway/commit/0e3f7d7f252519d0cc134c00ae269d9555937d53))

## [0.2.1](https://github.com/hyprmcp/mcp-gateway/compare/0.2.0...0.2.1) (2025-08-28)


### Other

* add config for DCR public/private client ([#36](https://github.com/hyprmcp/mcp-gateway/issues/36)) ([282db47](https://github.com/hyprmcp/mcp-gateway/commit/282db47db602f3f48fcf773af29f3bc96a71ef47))

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
