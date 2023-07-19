# jkscan

`jkscan` is a simple port scanner written in golang. It is designed to be fast and easy to use.

## Usage

```
docker-compose up
```

## Output

By default, `jkscan` will output a list of open ports (`outputs`) and logs (`outlogs`) in current directory.

## 运行示例

go run ./jkscan.go -t 211.22.90.0/24 -oJ outputs/211.22.90.0.json