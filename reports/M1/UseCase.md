
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
Đảm nhiệm bởi Cung Bùi

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
Đảm nhiệm bởi Trí Thành

|     **Field**      |                                                                 **Content**                                                                                                     |
|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **ID**             | UC-REVOKE                                                                                                                                                                  |
| **Name**           | Revoke Share Link (Thu hồi link chia sẻ)                                                                                                                                        |
| **Description**    | Sender gửi yêu cầu thu hồi link chia sẻ qua Telegram Bot. Bot chuyển tiếp yêu cầu đến Backend API để xác thực, kiểm tra quyền sở hữu và đánh dấu link bị thu hồi bằng cách cập nhật trường `revoked_at`. |
| **Actor**          | Sender (Primary), Telegram Bot, Backend Service                                                                                                                                |
| **Preconditions**  | B1: Sender đã được xác thực trong hệ thống (tồn tại user có `telegram_user_id`).                                                                                                |
|                    | B2: Bản ghi share có `share_id` tồn tại trong bảng `shares`.                                                                                                                    |
|                    | B3: Link đang hoạt động (`revoked_at IS NULL`).                                                                                                                                 |
|                    | B4: Sender là chủ sở hữu của link (`share.owner_user_id = user.id`).                                                                                                            |
| **Postconditions** | B1: Link được đánh dấu thu hồi (`revoked_at` có timestamp).                                                                                                                     |
|                    | B2: Người nhận không thể truy cập link nữa.                                                                                                                                      |
|                    | B3: Mọi request truy cập sau đó trả lỗi 410 hoặc 404.                                                                                                                            |
| **Triggers**       | Sender dùng lệnh `/myshare` → `/revoke` trên Telegram.                                                                                                                          |
| **Normal Flow**    | B1: Sender thực hiện lệnh `/revoke`.                                                                                                                                            |
|                    | B2: Bot gửi yêu cầu thu hồi kèm thông tin xác thực đến Backend.                                                                                                                 |
|                    | B3: Backend xác thực API Key và định danh user từ Telegram ID.                                                                                                                   |
|                    | B4: Backend kiểm tra quyền sở hữu và trạng thái link.                                                                                                                            |
|                    | B5: Nếu hợp lệ, Backend cập nhật DB để đánh dấu link đã thu hồi.                                                                                                                |
|                    | B6: Bot nhận phản hồi và thông báo cho Sender.                                                                                                                                   |
| **Alternative Flow** | Không có.                                                                                                                                                                     |
| **Exception Flow** | **Link không tồn tại (ở bước 3):**                                                                                                                                               |
|                    | B3a.1: Backend không tìm thấy bản ghi share.                                                                                                                                     |
|                    | B3a.2: Bot báo "Link không tồn tại."                                                                                                                                             |
|                    | **Người gửi không có quyền (ở bước 4):**                                                                                                                                         |
|                    | B4a.1: Sender không phải chủ sở hữu.                                                                                                                                             |
|                    | B4a.2: Bot báo "Bạn không có quyền thu hồi link này."                                                                                                                            |
|                    | **Link đã được thu hồi trước đó (ở bước 4):**                                                                                                                                    |
|                    | B4b.1: `revoked_at` đã có giá trị.                                                                                                                                               |
|                    | B4b.2: Bot báo "Link này đã được thu hồi trước đó."                                                                                                                              |
|                    | **Xác thực thất bại (ở bước 3):**                                                                                                                                                 |
|                    | B3b.1: API Key hoặc thông tin xác thực không hợp lệ.                                                                                                                             |
|                    | B3b.2: Bot báo "Không thể xác thực yêu cầu."                                                                                                                                     |

## Download File
### Nội dung:
Đảm nhiệm bởi Hùng Dũng

| **Field**          | **Content** |
|--------------------|-------------|
| **ID**             | UC_04 |
| **Name**           | Download File |
| **Description**    | Receiver truy cập vào một link chia sẻ, xác thực (nếu được yêu cầu), và tải tệp tin về thiết bị của mình thông qua Telegram Bot. |
| **Actor**          | Receiver (Primary), Telegram Bot, Backend Service |
| **Preconditions**  | B1: Sender đã upload tệp tin thành công.<br> B2: Sender đã tạo một link chia sẻ trỏ tới tệp tin đó.<br> B3: Receiver đã nhận được link chia sẻ. |
| **Postconditions** | B1: Receiver nhận được tệp tin qua tin nhắn Telegram.<br> B2: Một bản ghi truy cập thành công được lưu vào.<br> B3: Bộ đếm của link chia sẻ được tăng lên 1. |
| **Triggers**       | Receiver nhấn vào một link chia sẻ của bot. |
| **Normal Flow**    | B1: Receiver nhấn vào link chia sẻ.<br> B2: Bot nhận diện link, gọi BE để kiểm tra link hợp lệ, còn hạn và người này có trong danh sách được nhận không.<br> B3: Bot gửi yêu cầu nhập mật khẩu/TOTP.<br> B4: Receiver nhập và gửi mật khẩu/TOTP.<br> B5: Bot gửi toàn bộ thông tin xác thực đến BE.<br> B6: BE xác nhận mọi thứ hợp lệ và gửi lại access token tạm thời.<br> B7: Bot dùng token để lấy file từ BE và gửi cho Receiver. |
| **Alternative Flow** | **Trường hợp link công khai (không yêu cầu mật khẩu/TOTP):**<br> B2.1: Bot kiểm tra và thấy link không yêu cầu mật khẩu/TOTP.<br> B2.2: Bot gửi file cho Receiver. |
| **Exception Flow** | **Link hết hạn hoặc bị thu hồi:**<br> B2.1: BE phản hồi lỗi khi Bot kiểm tra link.<br> B2.2: Bot hiển thị thông báo link không còn tồn tại.<br><br> **Receiver nhập sai mật khẩu/TOTP:**<br> B6.1: BE phát hiện mật khẩu/TOTP không chính xác.<br> B6.2: Bot hiển thị thông báo mật khẩu/TOTP không đúng.<br> B6.3: Quay lại bước B3.<br><br> **Hết lượt tải (Max downloads):**<br> B6.1: BE xác thực thành công nhưng lượt tải đã hết.<br> B6.2: Bot hiển thị thông báo link đã đạt tối đa số lượt tải. |




