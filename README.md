# user api service

## API Endpoints

| Method | Path          | Description          |
|--------|---------------|----------------------|
| POST   | `/users`      | 새 사용자 생성       |
| GET    | `/users`      | 모든 사용자 조회     |
| GET    | `/users/{id}` | 특정 사용자 조회     |
| PUT    | `/users/{id}` | 특정 사용자 업데이트 |
| DELETE | `/users/{id}` | 특정 사용자 삭제     |
| GET    | `/health`     | 서버 상태 확인       |

## 설치 및 실행
### 1. 환경 변수 설정
```bash
export BD_USER={DB USER}
export DB_PASS={DB PASSWORD}
export DB_HOST={DB HOST}
export DB_PORT=3306(default)
export DB_NAME={DB NAME}
```

### 2. 소스 코드에서 직접 빌드
```bash
# 모듈 다운로드
go mod tidy

# 빌드
go build -o user-api-service

# 실행
./user-api-service
```

#### 플랫폼별 빌드 예시
```bash
# Linux 64-bit
GOOS=linux GOARCH=amd64 go build -o user-api-service

# Windows 64-bit
GOOS=Windows GOARCH=amd64 go build -o user-api-service.exe
```
