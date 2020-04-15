# Message Pusher

## Demo

- [x] Add User - 添加用户
- [x] Add Receive - 添加接收模板
- [x] Add Pusher - 添加推送
- [x] Add Relation - 绑定 [用户\接收模板\推送] 关系
- [x] WebHook - 接收任意 Json 推送并由[接收模板]解析及后续推送

## Principle

/webhook    data    <-    json
            |
            |
 [(keyword) || (header)] 
 判断对应 [Receive] 模板
            |
            |    ->    {{result json}}    ->               |
            |                                              |
        [Relation]                                         |    ->    pust to [Target]
            |                                              |
            ->   [Pusher]   ->   {{Pusher.url&var}}   ->   |
                                {{Pusher.template}}   ->   |

## TODO

- [ ] Admin
- [ ] Auth Login/Logout
  - [x] Session
- [ ] Add Template
- [ ] Add UI
- [ ] Add ......more and more