# todo-rank — Rendering Spec (Step 4)

이 문서는 `todo-rank` 프로젝트의 **4단계(출력/렌더링 규칙)**를
다른 에이전트나 자동화 도구에 그대로 전달하기 위한 **실행 가능한 스펙**이다.

이 단계의 목표는:

* 스코어링이 끝난 TODO 목록을
* **사람이 바로 행동할 수 있는 형태**로 출력하는 것이다.

---

## 1. 목표 (Goal)

입력으로 받은 `[]ScoredTodo`를 다음 기준에 따라 출력한다:

* 일관된 정렬 규칙
* 우선순위(Level) 기반 그룹핑
* 두 가지 출력 포맷 지원 (`text`, `md`)

렌더링 단계는 **데이터를 바꾸지 않는다**.
(정렬·그룹핑·포맷만 담당)

---

## 2. 입력 데이터 (Input)

```go
type ScoredTodo struct {
    Tag      string
    Priority *int
    Message  string
    Path     string
    Line     int

    Score int
    Level string // "P0", "P1", "P2", "P3"
}
```

입력은 이미:

* 파싱 완료
* 스코어 계산 완료
* Level 산출 완료

상태라고 가정한다.

---

## 3. 정렬 규칙 (Sorting Rules)

렌더링 전에 반드시 아래 규칙으로 정렬한다:

1. `Score` 내림차순 (desc)
2. `Path` 오름차순 (asc)
3. `Line` 오름차순 (asc)

의도:

* 가장 급한 항목이 위로
* 동일 점수에서는 파일 단위로 묶임

---

## 4. 그룹핑 규칙 (Grouping Rules)

정렬된 결과를 `Level` 기준으로 그룹핑한다.

그룹 순서는 고정:

```text
P0 → P1 → P2 → P3
```

각 그룹에 포함된 항목이 없으면
**해당 그룹은 출력하지 않는다**.

---

## 5. Text 출력 포맷

### 5.1 전체 구조

```text
P0 (Now)
<items>

P1 (Soon)
<items>
```

---

### 5.2 아이템 포맷

```text
(<Score>) <Path>:<Line>
  <Tag>[P?]: <Message>
```

* Priority가 없는 경우 `[P?]` 부분은 생략
* Message는 trim된 상태로 출력

---

### 5.3 Text 출력 예시

```text
P0 (Now)
(185) cmd/server/main.go:42
  FIXME[P0]: data race when shutting down

P1 (Soon)
(120) pkg/auth/token.go:88
  TODO[P1]: refresh token expiry logic
```

---

## 6. Markdown 출력 포맷

Markdown 출력은 **체크리스트 기반**으로 작성한다.

### 6.1 그룹 헤더

```md
## P0 (Now)
```

---

### 6.2 아이템 포맷

```md
- [ ] <Path>:<Line>  
  <Tag>[P?]: <Message>
```

* 줄바꿈을 위해 **두 칸 공백 + 개행**을 사용한다
* Priority가 없는 경우 `[P?]` 부분은 생략

---

### 6.3 Markdown 출력 예시

```md
## P0 (Now)
- [ ] cmd/server/main.go:42  
  FIXME[P0]: data race when shutting down

## P1 (Soon)
- [ ] pkg/auth/token.go:88  
  TODO[P1]: refresh token expiry logic
```

---

## 7. 파일 출력 규칙

* `--out` 옵션이 없으면 stdout으로 출력
* `--out` 옵션이 있으면:

  * 파일을 새로 생성하거나 덮어쓴다
  * 출력 인코딩은 UTF-8

렌더링 로직은 **writer(io.Writer)** 를 받아 처리하는 것을 권장한다.

---

## 8. 함수 시그니처 (권장)

```go
func RenderText(w io.Writer, items []ScoredTodo) error
func RenderMarkdown(w io.Writer, items []ScoredTodo) error
```

또는 통합형:

```go
func Render(w io.Writer, items []ScoredTodo, format string) error
```

---

## 9. 명시적 Non-Goals (Rendering 단계)

이 단계에서는 **하지 않는다**:

* 컬러 출력 (ANSI)
* 아이템 접기/펼치기
* 인터랙티브 UI
* 링크 자동 생성

---

## 10. 설계 철학 요약

* 출력은 **단순하고 예측 가능**해야 한다
* 텍스트는 터미널 친화적
* Markdown은 바로 TODO 리스트로 사용 가능

> 이 출력은 **보고 끝나는 리포트가 아니라**,
> 바로 행동으로 이어지는 목록이다.
