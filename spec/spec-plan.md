# todo-rank

> Scan TODOs. Rank what matters now.

`todo-rank`는 코드베이스에 흩어져 있는 `TODO / FIXME / @next` 주석을 스캔해
**“지금 무엇을 먼저 고칠지”**를 점수 기반으로 정렬해주는 CLI 도구다.

이 도구는 TODO를 **관리**하지 않는다.
대신 TODO를 **결정 가능한 목록(actionable shortlist)**으로 바꾼다.

---

## 1. 문제 정의

대부분의 레포에는 수십~수백 개의 TODO/FIXME가 존재한다.

하지만:

* 무엇이 급한지 알기 어렵고
* 개발자 의도가 코드에 흩어져 있으며
* 테스트/헬퍼/메인 경로가 동일하게 취급된다

`todo-rank`는 다음 질문에 답한다:

> “이 레포에서 **지금** 손대야 할 TODO는 무엇인가?”

---

## 2. 목표 (Goals)

* TODO를 **우선순위 순서로 정렬**
* 개발자의 **명시적 의도([P0~P3])를 최우선 존중**
* 복잡한 설정 없이 **바로 실행 가능**
* 주말/스프린트 단위로 **즉시 쓸 수 있는 출력**

---

## 3. 비목표 (Non-Goals)

* TODO 히스토리 추적 (git blame/log)
* TODO 생성/완료 관리
* 이슈 트래커 연동
* IDE 플러그인

---

## 4. 지원 주석 포맷

```go
// TODO: something
// TODO[P1]: improve error handling
// FIXME[P0]: data race here
// @next: refactor naming
```

### 규칙

* Tag: `FIXME | TODO | @next`
* Optional priority: `[P0] ~ [P3]`
* 메시지는 `:` 이후 전체

---

## 5. 우선순위 스코어링 규칙 (v1.0)

### 5.1 태그 기본 점수

| Tag   | Score |
| ----- | ----- |
| FIXME | +100  |
| TODO  | +50   |
| @next | +30   |

### 5.2 명시 우선순위 (있을 경우 추가)

| Priority | Score |
| -------- | ----- |
| P0       | +100  |
| P1       | +70   |
| P2       | +40   |
| P3       | +10   |

> 최종 점수 = 태그 점수 + 우선순위 점수 + 위치 보정

### 5.3 위치 보정

* `cmd/` 또는 `main.go` 경로 포함: `+20`
* 테스트 파일(`*_test.*`): `-20`

---

## 6. 출력 우선순위 레벨

| Score Range | Level        |
| ----------- | ------------ |
| ≥ 120       | P0 (Now)     |
| ≥ 80        | P1 (Soon)    |
| ≥ 40        | P2 (Later)   |
| < 40        | P3 (Cleanup) |

---

## 7. CLI 스펙

### 기본 사용

```bash
todo-rank
todo-rank scan
```

### 옵션

```text
--root <path>        스캔 시작 디렉토리 (default: .)
--ignore <csv>       무시할 디렉토리 (default: .git,node_modules,dist)
--format <text|md>   출력 포맷 (default: text)
--out <path>         결과 파일로 저장 (default: stdout)
--min-score <n>      최소 점수 필터
--tags <csv>         스캔할 태그 (default: TODO,FIXME,@next)
```

---

## 8. 출력 예시

### Text

```
P0 (185) cmd/server/main.go:42
  FIXME[P0]: data race when shutting down

P1 (120) pkg/auth/token.go:88
  TODO[P1]: refresh token expiry logic
```

### Markdown

```md
## P0 (Now)
- [ ] cmd/server/main.go:42  
  FIXME[P0]: data race when shutting down
```

---

## 9. 내부 구조 (초안)

```text
cmd/
  todo-rank/
    main.go

internal/
  scan/     // 디렉토리 순회, 파일 읽기
  parse/    // 주석 파싱 (regex)
  score/    // 스코어 계산
  render/   // 출력(text/md)
```

---

## 10. 구현 계획 (주말 기준)

### Day 1

* CLI 옵션 파싱
* 디렉토리 스캔
* 주석 파싱
* `TodoItem` 모델 완성

### Day 2

* 스코어 계산
* 정렬/필터
* text/md 출력
* README 정리

---

## 11. 향후 아이디어 (v1.1+)

* `--git` 옵션: 최근 수정 파일 보정
* JSON 출력
* 파일/디렉토리별 통계
* 커스텀 스코어 규칙 설정

---

## 12. 철학 요약

> `todo-rank`는 TODO를 관리하지 않는다.
> **결정 가능한 순서로 만든다.**
