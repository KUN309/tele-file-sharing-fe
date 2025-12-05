# Bài báo cáo lần 1
## 1. Requirements (Yêu cầu):
- Functional (Có chức năng)
> Thông tin xem [tại đây](FunctionalRequirement.md)  
- Non-functional (Phi chức năng)
> Thông tin xem [tại đây](NonFunctionalRequirement.md)  

## 2. Use cases (Nghiệp vụ):
Gồm có 4 workflow chính, đó là:
- Upload (Tải lên)
- Share (Chia sẻ)
- Revoke (Thu hồi)
- Download (Tải xuống)
> Thông tin xem [tại đây](UseCase.md)

## 3. API Specification (Đặc tả API)
**Tổng quan nội dung**:  
Dựa theo đặc tả backend API, file openapi.yaml mô tả toàn bộ đặc tả của File Sharing API, được sử dụng trong hệ thống chia sẻ file qua Telegram.
Tài liệu này tuân theo chuẩn OpenAPI 3.0.0, cung cấp cấu trúc đầy đủ của API để phục vụ phát triển backend, frontend, cũng như hỗ trợ sinh tài liệu tự động hoặc tạo client SDK.

### 3.1. Thông tin chung của API
**Tên API**: File Sharing API  
**Mô tả**: API cho dịch vụ chia sẻ file, bao gồm tải lên, chia sẻ, thu hồi và tải file.  
**Phiên bản**: 1.0.0-rc.1   
Tất cả các endpoint xác thực người dùng thông qua header Telegram:
- X-Telegram-User-Id
- X-Telegram-Username  
Middleware tạo user tự động nếu user chưa tồn tại.

### 3.2. Cấu trúc chức năng chính
**Ghi chú**: Chưa nhất quán giữa /api và /v1  
API được chia thành 4 nhóm chức năng lớn:
#### (A) Users
- `GET /api/me`
Lấy thông tin người dùng hiện tại dựa trên Telegram headers.
Nếu user chưa tồn tại → tạo mới.

#### (B) Files
Bao gồm toàn bộ vòng đời upload file:
- `POST /v1/files`
Khởi tạo upload → tạo file pending → trả upload URL hoặc nhận file trực tiếp.

- `GET /v1/files`
Liệt kê danh sách file thuộc user.

- `POST /v1/files/{file_id}/report-complete`
Người dùng báo cáo hoàn tất upload (thành công hoặc thất bại).

- `GET /v1/files/{file_id}/report`
Lấy báo cáo upload mới nhất của 1 file.

- `GET /v1/upload-reports`
Liệt kê toàn bộ báo cáo upload của user (có phân trang).

#### (C) Shares
- `GET /v1/shares`
Liệt kê các link chia sẻ do user tạo.

- `POST /v1/shares`
Tạo link chia sẻ mới cho file (có thể đặt mật khẩu, đặt hạn dùng).

- `POST /v1/shares/{id}/revoke`
Thu hồi link chia sẻ.

- `GET /v1/shares/{id}`
Lấy metadata của một share (không cần header Telegram trong bản này).

- `GET /v1/shares/{id}/download`
Tải file đã được chia sẻ (có kiểm tra: link còn hạn, chưa bị revoke, nếu có mật khẩu phải authorize trước).

#### (D) Authorization
**Chú thích**: Dành cho các link chia sẻ có mật khẩu.
- `POST /v1/shares/{id}/authorize`
Xác thực mật khẩu → trả về access token tạm thời để được phép tải file.

### 3.3. Các thành phần được định nghĩa
File YAML cũng mô tả toàn bộ schema dữ liệu:
- **User**: `User, UserMinimal`  
Thông tin người dùng: id, telegram_id, username, timestamps.
- **File**: `File, FileInitResponse, FileMetadata`  
Mô tả file, trạng thái, object key, metadata upload.
- **Upload Reports**: `ReportUploadCompleteRequest, ReportUploadCompleteResponse, UploadReport`  
Dùng để ghi nhận thông tin upload thành công hoặc thất bại.
- **Share**: `Share, CreateShareRequest, ShareRevokeResponse, ShareMetadataResponse`  
Mô tả link chia sẻ, trạng thái, bảo mật, mật khẩu, hạn dùng.
- **Error**:  
Schema phản hồi lỗi chuẩn.  

### 3.4. Cách sử dụng spec này
File openapi.yaml có thể sử dụng cho:
- Sinh tài liệu API tự động (Swagger UI, Redoc).
- Sinh code client → Go / TS / Python / Java…
- Kiểm thử API qua Postman/Insomnia bằng import spec.
- Chuẩn hóa giao tiếp giữa FE ↔ BE ↔ Bot Telegram.
> Thông tin xem [tại đây](https://github.com/dath-251-thuanle/tele-file-sharing-be/blob/MILESTONE-M1/api/openapi.yaml)
