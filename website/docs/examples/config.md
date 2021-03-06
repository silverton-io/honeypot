---
sidebar_position: 100
---


# Sample Configuration

The following is a sample configuration for collecting `snowplow`, `cloudevents`, `self-describing`, `webhook`, and `relay` events before validating them and sending them to `redpanda`, `mysql`, `postgres`, `materialize`, and `clickhouse`:

```
version: 1.1

app:
  env: development
  mode: debug
  port: 8080
  trackerDomain: trck.slvrtnio.com
  health:
    enabled: true
    path: /health
  stats:
    enabled: true
    path: /stats

middleware:
  timeout:
    enabled: false
    ms: 2000
  rateLimiter:
    enabled: false
    period: S
    limit: 10
  identity:
    cookie:
      enabled: true
      name: nuid
      secure: true
      ttlDays: 365
      domain: ""
      path: /
      sameSite: Lax
    fallback: 00000000-0000-4000-A000-000000000000
  cors:
    enabled: true
    allowOrigin:
      - "*"
    allowCredentials: true
    allowMethods:
      - POST
      - OPTIONS
      - GET
    maxAge: 86400
  requestLogger:
    enabled: true
  yeet:
    enabled: false

inputs:
  snowplow:
    enabled: true
    standardRoutesEnabled: true
    openRedirectsEnabled: true
    getPath: /plw/g
    postPath: /plw/p
    redirectPath: /plw/r
  cloudevents:
    enabled: true
    path: /ce/p
  generic:
    enabled: true
    path: /gen/p
    contexts:
      rootKey: contexts
      schemaKey: schema
      dataKey: data
    payload:
      rootKey: payload
      schemaKey: schema
      dataKey: data
  webhook:
    enabled: true
    path: /wb/hk
  relay:
    enabled: true
    path: /relay

schemaCache:
  backend:
    type: fs
    path: ./schemas/
  ttlSeconds: 300
  maxSizeBytes: 104857600 # 100mb -> 100 * 1024 * 1024
  purge:
    enabled: true
    path: /c/purge
  schemaDirectory:
    enabled: true

sinks:
  - name: nada
    type: blackhole
    deliveryRequired: false
  - name: redpanda
    type: kafka
    deliveryRequired: true
    kafkaBrokers:
      - 127.0.0.1:9092
    validTopic: honeypot-valid
    invalidTopic: honeypot-invalid
  - name: mysql
    type: mysql
    deliveryRequired: true
    mysqlHost: 127.0.0.1
    mysqlPort: 3306
    mysqlDbName: honeypot
    mysqlUser: honeypot
    mysqlPass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
  - name: postgres
    type: postgres
    deliveryRequired: true
    pgHost: 127.0.0.1
    pgPort: 5432
    pgDbName: honeypot
    pgUser: honeypot
    pgPass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
  - name: materialize
    type: materialize
    deliveryRequired: false
    mzHost: 127.0.0.1
    mzPort: 6875
    mzDbName: materialize
    mzUser: materialize
    mzPass: ""
    validTable: honeypot_valid
    invalidTable: honeypot_invalid
  - name: clickhouse
    type: clickhouse
    deliveryRequired: true
    clickhouseHost: 127.0.0.1
    clickhousePort: 9000
    clickhouseDbName: honeypot
    clickhouseUser: honeypot
    clickhousePass: honeypot
    validTable: honeypot_valid
    invalidTable: honeypot_invalid

squawkBox:
  enabled: true
  cloudeventsPath: /sqwk/ce
  snowplowPath: /sqwk/sp
  genericPath: /sqwk/gen

privacy:
  anonymize:
    device:
      ip: false
      useragent: false
    user:
      id: false

tele:
  enabled: true
```