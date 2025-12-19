# URL Shortener Service

Backend Service rút gọn link tương tự bit.ly, tinyurl - Golang.

**Live Demo:** https://shorten.quocbui.dev/swagger/index.html

## Mô tả bài toán

User có URL dài muốn rút gọn thành link ngắn. Khi truy cập link ngắn sẽ tự động redirect về URL gốc và hệ thống tracking số lượt click, thông tin thiết bị, vị trí.

**Ví dụ:**
```
Input:  https://example.com/very/long/path?param1=value1&param2=value2
Output: https://shorten.quocbui.dev/abc123
```

## Tính năng đã có

- Tạo link rút gọn (Đăng nhập/không đăng nhập)
- Custom alias (cho phép nhập: vd:  https://shorten.quocbui.dev/abc123)
- Link expiration (vd: 24 giờ)
- Analytics: click, browser, OS, device, country
- QR Code với logo
- Rate limiting (100 req/60s)
- Swagger API docs
- JWT

## Tech Stack

| Stack | Công nghệ | Ưu điểm |
|-----------|------------|------------|
| Language | Go 1.24 | Performance, concurrency tốt, compile nhẹ và đơn giản hơn so với Node.js hay Java |
| Framework | Gin | Framework nhẹ, phổ biến |
| Database | PostgreSQL | Tài liệu support tốt, ổn định với GORM |
| ORM | GORM | Auto migration, query builder |
| Auth | JWT | Hỗ trợ tốt cho RestAPI và scale sau này |
| Docs | Swagger | Hỗ trợ cho dev |
| Deploy | Docker + Caddy | Deploy dễ, config dễ hơn so với nginx |

## Cách chạy project

### Local

```bash
# 1. Clone repo
git clone https://github.com/quocbui2020/shorten_url.git
cd shorten_url

# 2. Copy env file
cp .env.example .env

# 3. Sửa .env với database credentials
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=123
DB_NAME=shorten_url

# 4. Chạy server
go run cmd/main.go

# 5. Truy cập
# API: http://localhost:8080
# Swagger: http://localhost:8080/swagger/index.html
```

### Docker

```bash
# 1. Sửa .env 
APP_DOMAIN=localhost:8080
JWT_SECRET=your-super-secret-key
DB_PASSWORD=123

# 2. Chạy
docker-compose up -d --build

```

## Endpoints đã có

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | /api/v1/auth/register | Đăng ký |
| POST | /api/v1/auth/login | Đăng nhập (email,password) |
| POST | /api/v1/shorten | Tạo link rút gọn (url gốc, alias, expired_in) |
| GET | /api/v1/me | Thông tin user (email+name) |
| GET | /api/v1/me/links | Danh sách links đã tạo của user |
| GET | /api/v1/me/links/:code | Chi tiết analytics của link (click, browser, device,...) |
| DELETE | /api/v1/me/links/:code | Xóa link (soft delete) |
| GET | /:code | Logic để redirect |

*Nếu không có token, tự động tạo guest account và trả về token.

## Thiết kế & Quyết định kỹ thuật

### Database Schema

```
users => 1:N => links => 1:N => clicks
```
![Database Entity Diagram](assets/readme/database-entity.png)

**Tại sao chọn PostgreSQL?**
- ACID compliance cho data integrity
- Index B-tree cho lookup short_code O(log n)
- Soft delete với deleted_at
- JSON support nếu cần mở rộng

### Thuật toán generate short code

```go
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 62^6 = 56.8 tỷ combinations
```

- Dùng `crypto/rand` thay vì `math/rand` để secure
- Retry 5 lần nếu bị trùng
- Custom alias validate: 3-20 ký tự, alphanumeric + `_` + `-`

### Xử lý Concurrency

**Vấn đề:** 2 requests cùng lúc tạo link với cùng alias?

**Giải pháp:**
- Database UNIQUE constraint trên `short_code`
- Check existence trước khi insert
- Nếu insert fail do duplicate → retry với code mới

```go
for i := 0; i < 5; i++ {
    shortCode = generateCode()
    if !exists(shortCode) {
        break
    }
}
```

### Click Tracking

```go
// Async tracking - không block redirect
go s.trackClick(linkID, clickInfo)
return originalURL
```

- Parse User-Agent → browser, OS, device
- GeoIP lookup → country, city (ip-api.com free tier)
- Aggregate stats trong bảng clicks

## Trade-offs

### 1. Guest User vs Anonymous Links

**Chọn:** Tự động tạo guest user khi shorten không có token

**Lý do:**
- User có thể xem lại links đã tạo
- Có analytics cho mọi link
- Dễ convert guest → registered user sau

**Nhược điểm:**
- Tạo nhiều guest users trong DB
- Cần cleanup job cho guest users cũ

### 2. GeoIP: External API vs Local Database

**Chọn:** External API (ip-api.com)

**Lý do:**
- Không cần maintain GeoIP database
- Free tier đủ dùng (45 req/min)
- Data luôn up-to-date

**Nhược điểm:**
- Dependency external service
- Rate limit
- Latency (nhưng async nên không ảnh hưởng)

### 3. QR Code: Generate on-demand vs Pre-generate

**Chọn:** Generate on-demand

**Lý do:**
- Không tốn storage
- Luôn có QR mới nhất

**Nhược điểm:**
- CPU cost mỗi request
- Có thể cache nếu cần optimize

## Challenges & Solutions

### 1. Go Version Mismatch

**Vấn đề:** Local Go 1.24, Docker image chỉ có 1.23

**Giải pháp:** Dùng `golang:latest` trong Dockerfile

### 2. QR Code Logo

**Vấn đề:** Logo không hiển thị, file format không đúng

**Giải pháp:** 
- Check file format thật (không chỉ extension)
- Resize logo < 1/3 QR size
- Dùng ErrorCorrectionHighest để QR vẫn scan được với logo

### 3. Swagger Host Hardcoded

**Vấn đề:** Swagger luôn gọi localhost

**Giải pháp:** Dynamic host từ env
```go
docs.SwaggerInfo.Host = cfg.App.Domain
```

## Limitations & Future Improvements

### Hiện tại còn thiếu

- [ ] Unit tests
- [ ] Integration tests
- [ ] Caching layer (Redis)
- [ ] Bulk URL shortening
- [ ] Link edit (update URL)
- [ ] Custom domain per user
- [ ] Webhook notifications

### Production-ready cần thêm

1. **Caching:** Redis cho hot links
2. **CDN:** Cache redirect responses
3. **Monitoring:** Prometheus + Grafana
4. **Logging:** Structured logging (Zap/Zerolog)
5. **Database:** Read replicas, connection pooling
6. **Security:** HTTPS only, CSP headers, input sanitization
7. **Cleanup Job:** Xóa expired links, guest users cũ

### Scalability Plan

```
                    ┌─────────────┐
                    │   Caddy/    │
                    │   Nginx     │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼─────┐ ┌────▼────┐ ┌─────▼─────┐
        │   App 1   │ │  App 2  │ │   App 3   │
        └─────┬─────┘ └────┬────┘ └─────┬─────┘
              │            │            │
              └────────────┼────────────┘
                           │
                    ┌──────▼──────┐
                    │    Redis    │
                    │   (Cache)   │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
        ┌─────▼─────┐ ┌────▼────┐
        │  Primary  │ │ Replica │
        │    DB     │ │   DB    │
        └───────────┘ └─────────┘
```

## Project Structure

```
shorten_url_go/
├── cmd/
│   └── main.go              # Entry point
├── internal/
│   ├── config/              # Configuration
│   ├── dto/                 # Request/Response DTOs
│   ├── handlers/            # HTTP handlers
│   ├── middleware/          # Auth, CORS, Rate limit
│   ├── models/              # Database models
│   ├── repository/          # Data access layer
│   ├── router/              # Route definitions
│   └── service/             # Business logic
├── pkg/
│   └── utils/               # Shared utilities
├── assets/                  # Static files (logo)
├── docs/                    # Swagger generated
├── Dockerfile
├── docker-compose.yml
├── Caddyfile
└── README.md
```

## Author

**Quoc Bui**
- GitHub: [@quocbui2020](https://github.com/quocbui2020)
- Email: quocbui26042005@gmail.com
