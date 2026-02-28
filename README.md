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
export DBUSER="{DB USER}"
export DB_PASS="{DB PASSWORD}"
export DB_HOST="{DB HOST}"
export DB_PORT="3306(default)"
export DB_NAME="{DB NAME}"
```

### 2. 빌드된 바이너리 사용
```bash
# Linux / macOS
./user-api-service

# Windows
user-api-service.exe
```

### 3. 소스 코드에서 직접 빌드
```bash
# 모듈 다운로드
go mod tidy

# 빌드
go build

# 실행
./user-api-service
```

#### 플랫폼별 빌드 예시
```bash
# Linux 64-bit
GOOS=linux GOARCH=amd64 go build

# macOS 64-bit (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build

# Windows 64-bit
GOOS=windows GOARCH=amd64 go build
```

## 데이터베이스 예시
```sql
CREATE TABLE users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    age INT UNSIGNED,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```
