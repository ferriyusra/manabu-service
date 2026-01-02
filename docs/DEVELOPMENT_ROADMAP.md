# Development Roadmap - Manabu App

## 🎯 Quick Overview

```
Overall Progress: ██░░░░░░░░░░░░░░ 13.3% (2/15 APIs)

MVP Progress:     ████████░░░░░░░░ 40% (2/5 APIs)
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

---

### 🔄 CURRENT WEEK

#### Week 1 (Foundation - Part 2)
- 🔴 **Vocabulary API** ← YOU ARE HERE
  - Priority: CRITICAL
  - Complexity: HIGH
  - 12 endpoints with advanced filtering
  - Estimated: 2-3 days

- 🔴 **Tags API**
  - Priority: MEDIUM
  - Complexity: MEDIUM
  - 6 endpoints
  - Estimated: 1 day

---

### 📋 UPCOMING

#### Week 2 (Learning Content - Part 1)
- 🔴 User Vocabulary Status API (Spaced Repetition)
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

### Milestone 1: MVP Backend ✅ 40%
**Target**: Week 2
**Status**: IN PROGRESS

**Requirements**:
- [x] User Auth
- [x] JLPT Levels
- [x] Categories
- [ ] Vocabulary CRUD
- [ ] User Vocabulary Status (Spaced Repetition)

**Blockers**: None
**Risks**: Vocabulary API complexity

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

**Dependencies**: Vocabulary API must be completed

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
| 4  | **Vocabulary**           | 🔴     | CRITICAL | 12        | **HIGH**   | **3 days** |
| 5  | Tags                     | 🔴     | MEDIUM   | 6         | MEDIUM     | 1 day     |
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
2. 🔄 **Start Vocabulary API development**
   - Use golang-pro agent
   - Code review with code-review-expert
   - Apply all best practices from Categories API

### This Week
3. Complete Tags API
4. Start User Vocabulary Status (Spaced Repetition)

### Next Week
5. Courses & Lessons APIs
6. Exercises & Questions APIs

---

## 💡 Development Best Practices (Learned)

### From Categories API Implementation:
✅ **DO**:
- Add unique constraints at database level
- Check RowsAffected in Update operations
- Use helper methods to avoid code duplication
- Implement proper HTTP status code mapping (404, 409, 422, 500)
- Add defensive validation in repository layer
- Optimize queries to avoid N+1 issues
- Create proper Swagger response types
- Use golang-pro agent for implementation
- Use code-review-expert for quality assurance

❌ **DON'T**:
- Rely only on service-layer validation (not atomic)
- Return generic 400 for all errors
- Duplicate model-to-DTO mapping code
- Skip defensive checks in repositories
- Create double-wrapped responses

---

## 📈 Velocity Tracking

### Week 1 Velocity
- **Planned**: 4 APIs (JLPT, Categories, Vocabulary, Tags)
- **Completed**: 2 APIs (JLPT, Categories)
- **Velocity**: 50%
- **Quality**: HIGH (all code reviewed & fixed)

### Adjustments
- Vocabulary API is more complex than estimated
- Add extra time for code review & fixes
- Reduce scope for Week 2 if needed

---

## 🎯 Success Criteria

### MVP Launch (Week 2)
- [ ] All CRITICAL APIs completed (5/5)
- [ ] All endpoints have Swagger docs
- [ ] All code reviewed by expert agent
- [ ] No critical or major issues
- [ ] Can create, learn, and review vocabulary
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
**Sprint Goal**: Complete MVP Backend APIs
**Last Updated**: 2026-01-02
**Next Review**: After Vocabulary API completion
