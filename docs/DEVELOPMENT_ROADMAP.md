# Development Roadmap - Manabu App

## 🎯 Quick Overview

```
Overall Progress: ████░░░░░░░░░░░░ 27% (4/15 APIs)

MVP Progress:     ████████████░░░░ 60% (3/5 APIs)
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

---

### 🔄 CURRENT WEEK

#### Week 2 (Learning Content - Part 1)
- 🔴 **User Vocabulary Status API** ← YOU ARE HERE
  - Priority: CRITICAL
  - Complexity: HIGH
  - 5 endpoints
  - Estimated: 2 days
  - **Note**: Final API to complete MVP Backend!

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

### Milestone 1: MVP Backend ✅ 60%
**Target**: Week 2
**Status**: IN PROGRESS

**Requirements**:
- [x] User Auth
- [x] JLPT Levels
- [x] Categories
- [x] Vocabulary CRUD
- [ ] User Vocabulary Status (Spaced Repetition)

**Blockers**: None
**Risks**: None - On track!

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
| 6  | User Vocab Status        | 🔴     | CRITICAL | 5         | HIGH       | 2 days    |
| 7  | Courses                  | 🔴     | HIGH     | 8         | MEDIUM     | 2 days    |
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

### This Week
4. 🔄 **Start User Vocabulary Status API** ← NEXT (CRITICAL - Final MVP API!)
5. Start Courses API

### Next Week
5. Courses & Lessons APIs
6. Exercises & Questions APIs

---

## 💡 Development Best Practices (Learned)

### From Categories, Vocabulary & Tags API Implementation:
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

### Week 2 Progress (Updated)
- **Completed**: Tags API (DONE!)
- **Current**: Ready for User Vocabulary Status API
- **Target**: Complete MVP Backend (5/5 APIs)
- **Status**: Excellent progress! On track for MVP completion.

### Lessons Learned
- Code review process adds quality but takes time
- golang-pro agent produces high-quality code from start (9.0/10 initial score)
- Fixing issues immediately prevents technical debt
- Database indexes should be added from the beginning
- Specific error constants improve API clarity (ErrInvalidTagName > ErrInvalidID)
- Case-insensitive validation crucial for user-facing name fields
- Remove dead code during review to maintain clean codebase

---

## 🎯 Success Criteria

### MVP Launch (Week 2)
- [x] All CRITICAL APIs completed (3/5) - In Progress
  - [x] User Management
  - [x] JLPT Levels
  - [x] Categories
  - [x] Vocabulary CRUD
  - [ ] User Vocabulary Status (Spaced Repetition)
- [x] All endpoints have Swagger docs
- [x] All code reviewed by expert agent
- [x] No critical or major issues
- [x] Can create, learn, and review vocabulary
- [ ] Basic spaced repetition working

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

**Current Sprint**: Foundation Phase (Week 1-2)
**Sprint Goal**: Complete MVP Backend APIs (80% Complete!)
**Last Updated**: 2026-01-05
**Last Completed**: Tags API ✅ (Quality: 9.5/10)
**Next Target**: User Vocabulary Status API (Final MVP API!)
