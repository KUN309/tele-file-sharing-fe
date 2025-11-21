# 📦 Telegram File Sharing — Frontend (Telegram Bot)  
Dự án **Tele File Sharing – Frontend** là Telegram Bot dùng để giao tiếp với người dùng, cho phép upload file, xem danh sách file, chia sẻ file và tương tác với Backend API.  
Bot được viết bằng **Go**, giao tiếp trực tiếp với **Telegram Bot API**, và kết nối đến **tele-file-sharing-be** thông qua HTTP.

---
## 👥 Thành viên nhóm Frontend
| Họ tên                      | Vai trò                                                                                   |
| --------------------------- | ----------------------------------------------------------------------------------------- |
| Nguyễn Nguyên Ngọc (Leader) | Tổng hợp và tạo file sườn, sync data giữa FE và BE, hỗ trợ thành viên, use case (share)   |
| Bùi Hoàng Cung              | Làm requirement (non-functional), use case (upload), xử lý command về upload (dự tính)    |
| Nguyễn Trí Thành            | Làm requirement (non-functional), use case (revoke), xử lý command về revoke (dự tính)    |
| Võ Hùng Dũng                | Làm requirement (non-functional), use case (download), xử lý command về download (dự tính)|

---
## 🚀 Hướng dẫn chạy dự án (Local Development)
### 1. Chuẩn bị môi trường
Cần cài đặt trước:  
- [Go](https://go.dev/doc/install) (phiên bản 1.25+)
- [Docker](https://www.docker.com/products/docker-desktop/)
- Docker Compose
- [Sops](https://github.com/getsops/sops) để mã hóa file dev.enc (Hiện tại, nhóm bị trục trặc sops nên chưa có encrypt dev.env thành dev.enc được)

### 2. Cấu hình môi trường
- Tạo file env:
```bash
cp env/example.env env/dev.env
```

- Nếu dùng file .enc:
```bash
sops -d env/dev.enc > env/dev.env
```

- Cấu hình trong env/example.env:
```bash
TELEGRAM_BOT_TOKEN=YOUR_TOKEN
BE_API_BASE=http://your_backend/
USE_WEBHOOK=true
WEBHOOK_URL=http://yourdomain.com/webhook
PORT=YOUR_PORT
API_TIMEOUT=10
TEMP_DIR=/tmp/bot
FILE_BASE_URL=http://your_backend/static
```

### 3. Chạy bot bằng Docker
```bash
docker compose up --build
```

### 4. Chạy bot trực tiếp bằng Go
```bash
go mod tidy
go run ./cmd/bot
```

---
## 🤖 Các lệnh Telegram hỗ trợ
| Lệnh                 | Chức năng                             |
| -------------------- | ------------------------------------- |
| `/me`                | Xem thông tin người dùng trên backend |
| `/files`             | Liệt kê các file đã upload            |
| (WIP) `/share`       | Tạo link chia sẻ file                 |
| (Chưa thử) upload    | Gửi file, bot tự tải file lên backend |
| (Chưa thử) download  | Tải file về từ đường dẫn được chia sẻ |

---
## 📂 Cấu trúc thư mục
```
tele-file-sharing-fe/
├── README.md                     # Giới thiệu & hướng dẫn
├── Dockerfile                    # Build Telegram Bot
├── compose.yml                   # Docker Compose dev environment
├── .github/
│   └── workflows/
│       └── ci.yml                # CI (build-only, chưa có test/lint/scan)
│
├── env/
│   ├── example.env               # File mẫu config
│   ├── dev.enc                   # File env mã hoá (WIP, lỗi với sops)
│
├── cmd/
│   └── bot/
│       └── main.go               # Entry point — chạy bot
│
├── internal/
│   ├── config/
│   │   └── config.go             # Load & validate env
│   │
│   ├── bot/
│   │   ├── handler.go            # Xử lý command & upload file
│   │   ├── router.go             # (Dự tính) Map command → handler
│   │   └── middleware.go         # (Dự tính) Logging, auth
│   │
│   └── api/
│       ├── client.go             # Gửi request đến backend API
│       └── models.go             # Structs API (dự tính tách thành dto/)
│
├── docs/                         # Tài liệu nội bộ
└── reports/                      # Báo cáo môn học (.md)
    ├── M1/
    ├── M2/
    └── M3/
```

---
## 🛠 Quy trình phát triển một tính năng mới
### 1. Tạo nhánh mới
Tạo nhánh riêng biệt trên repository chung và tự gửi Commit & Pull Request để chờ duyệt
```bash
git checkout -b <ten-tinh-nang>
git push origin dev:<ten-tinh-nang>
```

### 2. Cập nhật API (nếu backend thay đổi)
- Chỉnh sửa internal/api/client.go
- Cập nhật struct trong models.go
- Điều chỉnh command tương ứng trong handler.go

### 3. Viết code
- internal/api/* → Giao tiếp backend
- internal/bot/handler.go → Logic command
- internal/config/* → Thêm config nếu cần

### 4. Test thủ công qua Telegram
- Gửi file thử
- Test lệnh /me, /files
- Kiểm tra trường hợp backend lỗi

### 5. Tạo Pull Request
- Push code lên nhánh riêng
- Tạo PR vào main
- Mô tả rõ:
  - Mục đích thay đổi
  - Flow backend liên quan
  - Ảnh chụp bot (nếu cần)

### 6. Merge & dọn nhánh
```bash
git push origin --delete <ten-tinh-nang>

```

