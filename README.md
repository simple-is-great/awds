## awds 실행방법
### 1. 실행파일 생성
- Makefile로 실행파일을 생성합니다.   
```
    make
```
- /awds 경로에 /bin 폴더가 생성되고 내부에 실행파일(awds)가 생성됩니다.

### 2. awds 실행
```
    ./bin/awds
```
- awds/경로에서 awds를 실행합니다.

--- 
## 주요 기능
Body에 넣어야 하는 항목이 많으므로, curl 보다는 Postman 등의 환경에서 실행을 추천합니다. \
Body는 raw JSON 형태로 전송합니다. 테스트는 Postman 환경에서 진행하였습니다.

### device 등록(POST)
- 요청 경로
```
    http://155.230.36.27:10270/devices
```

- Body에 들어갈 내용
```
    {
        "end_point": {디바이스 엔드포인트},
        "description": {설명, 생략 가능}
    }
```

### pod 등록(POST)
- 요청 경로
```
    http://155.230.36.27:10270/pods
```

- Body에 들어갈 내용
```
    {
        "end_point": {파드 엔드포인트},
        "description": {설명, 생략 가능}
    }
```

### job 등록(POST)
- 요청 경로
```
    http://155.230.36.27:10270/jobs
```

- Body에 들어갈 내용
```
    {
        "device_id": {디바이스 엔드포인트},
        "pod_id": {설명, 생략 가능},
        "input_size": {입력 크기, 정수, 생략 가능},
    }
```

### schedule(POST임, GET 아님에 주의!)
- 요청 경로
```
    http://155.230.36.27:10270/schedules/{job_id}
```

- Body에 들어갈 내용 없음

### pod 조회(GET)
- 요청 경로
```
    http://155.230.36.27:10270/pods
```

- pod 리스트 반환함

### device 조회(GET)
- 요청 경로
```
    http://155.230.36.27:10270/devices
```

- device 리스트 반환함

### job 조회(GET)
- 요청 경로
```
    http://155.230.36.27:10270/jobs
```

- job 리스트 반환함

## 기타 편의 기능(PATCH, DELETE)
### 디바이스 업데이트(PATCH)
- 요청 경로
```
    http://155.230.36.27:10270/devices/{device_id}
```

- Body에 들어갈 내용
```
    {
        "end_point": {디바이스 엔드포인트},
        "description": {설명, 생략 가능}
    }
```

### 디바이스 삭제(DELETE)
- 요청 경로
```
    http://155.230.36.27:10270/devices/{device_id}
```

- Body에 들어갈 내용 없음

### 파드 업데이트(PATCH)
- 요청 경로
```
    http://155.230.36.27:10270/pods/{pod_id}
```

- Body에 들어갈 내용
```
    {
        "end_point": {디바이스 엔드포인트},
        "description": {설명, 생략 가능}
    }
```

### 파드 삭제(DELETE)
- 요청 경로
```
    http://155.230.36.27:10270/pods/{pod_id}
```

- Body에 들어갈 내용 없음

### job 업데이트(PATCH)
- 요청 경로
```
    http://155.230.36.27:10270/jobs/{job_id}
```

- 주로 사용할 Body
```
    {
        "device_id": {디바이스 엔드포인트},
        "pod_id": {설명, 생략 가능},
        "input_size": {입력 크기, 정수, 생략 가능},
        "completed": {완료 여부, 불리언, 생략 가능},
    }
```

- Body에 들어갈 수 있는 내용
```
    {
        "device_id": {디바이스 엔드포인트},
        "pod_id": {설명, 생략 가능},
        "input_size": {입력 크기, 정수, 생략 가능},
        "partition_rate": {분배 비율, 실수, 생략 가능},
        "completed": {완료 여부, 불리언, 생략 가능},
        "DeviceStartIndex": {디바이스 시작 인덱스, 정수, 생략 가능},
        "DeviceEndIndex": {디바이스 종료 인덱스, 정수, 생략 가능},
        "PodStartIndex": {파드 시작 인덱스, 정수, 생략 가능},
        "PodEndIndex": {파드 종료 인덱스, 정수, 생략 가능},
    }
```

### job 삭제(DELETE)
- 요청 경로
```
    http://155.230.36.27:10270/jobs/{job_id}
```

- Body에 들어갈 내용 없음