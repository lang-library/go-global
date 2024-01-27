package global

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	PrepareForFile(path)
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func Echo(x ...any) {
	len := len(x)
	if len == 0 || len > 2 {
		panic("global.Echo(): args out of range")
	}
	if len == 2 {
		if x[1] != nil {
			fmt.Printf("%s: ", x[1])
		}
	}
	fmt.Println(x[0])
}

func ExeDir() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := GetParent(ex)
	return exPath
}

func FromJson(jsonStr string) any {
	jsonBytes := []byte(jsonStr)
	var result any
	err := json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil
	}
	return result
}

func GetParent(path string) string {
	return filepath.Dir(path)
}

func PrettifyJson(jsonStr string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(jsonStr), "", "  ")
	if err != nil {
		return jsonStr
	}
	indentJson := buf.String()
	return indentJson
}

func Print(x ...any) {
	len := len(x)
	if len == 0 || len > 2 {
		panic("global.PrettyPrint(): args out of range")
	}
	if len == 2 {
		if x[1] != nil {
			fmt.Printf("%s: ", x[1])
		}
	}
	fmt.Println(ToPrettyJson(x[0]))
}

func Prepare(path string) {
	os.MkdirAll(path, os.ModePerm)
}

func PrepareForFile(path string) {
	parent := GetParent(path)
	Prepare(parent)
}

func ToJson(x any) string {
	jsonBytes, err := json.Marshal(x)
	if err != nil {
		return "null"
	}
	return string(jsonBytes)
}

func ToPrettyJson(x any) string {
	jsonStr := ToJson(x)
	return PrettifyJson(jsonStr)
}

// https://qiita.com/brushwood-field/items/417f7c07ee5813239ff3
// unZip zipファイルを展開する
func UnZip(src, destDir string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	//ext := filepath.Ext(src)
	//rep := regexp.MustCompile(ext + "$")
	//dir := filepath.Base(rep.ReplaceAllString(src, ""))
	//destDir := filepath.Join(dest, dir)
	// ファイル名のディレクトリを作成する
	if err := os.MkdirAll(destDir, os.ModeDir); err != nil {
		return err
	}

	for _, f := range r.File {
		if f.Mode().IsDir() {
			// ディレクトリは無視して構わない
			continue
		}
		if err := saveUnZipFile(destDir, *f); err != nil {
			return err
		}
	}

	return nil
}

// saveUnZipFile 展開したZipファイルをそのままローカルに保存する
func saveUnZipFile(destDir string, f zip.File) error {
	// 展開先のパスを設定する
	destPath := filepath.Join(destDir, f.Name)
	// 子孫ディレクトリがあれば作成する
	if err := os.MkdirAll(filepath.Dir(destPath), f.Mode()); err != nil {
		return err
	}
	// Zipファイルを開く
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	// 展開先ファイルを作成する
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()
	// 展開先ファイルに書き込む
	if _, err := io.Copy(destFile, rc); err != nil {
		return err
	}

	return nil
}
