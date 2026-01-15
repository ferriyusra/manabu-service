# Development Roadmap - Manabu App

## 🎯 Quick Overview

```
Overall Progress: ██████░░░░░░░░░░ 40% (6/15 APIs)

MVP Progress:     ████████████████ 100% (5/5 APIs) ✅ COMPLETE!
Learning Phase:   ████░░░░░░░░░░░░ 25% (1/4 APIs)
```

---

## 📅 Development Timeline

### ✅ COMPLETED

#### Week 0 (Initial Setup)
- ✅ User Management API (Auth)
  - Login, Register, Get User, Update User

#### Week 1 (Foundation - Part 1)
- ✅ JLPT Levels API
  - Full CRUD with Swagger
- ✅ Categories API
  - Full CRUD with pagination
  - Code review & optimizations applied
  - **Quality**: Production-ready

#### Week 1 (Foundation - Part 2)
- ✅ Vocabulary API
  - 5 CRUD endpoints with advanced filtering
  - Search, sort, pagination implemented
  - Code reviewed & all issues fixed
  - Database indexes optimized
  - URL validation & godoc comments added
  - **Quality**: Production-ready (9.5/10)

#### Week 2 (Foundation - Part 3)
- ✅ Tags API
  - 6 CRUD endpoints with advanced features
  - Case-insensitive search & duplicate check
  - Hex color validation with regex
  - Code reviewed & all major issues fixed
  - Proper error handling (ErrInvalidTagName)
  - **Quality**: Production-ready (9.0/10 → 9.5/10 after fixes)

#### Week 2 (MVP Completion)
- ✅ **User Vocabulary Status API** 🎉 **MVP COMPLETE!**
  - 5 endpoints with simple progress tracking
  - Simple boolean tracking (isCorrect: true/false)
  - 5 correct answers = completed status
  - User-controlled scheduling (no automatic intervals)
  - Removed SM-2 algorithm per user preference
  - Cleaned unused fields from response (easeFactor, interval, nextReviewDate)
  - Updated Status constraint (learning/completed)
  - Changed endpoint to use vocabulary_id (snake_case)
  - Created migration & comprehensive documentation
  - **Quality**: Production-ready (simplified & user-friendly)

#### Week 3 (Learning Content - Part 1)
- ✅ **Courses API** 🎉 **COMPLETED!**
  - 8 endpoints (CRUD + publish/unpublish + get published)
  - Advanced filtering (JLPT level, difficulty, published status, search)
  - Case-insensitive search & uniqueness check
  - Publish/Unpublish functionality with timestamp management
  - Code reviewed & all critical/major issues fixed
  - SQL injection prevention with whitelist validation
  - Proper error constants (including ErrInvalidCourseEstimatedHours)
  - Database migration created (004_create_courses_table.sql)
  - Complete Swagger documentation
  - **Quality**: Production-ready (8.5/10 → 9.3/10 after fixes)

---

### 🔄 CURRENT WEEK

#### Week 3 (Learning Content - Part 2)
- 🔴 **Lessons API** ← YOU ARE HERE
  - Priority: HIGH
  - Complexity: MEDIUM
  - 6 endpoints
  - Next target after Courses API

---

### 📋 UPCOMING

#### Week 2 (Learning Content - Part 2)
- 🔴 Courses API

#### Week 3 (Learning Content - Part 2)
- 🔴 Lessons API
- 🔴 Exercises API
- 🔴 Exercise Questions API

#### Week 4 (User Progress)
- 🔴 User Course Progress API
- 🔴 User Lesson Progress API
- 🔴 User Exercise Attempts API

#### Week 5-6 (Gamification)
- 🔴 Achievements API
- 🔴 Streaks API
- 🔴 Statistics API

---

## 🎯 Milestones

### Milestone 1: MVP Backend ✅ 100% 🎉 COMPLETE!
**Target**: Week 2
**Status**: ✅ COMPLETED

**Requirements**:
- [x] User Auth
- [x] JLPT Levels
- [x] Categories
- [x] Vocabulary CRUD
- [x] User Vocabulary Status (Simple Progress Tracking)

**Achievement**: All MVP APIs delivered with production-ready quality!

---

### Milestone 2: Learning Platform 🔴 0%
**Target**: Week 4
**Status**: NOT STARTED

**Requirements**:
- [ ] Tags
- [ ] Courses
- [ ] Lessons
- [ ] Exercises
- [ ] Exercise Questions

**Dependencies**: ✅ Vocabulary API completed!

---

### Milestone 3: Progress Tracking 🔴 0%
**Target**: Week 5
**Status**: NOT STARTED

**Requirements**:
- [ ] User Course Progress
- [ ] User Lesson Progress
- [ ] User Exercise Attempts

**Dependencies**: Learning Platform APIs

---

### Milestone 4: User Engagement 🔴 0%
**Target**: Week 6
**Status**: NOT STARTED

**Requirements**:
- [ ] Achievements
- [ ] Streaks
- [ ] Statistics

**Dependencies**: Progress Tracking APIs

---

## 📊 API Status Dashboard

| #  | API Name                  | Status | Priority | Endpoints | Complexity | ETA       |
|----|---------------------------|--------|----------|-----------|------------|-----------|
| 1  | User Management           | ✅     | CRITICAL | 5         | MEDIUM     | DONE      |
| 2  | JLPT Levels              | ✅     | HIGH     | 5         | LOW        | DONE      |
| 3  | Categories               | ✅     | HIGH     | 6         | MEDIUM     | DONE      |
| 4  | Vocabulary               | ✅     | CRITICAL | 5         | HIGH       | DONE      |
| 5  | Tags                     | ✅     | MEDIUM   | 6         | MEDIUM     | DONE      |
| 6  | User Vocab Status        | ✅     | CRITICAL | 5         | HIGH       | DONE      |
| 7  | Courses                  | ✅     | HIGH     | 8         | MEDIUM     | DONE      |
| 8  | Lessons                  | 🔴     | HIGH     | 6         | MEDIUM     | 1 day     |
| 9  | Exercises                | 🔴     | HIGH     | 6         | MEDIUM     | 2 days    |
| 10 | Exercise Questions       | 🔴     | HIGH     | 5         | MEDIUM     | 1 day     |
| 11 | User Course Progress     | 🔴     | CRITICAL | 4         | MEDIUM     | 1 day     |
| 12 | User Lesson Progress     | 🔴     | CRITICAL | 3         | LOW        | 1 day     |
| 13 | User Exercise Attempts   | 🔴     | CRITICAL | 5         | HIGH       | 2 days    |
| 14 | Achievements             | 🔴     | MEDIUM   | 3         | LOW        | 1 day     |
| 15 | Streaks                  | 🔴     | MEDIUM   | 2         | LOW        | 0.5 day   |
| 16 | Statistics               | 🔴     | MEDIUM   | 4         | MEDIUM     | 1 day     |

**Total Estimated Time**: ~20 days (4 weeks)

---

## 🚀 Next Actions

### Immediate (Today)
1. ✅ Complete Categories API fixes
2. ✅ **Complete Vocabulary API development**
   - ✅ Use golang-pro agent
   - ✅ Code review with code-review-expert
   - ✅ Fix all critical & major issues
   - ✅ Apply all best practices from Categories API
3. ✅ **Complete Tags API development**
   - ✅ Use golang-pro agent
   - ✅ Code review with code-review-expert (Score: 9.0/10)
   - ✅ Fix all major issues (Error constants, dead code)
   - ✅ Quality improved to 9.5/10

4. ✅ **Complete User Vocabulary Status API** 🎉 MVP COMPLETE!
   - ✅ Use golang-pro agent
   - ✅ Simple tracking system (boolean instead of SM-2)
   - ✅ Clean response (removed unused fields)
   - ✅ Fix Status constraint (learning/completed)
   - ✅ Update endpoint to use vocabulary_id (snake_case)
   - ✅ Create migration & documentation

5. ✅ **Complete Courses API** 🎉 Learning Platform Phase Started!
   - ✅ Use golang-pro agent (8 endpoints implemented)
   - ✅ Code review with code-review-expert (Initial score: 8.5/10)
   - ✅ Fix all critical issues (SQL injection, published_at clearing, error constants)
   - ✅ Fix all major issues (case-insensitive uniqueness check)
   - ✅ Create database migration (004_create_courses_table.sql)
   - ✅ Complete Swagger documentation
   - ✅ Quality improved to 9.3/10

### This Week
6. 🔄 **Start Lessons API** ← NEXT (Learning Platform Phase)

### Next Week
7. Exercises & Questions APIs
8. User Progress APIs

---

## 💡 Development Best Practices (Learned)

### From Categories, Vocabulary, Tags & Courses API Implementation:
✅ **DO**:
- Add unique constraints at database level
- Check RowsAffected in Update operations
- Use helper methods to avoid code duplication
- Implement proper HTTP status code mapping (404, 409, 422, 500)
- Add defensive validation in repository layer
- Optimize queries to avoid N+1 issues (use Preload)
- Create proper Swagger response types
- Use golang-pro agent for implementation
- Use code-review-expert for quality assurance
- Add database indexes for foreign keys
- Validate URLs with `url` validator tag
- Add godoc comments to all interfaces
- **Use whitelist validation for ORDER BY clauses** (prevent SQL injection)
- **Case-insensitive search and uniqueness checks** (LOWER() function)
- **Create specific error constants** for each validation type
- **Clear related timestamps** when unpublishing/deactivating
- Ensure validation logic is consistent (1-5, not 0-5)
- Create specific error constants (e.g., ErrInvalidTagName vs generic ErrInvalidID)
- Use case-insensitive search with LOWER() for name fields
- Validate special fields with regex (e.g., hex color patterns)

❌ **DON'T**:
- Rely only on service-layer validation (not atomic)
- Return generic 400 for all errors
- Duplicate model-to-DTO mapping code
- Skip defensive checks in repositories
- Create double-wrapped responses
- Allow validation inconsistencies between layers
- Reuse generic errors when specific ones are more appropriate
- Include redundant/dead code in conditional logic

---

## 📈 Velocity Tracking

### Week 1 Velocity (Final)
- **Planned**: 4 APIs (JLPT, Categories, Vocabulary, Tags)
- **Completed**: 3 APIs (JLPT, Categories, Vocabulary)
- **Velocity**: 75%
- **Quality**: EXCELLENT (all code reviewed & fixed, score 9.5/10)

### Week 2 Progress (FINAL) 🎉
- **Completed**: Tags API + User Vocabulary Status API
- **Achievement**: 🎉 **MVP COMPLETE!** (5/5 APIs)
- **Velocity**: 100% - All critical APIs delivered!
- **Quality**: Production-ready across all APIs

### New Lessons Learned (User Vocabulary Status API)
- **Simplicity over complexity**: User chose simple boolean tracking over SM-2 algorithm
- **Listen to user feedback**: Changed from complex spaced repetition to user-controlled review
- **Clean APIs**: Remove unused fields from responses (easeFactor, interval, nextReviewDate)
- **Intuitive endpoints**: Use vocabulary_id instead of status_id for better UX
- **Database constraints**: Update check constraints when business logic changes
- **Comprehensive docs**: Create feature-specific documentation when complexity requires understanding

### Previous Lessons
- Code review process adds quality but takes time
- golang-pro agent produces high-quality code from start (9.0/10 initial score)
- Fixing issues immediately prevents technical debt
- Database indexes should be added from the beginning
- Specific error constants improve API clarity (ErrInvalidTagName > ErrInvalidID)
- Case-insensitive validation crucial for user-facing name fields
- Remove dead code during review to maintain clean codebase

---

## 🎯 Success Criteria

### MVP Launch (Week 2) ✅ ACHIEVED!
- [x] All CRITICAL APIs completed (5/5) ✅
  - [x] User Management
  - [x] JLPT Levels
  - [x] Categories
  - [x] Vocabulary CRUD
  - [x] User Vocabulary Status (Simple Progress Tracking)
- [x] All endpoints have Swagger docs
- [x] All code reviewed by expert agent
- [x] No critical or major issues
- [x] Can create, learn, and review vocabulary
- [x] Simple progress tracking working

### Full Platform (Week 4)
- [ ] All HIGH priority APIs completed
- [ ] Course structure working
- [ ] Exercise system functional
- [ ] Progress tracking accurate

### Production Ready (Week 6)
- [ ] All 15 core APIs completed
- [ ] Gamification features live
- [ ] Performance optimized
- [ ] Security hardened
- [ ] Documentation complete

---

**Current Sprint**: Learning Platform Phase (Week 3)
**Sprint Goal**: Complete course/lesson structure APIs
**Last Updated**: 2026-01-14
**Last Completed**: Courses API ✅ 🎉 **Learning Platform Phase Started!**
**Next Target**: Lessons API (Learning Platform Phase)
