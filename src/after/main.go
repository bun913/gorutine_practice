package main

import (
	"fmt"
	"sync"
	"testing"
)

// Document 書類
// 営業が内容を書いた後であれば isWritedがtrueに
// チェッカーが内容を精査した後であれば isCheckedがtrueになる
type Document struct {
	isWrited  bool
	isChecked bool
}

// Employee 社員です。社員はみな家に帰れます。
type Employee struct {
}

// GoHome 家に帰ります。ビール飲みたい
func (employee Employee) GoHome() {
	fmt.Println("華金じゃあ!!!")
}

// Sales 営業です。資料に必要な情報を記載するまでが仕事です。
// Document束を持ちます
type Sales struct {
	Employee
	Documents []Document
	DocStream chan Document
}

// WriteDocs 営業が営業資料に必要事項を記載します。
func (sales *Sales) WriteDocs(wg *sync.WaitGroup) {
	// 処理終了したら、チャンネルを閉じてもう書類がこないことを伝えます
	defer close(sales.DocStream)
	// 処理終了したら、このスレッドの処理は終わったよと伝えます
	defer wg.Done()
	for _, doc := range sales.Documents {
		doc.isWrited = true
		sales.DocStream <- doc
	}
	sales.Employee.GoHome()
}

// Checker チェッカーです。営業が作った資料をチェックするのが仕事です。
// Documentの束を持ちます
type Checker struct {
	Employee
	DocStream chan Document
}

// CheckDocs チェッカーが営業が作った資料をチェックします。
func (checker *Checker) CheckDocs(wg *sync.WaitGroup) {
	defer wg.Done()
	for doc := range checker.DocStream {
		if doc.isWrited {
			doc.isChecked = true
		}
	}
	checker.Employee.GoHome()
}

func main() {
	result := testing.Benchmark(func(b *testing.B) {
		b.ResetTimer()
		// 営業とチェッカーの2つの並行処理が終わるのをmain関数で待ち受けるためのもの
		var wg sync.WaitGroup
		// 何も書いていない書類を打ち出します
		documents := []Document{}
		for range make([]int, 1000000) {
			newDoc := Document{}
			documents = append(documents, newDoc)
		}
		// 　資料受け渡しのためのチャンネルを作成します
		docStream := make(chan Document, 12)
		// 営業が資料に記載します
		sales := &Sales{Employee: Employee{},
			Documents: documents,
			DocStream: docStream,
		}
		wg.Add(1)
		go sales.WriteDocs(&wg)
		// 営業の資料を渡してチェッカーがチェックを始めます
		checker := &Checker{Employee: Employee{}, DocStream: docStream}
		wg.Add(1)
		go checker.CheckDocs(&wg)
		wg.Wait()
	})
	fmt.Printf("%s", result)
}
