# todo-rank — CLI Spec (Step 2)

이 문서는 `todo-rank` 프로젝트의 **2단계(CLI 옵션 파싱 및 실행 흐름)**를
다른 에이전트나 자동화 도구에 그대로 전달하기 위한 **실행 가능한 스펙**이다.

표준 Go `flag` 패키지를 사용하는 것을 전제로 하며,
**복붙 후 바로 구현 지시**로 사용할 수 있도록 작성되었다.

---

## 1. 목표 (Goal)

CLI 실행을 통해 다음 흐름을 완성한다:

```text
CLI 입력 → 옵션 파싱 → ScanOptions 구성 → scan.Run(options) 호출
```

이 단계의 목표는 **CLI UX를 고정**하고,
이후 로직(scan/parse/score/render)이 안정적으로 연결되도록 하는 것이다.

---

## 2. 기본 실행 규칙

### 2.1 실행 방식

```bash
todo-rank
todo-rank scan
```

* `todo-rank` 단독 실행 시 내부적으로 `scan`을 수행한다
* `scan`은 **논리적 서브커맨드**이며, cobra 등은 사용하지 않는다

---

## 3. 지원 옵션 (MVP)

모든 옵션은 **long flag**만 지원한다.

```text
--root <path>        스캔 시작 디렉토리 (default: .)
--ignore <csv>       무시할 디렉토리 목록
--format <text|md>   출력 포맷 (default: text)
--out <path>         결과 저장 경로 (default: stdout)
--min-score <n>      최소 점수 필터 (default: 0)
--tags <csv>         스캔할 태그 목록
```

### 3.1 기본값

| Option    | Default                  |
| --------- | ------------------------ |
| root      | `.`                      |
| ignore    | `.git,node_modules,dist` |
| format    | `text`                   |
| out       | (empty → stdout)         |
| min-score | `0`                      |
| tags      | `TODO,FIXME,@next`       |

---

## 4. 옵션 파싱 정책

### 4.1 `--ignore`

* CSV 문자열을 `,` 기준으로 split
* 디렉토리 이름 기준으로만 비교
* 경로에 포함되면 해당 디렉토리는 **Walk 자체를 skip**

예:

```text
--ignore .git,node_modules,dist
```

---

### 4.2 `--tags`

* CSV 문자열을 split
* 파서에 그대로 전달
* 파서 단계에서는 **허용 tag인지 검증하지 않는다**

예:

```text
--tags TODO,FIXME,@next,HACK
```

---

### 4.3 `--format`

* 허용 값: `text`, `md`
* 그 외 값이 들어오면 **즉시 에러 종료**

---

### 4.4 `--out`

* 값이 없으면 stdout
* 값이 있으면 해당 path에 파일 생성/덮어쓰기
* 디렉토리가 없으면 에러

---

## 5. ScanOptions 구조체 (공유 모델)

CLI 파싱 결과는 다음 구조체 하나로 모인다:

```go
type ScanOptions struct {
    Root      string
    Ignore    []string
    Format    string
    OutPath   string
    MinScore  int
    Tags      []string
}
```

---

## 6. 실행 흐름 (Control Flow)

### 6.1 main.go 책임

`cmd/todo-rank/main.go`의 역할은 **오직 다음만** 수행한다:

1. `os.Args`를 검사해 `scan` 여부 판단
2. `flag`로 옵션 파싱
3. `ScanOptions` 생성
4. `scan.Run(options)` 호출
5. 에러 발생 시 stderr 출력 후 `os.Exit(1)`

비즈니스 로직은 절대 포함하지 않는다.

---

## 7. 에러 처리 정책

* 잘못된 옵션 값 → 즉시 종료 + usage 출력
* `--format` invalid → 에러
* `--root` 존재하지 않음 → 에러
* `--out` 경로 생성 불가 → 에러

에러 메시지는 **짧고 직접적**이어야 한다.

---

## 8. 사용 예시

```bash
# 기본 스캔
todo-rank

# markdown 출력
todo-rank --format md --out TODOs.md

# 높은 우선순위만
todo-rank --min-score 80

# 특정 디렉토리 제외
todo-rank --ignore .git,node_modules,dist,vendor
```

---

## 9. 명시적 Non-Goals (CLI 단계)

이 단계에서는 **하지 않는다**:

* Short flag (`-f`, `-o`) 지원
* Subcommand 확장 (`list`, `config` 등)
* Config file 지원
* 인터랙티브 UI

---

## 10. 구현 체크리스트

* [ ] 표준 `flag` 패키지 사용
* [ ] 기본 실행 = scan
* [ ] 옵션 → `ScanOptions` 1곳에만 모으기
* [ ] 로직과 CLI 코드 분리

> 이 CLI는 **얇고 예측 가능**해야 한다.
> 똑똑함은 내부 로직의 몫이다.
