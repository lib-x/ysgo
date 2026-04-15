package ysgo

import "fmt"

func Example_buildDownloadURLHelpers() {
	fmt.Println(downloadHost("D", "ysepan.com"))
	fmt.Println(downloadHost("X", "ysepan.com"))
	fmt.Println(joinPath("nested", "demo"))
	// Output:
	// ys-d.ysepan.com
	// y.ys168.com:8000
	// nested/demo
}

func ExampleListSubdirectories() {
	files := []RemoteFile{
		{Subdirectory: "nested/a"},
		{Subdirectory: "nested/b"},
		{Subdirectory: "other"},
	}
	for _, s := range ListSubdirectories(files, "nested") {
		fmt.Println(s.Path)
	}
	// Output:
	// nested/a
	// nested/b
}
