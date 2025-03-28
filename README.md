# 🏁 Mini Feature Flag Service

A lightweight feature flag API written in Go, inspired by production-grade systems like Netflix's Fast Feature Work Delivery. Supports:

- In-memory flag storage with optional persistence to file
- User and region-based targeting
- CRUD operations for flags via REST API
- Auth-protected write/delete routes

---

## 🚀 Getting Started

### 1. Start the Server

```bash
go run ./cmd/server
```

## 2. API Endpoints (REST)

### Health Check

```
GET /healthz
```

Response:
```json
OK
```

### Check a Flag

GET /flags/:flagName?userId=FILL_ME_IN&region=FILL_ME_IN

Response:
```json
{
  "enabled": true,
  "reason": "User explicitly targeted"
}
```

### Get all Flags

```
GET /flags
```

Response:
```json
[
  {
    "Name": "beta-dashboard",
    "Enabled": true,
    "TargetUsers": null,
    "TargetRegions": null
  },
  {
    "Name": "new-homepage",
    "Enabled": false,
    "TargetUsers": [
      "123",
      "456"
    ],
    "TargetRegions": [
      "us",
      "ca"
    ]
  }
]
```


### Create a Flag

```
POST /flags/:flagName
Authorization: Bearer super-secret-token
Content-Type: application/json
```

Payload:
```json
{
  "enabled": true,
  "targetUsers": ["123", "456"],
  "targetRegions": ["us", "ca"]
}
```

### Update a Flag

```
PUT /flags/:flagName
Authorization: Bearer super-secret-token
Content-Type: application/json
```

Payload:
```json
{
  "enabled": false,
  "targetUsers": [],
  "targetRegions": ["uk"]
}
```

### Delete a Flag

```
DELETE /flags/:flagName
Authorization: Bearer super-secret-token
```

## 💾 Persistence

- All flags are saved to `flags.json` in the root directory (Temporary)
- Automatically loaded on startup and saved on create/update/delete with dummy data.

## 🔒 Security

All destructive actions (POST, PUT, DELETE) require the following header:

```
# Temporary token
Authorization: Bearer super-secret-token
```
