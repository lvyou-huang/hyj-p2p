package processfile

import (
	"errors"
	"fmt"
	"io"
	"log"
	"strings"

	"os"
)

func MakeSilce(filePath string) {
	// 将文件读入内存
	file, err := os.Open(filePath)
	split := strings.Split(filePath, ".")
	name := split[0]
	//种类
	tp := split[1]
	if err != nil {
		panic(err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	// 计算每个分片的大小
	chunkSize := len(data) / 10

	for i := 0; i < 10; i++ {
		// 计算每个分片的起始位置和结束位置
		start := i * chunkSize
		end := start + chunkSize

		if i == 9 {
			end = len(data)
		}

		// 将分片写入文件
		file, err := os.Create(fmt.Sprintf("%s-%d.%s", name, i, tp))
		if err != nil {
			log.Println(err)
			return
		}
		file.Write(data[start:end])
		file.Close()
	}
}

// c.docx c-0.docx,filePath为合并后的路径与名字,名字应与切片一致
func MakeFile(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	split := strings.Split(filePath, ".")
	name := split[0]
	//种类
	tp := split[1]
	for i := 0; i < 10; i++ {
		log.Println("导入第", i, "个切片")
		open, err := os.Open(fmt.Sprintf("%s-%d.%s", name, i, tp))
		if err != nil {
			log.Println(err)
			return errors.New("没有获取全部切片")
		}
		_, err = io.Copy(file, open)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
