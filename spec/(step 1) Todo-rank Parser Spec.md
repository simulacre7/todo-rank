# todo-rank — Parser Spec (Step 1)

이 문서는 `todo-rank` 프로젝트의 **1단계(정규식 / 파서 설계)**를
다른 에이전트나 자동화 도구에 그대로 전달하기 위한 **실행 가능한 스펙**이다.

복붙해서 바로 구현 지시로 사용할 수 있도록 **규칙·입력·출력·알고리즘** 중심으로 작성되었다.

---

## 1. 목표 (Goal)

**한 줄의 텍스트(line)**에서 다음 정보를 추출한다:

* `Tag`      : TODO / FIXME / @next
* `Priority` : P0 ~ P3 (옵션)
* `Message`  : 주석 메시지 본문

파싱에 실패한 라인은 무시한다.

---

## 2. 지원 입력 포맷 (Supported Inputs)

다음 형태만 **정식 지원**한다 (MVP 범위).

```text
// TODO: something
# TODO: something
/* TODO: something */
TODO: something

// TODO[P1]: improve error handling
// FIXME[P0]: data race here
// @next: refactor naming
// @next[P2]: cleanup naming

// TODO - alternative separator is allowed
```

### 지원 규칙 요약

* 주석 마커(`//`, `#`, `/* */`)는 **있어도 되고 없어도 된다**
* 반드시 **구분자(`:` 또는 `-`)**가 있어야 한다
* 구분자 뒤는 전부 메시지로 취급한다

❌ 아래는 MVP에서 **지원하지 않는다**:

```text
TODO something        // 구분자 없음
TODOS: something     // tag 불일치
```

---

## 3. Tag 정의

지원 Tag는 **대소문자 구분**하며, 정확히 다음만 허용한다:

```text
TODO
FIXME
@next
```

> 추후 확장은 CLI 옵션(`--tags`)에서 처리한다.

---

## 4. Priority 정의

* 형식: `[P0]`, `[P1]`, `[P2]`, `[P3]`
* **옵션**이며, tag 바로 뒤에만 올 수 있다

예:

```text
TODO[P1]: ...     // OK
TODO [P1]: ...    // ❌ (공백 허용 안 함)
```

---

## 5. 파싱 규칙 (Algorithm)

### 5.1 처리 순서

1. 입력 line 전체를 문자열로 받는다
2. 정규식으로 아래 요소를 순서대로 캡처한다

   * Tag
   * Priority (optional)
   * Separator (`:` or `-`)
   * Message
3. 매칭 실패 시 `ok = false`
4. 성공 시 각 필드를 trim 후 반환

---

## 6. 정규식 스펙 (Regex Spec)

### 요구사항

* 주석 마커는 optional
* Tag는 반드시 캡처
* Priority는 optional 캡처
* Separator는 `:` 또는 `-`
* Message는 **빈 문자열 불가**

### 개념적 패턴

```regex
(optional comment)
(tag)
(optional priority)
(separator)
(message)
```

### Go `regexp` 기준 권장 패턴

> ⚠️ Go는 named group을 지원하지 않으므로 **group index 기반**으로 파싱한다.

```regex
^\s*(?:\/\/|#|\/\*)?\s*(TODO|FIXME|@next)(?:\[(P[0-3])\])?\s*[:\-]\s*(.+?)\s*(?:\*\/)?$
```

### Group Index 매핑

| Index | 의미                  |
| ----: | ------------------- |
|     1 | Tag                 |
|     2 | Priority (optional) |
|     3 | Message             |

---

## 7. 출력 데이터 모델 (Expected Output)

파서 함수는 다음 형태를 만든다:

```go
type ParsedTodo struct {
    Tag      string  // TODO | FIXME | @next
    Priority *int    // nil if not specified
    Message  string
}
```

Priority 변환 규칙:

```text
"P0" -> 0
"P1" -> 1
"P2" -> 2
"P3" -> 3
```

---

## 8. 함수 시그니처 (권장)

```go
func ParseLine(line string) (todo ParsedTodo, ok bool)
```

### 동작 규칙

* 매칭 실패 → `ok = false`
* 매칭 성공 → `ok = true`, `ParsedTodo` 채움

---

## 9. 파싱 예시

### Input

```text
// TODO[P1]: improve error handling
```

### Output

```go
ParsedTodo{
    Tag: "TODO",
    Priority: ptr(1),
    Message: "improve error handling",
}
```

---

## 10. 명시적 Non-Goals (Parser 단계)

이 단계에서는 **하지 않는다**:

* 파일 경로 처리
* 라인 번호 처리
* 스코어 계산
* Tag alias 처리 (e.g. `TODO!`)
* 멀티라인 주석 파싱

---

## 11. 구현 우선순위 요약

1. 정규식 정확성
2. false positive 최소화
3. 파싱 실패 시 조용히 skip

> 이 파서는 **엄격함(strict)**을 우선한다.
> 놓치는 TODO보다 잘못 잡는 TODO가 더 위험하다.
