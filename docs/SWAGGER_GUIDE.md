# Swagger/OpenAPI Documentation Guide

## Overview

Manabu API menggunakan Swagger/OpenAPI untuk dokumentasi API yang interaktif dan mudah digunakan.

---

## Akses Swagger UI

Setelah aplikasi berjalan, akses Swagger UI di:

```
http://localhost:8001/swagger/index.html
```

---

## Generate Swagger Documentation

### Menggunakan Makefile (Recommended)

```bash
make swagger
```

### Menggunakan Swag CLI Langsung

```bash
swag init -g cmd/main.go -o docs
```

---

## File yang Dihasilkan

Setelah menjalankan perintah generate, file-file berikut akan dibuat di folder `docs/`:

- `docs.go` - Go code untuk Swagger
- `swagger.json` - OpenAPI spec dalam format JSON
- `swagger.yaml` - OpenAPI spec dalam format YAML

---

## Cara Menambahkan Dokumentasi untuk Endpoint Baru

### 1. Tambahkan Annotation di Controller

```go
// CreateVocabulary godoc
// @Summary      Create new vocabulary
// @Description  Create a new Japanese vocabulary entry
// @Tags         Vocabulary
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateVocabularyRequest true "Vocabulary details"
// @Success      200 {object} response.Response{data=dto.VocabularyResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Failure      422 {object} response.Response
// @Router       /vocabularies [post]
func (v *VocabularyController) CreateVocabulary(ctx *gin.Context) {
    // Implementation...
}
```

### 2. Re-generate Documentation

```bash
make swagger
```

### 3. Restart Aplikasi

Restart aplikasi Anda untuk melihat perubahan di Swagger UI.

---

## Annotation Reference

### Annotation Umum

| Annotation | Deskripsi | Contoh |
|------------|-----------|--------|
| `@Summary` | Ringkasan singkat endpoint | `@Summary Create user` |
| `@Description` | Deskripsi detail endpoint | `@Description Create a new user account` |
| `@Tags` | Group endpoint | `@Tags Users` |
| `@Accept` | Content-Type yang diterima | `@Accept json` |
| `@Produce` | Content-Type response | `@Produce json` |
| `@Security` | Security scheme | `@Security BearerAuth` |
| `@Param` | Parameter | `@Param id path string true "User ID"` |
| `@Success` | Response sukses | `@Success 200 {object} dto.UserResponse` |
| `@Failure` | Response error | `@Failure 400 {object} response.Response` |
| `@Router` | Route path & method | `@Router /users [post]` |

### Parameter Types

```go
// Path parameter
// @Param id path string true "User ID"

// Query parameter
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)

// Body parameter
// @Param request body dto.CreateRequest true "Request body"

// Header parameter
// @Param Authorization header string true "Bearer token"
```

### Response Types

```go
// Simple response
// @Success 200 {object} dto.UserResponse

// Response with nested data
// @Success 200 {object} response.Response{data=dto.UserResponse}

// Array response
// @Success 200 {array} dto.VocabularyResponse

// Array in Response wrapper
// @Success 200 {object} response.Response{data=[]dto.VocabularyResponse}
```

---

## Security Definitions

Security sudah dikonfigurasi di `cmd/main.go`:

```go
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

Untuk endpoint yang memerlukan authentication, tambahkan:

```go
// @Security BearerAuth
```

---

## Testing di Swagger UI

### 1. Login untuk Mendapatkan Token

1. Buka Swagger UI
2. Cari endpoint `POST /auth/login`
3. Klik "Try it out"
4. Masukkan credentials:
   ```json
   {
     "username": "your_username",
     "password": "your_password"
   }
   ```
5. Klik "Execute"
6. Copy token dari response

### 2. Authorize dengan Token

1. Klik tombol "Authorize" di kanan atas Swagger UI
2. Masukkan token dalam format: `Bearer your_token_here`
3. Klik "Authorize"
4. Klik "Close"

### 3. Test Protected Endpoints

Sekarang Anda bisa test semua endpoint yang memerlukan authentication.

---

## Best Practices

### 1. Konsisten dengan Tags

Gunakan tags yang konsisten untuk grouping:

```go
// @Tags Authentication  // untuk login, register
// @Tags Users           // untuk user management
// @Tags Vocabulary      // untuk vocabulary management
// @Tags Courses         // untuk course management
```

### 2. Deskripsi yang Jelas

```go
// ❌ Bad
// @Summary Get user
// @Description Get user

// ✅ Good
// @Summary Get user by UUID
// @Description Retrieve detailed user information by their unique identifier
```

### 3. Dokumentasikan Semua Status Code

```go
// @Success 200 {object} response.Response{data=dto.UserResponse} "Success"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 404 {object} response.Response "User Not Found"
// @Failure 422 {object} response.Response "Validation Error"
// @Failure 500 {object} response.Response "Internal Server Error"
```

### 4. Example Values di DTO

Tambahkan example values di struct DTO:

```go
type CreateVocabularyRequest struct {
    Kanji    string `json:"kanji" example:"食べる"`
    Hiragana string `json:"hiragana" example:"たべる"`
    Romaji   string `json:"romaji" example:"taberu"`
    Meaning  string `json:"meaning" example:"to eat"`
}
```

---

## Troubleshooting

### Error: "docs" package not found

**Solusi:**
```bash
make swagger
go mod tidy
```

### Swagger UI tidak menampilkan endpoint baru

**Solusi:**
1. Pastikan annotation sudah benar
2. Re-generate docs: `make swagger`
3. Restart aplikasi

### Security tidak bekerja

**Pastikan:**
1. Security definition ada di main.go
2. Annotation `@Security BearerAuth` ada di controller
3. Token format: `Bearer <token>`

---

## Contoh Lengkap

### Controller dengan Berbagai Parameter

```go
// GetVocabularies godoc
// @Summary      Get vocabularies list
// @Description  Get paginated list of vocabularies with optional filters
// @Tags         Vocabulary
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(20)
// @Param        jlpt_level query string false "JLPT Level filter" Enums(N5, N4, N3, N2, N1)
// @Param        word_type query string false "Word type filter" Enums(noun, verb, adjective)
// @Param        search query string false "Search keyword"
// @Success      200 {object} response.Response{data=[]dto.VocabularyResponse}
// @Failure      400 {object} response.Response
// @Failure      401 {object} response.Response
// @Router       /vocabularies [get]
func (v *VocabularyController) GetVocabularies(ctx *gin.Context) {
    // Implementation...
}
```

---

## API Metadata

Metadata API dikonfigurasi di `cmd/main.go`:

```go
// @title           Manabu API - Japanese Learning Application
// @version         1.0
// @description     API documentation for Manabu Japanese Learning Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.manabu.com/support
// @contact.email  support@manabu.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1
```

Ubah nilai-nilai ini sesuai kebutuhan project Anda.

---

## Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Gin Swagger](https://github.com/swaggo/gin-swagger)

---

## Summary

- ✅ Swagger UI: `http://localhost:8001/swagger/index.html`
- ✅ Generate docs: `make swagger`
- ✅ Security: Bearer Token authentication
- ✅ Interactive testing: Login → Authorize → Test endpoints
- ✅ Auto-generated from code annotations

Happy documenting! 📚
