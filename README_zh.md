# Ollama 性能测试工具
[English Version](https://github.com/SECSpell/Ollama_bench/blob/main/README.md)

这是一个用Go语言编写的性能测试工具，主要用于测试Ollama在本地的生成速度。该工具也支持测试其他兼容OpenAI API接口规范的服务。

## 测试原理

本工具通过并发发送多个请求到Ollama API服务，每个请求包含一个随机选择的问题。工具会记录总的token数量和请求时间，最后计算出平均每个请求的时间和每秒生成的token数量，从而评估该Ollama API的性能。

## 使用方法

1. 确保Ollama已经启动，并且已经加载了llama3.1模型。

2. 从[Releases](https://github.com/SECSpell/Ollama_bench/releases)页面下载适合您系统的二进制文件。

3. 运行下载的二进制文件：

   ```
   ./ollama_bench_darwin_arm64
   ```

   默认情况下，工具会使用1个并发请求，总共发送4个请求。

4. 您也可以指定并发数（C）和总请求数（N）：

   ```
   ./ollama_bench_darwin_arm64 <C> <N>
   ```

   例如：
   ```
   ./ollama_bench_darwin_arm64 5 20
   ```
   这将使用5个并发请求，总共发送20个请求。

## 配置

工具首次启动会在同一目录下自动创建一个`config.json`文件，您可以根据需要修改配置，从而支持更多的模型。

## 部分平台的测试结果，llama3.1:latest（8B）
| 测试平台 | 平均每秒生成token数量 | 备注 |
| --- | --- | --- |
| MacBook Pro M1 Pro | 30token/s |  |
| AMD Ryzen R7-7840H（780M GPU） | 11token/s | 纯 CPU 推理 |
| AMD Ryzen R7-7840H（780M GPU） | 17token/s | 核显 GPU 推理，[第三方编译OLLAMA](https://github.com/likelovewant/ollama-for-amd) |
| AMD Ryzen R7-8845HS（780M GPU） | 17token/s | 核显 GPU 推理，[第三方编译OLLAMA](https://github.com/likelovewant/ollama-for-amd)  |
| Nvidia P106-100(6G) | 25token/s | docker ollama 推理 |
| Nvidia 2080ti（22G） | 86token/s |  GPU 魔改过显存 |
|  Nvidia 3090（24G） | 124token/s |  |
|  Groq Cloud | 581token/s | 测试服务器为谷歌云美东节点，模型使用llama-3.1-8b-instant |
---