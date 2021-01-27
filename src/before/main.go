package main

import (
	"fmt"
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
}

// WriteDocs 営業が営業資料に必要事項を記載します。
func (sales *Sales) WriteDocs() {
	for _, doc := range sales.Documents {
		doc.isWrited = true
	}
	// 終わったら帰宅する
	sales.Employee.GoHome()
}

// Checker チェッカーです。営業が作った資料をチェックするのが仕事です。
// Documentの束を持ちます
type Checker struct {
	Employee
	Documents []Document
}

// CheckDocs チェッカーが営業が作った資料をチェックします。
func (checker *Checker) CheckDocs() {
	if len(checker.Documents) > 100 {
		fmt.Println("まとめて持ってくんなよ・・・")
	}
	for _, doc := range checker.Documents {
		if doc.isWrited {
			doc.isChecked = true
		}
	}
	// 終わったら帰宅する
	checker.Employee.GoHome()
}

func main() {
	result := testing.Benchmark(func(b *testing.B) {
		// 何も書いていない書類を打ち出します
		documents := []Document{}
		for range make([]int, 1000000) {
			newDoc := Document{}
			documents = append(documents, newDoc)
		}
		// 営業が資料に記載します
		sales := &Sales{Employee: Employee{}, Documents: documents}
		sales.WriteDocs()
		// 営業の資料を渡してチェッカーがチェックを始めます
		checker := &Checker{Employee: Employee{}, Documents: sales.Documents}
		checker.CheckDocs()
	})
	fmt.Printf("%s", result)
}
