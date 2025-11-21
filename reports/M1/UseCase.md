# User Flow
## Nội dung:
Hệ thống hỗ trợ một chu trình chia sẻ tệp tin an toàn và tiện lợi thông qua Telegram Bot.
Toàn bộ trải nghiệm của người dùng xoay quanh bốn luồng thao tác chính: **Upload**, **Share**, **Revoke** và **Download**.
Mỗi luồng thể hiện một bước trong hành trình quản lý và chia sẻ tệp, đảm bảo đáp ứng các yêu cầu về bảo mật, tốc độ và tính minh bạch của hệ thống.  
1. [**Upload**](#upload-file):
Người dùng gửi tệp tin hoặc đường dẫn đến Telegram Bot.
Bot xử lý thông tin, thực hiện kiểm tra kích thước và định dạng, sau đó tải tệp lên hệ thống Backend.
File được lưu trữ an toàn và sẵn sàng cho các thao tác tiếp theo.

2. [**Share**](#share-file):
Sau khi upload thành công, người dùng yêu cầu bot tạo liên kết chia sẻ.
Bot gửi metadata lên BE để tạo Share Record với các tùy chọn như thời gian hiệu lực, mật khẩu hoặc TOTP.
Kết quả là một đường dẫn bảo mật được gửi ngược lại cho người dùng.

3. [**Revoke**](#revoke-link):
Người dùng có thể thu hồi quyền truy cập một liên kết chia sẻ khi không còn muốn người khác tải xuống.
Bot gọi API revoke trên BE, đánh dấu share là invalid (hoặc revoked là true) và phản hồi cho người dùng biết rằng liên kết đã bị vô hiệu hóa.

4. [**Download**](#download-file):
Khi người nhận có liên kết chia sẻ, họ truy cập backend (thông qua bot hoặc web).
Nếu cần mật khẩu, hệ thống xác thực trước khi cấp token tạm thời để tải nội dung file.
File được trả về dưới dạng binary stream an toàn.

Nhờ thiết kế này, toàn bộ vòng đời chia sẻ tệp được tối ưu hóa — từ khi người dùng gửi tệp, chia sẻ với người khác, cho đến khi thu hồi hoặc tải xuống.
Mỗi luồng được xây dựng để đảm bảo **tính đơn giản, bảo mật, và ổn định**.

## Upload File
### Nội dung:

## Share File
### Nội dung:
Đảm nhiệm bởi Nguyên Ngọc

|     **Field**      |                                                                 **Content**                                                                                                     |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **ID**             | UC_02                                                                                                                                                                           |
| **Name**           | Share File                                                                                                                                                                      |
| **Description**    | Sender gửi file hoặc đường dẫn đến Telegram Bot để chia sẻ. Bot gửi dữ liệu lên hệ thống BE để lưu trữ và sinh liên kết chia sẻ bảo mật, sau đó phản hồi lại cho Sender.        |
| **Actor**          | Sender (Primary), Telegram Bot, Backend Service                                                                                                                                 |
|**Preconditions**   | B1: Bot hoạt động và kết nối được với Telegram API.                                                                                                                             |
|                    | B2: BE Service khả dụng (VD: API `/api/v1/share` online).                                                                                                                       |
|                    | B3: Sender đã bắt đầu phiên làm việc bằng lệnh `/start`.                                                                                                                        |
|**Postconditions**  | B1: File được tải lên và lưu trữ thành công trên BE.                                                                                                                            |
|                    | B2: Hệ thống tạo link chia sẻ.                                                                                                                                                  |
|                    | B3: Sender nhận link chia sẻ qua Telegram.                                                                                                                                      |
| **Triggers**       | Sender gửi file hoặc dùng lệnh `/share` trong Telegram.                                                                                                                         |
| **Normal Flow**    | B1: Sender gửi file hoặc lệnh `/share` cho bot.                                                                                                                                 |
|                    | B2: Bot xác định loại dữ liệu (file hoặc link).                                                                                                                                 |
|                    | B3: Bot kiểm tra dung lượng và định dạng hợp lệ.                                                                                                                                |
|                    | B4: Bot gửi file + metadata lên BE qua API `/api/v1/share`.                                                                                                                     |
|                    | B5: BE lưu file, sinh link và trả phản hồi.                                                                                                                                     |
|                    | B6: Bot gửi lại link cho Sender.                                                                                                                                                |
|                    | B7: Sender chia sẻ link với người khác.                                                                                                                                         |
|**Alternative Flow**| **Trường hợp file quá lớn (ở bước 3):**                                                                                                                                         |
|                    | B3.1: File vượt quá giới hạn (vd: 100MB).                                                                                                                                       |
|                    | B3.2: Bot báo *"File quá lớn, vui lòng chọn file nhỏ hơn."*                                                                                                                     |
|                    | **Sender muốn tăng bảo mật (ở bước 6):**                                                                                                                                        |
|                    | B6.1: Sender muốn đặt mật khẩu hoặc TOTP.                                                                                                                                       |
|                    | B6.2: Bot hỏi thêm thông tin và gọi API thiết lập bảo mật.                                                                                                                      |
| **Exception Flow** | **File sai định dạng (ở bước 3):**                                                                                                                                              |
|                    | B3.1: File sai định dạng.                                                                                                                                                       |
|                    | B3.2: Bot báo *"Định dạng file không được hỗ trợ."*                                                                                                                             |
|                    | **BE không phản hồi hoặc lỗi server (ở bước 4):**                                                                                                                               |
|                    | B4a.1: BE phản hồi chậm hoặc treo.                                                                                                                                              |
|                    | B4a.2: Bot báo *"Server đang bận, thử lại sau."*                                                                                                                                |
|                    | **Lỗi mạng (ở bước 4):**                                                                                                                                                        |
|                    | B4b.1: Mạng gián đoạn.                                                                                                                                                          |
|                    | B4b.2: Bot lưu tạm và gửi lại khi mạng ổn định.                                                                                                                                 |

## Revoke Link
### Nội dung:

## Download File
### Nội dung
