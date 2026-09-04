# 📐 System Architecture & Coding Standards

## 1. Core Principles
- **Clean Architecture & Separation of Concerns**: Keep business models decoupled from I/O and HTTP delivery layers.
- **Explicit Error Handling**: Always wrap errors with domain context; never suppress errors.
- **Zero Secrets in Code**: Secrets must only be passed through environment variables or secure key vaults.

## 2. Testing Standard
- All pull requests must include unit tests.
- Target minimum test coverage: 80% on domain and business logic packages.
- Follow table-driven test patterns in Go / TypeScript.
