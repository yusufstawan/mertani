# Mertani Backend Technical Test

REST API untuk manajemen data Device dan Sensor IoT. Project ini menggunakan Go, Echo, PostgreSQL, Docker Compose, dan struktur sederhana berbasis module feature.

## Quick Start

```bash
cp .env.example .env
make compose-up
make migrate-up
make run
```

API berjalan di `http://localhost:8080` dan Swagger UI tersedia di `http://localhost:8080/swagger`.

## Struktur Singkat

```text
cmd/api              Entry point aplikasi API
cmd/migrate          Entry point migration runner
internal/app         Server, middleware, route registration, OpenAPI route
internal/device      Entity, DTO, repository, service, handler Device
internal/sensor      Entity, DTO, repository, service, handler Sensor
internal/shared      Response, error mapping, validator, ID helper
migrations           SQL migration
docs                 OpenAPI dan dokumen studi kasus
```

## Environment

Buat file `.env` dari contoh:

```bash
cp .env.example .env
```

## Menjalankan Database

Start PostgreSQL dan pgweb:

```bash
make compose-up
```

pgweb dapat dibuka di:

```text
http://localhost:8081
```

## Menjalankan Migration

Apply semua migration:

```bash
make migrate-up
```

## Menjalankan Aplikasi

Pastikan database sudah berjalan dan migration sudah diterapkan, lalu jalankan:

```bash
make run
```

API akan berjalan di:

```text
http://localhost:8080
```

Swagger UI:

```text
http://localhost:8080/swagger
```

## Endpoint Device

| Method   | Endpoint              | Keterangan                                  |
| -------- | --------------------- | ------------------------------------------- |
| `POST`   | `/api/v1/devices`     | Membuat device.                             |
| `GET`    | `/api/v1/devices`     | Mengambil daftar device dengan pagination.  |
| `GET`    | `/api/v1/devices/:id` | Mengambil detail device.                    |
| `PATCH`  | `/api/v1/devices/:id` | Mengubah sebagian field device.             |
| `DELETE` | `/api/v1/devices/:id` | Menghapus device.                           |

Query list: `page`, `limit`, `search`. Contoh: `/api/v1/devices?page=1&limit=10&search=greenhouse`.

## Endpoint Sensor

| Method   | Endpoint              | Keterangan                                  |
| -------- | --------------------- | ------------------------------------------- |
| `POST`   | `/api/v1/sensors`     | Membuat sensor.                             |
| `GET`    | `/api/v1/sensors`     | Mengambil daftar sensor dengan pagination.  |
| `GET`    | `/api/v1/sensors/:id` | Mengambil detail sensor.                    |
| `PATCH`  | `/api/v1/sensors/:id` | Mengubah sebagian field sensor.             |
| `DELETE` | `/api/v1/sensors/:id` | Menghapus sensor.                           |

Query list: `page`, `limit`, `search`. Contoh: `/api/v1/sensors?page=1&limit=10&search=temperature`.
Untuk `PATCH /api/v1/sensors/:id`, `device_id` tidak wajib dikirim. Kirim `device_id` hanya jika sensor ingin dipindahkan ke device lain.

## Makefile Commands

Lihat semua command yang tersedia:

```bash
make help
```
