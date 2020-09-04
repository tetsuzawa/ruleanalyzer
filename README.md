# ruleanalyzer

exampleのようにルールを記述したコードから静的解析用のAnalyzerを生成するツール

# インストール

```shell script
go get -u github.com/tetsuzawa/ruleanalyzer/cmd/ruleanalyzer
```

# 使い方

1. ルールを記述したファイルをtestdata以下に作る。（コンパイルに含まれなければどこでも良い）  
    ```go
    func RuleOsOpen() {
        // step: call os.Open
        f, _ := os.Open("xxx")
        // step: call *File.close
        defer f.Close()
    }
    ```
   
1. ruleanalyzerを実行する
    ```shell script
    ruleanalyzer testdata/src/osopen/osopen_rule.go
    ```

1. カレントディレクトリ以下にルール名のディレクトリ作成され、中にAnalyzerが生成される。
    ```go
    var Analyzer = &analysis.Analyzer{
        Name: "OsOpen",
        Doc:  doc,
        Run:  run,
        Requires: []*analysis.Analyzer{
           buildssa.Analyzer,
        },
    }
    
    const ruleName = "OsOpen"
    
    func run(pass *analysis.Pass) (interface{}, error) {
       /* ...*/
    ```
   
1. ルールからはステップに応じたマイルストーンが生成されており、Analyzerではマイルストーンが全て達成されたかを調べる。マイルストーンが全て達成されていなければログが表示される。

# ルールのフォーマット

- 関数名は`RuleXxx`のように`Rule`で始める
- 各ステップの処理を次のフォーマットのコメントの次の行に書く   
  `// step: xxxxxxx`
- コメントの後に複数行の処理を書いても必要なステップとして認識されるのは直後の行のみである。


# コード生成までの原理


**ルール解析**

1. ルールが記述された関数をSSA形式に変換する
1. 特定フォーマットのコメント直後の命令から`*types.Object`を取り出し、`MilestoneQueue`に追加する。

**コード生成**

1. `*types.Object`からパッケージ名、型名、オブジェクト名などを取得する
1. これらの情報を元にAnalyzer内でオブジェクトを再取得するコードを構築する
1. テンプレートを使用し、ソースコードを組み立てる。

# 依存ライブラリ・コマンド

- golang.org/x/tools/go
- golang.org/x/tools/cmd/goimports
- github.com/gostaticanalysis/analysisutil

# TODO

- 対応する命令を増やす  
    現状、`*ssa.Call`、`*ssa.Alloc`、`*ssa.Defer`のみしか命令を処理できていないため、宣言をはじめとする他の命令にも対応したい
- 各ステップで作られた変数を追従し、その後のステップと結びつけられるようにする  
    現状、変数と関数呼び出しなどの対応付けができていない状態である。元のコードのSSA形式の構造を生成したAnalyzerに持ち込む方法を思いついていないので考えたい。
- インターフェースのメソッド呼び出しに対応する  
    インターフェースからの呼び出しはポインタ解析が必要となり、現状対応できていない。  
- コードのリファクタ  
    時間が限られていたためパッケージ構成やファイル分けが雑になってしまっている。また、テストも書けていないため追加したい。
    


