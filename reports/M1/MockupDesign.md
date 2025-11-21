# Mockup Design & User Flows 
## Tổng quan
### [User Flows](images/user_flow.png)
Mockup mô phỏng giao diện Telegram File Share Bot dựa trên UI Telegram.
Bộ thiết kế gồm 12 khung thể hiện đầy đủ luồng người dùng (xem bên tại mục High fidelity ở dưới):
- Owner: Upload file, xem file, chia sẻ có mật khẩu hoặc công khai, thu hồi link.
- Receiver: Truy cập link, nhập mật khẩu (nếu có), tải file.
- Edge Cases: Link bị thu hồi, truy cập bị từ chối.

## User Flows
### Low fidelity
**1. Flow /start**
- User: Nhập /start (hoặc /me)
- FE: Gọi API `GET /api/me`
- Bot: Gửi lời chào và danh sách các lệnh

**2. Flow /upload**
- User: Nhập /upload
- FE: Gọi API `POST /v1/files`
- BE: Tạo presigned URL
- Bot: Nhận presigned url và đồng thời trả về format upload (hiệu lực 15p, kích thước tối đa 100MB, kiểu file hỗ trợ pdf, doc, png,...)

**3. Flow /upload complete**
- Bot:
  - Hiện thị đang kiểm tra file
  - Trường hợp:
    - File đúng với format → Bot hiện thị file hợp lệ
    - File sai với format → Bot hiện thị vui lòng upload file khác
- User: Nhận /upload complete
- FE: Gọi API `POST /v1/files/:id/complete`
- Bot: Upload thành công, bạn muốn hiện thực gì tiếp theo, ví dụ: /myfiles hoặc /share

**4. Flow /files**
- User: Nhập /files
- FE: Gọi API `GET /v1/files`
- Bot: Bhận danh sách các file và hiển thị chúng

**5. Flow /shares**
- User: Nhập /shares
  - Trường hợp: file đã có từ upload trước đó
    - User: Chọn file muốn chia sẻ và nhập pass(nếu có)
    - Bot: Gửi file và metadata cho BE
    - BE: Gọi API `POST /v1/shares` trả về 1 link chia sẻ cho Bot
    - Bot: Gửi link cho User
    - User: chia sẻ cho người khác
  - Trường hợp: chưa có file
    - User: Gửi file mình muốn chia sẻ
    - Bot: Kiểm tra file và trả kết quả nếu hợp lệ thì gửi file và metadata cho BE
    - BE: Gọi API `POST /v1/shares` trả về 1 link chia sẻ cho Bot
    - Bot: Gửi link cho User

**6. Flow /shares/:id**
- User: Nhập /shares/:id (:id là ID của file mình đã chia sẻ) hoặc /shares
- FE: Gọi API `GET /v1/shares/:id` hoặc `GET /v1/shares`
- Bot:
  - Hiển thị link của file đã chia sẻ hoặc các link đã chia sẻ
  - Hiển thị gợi ý: ”Để thu hồi các link đã chia sẻ bạn có thể sử dụng lệnh /shares/:id/revoke (hoặc /revoke)”

**7. Flow /revoke**
- User: Nhập /shares/:id/revoke (hoặc /revoke)
- Bot: Hiển thị “Hãy gửi link bạn muốn thu hồi vào” (nếu chọn /revoke)
- FE: Gọi API `POST /v1/shares/:id/revoke`
  - Trường hợp: bình thường
    - User: Nhập pass(nếu có)
    - Bot: Hiển thị “ Đã thu hồi thành công”
  - Trường hợp: link đã được thu hồi
    - User: nhập pass(nếu có)
    - Bot: Hiển thị “Link đã được thu hồi rồi, vui lòng chọn link khác.”
  - Trường hợp: link không hợp lệ
    - Bot: Hiện thị “Link không đúng hoặc bạn không phải chủ sở hữu của link này, vui lòng nhập lại link”

**8. Flow tải file (metadata)**
- FE: Gọi API `GET /v1/shares/:id/download`
- Bot: Kiểm tra metadata
  - Trường hợp: Link đã bị thu hồi
    - Bot: Link hiện không tồn tại
  - Trường hợp: Link đã hết hạn
    - Bot: Hiển thị link đã hết hạn
  - Trường hợp: Link yêu cầu mật khẩu
    - Bot: Trả về (mẫu):
       - Tên file
       - Dung lượng
       - Có bảo mật
  - Trường hợp: Không có yêu cầu mật khẩu
    - FE: Gọi API `GET /v1/shares/:id/download` để lấy Presigned URL
    - Bot: Hiển thị link tải cho User

**9. Flow password + download**
- Trường hợp: link có yêu cầu pass
  - Bot: Hiển thị “Vui lòng nhập password”
- Trường hợp: Đúng
  - FE: gọi API `GET /v1/shares/:id/download` để lấy Presigned URL
  - Bot: Hiển thị link tải cho User
- Trường hợp: Sai
  - Bot: Hiển thị ”Mật khẩu không đúng, bạn còn 3 lần nhập, vui lòng nhập lại”.
    Nếu cả 3 lần đều nhập sai, Bot hiển thị bạn đã hết lần thử vui lòng đợi sau 5p, 10p,...

## High fidelity
> Các ảnh [xem tại thư mục images](images)
