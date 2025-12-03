# Functional Requirements (Yêu cầu có chức năng)
## 1. Tổng quan hệ thống
Hệ thống **Telegram File Sharing Bot** cho phép người dùng **tải lên, chia sẻ, thu hồi quyền truy cập, và tải xuống tệp tin** thông qua Telegram Bot.
Bot giao tiếp với Backend API để lưu trữ, quản lý quyền và tạo link chia sẻ.

---
## 2. Các chức năng (Features)
### 2.1. Quản lý người dùng (User Management)
- **Xác định người dùng qua Telegram UserID**
  - Hệ thống ghi nhận Telegram UserID của người dùng khi họ bắt đầu tương tác.
  - Không yêu cầu đăng ký tài khoản.
- **Lấy thông tin người dùng**
  - Lệnh: /me
  - Bao gồm (gợi ý mẫu):
    - ID
    - Telegram ID
    - Username
    - CreatedAt (optional)
    - UpdatedAt (optional)

### 2.2. Tải lên tệp tin (Upload File)
- **Người dùng gửi file trực tiếp cho bot**
  - Hỗ trợ: document, photo, video, audio.
- **Bot tải file từ Telegram về máy chủ tạm thời**
  - Sử dụng Telegram File API.
- **Bot gửi file lên Backend API để lưu trữ**
  - Bao gồm: metadata, kích thước, loại file, người sở hữu.
- **Bot trả về link chia sẻ hoặc mã file cho người dùng**.

### 2.3. Chia sẻ tệp tin (Share File)
- **Người dùng yêu cầu chia sẻ file đã upload**
  - Lệnh: /shares
- **Backend tạo link truy cập**.
- **Bot gửi link chia sẻ cho người dùng**.
- **Hỗ trợ chia sẻ nhiều chế độ** (nếu backend cho phép):
  - Công khai (public)
  - Chia sẻ riêng tư theo danh sách user

### 2.4. Thu hồi quyền truy cập (Revoke Share Link)
- **Người dùng thu hồi quyền của một file đã chia sẻ**
  - Lệnh: `/revoke <file_id>`
- **Backend cập nhật trạng thái file để vô hiệu hóa link cũ**.
- **Bot thông báo file đã được thu hồi quyền chia sẻ**.

### 2.5. Tải xuống tệp tin (Download File)
- **Người nhận sử dụng link được chia sẻ để tải file từ backend**.
- **Bot không trực tiếp gửi lại file qua Telegram**, trừ khi backend cho phép API gửi file.
- **Đảm bảo file chỉ khả dụng khi chưa bị revoke**.

### 2.6. Danh sách tệp tin của người dùng
- **Người dùng xem danh sách các file mình đã upload**
  - Lệnh: /files
  - Bao gồm (gợi ý mẫu):
    - ID
    - OwnerUserID
    - ObjectKey
    - Filename
    - Size
    - Mime
    - Status
    - CreatedAt
    - UpdatedAt
    - TelegramID
    - Username

---
## 3. Ràng buộc chức năng (Constraints)
- Chỉ chấp nhận file hợp lệ theo giới hạn Telegram.
- Backend phải hoạt động thì bot mới xử lý command (/me, /files, ...).
- Người dùng chỉ có thể thao tác với file họ sở hữu.

---
## 4. Luồng hoạt động chính (Main User Flow)
### 4.1. Upload
- Người dùng gửi file cho bot.
- Bot tải file từ Telegram.
- Bot upload file lên backend.
- Bot trả về mã file hoặc link.

### 4.2. Share
- Người dùng yêu cầu chia sẻ file bằng /shares.
- Backend tạo presigned link.
- Bot gửi link chia sẻ.

### 4.3. Revoke
- Người dùng gửi `/revoke <file_id>`.
- Backend hủy quyền truy cập.
- Bot xác nhận file đã được thu hồi.

### 4.4. Download
- Người nhận truy cập link.
- Backend phục vụ file cho người dùng.

---
## 5. Kết luận
Tài liệu này liệt kê toàn bộ yêu cầu chức năng cần thiết để xây dựng hệ thống Telegram File Sharing Bot phục vụ cho việc **tải lên, chia sẻ, thu hồi và tải xuống tệp tin**.

---
## 6 Danh sách chức năng theo Endpoint
Dưới đây là danh sách các chức năng được dẫn từ các endpoint của backend:

### **Nhóm: Người gửi (Uploader)**
| Endpoint                    | Vai trò                                 | Người dùng |
| --------------------------- | --------------------------------------- | ---------- |
| POST /v1/files              | Khởi tạo upload, trả URL hoặc nhận file | Người gửi  |
| POST /v1/files/:id/complete | Báo hoàn tất upload                     | Người gửi  |
| GET /v1/files               | Liệt kê file người dùng                 | Người gửi  |
| POST /v1/shares             | Tạo chia sẻ mới                         | Người gửi  |
| GET /v1/shares              | Liệt kê các chia sẻ của người gửi       | Người gửi  |
| POST /v1/shares/:id/revoke  | Thu hồi link                            | Người gửi  |

### **Nhóm: Người nhận (Receiver)**
| Endpoint                      | Vai trò                        | Người dùng |
| ----------------------------- | ------------------------------ | ---------- |
| GET /v1/shares/:id            | Lấy metadata của share         | Người nhận |
| POST /v1/shares/:id/authorize | Kiểm tra password/TOTP         | Người nhận |
| GET /v1/shares/:id/download   | Lấy link tải file (sau verify) | Người nhận |

### **Nhóm: Xác thực Telegram**
| Endpoint   | Vai trò                     | Người dùng  |
| ---------- | --------------------------- | ----------- |
| GET /v1/me | Lấy thông tin Telegram user | FE xác thực |
