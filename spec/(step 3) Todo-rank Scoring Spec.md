# todo-rank — Scoring Spec (Step 3)

이 문서는 `todo-rank` 프로젝트의 **3단계(우선순위 점수 계산 로직)**를
다른 에이전트나 자동화 도구에 그대로 전달하기 위한 **실행 가능한 스펙**이다.

이 단계의 목표는:

* TODO 하나당 **단일 정수 점수(score)**를 산출하고
* 그 점수가 **왜 그렇게 나왔는지 설명 가능**하게 만드는 것이다.

---

## 1. 목표 (Goal)

각 `TodoItem`에 대해 다음을 계산한다:

* `Score int`  : 정렬에 사용되는 최종 점수
* `Level`     : 사람이 이해하기 쉬운 우선순위 레벨 (`P0~P3`)

이 점수는 **결정 기준**이며, 절대적인 진실이 아니다.

---

## 2. 입력 데이터 (Input)

스코어 계산 함수는 다음 정보를 입력으로 받는다:

```go
type TodoItem struct {
    Tag      string   // TODO | FIXME | @next
    Priority *int     // 0~3 or nil
    Message  string
    Path     string   // file path (relative)
    Line     int
}
```

---

## 3. 출력 데이터 (Output)

```go
type ScoredTodo struct {
    TodoItem
    Score int
    Level string // "P0", "P1", "P2", "P3"
}
```

---

## 4. 점수 계산 개요

최종 점수는 다음 요소들의 **합(sum)** 으로 계산한다:

```text
Score = TagScore
      + PriorityScore
      + PathBonus
      + TestPenalty
```

각 항목은 독립적으로 계산되며,
어느 하나도 다른 항목을 덮어쓰지 않는다.

---

## 5. Tag 기본 점수 (TagScore)

| Tag   | Score |
| ----- | ----- |
| FIXME | +100  |
| TODO  | +50   |
| @next | +30   |

* Tag는 반드시 존재한다고 가정한다 (파서 단계에서 보장)
* 정의되지 않은 Tag는 **0점**으로 처리한다

---

## 6. 명시적 우선순위 점수 (PriorityScore)

Priority는 **개발자의 의도**를 직접 표현한 것이므로
Tag 점수와 **독립적으로 추가 가중치**를 부여한다.

| Priority | Score |
| -------- | ----- |
| P0       | +100  |
| P1       | +70   |
| P2       | +40   |
| P3       | +10   |

* Priority가 없으면 `0점`
* Priority 값은 파서에서 `0~3`으로 정규화되어 전달된다

---

## 7. 경로 기반 보정 (PathBonus)

### 7.1 메인 경로 보너스

다음 조건 중 하나라도 만족하면 `+20`을 부여한다:

* 파일 경로에 `cmd/`가 포함됨
* 파일명이 `main.go`임

의도:

> 실행 흐름의 시작점에 가까운 TODO는 더 위험하다.

---

## 8. 테스트 파일 패널티 (TestPenalty)

다음 조건을 만족하면 `-20`을 부여한다:

* 파일명이 `*_test.*` 패턴과 매칭됨

의도:

> 테스트 TODO는 즉시 장애를 유발할 가능성이 낮다.

---

## 9. 점수 계산 순서 (권장)

디버깅 및 설명 가능성을 위해 **아래 순서로 계산한다**:

1. TagScore
2. PriorityScore
3. PathBonus
4. TestPenalty

중간 점수 로그를 남길 수 있도록
각 항목은 **개별 함수**로 구현하는 것을 권장한다.

---

## 10. 우선순위 레벨(Level) 산출

최종 점수를 기준으로 다음 레벨을 부여한다:

| Score Range | Level        |
| ----------- | ------------ |
| ≥ 120       | P0 (Now)     |
| ≥ 80        | P1 (Soon)    |
| ≥ 40        | P2 (Later)   |
| < 40        | P3 (Cleanup) |

* Level은 **렌더링 및 그룹핑 전용 값**이다
* 정렬은 항상 `Score desc` 기준으로 한다

---

## 11. 함수 시그니처 (권장)

```go
func ScoreTodo(item TodoItem) ScoredTodo
```

또는 분해형:

```go
func CalcTagScore(tag string) int
func CalcPriorityScore(p *int) int
func CalcPathBonus(path string) int
func CalcTestPenalty(path string) int
func CalcLevel(score int) string
```

---

## 12. 계산 예시

### Input

```text
Tag: TODO
Priority: P1
Path: cmd/server/main.go
```

### Calculation

```text
TagScore        = 50
PriorityScore   = 70
PathBonus       = 20
TestPenalty     = 0
----------------------
Total Score     = 140
Level           = P0
```

---

## 13. 명시적 Non-Goals (Scoring 단계)

이 단계에서는 **하지 않는다**:

* Git 히스토리 기반 가중치
* 동일 파일 내 TODO 개수 보정
* 시간 기반 decay
* 사용자 정의 수식

---

## 14. 설계 철학 요약

* 점수는 **설명 가능해야 한다**
* 규칙은 **간단하고 고정적**이어야 한다
* 개발자의 명시적 의도를 **자동 추론보다 우선**한다

> 이 스코어는 AI가 아니라 **결정을 돕는 규칙**이다.
