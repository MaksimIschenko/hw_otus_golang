befote optimization
```
=== RUN   TestGetUsers_Time_And_Memory
    stats_optimization_test.go:44: time used: 206.389292ms / 300ms
    stats_optimization_test.go:45: memory used: 174Mb / 30Mb
    assertion_compare.go:332: 
                Error Trace:    stats_optimization_test.go:48
                Error:          "183365256" is not less than "31457280"
                Test:           TestGetUsers_Time_And_Memory
                Messages:       [the program is too greedy]
--- FAIL: TestGetUsers_Time_And_Memory (3.12s)
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:76: time used: 307.038959ms / 300ms
    stats_optimization_test.go:77: memory used: 308Mb / 30Mb
    assertion_compare.go:332: 
                Error Trace:    stats_optimization_test.go:79
                Error:          "307038959" is not less than "300000000"
                Test:           TestGetDomainStat_Time_And_Memory
                Messages:       [the program is too slow]
--- FAIL: TestGetDomainStat_Time_And_Memory (5.07s
```

after optimization
```
=== RUN   TestGetUsers_Time_And_Memory
    stats_optimization_test.go:44: time used: 92.400458ms / 300ms
    stats_optimization_test.go:45: memory used: 3Mb / 30Mb
--- PASS: TestGetUsers_Time_And_Memory (0.98s)
=== RUN   TestGetDomainStat_Time_And_Memory
    stats_optimization_test.go:76: time used: 103.076834ms / 300ms
    stats_optimization_test.go:77: memory used: 19Mb / 30Mb
--- PASS: TestGetDomainStat_Time_And_Memory (1.04s)
PASS
```

