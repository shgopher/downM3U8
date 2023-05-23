/*
 * @Author: shgopher shgopher@gmail.com
 * @Date: 2023-05-23 22:09:13
 * @LastEditors: shgopher shgopher@gmail.com
 * @LastEditTime: 2023-05-23 22:13:43
 * @FilePath: /downM3U8/main.go
 * @Description:
 *
 * Copyright (c) 2023 by shgopher, All Rights Reserved.
 */
 package main

 import (
	 "fmt"
	 "io"
	 "net/http"
	 "net/url"
	 "os"
	 "os/exec"
	 "path/filepath"
	 "regexp"
	 "sort"
	 "strconv"
	 "strings"
	 "sync"
 )
 
 var (
	 URL = ""
 )
 
 func init() {
	 URL = os.Args[1]
 }
 
 func main() {
	 HandleM3U8(URL)
 }
 
 func HandleM3U8(u string) {
	 fmt.Println("开始执行，请稍等")
	 // 1. 打开 M3U8 文件
	 resp, err := http.Get(u)
	 if err != nil {
		 fmt.Println(err)
	 }
	 burl := dealWithUrl(URL)
	 baseUrl, err := url.Parse(burl)
	 if err != nil {
		 panic(err)
	 }
	 defer resp.Body.Close()
 
	 // 2. 读取 M3U8 文件内容
	 data, err := io.ReadAll(resp.Body)
	 if err != nil {
		 panic(err)
	 }
 
	 // 3. 解析 M3U8 文件内容，获取所有 .ts 文件链接
	 tsUrls := []string{}
	 lines := strings.Split(string(data), "\n")
	 for _, line := range lines {
		 if strings.HasPrefix(line, "#") {
			 continue
		 }
		 tsUrl, err := url.Parse(line)
		 if err != nil {
			 panic(err)
		 }
		 tsUrls = append(tsUrls, baseUrl.ResolveReference(tsUrl).String())
	 }
 
	 // 4. 下载所"有 .ts 文件并保存到本地
	 tsFiles := []string{}
	 wg := new(sync.WaitGroup)
	 wg.Add(len(tsUrls))
	 sc := make(chan string)
	 os.Mkdir("./poolHi", os.ModePerm)
	 fmt.Println("读取原始m3u8数据完毕，开始对分段进行处理，分段数量是", len(tsUrls))
	 for i, tsUrl := range tsUrls {
		 readTsFile(i, tsUrl, sc, wg)
	 }
	 go func() {
		 wg.Wait()
		 fmt.Println("全部分段数据获取完毕！")
		 close(sc)
	 }()
 
	 for v := range sc {
		 tsFiles = append(tsFiles, v)
	 }
	 sortString(tsFiles)
	 // 5. 合并所有 .ts 文件为一个 MP4 文件
	 fmt.Println("开始合并为mp4文件")
	 mp4File := "video.mp4"
	 cmd := exec.Command("ffmpeg", "-i", "concat:"+strings.Join(tsFiles, "|"), "-c", "copy", mp4File)
	 err = cmd.Start()
	 if err != nil {
		 fmt.Println(err)
	 }
	 err = cmd.Wait()
 
	 if err != nil {
		 fmt.Println(err)
	 }
	 removeAllFiles("./poolHi")
	 fmt.Println("任务已经执行完毕，并且删除临时文件，请在本目录下查看video.mp4")
 }
 
 // 删除中间ts文件
 func removeAllFiles(dir string) error {
	 d, err := os.Open(dir)
	 if err != nil {
		 return err
	 }
	 defer d.Close()
 
	 files, err := d.Readdir(-1)
	 if err != nil {
		 return err
	 }
 
	 for _, file := range files {
		 filePath := filepath.Join(dir, file.Name())
		 if !file.IsDir() {
			 err = os.Remove(filePath)
			 if err != nil {
				 return err
			 }
		 }
	 }
 
	 return nil
 }
 
 // 并发读取ts文件
 func readTsFile(i int, url string, sc chan string, wg *sync.WaitGroup) {
	 go func(i int, url string) {
		 defer func() {
			 fmt.Printf("第%d个分段执行完毕\n", i)
		 }()
		 defer wg.Done()
		 resp, err := http.Get(url)
		 if err != nil {
			 fmt.Println("读取ts文件错误，", url, err)
		 }
		 defer resp.Body.Close()
 
		 // 创建 .ts 文件
		 tsFile := fmt.Sprintf("./poolHi/%d.ts", i)
		 file, err := os.Create(tsFile)
		 if err != nil {
			 fmt.Println("创建文件错误", err)
		 }
		 defer file.Close()
 
		 // 将 .ts 文件内容写入本地文件
		 _, err = io.Copy(file, resp.Body)
		 if err != nil {
			 panic(err)
		 }
		 sc <- tsFile
		 //tsFiles = append(tsFiles, tsFile)
	 }(i, url)
 }
 
 func sortString(strs []string) {
 
	 // 编译正则表达式
	 re := regexp.MustCompile(`\d+`)
 
	 // 定义一个排序函数，将字符串中的数字提取出来并进行比较
	 sortFunc := func(i, j int) bool {
		 num1, _ := strconv.Atoi(re.FindString(strs[i]))
		 num2, _ := strconv.Atoi(re.FindString(strs[j]))
		 return num1 < num2
	 }
 
	 // 使用排序函数对字符串切片进行排序
	 sort.Slice(strs, sortFunc)
 }
 
 func dealWithUrl(url string) (nurl string) {
	 // 编译正则表达式，匹配最后一个 / 之前的所有字符
	 re := regexp.MustCompile(`(.*/)[^/]*$`)
 
	 // 使用正则表达式提取路径部分
	 match := re.FindStringSubmatch(url)
 
	 // 如果能够匹配，则输出路径部分
	 if len(match) > 1 {
		 return match[1]
	 }
	 return ""
 }
 
