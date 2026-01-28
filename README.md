# Hi There 👋 !!

## requests
This project using rest api request tool called bruno  
Bruno project directory: `api/contract`

## point service

### run
```bash
make run
```

### down
```bash
make down
```

### debug
```bash
make run-debug
```

## for keploy

### trace
```bash
make run-keploy
```

<details>
    <summary>example</summary>

```bash

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.4
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
{"time":"2025-12-22T14:45:33.162011341Z","level":"ERROR","msg":"REQUEST_ERROR","uri":"/point/confirm","status":400,"err":"code=400, message=&{[user not found]}"}
🐰 Keploy: 2025-12-22T14:45:35.186510925Z       INFO    🟠 Keploy has captured test cases for the user's application.   {"path": "/Users/harukisugiyama/GolandProjects/pointservice/keploy/test-set-2/tests", "testcase name": "test-1"}
{"time":"2025-12-22T14:45:41.982762386Z","level":"INFO","msg":"REQUEST","uri":"/point/add","status":200}
🐰 Keploy: 2025-12-22T14:45:44.085712804Z       INFO    🟠 Keploy has captured test cases for the user's application.   {"path": "/Users/harukisugiyama/GolandProjects/pointservice/keploy/test-set-2/tests", "testcase name": "test-2"}
{"time":"2025-12-22T14:45:47.475678166Z","level":"INFO","msg":"REQUEST","uri":"/point/confirm","status":200}
🐰 Keploy: 2025-12-22T14:45:49.482654167Z       INFO    🟠 Keploy has captured test cases for the user's application.   {"path": "/Users/harukisugiyama/GolandProjects/pointservice/keploy/test-set-2/tests", "testcase name": "test-3"}
```

</details>


### test
```bash
make test-keploy
```

<details>
<summary>example</summary>

```bash
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.4
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:1323
🐰 Keploy: 2025-12-22T14:53:55.811711212Z       INFO    starting test for       {"test case": "[test-1]", "test set": "[test-set-3]"}
Testrun passed for testcase with id: "test-1"

--------------------------------------------------------------------

🐰 Keploy: 2025-12-22T14:53:55.815419795Z       INFO    result  {"testcase id": "[test-1]", "testset id": "[test-set-3]", "passed": "[true]"}
🐰 Keploy: 2025-12-22T14:53:55.815692712Z       INFO    starting test for       {"test case": "[test-2]", "test set": "[test-set-3]"}
{"time":"2025-12-22T14:53:55.814657462Z","level":"ERROR","msg":"REQUEST_ERROR","uri":"/point/confirm","status":400,"err":"code=400, message=&{[user not found]}"}
Testrun passed for testcase with id: "test-2"

--------------------------------------------------------------------

🐰 Keploy: 2025-12-22T14:53:55.817695504Z       INFO    result  {"testcase id": "[test-2]", "testset id": "[test-set-3]", "passed": "[true]"}
🐰 Keploy: 2025-12-22T14:53:55.81819617Z        INFO    starting test for       {"test case": "[test-3]", "test set": "[test-set-3]"}
{"time":"2025-12-22T14:53:55.817498337Z","level":"INFO","msg":"REQUEST","uri":"/point/add","status":200}
Testrun passed for testcase with id: "test-3"

--------------------------------------------------------------------

🐰 Keploy: 2025-12-22T14:53:55.81965017Z        INFO    result  {"testcase id": "[test-3]", "testset id": "[test-set-3]", "passed": "[true]"}
{"time":"2025-12-22T14:53:55.819442379Z","level":"INFO","msg":"REQUEST","uri":"/point/confirm","status":200}

 <=========================================> 
  TESTRUN SUMMARY. For test-set: "test-set-3"
        Total tests: 3
        Total test passed: 3
        Total test failed: 0
        Time Taken: "11.03 s"
 <=========================================> 

🐰 Keploy: 2025-12-22T14:53:56.935187213Z       WARN    To enable storing mocks in cloud, please use --disableMockUpload=false flag or test:disableMockUpload:false in config file

 <=========================================> 
  COMPLETE TESTRUN SUMMARY. 
        Total tests: 3
        Total test passed: 3
        Total test failed: 0
        Total time taken: "11.03 s"

        Test Suite Name         Total Test      Passed          Failed          Time Taken      

        "test-set-3"            3               3               0               "11.03 s"
<=========================================> 

🐰 Keploy: 2025-12-22T14:53:56.935641754Z       INFO    stopping Keploy {"reason": "replay completed successfully"}
🐰 Keploy: 2025-12-22T14:53:56.936717713Z       INFO    proxy stopped...
🐰 Keploy: 2025-12-22T14:53:57.668102005Z       INFO    eBPF resources released successfully...
🐰 Keploy: 2025-12-22T23:53:57.845781+09:00     INFO    exiting the current process as the command is moved to docker
```

```

</details>

## Spec と実装の対応

| Spec | Go 実装 |
|----|----|
| openapi/paths/point/add.yaml | internal/presentation/point_handler.go (PointAdd) |
| openapi/paths/point/sub.yaml | internal/presentation/point_handler.go (PointSub) |
| openapi/paths/point/confirm.yaml | internal/presentation/point_handler.go (PointConfirm) |
| openapi/paths/reservation/create.yaml | internal/presentation/point_handler.go (PointReserve) |
| asyncapi/channels/point/updated.yaml | internal/usecase/point_upsert.go (Publisher) |
| asyncapi/channels/reservation/created.yaml | cmd/worker/main.go (Consumer) |

# 新しく追加する時のルール

### REST API 追加手順

1. `openapi/paths/<domain>/<usecase>.yaml` を追加
2. `openapi/paths/<domain>/index.yaml` に追記
3. schema が必要なら `openapi/components/schemas` に追加し index.yaml 更新
4. handler を実装
5. CI で OpenAPI validate

### Event 追加手順

1. `asyncapi/channels/<domain>/<event>.yaml` を追加
2. `asyncapi/channels/index.yaml` に追記
3. `asyncapi/components/messages/<domain>.yaml` (or separate) 追加
4. schema が必要なら `asyncapi/components/schemas`
5. publisher / consumer 実装