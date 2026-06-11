# World Cup Predictor

Web app tự chơi đoán tỉ số các trận đấu World Cup 2026. Đăng ký/đăng nhập, dự đoán
tỉ số từng trận, và xem bảng xếp hạng (leaderboard) ai đoán đúng nhiều nhất.

## Tech stack

- **Frontend**: Next.js (App Router), chạy trên cổng `3000`
- **Backend**: Go (chuẩn `net/http`, không dùng framework ngoài), JWT auth, cổng `8080`
- **Database**: PostgreSQL 16
- Tất cả chạy qua **Docker Compose**

## Cách chạy

```bash
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432 (user `wcuser` / pass `wcpass` / db `worldcup`)

Lần chạy đầu tiên, Postgres sẽ tự chạy `db/init.sql` để tạo bảng và seed dữ liệu
(đội bóng + lịch vòng bảng).

## Tài khoản admin mặc định

Backend tự tạo 1 tài khoản admin khi khởi động lần đầu (cấu hình trong
`docker-compose.yml`):

- Username: `admin`
- Password: `admin123`

> ⚠️ Hãy đổi `ADMIN_USERNAME` / `ADMIN_PASSWORD` và `JWT_SECRET` trong
> `docker-compose.yml` trước khi public ra ngoài.

## Cách tính điểm

Khi admin nhập kết quả cuối cùng của 1 trận đấu:

- Đoán **đúng tỉ số chính xác** → **3 điểm**
- Đoán **đúng kết quả** (thắng/thua/hòa) nhưng sai tỉ số → **1 điểm**
- Sai cả hai → **0 điểm**

## Cập nhật lịch thi đấu / dữ liệu đội bóng

Dữ liệu seed trong `db/init.sql` hiện có đầy đủ tên đội cho các bảng **A, B, D, G, I**
(lấy từ kết quả bốc thăm thực tế). Các bảng còn lại (**C, E, F, H, J, K, L**) đang để
tên placeholder dạng `TBD C1`, `TBD C2`, ... — bạn có thể:

1. Sửa trực tiếp trong `db/init.sql` rồi `docker compose down -v && docker compose up --build`
   (xóa volume để chạy lại seed), hoặc
2. Sửa trực tiếp trong bảng `teams` của Postgres bằng `psql`/DBeaver sau khi đã chạy.

### Thêm trận đấu mới (vòng knock-out, v.v.) — cần token admin

```bash
# 1. Đăng nhập lấy token
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'

# 2. Tạo trận đấu mới
curl -X POST http://localhost:8080/api/admin/matches \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "stage": "round_of_32",
    "group_name": "",
    "home_team_id": 1,
    "away_team_id": 5,
    "match_time": "2026-06-30T18:00:00Z"
  }'
```

### Nhập kết quả trận đấu (để chấm điểm dự đoán)

```bash
curl -X PUT http://localhost:8080/api/admin/matches/1/result \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"home_score": 2, "away_score": 1}'
```

Sau khi gọi API này, hệ thống sẽ tự động chấm điểm cho tất cả dự đoán của trận đó
và cập nhật bảng xếp hạng.

## API tổng quan

| Method | Endpoint                          | Auth        | Mô tả                              |
|--------|------------------------------------|-------------|--------------------------------------|
| POST   | `/api/register`                   | -           | Đăng ký tài khoản                   |
| POST   | `/api/login`                      | -           | Đăng nhập, trả về JWT                |
| GET    | `/api/me`                         | user        | Thông tin tài khoản hiện tại         |
| GET    | `/api/teams`                      | -           | Danh sách đội bóng                   |
| GET    | `/api/matches`                    | optional    | Danh sách trận đấu (kèm dự đoán của bạn nếu đăng nhập) |
| POST   | `/api/predictions`                | user        | Lưu/sửa dự đoán cho 1 trận           |
| GET    | `/api/predictions/me`             | user        | Danh sách dự đoán của bạn            |
| GET    | `/api/leaderboard`                | -           | Bảng xếp hạng                        |
| POST   | `/api/admin/matches`              | admin       | Tạo trận đấu mới                     |
| PUT    | `/api/admin/matches/{id}/result`  | admin       | Nhập kết quả & chấm điểm             |

## Cấu trúc thư mục

```
.
├── docker-compose.yml
├── db/
│   └── init.sql        # schema + seed data
├── backend/             # Go API
│   ├── Dockerfile
│   ├── go.mod
│   └── *.go
└── frontend/             # Next.js app
    ├── Dockerfile
    ├── package.json
    └── app/
```
