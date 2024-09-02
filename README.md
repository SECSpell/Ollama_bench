# Ollama Performance Testing Tool
[中文版本](https://github.com/SECSpell/Ollama_bench/blob/main/README_zh.md
This is a performance testing tool written in Go, primarily designed to test the local generation speed of Ollama. The tool also supports testing other services compatible with the OpenAI API interface specification.

## Testing Principle

This tool sends multiple concurrent requests to the Ollama API service, with each request containing a randomly selected question. The tool records the total number of tokens and request time, then calculates the average time per request and the number of tokens generated per second, thereby evaluating the performance of the Ollama API.

## Usage

1. Ensure that Ollama is running and has loaded the llama3.1 model.

2. Download the appropriate binary file for your system from the [Releases](https://github.com/SECSpell/Ollama_bench/releases) page.

3. Run the downloaded binary file:

   ```
   ./ollama_bench_darwin_arm64
   ```

   By default, the tool will use 1 concurrent request and send a total of 4 requests.

4. You can also specify the number of concurrent requests (C) and the total number of requests (N):

   ```
   ./ollama_bench_darwin_arm64 <C> <N>
   ```

   For example:
   ```
   ./ollama_bench_darwin_arm64 5 20
   ```
   This will use 5 concurrent requests and send a total of 20 requests.

## Configuration

Upon first launch, the tool will automatically create a `config.json` file in the same directory. You can modify this configuration as needed to support more models.