# shortUrl


## 生成短链接
```
curl  POST 'http://url.aliyundevops.com/api/shorten' \
--header 'Content-Type: application/json' \
-d '{
    "url": "https://www.baidu.com",
    "expiration_in_minutes": 160
}'
```
## 查看链接信息
```
curl 'http://url.aliyundevops.com/api/info?shortlink=1&='
```

## 跳转
```
curl 'http://url.aliyundevops.com/1'
```
