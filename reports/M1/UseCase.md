| **Field** | **Content** |
|-------|---------|
| **ID** | UC_01 |
| **Name** | Upload File |
| **Description** | Upload file lên hệ thống |
| **Actor** | Sender (Primary), Telegram Bot, Backend System, Storage Service |
| **Precondition** | Bot đã được khởi động và Sender có file muốn gửi |
| **Postcondition** | File được lưu trữ và sẵn sàng để gửi |
| **Trigger** | Không có |
| **Normal Flow** | B1: Sender gọi /upload. |
|                 | B2: Bot xử lý lệnh và gọi API Request Upload URL. |
|                 | B3: Bot gọi POST /v1/files. |
|                 | B4: Backend tạo bản ghi file tạm thời + presigned URL. |
|                 | B5.1: Bot nhận link upload. |
|                 | B5.2: Bot tạo cấu hình upload (thời hạn 15 phút, tối đa 100MB, pdf/doc/png...). |
|                 | B5.3: Bot gửi presigned URL cho Sender. |
|                 | B6: Sender click URL để upload file. |
|                 | B7: Bot gọi PUT upload file lên Storage Service. |
|                 | B8: Sender ấn nút xác nhận. |
|                 | B9: Bot gọi POST /v1/files/:id/complete. |
|                 | B10: Backend đánh dấu hoàn tất. |
|                 | B11: Bot hiển thị thông báo upload thành công. |
| **Alternative Flow** | **Trường hợp link hết hạn (B6):** |
|                      | B6.1: Sender click vào link hết hạn. |
|                      | B6.2: Bot hiển thị "Link đã hết hiệu lực". |
|                      | B6.3: Sender quay lại bước B1. |
|                      | **Trường hợp file vượt quá giới hạn (B7):** |
|                      | B7.1: Sender chọn file lớn hơn 100MB. |
|                      | B7.2: Storage Service từ chối lưu trữ. |
|                      | B7.3: Bot hiển thị "File quá lớn, vui lòng chọn file nhỏ hơn". |
|                      | B7.4: Sender quay lại B1. |
|                      | **Trường hợp mất kết nối khi upload (B9):** |
|                      | B9.1: Sender bị mất kết nối trong lúc upload. |
|                      | B9.2: Sender chọn reload. |
|                      | B9.3: Backend cleanup sau 5 phút nếu không hoàn tất. |
|                      | B9.4: Sender quay lại bước B6. |
| **Exception Flow** | Không có |
