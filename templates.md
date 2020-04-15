# Templates

## FeiShu

Custom Bot

`https://open.feishu.cn/open-apis/bot/hook/xxxxxxxxxxxxxxxxxxxxxxxxxxx`

```json
{
    "title": "Hello Feishu",  # 选填
    "text": "Good Feishu"  # 必填
}
```


`curl -X POST -H "Content-Type: application/json" -d '{"title": "Hello Feishu", "text": "Good Feishu"}' https://open.feishu.cn/open-apis/bot/hook/xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx `

## DingDing

https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq

`https://oapi.dingtalk.com/robot/send?access_token=XXXXXX`

```json
{
    "msgtype": "text", 
    "text": {
        "content": "我就是我, 是不一样的烟火@156xxxx8827"
    }
}
```

