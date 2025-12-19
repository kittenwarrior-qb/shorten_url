# URL Shortener Service

Backend Service rút gọn link tương tự bit.ly, tinyurl - Golang.

**Live Demo:** https://shorten.quocbui.dev/swagger/index.html

## Mô tả bài toán

User có URL dài muốn rút gọn thành link ngắn. Khi truy cập link ngắn sẽ tự động redirect về URL gốc và hệ thống tracking số lượt click, thông tin thiết bị, vị trí.

```
Input:  https://example.com/very/long/path?param1=value1&param2=value2
Output: https://shorten.quocbui.dev/abc123
```

## Tính năng

- Tạo link rút gọn (có/không đăng nhập)
- Custom alias (vd: `/my-link`)
- Link expiration
- Analytics: click count, browser, OS, device, country
- QR Code với logo
- Rate limiting (100 req/60s)
- JWT Authentication
- Swagger API docs

## Tech Stack

| Component | Choice | Lý do |
|-----------|--------|-------|
| Language | Go 1.24 | Performance cao, concurrency native, binary nhẹ |
| Framework | Gin | Lightweight, ecosystem tốt |
| Database | PostgreSQL | ACID, indexing mạnh, GORM support tốt |
| Auth | JWT | Stateless, scale horizontal dễ |
| Deploy | Docker + Caddy | Auto HTTPS, config đơn giản hơn Nginx |

## Cách chạy

### Local
```bash
git clone https://github.com/quocbui2020/shorten_url.git
cd shorten_url
cp .env.example .env
# Sửa DB credentials trong .env
go run cmd/main.go
# → http://localhost:8080/swagger/index.html
```

### Docker
```bash
cp .env.example .env
docker-compose up -d --build
```

## API Endpoints

| Method | Endpoint | Mô tả |
|--------|----------|-------|
| POST | `/api/v1/auth/register` | Đăng ký |
| POST | `/api/v1/auth/login` | Đăng nhập |
| POST | `/api/v1/shorten` | Tạo link rút gọn |
| GET | `/api/v1/me/links` | Danh sách links của user |
| GET | `/api/v1/me/links/:code` | Chi tiết + analytics |
| DELETE | `/api/v1/me/links/:code` | Xóa link (soft delete) |
| GET | `/:code` | Redirect về URL gốc |

## Thiết kế Database

```
users (1) ──→ (N) links (1) ──→ (N) clicks
```

**Indexes:**
- `links.short_code` - UNIQUE, lookup O(1)
- `links.user_id` - Query links theo user
- `links.expires_at` - Filter expired links
- `clicks.link_id` - Aggregate analytics
- `clicks.clicked_at` - Time-series queries

**Tại sao PostgreSQL thay vì NoSQL?**
- Cần ACID cho việc tạo short code unique
- Foreign key đảm bảo data integrity
- Aggregate queries cho analytics phức tạp
- GORM migration tự động

## Thuật toán Generate Short Code

```go
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// 62^6 = 56.8 tỷ combinations
```

- Dùng `crypto/rand`
- Retry tối đa 5 lần 
- Custom alias: 3-20 ký tự

**Tại sao không dùng hash (MD5/SHA)**
- Random 6 chars với 62 charset đủ cho scale vừa
- Đơn giản, dễ debug

## Xử lý Concurrency

**Vấn đề:** 2 requests cùng tạo link với alias "my-link" cùng lúc?

**Giải pháp:**
```go
// Transaction
err = s.txManager.ExecuteInTransaction(func(tx *gorm.DB) error {
    existing, _ := s.linkRepo.GetByShortCodeForUpdate(tx, *customAlias)
    if existing != nil {
        return ErrAliasAlreadyExists
    }
    return s.linkRepo.CreateWithTx(tx, link)
})
```

- Row-level locking ngăn race condition

## Trade-offs

### 1. Guest User vs Anonymous Links
**Chọn:** Tự động tạo guest account khi shorten không có token

**Ưu điểm:**
- User có thể xem lại links đã tạo
- Mọi link đều có analytics
- Dễ convert guest thành registered user sau này

**Nhược điểm:**
- Tạo nhiều records, cần cleanup job

### 2. QR Code Generation
**Chọn:** Generate on-the-fly, trả về base64

**Ưu điểm:**
- Không tốn storage
- Luôn up-to-date

**Nhược điểm:**
- Tốn CPU mỗi request
- *Cải thiện:* Cache với Redis hoặc pre-generate

### 3. Click Tracking
**Chọn:** Async goroutine + transaction

```go
go s.trackClick(link.ID, clickInfo)  // Không block redirect
```

**Ưu điểm:**
- Redirect response nhanh 
- Transaction đảm bảo click record + click_count

**Nhược điểm:**
- Có thể mất click nếu server crash giữa chừng
- *Cải thiện:* Message queue (Kafka/RabbitMQ)

## Challenges & Solutions

### 1. Race Condition khi tạo Custom Alias
**Vấn đề:** Check exist → Insert có gap, 2 requests có thể cùng pass check

**Giải pháp:** `SELECT ... FOR UPDATE` lock row trong transaction

### 2. N+1 Query trong Analytics
**Vấn đề:** Aggregate browser, OS, device, country riêng lẻ

**Giải pháp:** Batch queries với GROUP BY

## Security Considerations

- **Rate Limiting:** IP-based, 100 req/60s
- **URL Validation:** Chỉ accept http/https, max 2048 chars
- **Soft Delete:** Links không bị xóa khỏi db

## Performance Considerations

### Nếu có 1 triệu links

| Vấn đề | Giải pháp | Giải thích |
|--------|-----------|------------|
| query chậm | `short_code` UNIQUE INDEX GORM tích hợp sẵn | B-tree index giúp tìm kiếm O(log n) thay vì O(n). Với 1M records, chỉ cần ~20 comparisons thay vì scan toàn bộ table |
| response lớn | Đã có Pagination | Không load hết 1M links vào memory. Dùng OFFSET/LIMIT, trả về 10-50 records/page |
| Connection limit | Connection pooling | GORM mặc định pool connections, reuse thay vì tạo mới mỗi request.  |

### Nếu traffic tăng 100x

| Vấn đề | Giải pháp | Giải thích |
|--------|-----------|------------|
| DB quá tải khi redirect | **Redis cache** | Lưu link hay dùng vào bộ nhớ tạm. Truy cập nhanh hơn nhiều so với query DB mỗi lần |
| Query analytics chậm | **Read replicas** | Có thể scale DB phụ để đọc báo cáo, DB chính chỉ lo ghi. Hai việc không chặn nhau |
| Mất click khi traffic cao | **Message queue** | Gửi tất cả click vào hàng đợi trước, xử lý sau (Job). Không sợ mất dữ liệu khi server bận |

## Limitations & Future Improvements

### Hiện tại còn thiếu
- Custom domain cho user
- Hệ thống trả phí cho user
- Quản lý user auth đầy đủ
- Quản lý realtime cho analytics overview
- Phân quyền

### Production-ready cần thêm
- Distributed rate limiting scale lên Redis
- Unit tests, Integration tests và CI/CD actions
- Metrics & monitoring 
- Structured logging 
- Database backup 
- Thiếu tính năng nghiệp vụ