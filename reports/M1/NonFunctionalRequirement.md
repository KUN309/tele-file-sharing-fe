# Non-Functional Requirements (NFR)

Bảng dưới đây liệt kê các yêu cầu phi chức năng (NFR) của hệ thống. Các chỉ số này đóng vai trò là cam kết về chất lượng dịch vụ (QoS) của Bot đối với người dùng.

| **Loại yêu cầu**          | **Mô tả chi tiết** |
|---------------------------|---------------------|
| **Hiệu năng (Performance)** | Bot phải phản hồi tin nhắn Telegram trong ≤2s kể từ khi nhận. API Backend phản hồi request (tạo link, check pass) trong ≤ 1s (không tính thời gian tải file). |
| **Độ sẵn sàng (Availability)** | Hệ thống phải hoạt động ≥99% uptime, hỗ trợ restart tự động qua Docker. |
| **Khả năng mở rộng (Scalability)** | Có thể triển khai nhiều bot instance bằng `docker-compose scale` hoặc Kubernetes mà không cần chỉnh code. |
| **Khả năng phục hồi (Reliability)** | Nếu BE tạm thời không phản hồi, bot sẽ retry tối đa 3 lần với độ trễ exponential backoff. |
| **Bảo mật (Security)** | Không lưu thông tin nhạy cảm (token, API key) trong source; tất cả biến môi trường được mã hóa bằng SOPS. Giới hạn 20 req/phút/user để chống Spam. |
| **Ràng buộc tài nguyên (Constraints)** | Giới hạn upload tối đa 100MB/file (do giới hạn của Bot API và băng thông). File hết hạn hoặc bị thu hồi sẽ được xóa vĩnh viễn sau 24-48h. |
| **Khả năng bảo trì (Maintainability)** | Code tuân thủ Go convention, có linter CI (golangci-lint), và test coverage ≥70%. |
| **Khả năng triển khai (Deployability)** | Triển khai bằng Docker Compose. Chỉ cần chạy `docker compose up -d` để khởi chạy bot. |
| **Khả năng tương thích (Compatibility)** | Bot tương thích với Telegram API ≥v6; Backend API tuân theo OpenAPI spec v1. |
| **Khả năng giám sát (Observability)** | Ghi log mức INFO và ERROR, xuất ra file hoặc stdout để tích hợp CI/CD. |
| **Dễ sử dụng (Usability)** | Bot có command dễ hiểu, hỗ trợ hướng dẫn tự động `/help`, phản hồi rõ ràng bằng Markdown. |
