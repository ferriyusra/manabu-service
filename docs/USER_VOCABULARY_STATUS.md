# User Vocabulary Status API Documentation

## 📚 Table of Contents

1. [Overview](#overview)
2. [Simple Tracking Concept](#simple-tracking-concept)
3. [API Endpoints](#api-endpoints)
4. [Usage Flow](#usage-flow)
5. [Examples](#examples)
6. [Status Progression](#status-progression)
7. [Troubleshooting](#troubleshooting)

---

## Overview

User Vocabulary Status API adalah sistem manajemen pembelajaran vocabulary berbasis **Simple Progress Tracking**. Sistem ini membantu user melacak progress pembelajaran vocabulary dengan cara yang sederhana dan mudah dipahami.

### Key Features

- ✅ **Simple Progress Tracking** - Track berapa kali vocabulary berhasil dijawab benar
- ✅ **User-Controlled Review** - User sendiri yang memutuskan kapan mau review
- ✅ **Clear Completion Goal** - 5 jawaban benar = completed
- ✅ **Status Monitoring** - Monitor status pembelajaran (learning → completed)
- ✅ **Reset on Failure** - Salah jawab akan reset progress ke 0

### Business Value

- **Simplicity**: Mudah dipahami dan digunakan
- **User Control**: User punya kontrol penuh kapan mau belajar
- **Clear Goals**: Target yang jelas (5x benar = selesai)
- **Motivation**: Progress yang terlihat jelas mendorong belajar

---

## Simple Tracking Concept

### How It Works

Sistem ini menggunakan simple counter untuk tracking progress:

```
Start Learning:  ●────────────────────  repetitions = 0, status = "learning"
                 ↓ Review (correct)

Review 1:        ──●──────────────────  repetitions = 1, status = "learning"
                 ↓ Review (correct)

Review 2:        ────●────────────────  repetitions = 2, status = "learning"
                 ↓ Review (correct)

Review 3:        ──────●──────────────  repetitions = 3, status = "learning"
                 ↓ Review (correct)

Review 4:        ────────●────────────  repetitions = 4, status = "learning"
                 ↓ Review (correct)

Review 5:        ──────────●──────────  repetitions = 5, status = "completed" ✓
```

**Key Principle**: Jawab benar 5 kali berturut-turut untuk menyelesaikan vocabulary.

### Review Logic

#### Correct Answer (isCorrect = true)
```
repetitions = repetitions + 1

if repetitions >= 5:
    status = "completed"  ← Goal reached!
```

#### Incorrect Answer (isCorrect = false)
```
repetitions = 0          ← Reset to start
status = "learning"
```

### Simple & Clear

```
┌─────────────────────────────┐
│  User Reviews Vocabulary    │
└──────────┬──────────────────┘
           │
           ▼
    ┌──────────────┐
    │  Is Correct? │
    └──────┬───────┘
           │
      Yes  │  No
    ┌──────▼──┐  ┌────────┐
    │ Add +1  │  │ Reset  │
    │ to Rep  │  │ to 0   │
    └──────┬──┘  └────┬───┘
           │          │
           └────┬─────┘
                │
                ▼
         ┌──────────────┐
         │ Rep >= 5?    │
         └──────┬───────┘
                │
           Yes  │  No
         ┌──────▼──┐  ┌───────────┐
         │Complete!│  │Keep Going │
         │Status=  │  │Status=    │
         │completed│  │learning   │
         └─────────┘  └───────────┘
```

---

## API Endpoints

### 1. Start Learning Vocabulary

**Endpoint**: `POST /api/v1/user-vocabulary-status`

**Purpose**: Mulai belajar vocabulary baru

**Request Body**:
```json
{
  "vocabularyId": 1
}
```

**Response (201 Created)**:
```json
{
  "status": "success",
  "message": "Created",
  "data": {
    "id": 1,
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "vocabularyId": 1,
    "vocabulary": {
      "id": 1,
      "word": "こんにちは",
      "reading": "konnichiwa",
      "meaning": "Hello"
    },
    "status": "learning",
    "repetitions": 0,
    "easeFactor": 2.5,
    "interval": 0,
    "nextReviewDate": "2026-01-11T10:00:00Z",
    "lastReviewedAt": null
  }
}
```

**Note**: Fields `easeFactor`, `interval`, dan `nextReviewDate` tidak digunakan dalam simple tracking (set ke default values untuk backward compatibility).

---

### 2. Get Vocabulary Status

**Endpoint**: `GET /api/v1/user-vocabulary-status/:id`

**Purpose**: Cek progress satu vocabulary

**Example Request**:
```
GET /api/v1/user-vocabulary-status/1
Authorization: Bearer <token>
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "id": 1,
    "userId": "550e8400-e29b-41d4-a716-446655440000",
    "vocabularyId": 1,
    "vocabulary": {
      "word": "こんにちは",
      "reading": "konnichiwa",
      "meaning": "Hello"
    },
    "status": "learning",
    "repetitions": 3,
    "lastReviewedAt": "2026-01-11T10:30:00Z"
  }
}
```

---

### 3. List All Vocabulary

**Endpoint**: `GET /api/v1/user-vocabulary-status`

**Purpose**: List semua vocabulary yang sedang dipelajari

**Query Parameters**:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `sort`: Sort field - `next_review_date`, `created_at`, `status` (default: `next_review_date`)
- `order`: `asc` or `desc` (default: `asc`)
- `status`: Filter by status - `learning`, `completed`

**Example Request**:
```
GET /api/v1/user-vocabulary-status?status=learning&page=1&limit=10
Authorization: Bearer <token>
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": [
    {
      "id": 1,
      "vocabularyId": 1,
      "vocabulary": {
        "word": "こんにちは",
        "meaning": "Hello"
      },
      "status": "learning",
      "repetitions": 3
    },
    {
      "id": 2,
      "vocabularyId": 2,
      "vocabulary": {
        "word": "ありがとう",
        "meaning": "Thank you"
      },
      "status": "completed",
      "repetitions": 5
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "totalPages": 1,
    "totalItems": 2
  }
}
```

---

### 4. Get Due Vocabularies

**Endpoint**: `GET /api/v1/user-vocabulary-status/due`

**Purpose**: Get vocabulary yang perlu direview

**Note**: Karena sistem ini user-controlled, endpoint ini mengembalikan semua vocabulary yang status-nya "learning" (belum completed). User yang menentukan mana yang mau direview.

**Example Request**:
```
GET /api/v1/user-vocabulary-status/due
Authorization: Bearer <token>
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": [
    {
      "id": 1,
      "vocabulary": {
        "word": "こんにちは",
        "reading": "konnichiwa",
        "meaning": "Hello"
      },
      "status": "learning",
      "repetitions": 3,
      "lastReviewedAt": "2026-01-10T10:00:00Z"
    },
    {
      "id": 3,
      "vocabulary": {
        "word": "さようなら",
        "reading": "sayounara",
        "meaning": "Goodbye"
      },
      "status": "learning",
      "repetitions": 1,
      "lastReviewedAt": "2026-01-09T10:00:00Z"
    }
  ]
}
```

---

### 5. Submit Review ✨

**Endpoint**: `POST /api/v1/user-vocabulary-status/:vocabulary_id/review`

**Purpose**: Submit hasil review vocabulary

**Request Body**:
```json
{
  "isCorrect": true
}
```

**Parameters**:
- `isCorrect` (boolean, required): Apakah user menjawab benar atau salah
  - `true`: Jawaban benar → increment repetitions
  - `false`: Jawaban salah → reset ke 0

**Example Requests**:

#### Correct Answer
```json
POST /api/v1/user-vocabulary-status/1/review
Authorization: Bearer <token>

{
  "isCorrect": true
}
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "id": 1,
    "vocabularyId": 1,
    "vocabulary": {
      "word": "こんにちは",
      "meaning": "Hello"
    },
    "status": "learning",
    "repetitions": 4,
    "lastReviewedAt": "2026-01-11T11:00:00Z"
  }
}
```

#### Incorrect Answer
```json
POST /api/v1/user-vocabulary-status/1/review
Authorization: Bearer <token>

{
  "isCorrect": false
}
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "id": 1,
    "vocabularyId": 1,
    "vocabulary": {
      "word": "こんにちは",
      "meaning": "Hello"
    },
    "status": "learning",
    "repetitions": 0,
    "lastReviewedAt": "2026-01-11T11:00:00Z"
  }
}
```

#### Completion (5th Correct Answer)
```json
POST /api/v1/user-vocabulary-status/1/review
Authorization: Bearer <token>

{
  "isCorrect": true
}
```

**Response (200 OK)**:
```json
{
  "status": "success",
  "message": "OK",
  "data": {
    "id": 1,
    "vocabularyId": 1,
    "vocabulary": {
      "word": "こんにちは",
      "meaning": "Hello"
    },
    "status": "completed",
    "repetitions": 5,
    "lastReviewedAt": "2026-01-11T11:00:00Z"
  }
}
```

---

## Usage Flow

### Complete Learning Cycle

```
1. START LEARNING
   POST /user-vocabulary-status { vocabularyId: 1 }
   ↓
   Status: learning
   Repetitions: 0

2. FIRST REVIEW
   POST /user-vocabulary-status/1/review { isCorrect: true }
   ↓
   Status: learning
   Repetitions: 1

3. SECOND REVIEW
   POST /user-vocabulary-status/1/review { isCorrect: true }
   ↓
   Status: learning
   Repetitions: 2

4. FAILED REVIEW (oops!)
   POST /user-vocabulary-status/1/review { isCorrect: false }
   ↓
   Status: learning
   Repetitions: 0 ← Reset!

5. TRY AGAIN - Review 1
   POST /user-vocabulary-status/1/review { isCorrect: true }
   ↓
   Status: learning
   Repetitions: 1

6. Continue... Reviews 2, 3, 4
   POST /user-vocabulary-status/1/review { isCorrect: true }
   ↓
   Status: learning
   Repetitions: 2, 3, 4

7. FINAL REVIEW - 5th Correct
   POST /user-vocabulary-status/1/review { isCorrect: true }
   ↓
   Status: completed ✓
   Repetitions: 5
```

### Daily Study Routine

```javascript
// Study Session Example
async function studySession() {
  // 1. Get vocabularies to review
  const response = await fetch('/user-vocabulary-status/due', {
    headers: { 'Authorization': `Bearer ${token}` }
  })

  const vocabularies = await response.json()

  // 2. For each vocabulary
  for (const vocab of vocabularies.data) {
    // 3. Show word to user
    showWord(vocab.vocabulary.word)  // "こんにちは"

    // 4. User tries to recall meaning
    const userAnswer = getUserInput()

    // 5. Check if correct
    const isCorrect = userAnswer === vocab.vocabulary.meaning

    // 6. Submit review
    await fetch(`/user-vocabulary-status/${vocab.id}/review`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ isCorrect })
    })

    // 7. Show feedback
    if (isCorrect) {
      showFeedback(`Correct! ${vocab.repetitions + 1}/5`)
    } else {
      showFeedback(`Wrong! Reset to 0/5. Try again!`)
    }
  }
}
```

---

## Examples

### Example 1: Perfect Learning Path

User jawab semua benar tanpa kesalahan:

```
Day 1:
POST /user-vocabulary-status { vocabularyId: 1 }
→ repetitions: 0, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 1, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 2, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 3, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 4, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 5, status: completed ✓

Timeline: 5 reviews → COMPLETED!
```

---

### Example 2: Struggling Learner

User beberapa kali salah, harus mengulang:

```
POST /user-vocabulary-status { vocabularyId: 1 }
→ repetitions: 0, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 1, status: learning

POST /user-vocabulary-status/1/review { isCorrect: false }
→ repetitions: 0, status: learning ← RESET!

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 1, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 2, status: learning

POST /user-vocabulary-status/1/review { isCorrect: false }
→ repetitions: 0, status: learning ← RESET AGAIN!

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 1, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 2, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 3, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 4, status: learning

POST /user-vocabulary-status/1/review { isCorrect: true }
→ repetitions: 5, status: completed ✓

Timeline: 12 reviews (with 2 failures) → COMPLETED!
```

**Key Points**:
- Failed reviews reset progress to 0
- User needs to answer correctly 5 times in a row
- No time pressure, user reviews whenever ready

---

## Status Progression

### Learning Status States

| Status | Repetitions | Meaning |
|--------|-------------|---------|
| **learning** | 0-4 | Still learning, need more practice |
| **completed** | 5+ | Successfully mastered! |

### State Transitions

```
┌──────────┐
│ learning │ ← Start here (repetitions: 0)
└────┬─────┘
     │
     │ Answer correctly 5 times
     ↓
┌──────────┐
│completed │ ← Goal! (repetitions: 5)
└──────────┘

Note: Wrong answer resets to learning (repetitions: 0)
```

### Visual Progress

```
Repetitions:     0    1    2    3    4    5
Status:       ┌────────────────────────┬────┐
              │       learning         │comp│
              └────────────────────────┴────┘

Progress:     [    ][    ][    ][    ][    ] ← Empty (0/5)
              [████][    ][    ][    ][    ] ← 1/5
              [████][████][    ][    ][    ] ← 2/5
              [████][████][████][    ][    ] ← 3/5
              [████][████][████][████][    ] ← 4/5
              [████][████][████][████][████] ← 5/5 DONE!
```

---

## Troubleshooting

### Common Issues

#### 1. Vocabulary Already Being Learned

**Error**:
```json
{
  "status": "error",
  "message": "vocabulary already being learned by user"
}
```

**Cause**: Trying to start learning vocabulary yang sudah ada di learning list.

**Solution**: Check existing status terlebih dahulu.

---

#### 2. Status Not Found

**Error**:
```json
{
  "status": "error",
  "message": "user vocabulary status not found"
}
```

**Cause**: Status ID tidak exist atau sudah dihapus.

**Solution**: Verify ID dari list endpoint.

---

#### 3. Forbidden Access

**Error**:
```json
{
  "status": "error",
  "message": "forbidden - this status belongs to another user"
}
```

**Cause**: Trying to access user lain punya learning status.

**Solution**: Check authentication token. Setiap user hanya bisa access data mereka sendiri.

---

### Best Practices

#### ✅ DO

1. **Review Regularly**
   - Practice vocabulary setiap hari untuk consistency
   - Use `/due` endpoint untuk list vocabulary yang perlu direview

2. **Be Honest**
   - Jawab `isCorrect: true` hanya jika benar-benar ingat
   - Jangan curang, itu hanya merugikan diri sendiri

3. **Track Progress**
   ```
   GET /user-vocabulary-status?status=completed
   ```
   Monitor completed vocabularies untuk motivasi

4. **Focus on Learning**
   - Jangan terburu-buru mencapai "completed"
   - Focus pada understanding, bukan speed

#### ❌ DON'T

1. **Don't Skip Learning**
   - Jangan langsung mark `isCorrect: true` tanpa benar-benar review
   - Defeats the purpose of learning

2. **Don't Give Up**
   - Reset ke 0 itu normal, bagian dari learning process
   - Keep trying sampai berhasil 5x berturut-turut

3. **Don't Delete & Re-add**
   - Keep failed items dan retry
   - Progress history penting untuk tracking

---

## API Reference Summary

| Endpoint | Method | Purpose | Auth Required |
|----------|--------|---------|---------------|
| `/user-vocabulary-status` | POST | Start learning | ✅ |
| `/user-vocabulary-status/:id` | GET | Get status | ✅ |
| `/user-vocabulary-status` | GET | List all | ✅ |
| `/user-vocabulary-status/due` | GET | Get due items | ✅ |
| `/user-vocabulary-status/:vocabulary_id/review` | POST | Submit review | ✅ |

---

## Comparison: Simple vs SM-2

| Aspect | Simple Tracking | SM-2 Algorithm |
|--------|----------------|----------------|
| **Complexity** | Very simple | Complex |
| **Input** | Boolean (true/false) | Quality scale (0-5) |
| **Completion** | 5 correct answers | 6+ reviews + EF ≥ 2.3 |
| **Scheduling** | User decides | Automatic intervals |
| **Reset Logic** | Simple (0 on fail) | Complex (EF adjustment) |
| **Status States** | 2 (learning, completed) | 3 (learning, reviewing, mastered) |
| **User Control** | Full control | Algorithm-controlled |
| **Learning Curve** | Easy to understand | Requires explanation |

**Why Simple Tracking?**
- ✅ Easier to understand for beginners
- ✅ User has full control
- ✅ Clear, measurable goals (5 = done)
- ✅ No complex calculations
- ✅ Suitable for casual learners

---

**Last Updated**: 2026-01-11
**API Version**: v1
**Status**: Production Ready ✅
