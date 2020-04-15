const request = require('request');
const random = require('string-random');
const util = require('util');

const API = "http://127.0.0.1:8006";

const randomNum = (i, a) => {
    return (i + Math.round(Math.random() * (a - i)));
}
class Mocks {
    cookie = ""

    getReq(uri, params) {
        return request({
            url: `${API}${uri}`,
            method: "GET",
            headers: {
                cookie: this.cookie
            },
            qs: params
        }, (err, resp, body) => {
            if (!err && resp.statusCode == 200) {
                console.log(body)
            }
        });
    }
    postReq(uri, data, cb = null) {
        return request({
            url: `${API}${uri}`,
            method: "POST",
            json: true,
            headers: {
                "content-type": "application/json",
                "X-Type": "custom",
                cookie: this.cookie
            },
            body: data
        }, (err, resp, body) => {
            if (!err && resp.statusCode == 200) {
                console.log(body)
                if (uri == '/login') {
                    this.cookie = resp.headers['set-cookie'].toString()
                    console.debug(this.cookie)
                }
                if (typeof cb == "function") {
                    cb(body)
                }
            } else {
                console.error(err)
                // console.error(resp)
            }
        });
    }
    putReq(uri, data, cb = null) {
        return request({
            url: `${API}${uri}`,
            method: "PUT",
            json: true,
            headers: {
                "content-type": "application/json",
                cookie: this.cookie
            },
            body: data
        }, (err, resp, body) => {
            if (!err && resp.statusCode == 200) {
                console.log(body)
                if (typeof cb == "function") {
                    cb(body)
                }
            }
        });
    }
    // register(username, password = '123456', cb = null) {
    //     let data = {
    //         username: username || random(12, { numbers: false }).toLowerCase(),
    //         password: password,
    //         cpassword: password
    //     }
    //     return this.postReq('/auth/register', data, cb)
    // }
    login(username, password = '123456', cb = null) {
        let data = {
            username: username,
            password: password
        }
        return this.postReq('/login', data, cb)
    }
    async addUser(username, password = '123456', cb = null) {
        let data = {
            "username": username,
            "password": password
        }
        return this.postReq('/debug/user', data, cb)
    }
    async addReceive(d = null, cb = null) {
        if (d == null)
            d = {
                "name": random(12, {
                    numbers: false
                }).toLowerCase(),
                "type": "custom",
                "header": "custom",
                "body": "{\"title\":[\"Github\"],\"text\":[\"repository.full_name\",\"repository.url\",\"commits.0.message\",\"commits.0.timestamp\",\"commits.0.author.name\"]}",
                "keyword": ""
            }
        return this.postReq('/api/receive', d, cb)
    }
    async addTemplate(url, vendor, body, cb = null) {
        let data = {
            // URL string
            "name": random(12, {
                numbers: false
            }).toLowerCase(),
            "url": url,
            "vendor": vendor,
            "body": body
        }
        return this.postReq('/api/template', data, cb)
    }
    async addPusher(data = null, cb = null) {
        if (data == null)
            data = {
                "name": random(12, {
                    numbers: false
                }).toLowerCase(),
                "url": "https://open.feishu.cn/open-apis/bot/hook/xxxxxxx",
                "vendor": "feishu",
                "template_id": 1
            }
        return this.postReq('/api/pusher', data, cb)
    }
    async addRelation(data = null, cb = null) {
        if (data == null)
            data = {
                "status": true,
                "user_id": randomNum(1, 10),
                "pusher_id": randomNum(1, 10),
                "receive_id": randomNum(1, 10)
            }
        return this.postReq('/api/relation', data, cb)
    }
    async testWebhook(cb = null) {
        let data = {
            "repository": {
                "name": "hongyan",
                "full_name": "virink/hongyan",
                "url": "https://github.com/virink/hongyan"
            },
            "commits": [{
                "id": "211df76b20ce909458a771d6dc5d8e1ef7c54b9b",
                "message": "Update Fuck thing",
                "timestamp": "2020-03-23T21:22:07+08:00",
                "author": {
                    "name": "Virink",
                    "email": "virink@outlook.com",
                    "username": "virink"
                }
            }]
        }
        return this.postReq('/webhook', data, cb)
    }
}

let m = new Mocks()
async function doit(u = '', cb = null) {
    var username = u || random(12, {
        numbers: false
    }).toLowerCase()
    await m.addUser(username, '123456', async (body) => {
        await m.login(username, '123456', async (body) => {
            await m.addReceive()
            await m.addPusher()
            await m.addPusher()
            cb && await cb()
        })
    })

}

async function main() {
    await doit('virink', async () => {
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key\"]}")
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key.subkey\"]}")
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key.subkey.subkey1\"]}")
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key\"]}")
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key.subkey\"]}")
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key.subkey.subkey1\"]}")
    });
    await doit('test');
    await doit();
    await doit();
    await doit();
    await doit();
    await doit();
    await doit();
    await doit();
    await doit();
}

async function chamd5(username = 'virink', password = '123456') {
    await m.login(username, username, async (_) => {
        await m.addTemplate("https://open.feishu.cn/open-apis/bot/hook/$key", "feishu", "{\"title\":[\"Title\"],\"text\":[\"key\"]}")
        await m.addTemplate("https://oapi.dingtalk.com/robot/send?access_token=$key", "dingding_text", "{\"msgtype\":[\"text\"],\"text\":{\"content\":[\"text.content\"]}}")
        let r = {
            "name": "chamd5_ti_dingding",
            "type": "dingding",
            "header": "",
            "body": "{\"title\":[\"ChaMD5\"],\"text\":[\"text.ccontent\"]}",
            "keyword": ""
        }
        await m.addReceive(r)
        var p = {
            "name": "chamd5_ti_feishu",
            "url": "https://open.feishu.cn/open-apis/bot/hook/e3fe90555f7b47289cabb3ab2b5a239f",
            "vendor": "feishu",
            "template_id": 0
        }
        await m.addPusher(p)
        p = {
            "name": "chamd5_ti_dingding",
            "url": "https://oapi.dingtalk.com/robot/send?access_token=e38ed299a78666082c96034ed47efae6aa8e812ed3aa9b4264d26a89ec534288",
            "vendor": "dingding",
            "template_id": 0
        }
        await m.addPusher(p)
        // let rr = {
        //     "status": true,
        //     "user_id": 1,
        //     "pusher_id": randomNum(1, 10),
        //     "receive_id": randomNum(1, 10)
        // }
        // await m.addRelation()
    })
}

// chamd5()
async function chamd5_r(username = 'virink', password = '123456') {
    await m.login(username, username, async (_) => {
        var rr = {
            "status": true,
            "user_id": 1,
            "pusher_id": 11,
            "receive_id": 12
        }
        await m.addRelation(rr)
        rr["pusher_id"] = 12
        await m.addRelation(rr)

    })
}

chamd5_r()


// main()
// setTimeout(() => {
//     m.testWebhook()
// }, 1000);