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
    async addReceive(cb = null) {
        let data = {
            "name": random(12, {
                numbers: false
            }).toLowerCase(),
            "type": "custom",
            "header": "custom",
            "body": "{\"title\":[\"Github\"],\"text\":[\"repository.full_name\",\"repository.url\",\"commits.0.message\",\"commits.0.timestamp\",\"commits.0.author.name\"]}",
            "keyword": ""
        }
        return this.postReq('/api/receive', data, cb)
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
    async addPusher(cb = null) {
        let data = {
            "name": random(12, {
                numbers: false
            }).toLowerCase(),
            "url": "https://open.feishu.cn/open-apis/bot/hook/xxxxxxx",
            "vendor": "feishu",
            "template_id": 1
        }
        return this.postReq('/api/pusher', data, cb)
    }
    async addRelation(cb = null) {
        let data = {
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
            await m.addRelation()
            await cb()
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

main()
// setTimeout(() => {
//     m.testWebhook()
// }, 1000);