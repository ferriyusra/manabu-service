# API Documentation - User Registration

## Endpoint

```
POST /auth/register
```

Endpoint untuk registrasi user baru dalam aplikasi Manabu (Japanese Learning App).

---

## Request

### Headers

```
Content-Type: application/json
```

### Request Body

```json
{
  "name": "string (required)",
  "username": "string (required)",
  "email": "string (required, valid email)",
  "password": "string (required)",
  "confirmPassword": "string (required)"
}
```

### Field Descriptions

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| `name` | string | Yes | Not empty | Nama lengkap user |
| `username` | string | Yes | Not empty, unique | Username untuk login |
| `email` | string | Yes | Valid email format, unique | Email user |
| `password` | string | Yes | Minimum 6 characters (recommended) | Password user |
| `confirmPassword` | string | Yes | Must match `password` | Konfirmasi password |

### Business Rules

1. **Username** harus unique (tidak boleh duplikat)
2. **Email** harus unique dan valid format
3. **Password** dan **confirmPassword** harus sama
4. **RoleID** akan di-set otomatis sebagai "student" (role default)
5. Password akan di-hash menggunakan bcrypt sebelum disimpan

---

## Response

### Success Response

**HTTP Status:** `200 OK`

```json
{
  "code": 200,
  "message": "OK",
  "data": {
    "user": {
      "uuid": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "username": "johndoe",
      "email": "john.doe@example.com",
      "role": "student"
    }
  }
}
```

### Error Responses

#### 1. Validation Error

**HTTP Status:** `422 Unprocessable Entity`

```json
{
  "code": 422,
  "message": "Unprocessable Entity",
  "data": [
    {
      "field": "email",
      "message": "email must be a valid email address"
    },
    {
      "field": "confirmPassword",
      "message": "confirmPassword is required"
    }
  ]
}
```

#### 2. Duplicate Username/Email

**HTTP Status:** `400 Bad Request`

```json
{
  "code": 400,
  "message": "Username already exists",
  "data": null
}
```

atau

```json
{
  "code": 400,
  "message": "Email already exists",
  "data": null
}
```

#### 3. Password Mismatch

**HTTP Status:** `400 Bad Request`

```json
{
  "code": 400,
  "message": "Password and confirm password do not match",
  "data": null
}
```

#### 4. Invalid Request Body

**HTTP Status:** `400 Bad Request`

```json
{
  "code": 400,
  "message": "Invalid request body",
  "data": null
}
```

---

## Example Requests

### cURL

```bash
curl -X POST http://localhost:8001/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "SecurePass123",
    "confirmPassword": "SecurePass123"
  }'
```

### JavaScript (Fetch API)

```javascript
const registerUser = async () => {
  try {
    const response = await fetch('http://localhost:8001/auth/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: 'John Doe',
        username: 'johndoe',
        email: 'john.doe@example.com',
        password: 'SecurePass123',
        confirmPassword: 'SecurePass123'
      })
    });

    const data = await response.json();

    if (response.ok) {
      console.log('Registration successful:', data.data.user);
    } else {
      console.error('Registration failed:', data.message);
    }
  } catch (error) {
    console.error('Network error:', error);
  }
};

registerUser();
```

### JavaScript (Axios)

```javascript
import axios from 'axios';

const registerUser = async () => {
  try {
    const response = await axios.post('http://localhost:8001/auth/register', {
      name: 'John Doe',
      username: 'johndoe',
      email: 'john.doe@example.com',
      password: 'SecurePass123',
      confirmPassword: 'SecurePass123'
    });

    console.log('Registration successful:', response.data.data.user);
  } catch (error) {
    if (error.response) {
      console.error('Registration failed:', error.response.data.message);
    } else {
      console.error('Network error:', error.message);
    }
  }
};

registerUser();
```

### Python (requests)

```python
import requests
import json

url = "http://localhost:8001/auth/register"
payload = {
    "name": "John Doe",
    "username": "johndoe",
    "email": "john.doe@example.com",
    "password": "SecurePass123",
    "confirmPassword": "SecurePass123"
}
headers = {
    "Content-Type": "application/json"
}

response = requests.post(url, json=payload, headers=headers)

if response.status_code == 200:
    user_data = response.json()["data"]["user"]
    print(f"Registration successful: {user_data}")
else:
    print(f"Registration failed: {response.json()['message']}")
```

### Go (net/http)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

type RegisterRequest struct {
    Name            string `json:"name"`
    Username        string `json:"username"`
    Email           string `json:"email"`
    Password        string `json:"password"`
    ConfirmPassword string `json:"confirmPassword"`
}

func main() {
    url := "http://localhost:8001/auth/register"

    payload := RegisterRequest{
        Name:            "John Doe",
        Username:        "johndoe",
        Email:           "john.doe@example.com",
        Password:        "SecurePass123",
        ConfirmPassword: "SecurePass123",
    }

    jsonData, _ := json.Marshal(payload)

    resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode == 200 {
        fmt.Println("Registration successful:", string(body))
    } else {
        fmt.Println("Registration failed:", string(body))
    }
}
```

### Postman Collection

```json
{
  "info": {
    "name": "Manabu API - Registration",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Register User",
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"John Doe\",\n  \"username\": \"johndoe\",\n  \"email\": \"john.doe@example.com\",\n  \"password\": \"SecurePass123\",\n  \"confirmPassword\": \"SecurePass123\"\n}"
        },
        "url": {
          "raw": "http://localhost:8001/auth/register",
          "protocol": "http",
          "host": ["localhost"],
          "port": "8001",
          "path": ["auth", "register"]
        }
      }
    }
  ]
}
```

---

## Test Cases

### Valid Registration

**Request:**
```json
{
  "name": "Tanaka Yuki",
  "username": "tanakayuki",
  "email": "tanaka.yuki@example.com",
  "password": "MyPassword123!",
  "confirmPassword": "MyPassword123!"
}
```

**Expected Response:** `200 OK` with user data

---

### Invalid Email Format

**Request:**
```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "invalid-email",
  "password": "SecurePass123",
  "confirmPassword": "SecurePass123"
}
```

**Expected Response:** `422 Unprocessable Entity`

---

### Password Mismatch

**Request:**
```json
{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john.doe@example.com",
  "password": "SecurePass123",
  "confirmPassword": "DifferentPassword"
}
```

**Expected Response:** `400 Bad Request` - Password mismatch

---

### Missing Required Fields

**Request:**
```json
{
  "name": "John Doe",
  "email": "john.doe@example.com"
}
```

**Expected Response:** `422 Unprocessable Entity`

---

### Duplicate Username

**Request:**
```json
{
  "name": "Jane Doe",
  "username": "johndoe",
  "email": "jane.doe@example.com",
  "password": "SecurePass123",
  "confirmPassword": "SecurePass123"
}
```

**Expected Response:** `400 Bad Request` - Username already exists (if `johndoe` already registered)

---

## Notes

1. **Password Security:**
   - Password di-hash menggunakan bcrypt
   - Tidak disimpan dalam bentuk plain text
   - Minimum panjang password: 6 karakter (disarankan lebih panjang)

2. **Default Role:**
   - Semua user baru otomatis mendapat role "student"
   - Role admin harus di-set manual melalui database

3. **UUID:**
   - Setiap user mendapat UUID unik
   - UUID digunakan untuk public-facing ID
   - Mencegah enumeration attack

4. **Response:**
   - Password TIDAK dikembalikan dalam response
   - Token TIDAK dikembalikan (user harus login setelah register)

5. **Rate Limiting:**
   - Endpoint ini di-protect dengan rate limiter
   - Default: 1000 requests per 60 seconds (configurable di config.json)

---

## Related Endpoints

- `POST /auth/login` - Login user
- `GET /auth/user` - Get logged in user info
- `GET /auth/:uuid` - Get user by UUID
- `PUT /auth/:uuid` - Update user profile

---

## Implementation Details

### DTO Location
`domain/dto/user.go` - RegisterRequest struct

### Controller
`controllers/user/user.go` - Register method

### Service
`services/user/user.go` - Register business logic

### Repository
`repositories/user/user.go` - Database operations

### Model
`domain/models/user.go` - User entity
