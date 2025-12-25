# ERD - Manabu Japanese Learning Application

## Database Schema Design

### Overview
This ERD represents the complete database structure for a Japanese language learning application with support for vocabulary, lessons, exercises, progress tracking, and JLPT level classification.

---

## Entity Relationship Diagram

```
┌─────────────────┐
│     users       │
├─────────────────┤
│ PK id           │
│    uuid         │
│    name         │
│    username     │
│    password     │
│    phone_number │
│    email        │
│ FK role_id      │
│    created_at   │
│    updated_at   │
└────────┬────────┘
         │
         │ 1:N
         │
┌────────┴────────┐
│     roles       │
├─────────────────┤
│ PK id           │
│    code         │
│    name         │
│    created_at   │
│    updated_at   │
└─────────────────┘


┌──────────────────────┐
│   jlpt_levels        │
├──────────────────────┤
│ PK id                │
│    code              │  (N5, N4, N3, N2, N1)
│    name              │
│    description       │
│    level_order       │  (5, 4, 3, 2, 1)
│    created_at        │
│    updated_at        │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────┴───────────┐
│   categories         │
├──────────────────────┤
│ PK id                │
│    code              │
│    name              │
│    description       │
│    icon              │
│ FK jlpt_level_id     │
│    created_at        │
│    updated_at        │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────┴───────────────────┐
│   vocabularies               │
├──────────────────────────────┤
│ PK id                        │
│    uuid                      │
│    kanji                     │
│    hiragana                  │
│    romaji                    │
│    meaning_id                │  (Indonesia)
│    meaning_en                │  (English)
│    word_type                 │  (noun, verb, adjective, etc.)
│    example_sentence_jp       │
│    example_sentence_id       │
│    example_sentence_en       │
│    audio_url                 │
│    image_url                 │
│ FK category_id               │
│ FK jlpt_level_id             │
│    is_active                 │
│    created_at                │
│    updated_at                │
└──────────┬───────────────────┘
           │
           │ N:M (through vocabulary_tags)
           │
┌──────────┴───────────┐
│   tags               │
├──────────────────────┤
│ PK id                │
│    name              │
│    slug              │
│    created_at        │
│    updated_at        │
└──────────────────────┘


┌──────────────────────┐
│   vocabulary_tags    │
├──────────────────────┤
│ PK id                │
│ FK vocabulary_id     │
│ FK tag_id            │
│    created_at        │
└──────────────────────┘


┌──────────────────────┐
│   courses            │
├──────────────────────┤
│ PK id                │
│    uuid              │
│    title             │
│    description       │
│    thumbnail_url     │
│ FK jlpt_level_id     │
│    level_order       │
│    is_active         │
│    created_at        │
│    updated_at        │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────┴───────────┐
│   lessons            │
├──────────────────────┤
│ PK id                │
│    uuid              │
│ FK course_id         │
│    title             │
│    description       │
│    content           │  (Rich text/Markdown)
│    lesson_order      │
│    duration_minutes  │
│    is_active         │
│    created_at        │
│    updated_at        │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────┴───────────┐
│   exercises          │
├──────────────────────┤
│ PK id                │
│    uuid              │
│ FK lesson_id         │
│    title             │
│    description       │
│    exercise_type     │  (quiz, flashcard, writing, listening, matching)
│    difficulty_level  │  (easy, medium, hard)
│    exercise_order    │
│    time_limit_sec    │
│    is_active         │
│    created_at        │
│    updated_at        │
└──────────┬───────────┘
           │
           │ 1:N
           │
┌──────────┴───────────────┐
│   exercise_questions     │
├──────────────────────────┤
│ PK id                    │
│    uuid                  │
│ FK exercise_id           │
│ FK vocabulary_id (null)  │
│    question_text         │
│    question_type         │  (multiple_choice, fill_blank, true_false, matching)
│    correct_answer        │
│    options               │  (JSON: ["opt1", "opt2", ...])
│    hint                  │
│    explanation           │
│    points                │
│    question_order        │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_course_progress   │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│ FK course_id             │
│    enrollment_date       │
│    completion_date       │
│    status                │  (enrolled, in_progress, completed)
│    progress_percentage   │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_lesson_progress   │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│ FK lesson_id             │
│    status                │  (not_started, in_progress, completed)
│    started_at            │
│    completed_at          │
│    time_spent_minutes    │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_exercise_attempts │
├──────────────────────────┤
│ PK id                    │
│    uuid                  │
│ FK user_id               │
│ FK exercise_id           │
│    attempt_number        │
│    started_at            │
│    completed_at          │
│    score                 │
│    max_score             │
│    percentage            │
│    time_taken_seconds    │
│    is_passed             │
│    created_at            │
│    updated_at            │
└──────────┬───────────────┘
           │
           │ 1:N
           │
┌──────────┴───────────────────┐
│   user_question_answers      │
├──────────────────────────────┤
│ PK id                        │
│ FK attempt_id                │
│ FK question_id               │
│    user_answer               │
│    is_correct                │
│    points_earned             │
│    time_taken_seconds        │
│    created_at                │
└──────────────────────────────┘


┌──────────────────────────┐
│   user_vocabulary_status │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│ FK vocabulary_id         │
│    status                │  (learning, reviewing, mastered)
│    mastery_level         │  (0-100)
│    review_count          │
│    correct_count         │
│    incorrect_count       │
│    last_reviewed_at      │
│    next_review_at        │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_achievements      │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│ FK achievement_id        │
│    earned_at             │
│    created_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   achievements           │
├──────────────────────────┤
│ PK id                    │
│    code                  │
│    name                  │
│    description           │
│    icon_url              │
│    badge_url             │
│    criteria              │  (JSON)
│    points                │
│    is_active             │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_streaks           │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│    current_streak        │
│    longest_streak        │
│    last_activity_date    │
│    created_at            │
│    updated_at            │
└──────────────────────────┘


┌──────────────────────────┐
│   user_daily_stats       │
├──────────────────────────┤
│ PK id                    │
│ FK user_id               │
│    date                  │
│    words_learned         │
│    lessons_completed     │
│    exercises_completed   │
│    time_spent_minutes    │
│    points_earned         │
│    created_at            │
│    updated_at            │
└──────────────────────────┘
```

---

## Relationships Summary

### Core User Management
- **users** → **roles** (Many-to-One)
  - Each user has one role (student, teacher, admin)

### Learning Content Hierarchy
- **jlpt_levels** → **categories** (One-to-Many)
  - Each JLPT level has multiple categories
- **categories** → **vocabularies** (One-to-Many)
  - Each category contains multiple vocabularies
- **jlpt_levels** → **vocabularies** (One-to-Many)
  - Each vocabulary is classified by JLPT level
- **vocabularies** ↔ **tags** (Many-to-Many through vocabulary_tags)
  - Vocabularies can have multiple tags for better organization

### Course Structure
- **jlpt_levels** → **courses** (One-to-Many)
  - Each course is aligned with a JLPT level
- **courses** → **lessons** (One-to-Many)
  - Each course contains multiple lessons
- **lessons** → **exercises** (One-to-Many)
  - Each lesson has multiple exercises
- **exercises** → **exercise_questions** (One-to-Many)
  - Each exercise contains multiple questions
- **vocabularies** → **exercise_questions** (One-to-Many, optional)
  - Questions can reference specific vocabulary items

### User Progress Tracking
- **users** → **user_course_progress** (One-to-Many)
  - Track which courses a user is enrolled in
- **users** → **user_lesson_progress** (One-to-Many)
  - Track lesson completion status
- **users** → **user_exercise_attempts** (One-to-Many)
  - Record all exercise attempts
- **user_exercise_attempts** → **user_question_answers** (One-to-Many)
  - Detailed answers for each question in an attempt

### Vocabulary Learning
- **users** → **user_vocabulary_status** (One-to-Many)
  - Track mastery level and review schedule for each vocabulary

### Gamification
- **users** → **user_achievements** (One-to-Many)
  - Track earned achievements
- **achievements** → **user_achievements** (One-to-Many)
  - Reference to achievement definitions
- **users** → **user_streaks** (One-to-One)
  - Track daily learning streaks
- **users** → **user_daily_stats** (One-to-Many)
  - Daily learning statistics

---

## Key Features Supported

### 1. Vocabulary Management
- ✅ Kanji, Hiragana, Romaji representation
- ✅ Multi-language meanings (ID, EN)
- ✅ Example sentences
- ✅ Audio pronunciation
- ✅ Images for visual learning
- ✅ JLPT level classification
- ✅ Category organization
- ✅ Tagging system

### 2. Course & Lesson System
- ✅ JLPT-aligned courses
- ✅ Sequential lesson structure
- ✅ Rich content support
- ✅ Duration tracking

### 3. Exercise Types
- ✅ Multiple choice
- ✅ Fill in the blank
- ✅ True/False
- ✅ Matching
- ✅ Flashcards
- ✅ Writing practice
- ✅ Listening comprehension

### 4. Progress Tracking
- ✅ Course enrollment & completion
- ✅ Lesson progress
- ✅ Exercise attempts with scoring
- ✅ Detailed question-level answers
- ✅ Time tracking

### 5. Spaced Repetition
- ✅ Vocabulary mastery levels
- ✅ Review scheduling
- ✅ Success/failure tracking

### 6. Gamification
- ✅ Achievements & badges
- ✅ Daily streaks
- ✅ Points system
- ✅ Daily statistics

### 7. Analytics
- ✅ Learning time tracking
- ✅ Performance metrics
- ✅ Progress percentage
- ✅ Daily activity logs

---

## Indexes Recommendations

### High Priority Indexes
```sql
-- User lookups
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_uuid ON users(uuid);

-- Vocabulary searches
CREATE INDEX idx_vocabularies_jlpt_level ON vocabularies(jlpt_level_id);
CREATE INDEX idx_vocabularies_category ON vocabularies(category_id);
CREATE INDEX idx_vocabularies_word_type ON vocabularies(word_type);
CREATE INDEX idx_vocabularies_active ON vocabularies(is_active);

-- Course/Lesson ordering
CREATE INDEX idx_courses_level_order ON courses(jlpt_level_id, level_order);
CREATE INDEX idx_lessons_course_order ON lessons(course_id, lesson_order);
CREATE INDEX idx_exercises_lesson_order ON exercises(lesson_id, exercise_order);

-- Progress tracking
CREATE INDEX idx_user_course_progress ON user_course_progress(user_id, course_id);
CREATE INDEX idx_user_lesson_progress ON user_lesson_progress(user_id, lesson_id);
CREATE INDEX idx_user_vocabulary_status ON user_vocabulary_status(user_id, vocabulary_id);
CREATE INDEX idx_user_vocabulary_next_review ON user_vocabulary_status(user_id, next_review_at);

-- Exercise attempts
CREATE INDEX idx_exercise_attempts_user ON user_exercise_attempts(user_id, exercise_id);
CREATE INDEX idx_question_answers_attempt ON user_question_answers(attempt_id);

-- Daily stats
CREATE INDEX idx_daily_stats_user_date ON user_daily_stats(user_id, date);
```

---

## Notes

1. **UUID Fields**: Used for public-facing IDs to prevent enumeration attacks
2. **Soft Deletes**: Consider adding `deleted_at` fields for important entities
3. **Audit Fields**: All tables include `created_at` and `updated_at`
4. **JSON Fields**: `options`, `criteria` use JSON for flexibility
5. **Status Enums**: Use constants in code for status values
6. **Time Tracking**: Multiple granularities (minutes, seconds) for different contexts
7. **Multi-language**: Support for Indonesian and English (expandable)

---

## Next Steps

1. Create migration files for each entity
2. Define GORM models in `domain/models/`
3. Create DTOs in `domain/dto/`
4. Implement repositories
5. Create seed data for JLPT levels and basic vocabulary
6. Build API endpoints
