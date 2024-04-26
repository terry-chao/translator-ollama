# Translator-ollama

# Quickstart

## requirements
- [Go](https://go.dev/)
- [ollama](https://ollama.ai/)

## Installation

```shell
go build .
```

## Usage
```shell
./translator-ollama
```

Then you can translate text by sending a POST request to `/translate` with the following JSON body:
```json
{
    "output_lang": "英语",
    "text": "你好"
}
```