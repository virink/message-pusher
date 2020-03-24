const request = require('request');
const random = require('string-random');
const util = require('util');

const API = "http://127.0.0.1:8006";

const randomNum = (i, a) => {
    return (i + Math.round(Math.random() * (a - i)));
}
class Mocks {
    getReq(uri, params) {
        return request({
            url: `${API}${uri}`,
            method: "GET",
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
                "X-Type": "custom"
            },
            body: data
        }, (err, resp, body) => {
            if (!err && resp.statusCode == 200) {
                console.log(body)
                if (uri == '/auth/login') {
                    this.jwt = typeof resp.headers.token == "string" ? resp.headers.token : resp.headers.token[0]
                }
                if (typeof cb == "function") {
                    cb(body)
                }
            } else {
                console.error(err)
                console.error(resp)
            }
        });
    }
    putReq(uri, data, cb = null) {
        return request({
            url: `${API}${uri}`,
            method: "PUT",
            json: true,
            headers: {
                "content-type": "application/json"
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
        return this.postReq('/auth/login', data, cb)
    }
    async addUser(username, password = '123456', cb = null) {
        let data = {
            "username": username,
            "password": password
        }
        return this.postReq('/admin/adduser', data, cb)
    }
    async addReceive(cb = null) {
        let data = {
            "name": "test",
            "type": "custom",
            "header": "custom",
            "body": "{\"title\":[\"Github\"],\"text\":[\"repository.full_name\",\"repository.url\",\"commits.0.message\",\"commits.0.timestamp\",\"commits.0.author.name\"]}",
            "keyword": ""
        }
        return this.postReq('/user/addreceive', data, cb)
    }
    async addPusher(cb = null) {
        let data = {
            "name": "feishu_sec",
            "url": "https://open.feishu.cn/open-apis/bot/hook/xxxxxxx",
            "vendor": "feishu",
            "template_id": 1
        }
        return this.postReq('/user/addpusher', data, cb)
    }
    async addRelation(cb = null) {
        let data = {
            "status": true,
            "user_id": 1,
            "pusher_id": 1,
            "receive_id": 1
        }
        return this.postReq('/user/addrelation', data, cb)
    }
    async testWebhook(cb = null) {
        let data = {
            "repository": {
                "name": "hongyan",
                "full_name": "virink/hongyan",
                "url": "https://github.com/virink/hongyan"
            },
            "commits": [
                {
                    "id": "211df76b20ce909458a771d6dc5d8e1ef7c54b9b",
                    "message": "Update Fuck thing",
                    "timestamp": "2020-03-23T21:22:07+08:00",
                    "author": {
                        "name": "Virink",
                        "email": "virink@outlook.com",
                        "username": "virink"
                    }
                }
            ]
        }
        return this.postReq('/webhook', data, cb)
    }
}

let m = new Mocks()
async function doit(u = '') {
    var username = u || random(12, { numbers: false }).toLowerCase()
    await m.addUser(username, '123456', async (body) => {
        // console.log(body)
        await m.addReceive()
        await m.addPusher()
        await m.addRelation()
    })
}

doit()
setTimeout(() => {
    m.testWebhook()
}, 1000);