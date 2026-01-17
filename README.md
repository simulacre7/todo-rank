# todo-rank

> Scan TODOs. Rank what matters now.

`todo-rank`는 코드베이스에 흩어진 `TODO`, `FIXME`, `@next` 주석을 스캔하여 **우선순위 점수 기반으로 정렬**해주는 CLI 도구입니다.

## 설치

```bash
go install github.com/simulacre7/todo-rank/cmd/todo-rank@latest
```

> **참고**: `go install`로 설치한 바이너리는 `$HOME/go/bin`에 위치합니다.
> `command not found` 오류가 발생하면 PATH에 Go bin 디렉토리를 추가하세요.
>
> ```bash
> # ~/.zshrc 또는 ~/.bashrc에 추가
> export PATH="$HOME/go/bin:$PATH"
> ```

또는 소스에서 빌드:

```bash
git clone https://github.com/simulacre7/todo-rank.git
cd todo-rank
go build -o todo-rank ./cmd/todo-rank
```

### 쉘 Alias 설정 (선택)

자주 사용하는 옵션을 alias로 등록하면 편리합니다.

**Bash** (`~/.bashrc`):

```bash
alias tr='todo-rank'
alias trp='todo-rank --min-score 80'  # P0, P1만 보기
alias trmd='todo-rank --format md'    # Markdown 출력
```

**Zsh** (`~/.zshrc`):

```bash
alias tr='todo-rank'
alias trp='todo-rank --min-score 80'
alias trmd='todo-rank --format md'
```

**Fish** (`~/.config/fish/config.fish`):

```fish
alias tr 'todo-rank'
alias trp 'todo-rank --min-score 80'
alias trmd 'todo-rank --format md'
```

설정 후 쉘을 재시작하거나 `source` 명령어로 적용:

```bash
source ~/.zshrc  # 또는 ~/.bashrc
```

## 사용법

### 기본 실행

```bash
# 현재 디렉토리 스캔
todo-rank

# scan 서브커맨드 (동일한 동작)
todo-rank scan
```

### 옵션

| 옵션 | 설명 | 기본값 |
|------|------|--------|
| `--root` | 스캔 시작 디렉토리 | `.` |
| `--ignore` | 무시할 디렉토리 (쉼표 구분) | `.git,node_modules,dist` |
| `--format` | 출력 포맷 (`text` 또는 `md`) | `text` |
| `--out` | 결과 저장 경로 | stdout |
| `--min-score` | 최소 점수 필터 | `0` |
| `--tags` | 스캔할 태그 (쉼표 구분) | `TODO,FIXME,@next` |

### 예시

```bash
# Markdown 파일로 출력
todo-rank --format md --out TODOs.md

# 높은 우선순위만 필터링 (P0, P1)
todo-rank --min-score 80

# 특정 디렉토리 제외
todo-rank --ignore .git,node_modules,dist,vendor

# 특정 디렉토리에서 스캔
todo-rank --root ./src
```

## 지원 주석 포맷

```go
// TODO: something
// TODO[P1]: improve error handling
// FIXME[P0]: data race here
// @next: refactor naming

# TODO: Python style comment
/* TODO: block comment */
TODO - alternative separator
```

### 규칙

- **태그**: `TODO`, `FIXME`, `@next` (대소문자 구분)
- **우선순위**: `[P0]` ~ `[P3]` (선택사항, 태그 바로 뒤에 위치)
- **구분자**: `:` 또는 `-` 필수
- 주석 마커(`//`, `#`, `/* */`)는 선택사항

## 스코어링 규칙

### 태그 기본 점수

| 태그 | 점수 |
|------|------|
| FIXME | +100 |
| TODO | +50 |
| @next | +30 |

### 명시적 우선순위 (추가 점수)

| 우선순위 | 점수 |
|----------|------|
| P0 | +100 |
| P1 | +70 |
| P2 | +40 |
| P3 | +10 |

### 경로 보정

| 조건 | 점수 |
|------|------|
| `cmd/` 포함 또는 `main.go` | +20 |
| `*_test.*` 파일 | -20 |

### 최종 레벨

| 점수 범위 | 레벨 |
|-----------|------|
| >= 120 | P0 (Now) |
| >= 80 | P1 (Soon) |
| >= 40 | P2 (Later) |
| < 40 | P3 (Cleanup) |

## 출력 예시

### Text (기본)

```
P0 (Now)
(185) cmd/server/main.go:42
  FIXME[P0]: data race when shutting down

P1 (Soon)
(120) pkg/auth/token.go:88
  TODO[P1]: refresh token expiry logic
```

### Markdown

```md
## P0 (Now)
- [ ] cmd/server/main.go:42
  FIXME[P0]: data race when shutting down

## P1 (Soon)
- [ ] pkg/auth/token.go:88
  TODO[P1]: refresh token expiry logic
```

## 프로젝트 구조

```
todo-rank/
├── cmd/
│   └── todo-rank/
│       └── main.go          # CLI 진입점, 옵션 파싱
├── internal/
│   ├── parse/
│   │   ├── parse.go         # TODO 주석 파싱 (정규식)
│   │   └── parse_test.go    # 파서 테스트
│   ├── scan/
│   │   ├── options.go       # ScanOptions 구조체
│   │   └── scan.go          # 디렉토리 순회, 파일 읽기
│   ├── score/
│   │   └── score.go         # 점수 계산 로직
│   └── render/
│       └── render.go        # 출력 포맷팅 (text/md)
├── spec/                    # 설계 문서
├── go.mod
└── README.md
```

### 모듈 책임

| 모듈 | 책임 |
|------|------|
| `cmd/todo-rank` | CLI 옵션 파싱, 진입점 |
| `internal/parse` | 한 줄에서 TODO 정보 추출 |
| `internal/scan` | 디렉토리 순회, 파일 읽기, 결과 수집 |
| `internal/score` | 점수 계산, 레벨 산출 |
| `internal/render` | 정렬, 그룹핑, 텍스트/마크다운 출력 |

## 개발 기여

### 요구사항

- Go 1.21 이상

### 빌드 및 테스트

```bash
# 빌드
go build ./...

# 테스트
go test ./...

# 실행
go run ./cmd/todo-rank
```

### 코드 스타일

- 표준 Go 포맷팅 (`gofmt`)
- 각 모듈은 단일 책임 원칙을 따름
- 외부 의존성 최소화 (표준 라이브러리만 사용)

### 기여 절차

1. 이슈 생성 또는 기존 이슈 확인
2. 브랜치 생성 (`feature/xxx` 또는 `fix/xxx`)
3. 변경사항 구현 및 테스트
4. PR 생성

## 설계 철학

- TODO를 **관리**하지 않고, **결정 가능한 순서**로 만든다
- 개발자의 **명시적 의도([P0~P3])를 최우선 존중**
- 복잡한 설정 없이 **바로 실행 가능**
- 점수는 **설명 가능하고 예측 가능**해야 한다

## 라이선스

MIT
