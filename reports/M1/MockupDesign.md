# Mockup Design & User Flows 
A. Tổng quan
 Mockup mô phỏng giao diện FileShareBot dựa trên UI Telegram. Bộ thiết kế gồm 12 khung thể hiện đầy đủ luồng người dùng:
- Owner: Upload file, xem file, chia sẻ có mật khẩu hoặc công khai, thu hồi link.
- Receiver: Truy cập link, nhập mật khẩu (nếu có), tải file.
- Edge Cases: Link bị thu hồi, truy cập bị từ chối.

B. User Flows

1. Flow /start
User: nhập /start
FE: gọi API: GET /v1/api/me
Bot:
+ Gửi lời chào và danh sách các lệnh

2. Flow /upload
User: nhập /upload
FE gọi API: gọi API POST /v1/files
BE: tạo presigned URL
Bot:
+ Nhận presigned url và đồng thời trả về format upload (hiệu lực
15p, kích thước tối đa 100MB, kiểu file hỗ trợ pdf, doc, png,...)

3. Flow /upload complete
User: push file lên
Bot:
+ Hiện thị đang kiểm tra file
+ TH:
- file đúng với format → Bot hiện thị file hợp lệ
- file sai với format → Bot hiện thị vui lòng upload file khác
User: nhập /upload complete
FE: gọi API POST /v1/files/:id/complete
Bot: Upload thành công, bạn muốn hiện thực gì tiếp theo, ví dụ: /myfiles
hoặc /share

4. Flow /myfiles
User: nhập /myfiles
FE: gọi API GET /v1/files
Bot: nhận danh sách các file và hiển thị chúng

5. Flow /share
User: nhập /share
TH: file đã có trong upload
User: chọn file muốn chia sẻ và nhập pass(nếu có)
Bot: gửi file và metadata cho BE
BE: gọi API POST /v1/shares trả về 1 link chia sẻ cho Bot
Bot: gửi link cho User
User: chia sẻ cho người khác
TH: chưa có file
User: gửi file mình muốn chia sẻ
Bot: kiểm tra file và trả kết quả nếu hợp lệ thì gửi file và metadata cho BE
BE: gọi API POST /v1/shares trả về 1 link chia sẻ cho Bot
Bot: gửi link cho User

6. Flow /myshares
User: nhập /myshares
FE: gọi API GET /v1/shares
Bot:
+ Hiển thị các link đã chia sẻ
+ Hiển thị gợi ý:” Để thu hồi các link đã chia sẻ bạn có thể sử dụng
lệnh /revoke”

7. low /revoke
UserF: nhập /revoke
Bot: Hiển thị “Hãy gửi link bạn muốn thu hồi vào”
FE: gọi API POST /v1/shares/:id/revoke
TH: bình thường
User: nhập pass(nếu có)
Bot: Hiển thị “ Đã thu hồi thành công”
TH: link đã được thu hồi
User: nhập pass(nếu có)
Bot: Hiển thị “Link đã được thu hồi rồi, vui lòng chọn link khác.”
TH: link không hợp lệ
Bot: Hiện thị: “Link không đúng hoặc bạn không phải chủ sở hữu của
link này, vui lòng nhập lại link”

8. Flow tải file (metadata)
FE: gọi API GET /v1/shares/:id
Bot: Kiểm tra metadata
TH: link đã bị thu hồi
Bot: Link hiện không tồn tại
TH: Link đã hết hạn
Bot: Hiển thị link đã hết hạn
TH: Link yêu cầu mật khẩu
Bot: trả về:
- Tên file
- Dung lượng
- Có bảo mật
TH: Không có yêu cầu mật khẩu
FE: gọi API GET /v1/shares/:id/download để lấy Presigned URL
Bot: Hiển thị link tải cho User

9. Flow password + download
TH: link có yêu cầu pass
Bot: Hiển thị “Vui lòng nhập password”
TH: Đúng
FE: gọi API GET /v1/shares/:id/download để lấy Presigned URL
Bot: Hiển thị link tải cho User
TH: Sai
Bot: Hiển thị” Mật khẩu không đúng, bạn còn 3 lần nhập,vui lòng nhập
lại” Nếu cả 3 lần đều nhập sai, Bot hiển thị bạn không đã hết lần thử vui
lòng đợi sau 5p, 10p,...
