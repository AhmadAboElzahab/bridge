---
layout: home

hero:
  name: Bridge API
  text: UAE Domestic Service Marketplace
  tagline: Go 1.24 · Gin · GORM · PostgreSQL · JWT · Docker
  actions:
    - theme: brand
      text: Get Started
      link: /getting-started
    - theme: alt
      text: API Reference
      link: /api-reference
    - theme: alt
      text: Swagger UI ↗
      link: http://localhost:8080/swagger/index.html

features:
  - icon: 🚀
    title: Getting Started
    details: Prerequisites, local setup, Docker setup, and how to regenerate Swagger docs.
    link: /getting-started

  - icon: ⚙️
    title: Environment Variables
    details: Every env var explained — database, auth, CORS, and all three storage drivers.
    link: /environment-variables

  - icon: 🏗️
    title: Architecture
    details: Layer responsibilities, directory map, route map, and known issues.
    link: /architecture

  - icon: 📡
    title: API Reference
    details: All endpoints with request/response shapes, error codes, and pagination.
    link: /api-reference

  - icon: 🗄️
    title: Storage
    details: Local disk, Cloudflare R2, and AWS S3 — configure with a single env var.
    link: /storage

  - icon: 🔍
    title: Tab & Filter System
    details: How per-user views, column configs, and recursive filters work end-to-end.
    link: /tab-filter-system

  - icon: 🐳
    title: Deployment
    details: Docker dev setup, production with Nginx + SSL, and operations commands.
    link: /deployment
---
