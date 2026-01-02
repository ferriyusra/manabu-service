# Feature Checklist - Manabu Japanese Learning App

## 📊 Overall Progress

**MVP Phase**: 2/5 (40%) ✅
**Phase 1 (Foundation)**: 2/4 (50%) ✅
**Phase 2 (Learning Content)**: 0/4 (0%)
**Phase 3 (User Progress)**: 0/4 (0%)
**Phase 4 (Gamification)**: 0/3 (0%)

**Total Progress**: 2/15 APIs (13.3%)

---

## 🎯 MVP Features (Must Have for Launch)

### ✅ 1. User Management (COMPLETED)
**Status**: ✅ Ready for Production
**Priority**: CRITICAL

- [x] POST `/auth/register` - User registration
- [x] POST `/auth/login` - User login
- [x] GET `/auth/user` - Get logged in user (authenticated)
- [x] GET `/auth/:uuid` - Get user by UUID (authenticated)
- [x] PUT `/auth/:uuid` - Update user (authenticated)

**Last Updated**: Initial implementation
**Notes**: Authentication & authorization working correctly

---

### ✅ 2. JLPT Levels Management (COMPLETED)
**Status**: ✅ Ready for Production
**Priority**: HIGH - Foundation for content classification

- [x] GET `/api/v1/jlpt-levels` - List all JLPT levels
- [x] GET `/api/v1/jlpt-levels/:id` - Get JLPT level detail
- [x] POST `/api/v1/jlpt-levels` - Create JLPT level (admin)
- [x] PUT `/api/v1/jlpt-levels/:id` - Update JLPT level (admin)
- [x] DELETE `/api/v1/jlpt-levels/:id` - Delete JLPT level (admin)

**Last Updated**: 2026-01-02
**Features**:
- ✅ Auto-generated Swagger documentation
- ✅ Proper error handling
- ✅ Admin authentication required for CUD operations
- ✅ Validation & business logic

---

### ✅ 3. Categories Management (COMPLETED)
**Status**: ✅ Ready for Production
**Priority**: HIGH - For vocabulary organization

- [x] GET `/api/v1/categories` - List all categories (paginated)
- [x] GET `/api/v1/categories/:id` - Get category detail
- [x] GET `/api/v1/categories/jlpt/:jlpt_level_id` - Get categories by JLPT level (paginated)
- [x] POST `/api/v1/categories` - Create category (admin)
- [x] PUT `/api/v1/categories/:id` - Update category (admin)
- [x] DELETE `/api/v1/categories/:id` - Delete category (admin)

**Last Updated**: 2026-01-02
**Features**:
- ✅ Unique constraint on (name, jlpt_level_id)
- ✅ Pagination support (default: page=1, limit=10, max=100)
- ✅ Proper HTTP status codes (404, 409, 422, 500)
- ✅ Optimized queries (no N+1 issues)
- ✅ Clean response structure with data array
- ✅ Defensive validation
- ✅ DRY code with helper methods
- ✅ Complete Swagger documentation

**Code Quality**:
- ✅ Reviewed by code-review-expert agent
- ✅ All critical & major issues fixed
- ✅ Follows clean architecture pattern
- ✅ Proper error handling & validation

---

### ⏳ 4. Vocabulary Management (IN PROGRESS)
**Status**: 🔴 Not Started
**Priority**: CRITICAL - Core feature

**Endpoints Needed** (10 total):
- [ ] GET `/api/v1/vocabularies` - List vocabularies (paginated, filterable)
- [ ] GET `/api/v1/vocabularies/:id` - Get vocabulary detail
- [ ] GET `/api/v1/vocabularies/search` - Search vocabularies
- [ ] GET `/api/v1/vocabularies/random` - Get random vocabulary (for practice)
- [ ] GET `/api/v1/vocabularies/jlpt/:level` - Get vocabularies by JLPT level
- [ ] GET `/api/v1/vocabularies/category/:id` - Get vocabularies by category
- [ ] GET `/api/v1/vocabularies/word-type/:type` - Get vocabularies by word type
- [ ] POST `/api/v1/vocabularies` - Create vocabulary (admin)
- [ ] PUT `/api/v1/vocabularies/:id` - Update vocabulary (admin)
- [ ] DELETE `/api/v1/vocabularies/:id` - Delete vocabulary (admin)
- [ ] PATCH `/api/v1/vocabularies/:id/activate` - Activate vocabulary (admin)
- [ ] PATCH `/api/v1/vocabularies/:id/deactivate` - Deactivate vocabulary (admin)

**Query Parameters Required**:
- `?page=1&limit=20` - Pagination
- `?jlpt_level=N5` - Filter by JLPT level
- `?category_id=1` - Filter by category
- `?word_type=noun` - Filter by word type (noun, verb, adjective, etc.)
- `?search=こんにちは` - Search kanji/hiragana/romaji/meaning

**Database Schema**:
- Kanji (nullable)
- Hiragana (required)
- Romaji (required)
- Meaning (required)
- WordType (enum: noun, verb, adjective, adverb, etc.)
- JlptLevelID (foreign key)
- CategoryID (foreign key)
- AudioURL (optional)
- ImageURL (optional)
- ExampleSentence (optional)
- IsActive (default: true)

**Related Tables**: vocabularies, categories, jlpt_levels, tags (many-to-many)

**Complexity**: HIGH - Most complex CRUD with multiple filters & search

---

### ⏳ 5. User Vocabulary Status (Spaced Repetition)
**Status**: 🔴 Not Started
**Priority**: CRITICAL - Core learning feature

**Endpoints Needed**:
- [ ] GET `/api/v1/my/vocabulary/status` - Get all vocabulary status
- [ ] GET `/api/v1/my/vocabulary/review` - Get vocabularies due for review
- [ ] POST `/api/v1/my/vocabulary/:id/review` - Mark vocabulary as reviewed
- [ ] PATCH `/api/v1/my/vocabulary/:id/mastery` - Update mastery level
- [ ] GET `/api/v1/my/vocabulary/stats` - Get vocabulary learning stats

**Spaced Repetition Algorithm**:
- Mastery levels: 0 (new), 1 (learning), 2 (familiar), 3 (mastered)
- Review intervals: 1 day, 3 days, 7 days, 14 days, 30 days
- Adjust intervals based on correctness

**Database Schema**:
- UserID (foreign key)
- VocabularyID (foreign key)
- MasteryLevel (0-3)
- ReviewCount
- CorrectCount
- LastReviewedAt
- NextReviewAt
- CreatedAt

---

## 🏗️ Phase 1: Foundation (Week 1-2)

### Progress: 2/4 (50%) ✅

1. ✅ **JLPT Levels API** - COMPLETED
2. ✅ **Categories API** - COMPLETED
3. 🔴 **Vocabulary CRUD API** - Not Started
4. 🔴 **Tags API (basic)** - Not Started

---

## 📚 Phase 2: Learning Content (Week 3-4)

### Progress: 0/4 (0%)

### ⏳ 1. Tags Management
**Status**: 🔴 Not Started
**Priority**: MEDIUM - Flexible organization

**Endpoints Needed**:
- [ ] GET `/api/v1/tags` - List all tags
- [ ] GET `/api/v1/tags/:id` - Get tag detail
- [ ] GET `/api/v1/tags/:id/vocabularies` - Get vocabularies by tag
- [ ] POST `/api/v1/tags` - Create tag (admin)
- [ ] PUT `/api/v1/tags/:id` - Update tag (admin)
- [ ] DELETE `/api/v1/tags/:id` - Delete tag (admin)

**Database Schema**:
- Name (unique)
- Color (optional, for UI)
- Many-to-many with vocabularies via `vocabulary_tags` table

---

### ⏳ 2. Courses Management
**Status**: 🔴 Not Started
**Priority**: HIGH - Learning structure

**Endpoints Needed**:
- [ ] GET `/api/v1/courses` - List all courses
- [ ] GET `/api/v1/courses/:id` - Get course detail with lessons
- [ ] GET `/api/v1/courses/jlpt/:level` - Get courses by JLPT level
- [ ] POST `/api/v1/courses` - Create course (admin)
- [ ] PUT `/api/v1/courses/:id` - Update course (admin)
- [ ] DELETE `/api/v1/courses/:id` - Delete course (admin)
- [ ] PATCH `/api/v1/courses/:id/activate` - Activate course (admin)
- [ ] PATCH `/api/v1/courses/:id/deactivate` - Deactivate course (admin)

**Database Schema**:
- Title
- Description
- JlptLevelID (foreign key)
- ThumbnailURL (optional)
- EstimatedHours
- OrderIndex
- IsActive (default: true)

---

### ⏳ 3. Lessons Management
**Status**: 🔴 Not Started
**Priority**: HIGH - Learning content

**Endpoints Needed**:
- [ ] GET `/api/v1/lessons/:id` - Get lesson detail
- [ ] GET `/api/v1/courses/:course_id/lessons` - Get lessons in a course
- [ ] POST `/api/v1/courses/:course_id/lessons` - Create lesson (admin)
- [ ] PUT `/api/v1/lessons/:id` - Update lesson (admin)
- [ ] DELETE `/api/v1/lessons/:id` - Delete lesson (admin)
- [ ] PATCH `/api/v1/lessons/:id/reorder` - Reorder lessons (admin)

**Database Schema**:
- CourseID (foreign key)
- Title
- Description
- Content (markdown/HTML)
- OrderIndex
- EstimatedMinutes
- CreatedAt

---

### ⏳ 4. Exercises Management
**Status**: 🔴 Not Started
**Priority**: HIGH - Practice exercises

**Endpoints Needed**:
- [ ] GET `/api/v1/exercises/:id` - Get exercise detail
- [ ] GET `/api/v1/lessons/:lesson_id/exercises` - Get exercises in a lesson
- [ ] POST `/api/v1/lessons/:lesson_id/exercises` - Create exercise (admin)
- [ ] PUT `/api/v1/exercises/:id` - Update exercise (admin)
- [ ] DELETE `/api/v1/exercises/:id` - Delete exercise (admin)
- [ ] GET `/api/v1/exercises/:id/questions` - Get questions in exercise

**Exercise Types**:
- quiz (multiple choice)
- flashcard (show/hide answer)
- writing (free text input)
- listening (audio playback)
- matching (drag & drop)

**Database Schema**:
- LessonID (foreign key)
- Title
- Type (enum)
- Instructions
- OrderIndex
- PassingScore (percentage)

---

## 📈 Phase 3: User Progress (Week 5-6)

### Progress: 0/4 (0%)

### ⏳ 1. Exercise Questions Management
**Status**: 🔴 Not Started
**Priority**: HIGH

**Endpoints Needed**:
- [ ] GET `/api/v1/questions/:id` - Get question detail
- [ ] GET `/api/v1/exercises/:exercise_id/questions` - Get questions in exercise
- [ ] POST `/api/v1/exercises/:exercise_id/questions` - Create question (admin)
- [ ] PUT `/api/v1/questions/:id` - Update question (admin)
- [ ] DELETE `/api/v1/questions/:id` - Delete question (admin)

**Database Schema**:
- ExerciseID (foreign key)
- QuestionText
- QuestionType (multiple_choice, true_false, fill_blank, etc.)
- Options (JSON array for multiple choice)
- CorrectAnswer
- VocabularyID (foreign key, optional)
- Points
- OrderIndex

---

### ⏳ 2. User Course Progress
**Status**: 🔴 Not Started
**Priority**: CRITICAL

**Endpoints Needed**:
- [ ] GET `/api/v1/my/courses` - Get enrolled courses
- [ ] POST `/api/v1/courses/:id/enroll` - Enroll in a course
- [ ] GET `/api/v1/my/courses/:id/progress` - Get course progress
- [ ] PATCH `/api/v1/my/courses/:id/complete` - Mark course as completed

**Database Schema**:
- UserID (foreign key)
- CourseID (foreign key)
- Status (not_started, in_progress, completed)
- ProgressPercentage (calculated)
- EnrolledAt
- CompletedAt (nullable)

---

### ⏳ 3. User Lesson Progress
**Status**: 🔴 Not Started
**Priority**: CRITICAL

**Endpoints Needed**:
- [ ] GET `/api/v1/my/lessons/:id/progress` - Get lesson progress
- [ ] POST `/api/v1/my/lessons/:id/start` - Start lesson (mark in_progress)
- [ ] PATCH `/api/v1/my/lessons/:id/complete` - Complete lesson

**Database Schema**:
- UserID (foreign key)
- LessonID (foreign key)
- Status (not_started, in_progress, completed)
- StartedAt
- CompletedAt (nullable)

---

### ⏳ 4. User Exercise Attempts
**Status**: 🔴 Not Started
**Priority**: CRITICAL

**Endpoints Needed**:
- [ ] POST `/api/v1/exercises/:id/attempt` - Start exercise attempt
- [ ] POST `/api/v1/attempts/:id/submit` - Submit exercise attempt
- [ ] GET `/api/v1/my/attempts` - Get user's attempt history
- [ ] GET `/api/v1/my/attempts/:id` - Get attempt detail with answers
- [ ] GET `/api/v1/exercises/:id/attempts` - Get attempts for specific exercise

**Submit Payload Example**:
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

**Database Schema**:
- UserID (foreign key)
- ExerciseID (foreign key)
- Score (percentage)
- TotalQuestions
- CorrectAnswers
- TimeSpentSeconds
- StartedAt
- SubmittedAt

---

## 🎮 Phase 4: Gamification (Week 7-8)

### Progress: 0/3 (0%)

### ⏳ 1. Achievements & Gamification
**Status**: 🔴 Not Started
**Priority**: MEDIUM - User engagement

**Endpoints Needed**:
- [ ] GET `/api/v1/achievements` - List all achievements
- [ ] GET `/api/v1/my/achievements` - Get user's achievements
- [ ] GET `/api/v1/my/achievements/progress` - Get achievement progress

**Achievement Types**:
- First Step (Complete first lesson)
- Vocabulary Master (Learn 100 words)
- Perfect Score (100% on any exercise)
- Week Streak (7 days in a row)
- JLPT N5 Ready (Complete all N5 content)

**Database Schema**:
- Name
- Description
- IconURL
- Criteria (JSON)
- Points
- BadgeType (bronze, silver, gold)

---

### ⏳ 2. User Streaks
**Status**: 🔴 Not Started
**Priority**: MEDIUM - Daily engagement

**Endpoints Needed**:
- [ ] GET `/api/v1/my/streak` - Get current streak
- [ ] POST `/api/v1/my/streak/checkin` - Daily check-in

**Database Schema**:
- UserID (foreign key)
- CurrentStreak
- LongestStreak
- LastCheckInDate
- TotalCheckIns

---

### ⏳ 3. User Statistics & Analytics
**Status**: 🔴 Not Started
**Priority**: MEDIUM - User dashboard

**Endpoints Needed**:
- [ ] GET `/api/v1/my/stats/daily` - Get daily stats
- [ ] GET `/api/v1/my/stats/summary` - Get overall stats summary
- [ ] GET `/api/v1/my/stats/weekly` - Get weekly stats
- [ ] GET `/api/v1/my/stats/monthly` - Get monthly stats

**Response Example**:
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

**Database Schema**:
- UserID (foreign key)
- Date
- WordsLearned
- LessonsCompleted
- ExercisesCompleted
- TimeSpentMinutes
- AverageScore

---

## 🔧 Additional Features (Post-MVP)

### File Upload Support
**Status**: 🔴 Not Planned Yet
**Priority**: LOW

**Requirements**:
- Vocabulary: audio files (.mp3), images (.jpg, .png)
- Courses: thumbnail images
- Need file storage service (AWS S3 or local storage)

**Endpoints Needed**:
- [ ] POST `/api/v1/upload/audio` - Upload audio file
- [ ] POST `/api/v1/upload/image` - Upload image file
- [ ] DELETE `/api/v1/upload/:id` - Delete uploaded file

---

### Admin Dashboard APIs
**Status**: 🔴 Not Planned Yet
**Priority**: MEDIUM

**Endpoints Needed**:
- [ ] GET `/api/v1/admin/stats/overview` - Admin dashboard overview
- [ ] GET `/api/v1/admin/users` - List all users (paginated)
- [ ] GET `/api/v1/admin/users/:id` - Get user detail
- [ ] PATCH `/api/v1/admin/users/:id/ban` - Ban user
- [ ] PATCH `/api/v1/admin/users/:id/unban` - Unban user

---

## 📋 Priority Summary

### 🔴 CRITICAL (Must-have for MVP)
1. ✅ User Management (Auth)
2. ✅ JLPT Levels
3. ✅ Categories
4. 🔴 **Vocabulary Management** ← NEXT
5. 🔴 **User Vocabulary Status (Spaced Repetition)** ← NEXT

### 🟡 HIGH (Important for full experience)
6. 🔴 Tags
7. 🔴 Courses
8. 🔴 Lessons
9. 🔴 Exercises
10. 🔴 Exercise Questions
11. 🔴 User Course Progress
12. 🔴 User Lesson Progress
13. 🔴 User Exercise Attempts

### 🟢 MEDIUM (Can be added later)
14. 🔴 Achievements
15. 🔴 Streaks
16. 🔴 Statistics

---

## 🎯 Recommended Next Steps

### Immediate (This Week)
1. **Develop Vocabulary API** (CRITICAL - Core feature)
   - Most complex CRUD with filters & search
   - Foundation for all learning features

2. **Develop User Vocabulary Status API** (CRITICAL - Spaced repetition)
   - Core learning mechanism
   - Implements spaced repetition algorithm

### Week 2
3. **Develop Tags API** (MEDIUM - Flexible organization)
4. **Develop Courses API** (HIGH - Learning structure)

### Week 3-4
5. **Develop Lessons & Exercises APIs** (HIGH - Practice content)
6. **Develop User Progress APIs** (CRITICAL - Track learning)

### Week 5+
7. **Gamification Features** (MEDIUM - Engagement)
8. **Admin Dashboard** (MEDIUM - Management)
9. **File Upload** (LOW - Enhanced content)

---

## 📊 Definition of Done

For each API endpoint to be marked as ✅ COMPLETED, it must have:

1. **Code Implementation**
   - [ ] Model with proper GORM tags
   - [ ] DTO (Request/Response)
   - [ ] Repository layer with CRUD methods
   - [ ] Service layer with business logic
   - [ ] Controller with all endpoints
   - [ ] Routes registered

2. **Quality Assurance**
   - [ ] Error handling with proper status codes
   - [ ] Input validation
   - [ ] Unique constraints where needed
   - [ ] Optimized queries (no N+1)
   - [ ] Security (authentication/authorization)

3. **Documentation**
   - [ ] Swagger/OpenAPI annotations
   - [ ] Auto-generated API docs
   - [ ] Code comments for complex logic

4. **Testing** (optional for MVP)
   - [ ] Unit tests for services
   - [ ] Integration tests for endpoints
   - [ ] Manual testing via Swagger UI

5. **Code Review**
   - [ ] Reviewed by code-review-expert agent
   - [ ] All critical issues fixed
   - [ ] All major issues fixed

---

**Last Updated**: 2026-01-02
**Next Review**: After completing Vocabulary API
