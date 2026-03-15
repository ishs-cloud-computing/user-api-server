# User API

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
| Name | Description           |
|------|-----------------------|
| DB_USER | 데이터베이스 유저 이름   |
| DB_PASS | 데이터베이스 비밀번호    |
| DB_HOST | 데이터베이스 접속 주소   |
| DB_PORT | 데이터베이스 접속 포트   |
| DB_NAME | 데이터베이스 이름       |

> DB_PORT의 기본 값은 3306입니다.

### 2. 빌드된 바이너리 사용
```bash
# 압축 해제
unxz ./user-api-server.xz

# 실행
./user-api-server
```

### 3. 소스 코드에서 직접 빌드
```bash
# 모듈 다운로드
go mod tidy

# 빌드
GOOS=linux GOARCH=amd64 go build

# 실행
./user-api-server
```

## 데이터베이스
이 서버는 MySQL 또는 MySQL 호환 DB와 동작합니다.

####  users TABLE 예제
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

## Third-Party Libraries
This project uses the following open-source libraries:

- [github.com/go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) — Mozilla Public License 2.0
- [github.com/gorilla/mux](https://github.com/gorilla/mux) — BSD 3-Clause License

All third-party libraries are used according to their license terms.

## License
This project is licensed under [Zero-Clause BSD (0BSD)](LICENSE).
