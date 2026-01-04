# API Gaps Analysis - Manabu Japanese Learning App

## Current API Status

### ✅ Existing APIs (User Management Only)

**Auth Routes** (`/auth`)
- `POST /auth/login` - User login
- `POST /auth/register` - User registration
- `GET /auth/user` - Get logged in user (authenticated)
- `GET /auth/:uuid` - Get user by UUID (authenticated)
- `PUT /auth/:uuid` - Update user (authenticated)

---

## 🚨 Missing APIs (Berdasarkan ERD)

Berdasarkan ERD yang sudah dibuat, berikut adalah API yang **BELUM ADA** dan perlu dibangun:

---

### 📚 1. JLPT Levels Management

**Priority: HIGH** - Fundamental untuk klasifikasi konten

```
GET    /api/v1/jlpt-levels              - List all JLPT levels
GET    /api/v1/jlpt-levels/:id          - Get JLPT level detail
POST   /api/v1/jlpt-levels              - Create JLPT level (admin)
PUT    /api/v1/jlpt-levels/:id          - Update JLPT level (admin)
DELETE /api/v1/jlpt-levels/:id          - Delete JLPT level (admin)
```

**Related Tables:** `jlpt_levels`

---

### 📂 2. Categories Management

**Priority: HIGH** - Untuk organisasi vocabulary

```
GET    /api/v1/categories                      - List all categories
GET    /api/v1/categories/:id                  - Get category detail
GET    /api/v1/categories/jlpt/:jlpt_level_id  - Get categories by JLPT level
POST   /api/v1/categories                      - Create category (admin)
PUT    /api/v1/categories/:id                  - Update category (admin)
DELETE /api/v1/categories/:id                  - Delete category (admin)
```

**Related Tables:** `categories`, `jlpt_levels`

---

### 📖 3. Vocabulary Management

**Priority: CRITICAL** - Core feature aplikasi belajar bahasa

```
GET    /api/v1/vocabularies                    - List vocabularies (with pagination & filters)
GET    /api/v1/vocabularies/:id                - Get vocabulary detail
GET    /api/v1/vocabularies/search             - Search vocabularies
GET    /api/v1/vocabularies/random             - Get random vocabulary (for practice)
GET    /api/v1/vocabularies/jlpt/:level        - Get vocabularies by JLPT level
GET    /api/v1/vocabularies/category/:id       - Get vocabularies by category
GET    /api/v1/vocabularies/word-type/:type    - Get vocabularies by word type
POST   /api/v1/vocabularies                    - Create vocabulary (admin)
PUT    /api/v1/vocabularies/:id                - Update vocabulary (admin)
DELETE /api/v1/vocabularies/:id                - Delete vocabulary (admin)
PATCH  /api/v1/vocabularies/:id/activate       - Activate vocabulary (admin)
PATCH  /api/v1/vocabularies/:id/deactivate     - Deactivate vocabulary (admin)
```

**Query Parameters:**
- `?page=1&limit=20` - Pagination
- `?jlpt_level=N5` - Filter by JLPT level
- `?category_id=1` - Filter by category
- `?word_type=noun` - Filter by word type
- `?search=こんにちは` - Search kanji/hiragana/romaji/meaning

**Related Tables:** `vocabularies`, `categories`, `jlpt_levels`, `tags`

---

### 🏷️ 4. Tags Management

**Priority: MEDIUM** - Untuk organisasi vocabulary lebih fleksibel

```
GET    /api/v1/tags                     - List all tags
GET    /api/v1/tags/:id                 - Get tag detail
GET    /api/v1/tags/:id/vocabularies    - Get vocabularies by tag
POST   /api/v1/tags                     - Create tag (admin)
PUT    /api/v1/tags/:id                 - Update tag (admin)
DELETE /api/v1/tags/:id                 - Delete tag (admin)
```

**Related Tables:** `tags`, `vocabulary_tags`, `vocabularies`

---

### 🎓 5. Courses Management

**Priority: HIGH** - Struktur pembelajaran

```
GET    /api/v1/courses                  - List all courses
GET    /api/v1/courses/:id              - Get course detail with lessons
GET    /api/v1/courses/jlpt/:level      - Get courses by JLPT level
POST   /api/v1/courses                  - Create course (admin)
PUT    /api/v1/courses/:id              - Update course (admin)
DELETE /api/v1/courses/:id              - Delete course (admin)
PATCH  /api/v1/courses/:id/activate     - Activate course (admin)
PATCH  /api/v1/courses/:id/deactivate   - Deactivate course (admin)
```

**Related Tables:** `courses`, `jlpt_levels`, `lessons`

---

### 📝 6. Lessons Management

**Priority: HIGH** - Konten pembelajaran

```
GET    /api/v1/lessons/:id              - Get lesson detail
GET    /api/v1/courses/:course_id/lessons  - Get lessons in a course
POST   /api/v1/courses/:course_id/lessons  - Create lesson (admin)
PUT    /api/v1/lessons/:id              - Update lesson (admin)
DELETE /api/v1/lessons/:id              - Delete lesson (admin)
PATCH  /api/v1/lessons/:id/reorder      - Reorder lessons (admin)
```

**Related Tables:** `lessons`, `courses`, `exercises`

---

### ✍️ 7. Exercises Management

**Priority: HIGH** - Latihan soal

```
GET    /api/v1/exercises/:id                    - Get exercise detail
GET    /api/v1/lessons/:lesson_id/exercises     - Get exercises in a lesson
POST   /api/v1/lessons/:lesson_id/exercises     - Create exercise (admin)
PUT    /api/v1/exercises/:id                    - Update exercise (admin)
DELETE /api/v1/exercises/:id                    - Delete exercise (admin)
GET    /api/v1/exercises/:id/questions          - Get questions in exercise
```

**Exercise Types:**
- quiz
- flashcard
- writing
- listening
- matching

**Related Tables:** `exercises`, `lessons`, `exercise_questions`

---

### ❓ 8. Exercise Questions Management

**Priority: HIGH** - Pertanyaan dalam exercise

```
GET    /api/v1/questions/:id                      - Get question detail
GET    /api/v1/exercises/:exercise_id/questions   - Get questions in exercise
POST   /api/v1/exercises/:exercise_id/questions   - Create question (admin)
PUT    /api/v1/questions/:id                      - Update question (admin)
DELETE /api/v1/questions/:id                      - Delete question (admin)
```

**Related Tables:** `exercise_questions`, `exercises`, `vocabularies`

---

### 📊 9. User Course Progress

**Priority: CRITICAL** - Tracking progress user

```
GET    /api/v1/my/courses                      - Get enrolled courses
POST   /api/v1/courses/:id/enroll              - Enroll in a course
GET    /api/v1/my/courses/:id/progress         - Get course progress
PATCH  /api/v1/my/courses/:id/complete         - Mark course as completed
```

**Related Tables:** `user_course_progress`, `courses`, `users`

---

### 📖 10. User Lesson Progress

**Priority: CRITICAL** - Tracking lesson progress

```
GET    /api/v1/my/lessons/:id/progress         - Get lesson progress
POST   /api/v1/my/lessons/:id/start            - Start lesson (mark in_progress)
PATCH  /api/v1/my/lessons/:id/complete         - Complete lesson
```

**Related Tables:** `user_lesson_progress`, `lessons`, `users`

---

### ✅ 11. User Exercise Attempts

**Priority: CRITICAL** - Submit & track jawaban exercise

```
POST   /api/v1/exercises/:id/attempt           - Start exercise attempt
POST   /api/v1/attempts/:id/submit             - Submit exercise attempt
GET    /api/v1/my/attempts                     - Get user's attempt history
GET    /api/v1/my/attempts/:id                 - Get attempt detail with answers
GET    /api/v1/exercises/:id/attempts          - Get attempts for specific exercise
```

**Submit Payload Example:**
```json
{
  "attempt_id": "uuid",
  "answers": [
    {
      "question_id": 1,
      "user_answer": "こんにちは",
      "time_taken_seconds": 15
    }
  ]
}
```

**Related Tables:** `user_exercise_attempts`, `user_question_answers`, `exercises`

---

### 🔄 12. User Vocabulary Status (Spaced Repetition)

**Priority: CRITICAL** - Sistem review vocabulary

```
GET    /api/v1/my/vocabulary/status            - Get all vocabulary status
GET    /api/v1/my/vocabulary/review            - Get vocabularies due for review
POST   /api/v1/my/vocabulary/:id/review        - Mark vocabulary as reviewed
PATCH  /api/v1/my/vocabulary/:id/mastery       - Update mastery level
GET    /api/v1/my/vocabulary/stats             - Get vocabulary learning stats
```

**Review Payload:**
```json
{
  "vocabulary_id": 123,
  "is_correct": true,
  "time_taken_seconds": 10
}
```

**Related Tables:** `user_vocabulary_status`, `vocabularies`, `users`

---

### 🏆 13. Achievements & Gamification

**Priority: MEDIUM** - Motivasi belajar

```
GET    /api/v1/achievements                    - List all achievements
GET    /api/v1/my/achievements                 - Get user's achievements
GET    /api/v1/my/achievements/progress        - Get achievement progress
```

**Related Tables:** `achievements`, `user_achievements`

---

### 🔥 14. User Streaks

**Priority: MEDIUM** - Daily streak tracking

```
GET    /api/v1/my/streak                       - Get current streak
POST   /api/v1/my/streak/checkin               - Daily check-in
```

**Related Tables:** `user_streaks`

---

### 📈 15. User Statistics & Analytics

**Priority: MEDIUM** - Dashboard user

```
GET    /api/v1/my/stats/daily                  - Get daily stats
GET    /api/v1/my/stats/summary                - Get overall stats summary
GET    /api/v1/my/stats/weekly                 - Get weekly stats
GET    /api/v1/my/stats/monthly                - Get monthly stats
```

**Response Example:**
```json
{
  "total_words_learned": 245,
  "total_lessons_completed": 12,
  "total_exercises_completed": 45,
  "total_time_spent_minutes": 1230,
  "current_streak": 7,
  "average_score": 85.5
}
```

**Related Tables:** `user_daily_stats`, `user_vocabulary_status`, `user_exercise_attempts`

---

## Priority Summary

### 🔴 CRITICAL (Harus dibuat pertama)
1. **Vocabulary Management** - Core feature
2. **User Vocabulary Status** - Spaced repetition system
3. **User Course/Lesson Progress** - Progress tracking
4. **User Exercise Attempts** - Exercise submission

### 🟡 HIGH (Prioritas tinggi)
5. **JLPT Levels** - Fundamental classification
6. **Categories** - Content organization
7. **Courses & Lessons** - Learning structure
8. **Exercises & Questions** - Practice content

### 🟢 MEDIUM (Bisa ditambahkan kemudian)
9. **Tags** - Flexible organization
10. **Achievements** - Gamification
11. **Streaks** - Engagement
12. **Statistics** - Analytics

---

## Recommended Implementation Order

### Phase 1: Foundation (Week 1-2)
1. JLPT Levels API
2. Categories API
3. Vocabulary CRUD API
4. Tags API (basic)

### Phase 2: Learning Content (Week 3-4)
5. Courses API
6. Lessons API
7. Exercises API
8. Questions API

### Phase 3: User Progress (Week 5-6)
9. Course Progress API
10. Lesson Progress API
11. Exercise Attempts API
12. Vocabulary Status/Review API

### Phase 4: Gamification (Week 7-8)
13. Achievements API
14. Streaks API
15. Statistics & Analytics API

---

## Additional Considerations

### Authentication & Authorization
- Semua endpoint `/my/*` butuh authentication
- Admin endpoints butuh role-based authorization
- Rate limiting untuk public endpoints

### Pagination
- Semua list endpoints harus support pagination
- Default: `page=1&limit=20`
- Max limit: 100

### Filtering & Sorting
- Support query parameters untuk filter
- Support sorting: `?sort_by=created_at&order=desc`

### Response Format
- Consistent response structure
- Include metadata (pagination, total count)
- Proper error handling

### File Upload
- Vocabulary: audio files, images
- Course: thumbnails
- Need file storage service (S3/local)

---

## Next Steps

Apakah Anda ingin saya mulai implement API dari Phase 1 (Foundation)?
