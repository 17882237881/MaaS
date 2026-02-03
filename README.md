# MaaS

MaaS platform scaffold (CloudWeGo-first): Hertz for HTTP, Kitex for internal RPC, Protobuf for IDL.

## Repo Layout
- cmd/      : service entrypoints
- internal/ : core business logic
- api/      : OpenAPI/Proto definitions
- deploy/   : Helm charts
- infra/    : Terraform
- docs/     : node/step documentation
- scripts/  : local helper scripts

## Quick Start
- make tidy
- make test
- make vet
