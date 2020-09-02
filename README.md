# ruleanalyzer
ruleanalyzer is static analyzer generater based on `func Rule_xxx()` like `func Test_xxx()`

# やりたいこと

```go
func Rule_xxx(){
    // 1:xxx
    f, _ := os.Open("_xxx")
    // 2:xxx
    f.Close()
    // 3:xxx
    return
}
```
のように`xxx_test.go`ファイルにルールを記述する。

- ruleを解析するanalyzerを作って、analyzerのコードを生成するイメージ
- 有向グラフの話を考慮しないとだめかも
- go generateで生成する
  - つまり、作るツールはコマンドラインで実行されるイメージ
  
# MVP

1. 次のコードを解析してルールを構築する
    ```go
    func Rule_xxx(){
        // 1:xxx
        f, _ := os.Open("_xxx")
        // 2:xxx
        f.Close()
        // 3:xxx
        return
    }
    ```
2. ルールを解析して
    

  





