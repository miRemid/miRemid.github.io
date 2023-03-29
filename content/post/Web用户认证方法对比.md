---
title: "Web用户认证方法对比"
date: 2021-01-08T21:33:13Z
draft: false
toc: false
images:
tags: 
  - web
  - python
  - 翻译
---

> 本文翻译自[testdriven.io](https://testdriven.io/blog/web-authentication-methods/#authentication-vs-authorization)

在这篇文章，我们将从一名Python开发者的视角来观察目前最常见的几种处理Web认证的方式
> 尽管本片文章所有的代码是面向Python开发者的，但是实际上对所有的Web开发者，每种认证方法实际上都是差不多的

# 认证vs授权(Authentication vs Authorization)
认证是一种处理用户或设备尝试通过凭证来访问受限系统的过程。而授权则是验证用户或设备是否有权限来访问系统提供的确切服务

简单来讲就是
- 认证：你是谁？
- 授权：你可以做什么？

认证比授权出现的更早，用户必须在获得访问资源权限之前经过合法验证。最常见的用户认证方法就是`username`和`password`。一旦认证完成，不同身份例如`admin`、`moderator`等等，都将被附加在用户身上用于提供访问系统的身份信息。

有了上述的解释，让我们来看一看验证用户的不同方法吧

<!--more-->

# HTTP基本认证(HTTP Basic Authentication)
基本认证是一种建立在HTTP协议上的最基础的验证方式。在基本认证中，每一个请求`Header`都要携带登录凭证信息
```
"Authorization: Basic dXNlcm5hbWU6cGFzc3dvcmQ=" your-website.com
```

用户名和密码都没有被加密，并且用户名和密码通过`:`组合在一起`username:password`。这个字符串通常需要通过`Base64`编码生成
```python
>>> import base64
>>>
>>> auth = "username:password"
>>> auth_bytes = auth.encode('ascii') # convert to bytes
>>> auth_bytes
b'username:password'
>>>
>>> encoded = base64.b64encode(auth_bytes) # base64 encode
>>> encoded
b'dXNlcm5hbWU6cGFzc3dvcmQ='
>>> base64.b64decode(encoded) # base64 decode
b'username:password'.
```

这种方法是无状态的因此客户端必须在每个请求中携带次凭证，这种方法非常适合单一的简单的不需要持久`sessions`的API的调用

## 流程
1. 未认证的客户端请求一个受限的资源
2. HTTP返回401状态码并携带`Basic`值的名为`WWW-Authenticate`的header
3. `WWW-Authenticate: Basic`让浏览器提示用户密码输入
4. 在输入凭证之后，在每一个请求的Header中都将携带该凭证`Authorization: Basic dcdvcmQ=`
![基本认证](https://testdriven.io/static/images/blog/web-authentication-methods/basic_auth.png)

## 优点
1. 步骤少，认证过程快
2. 实现简单
3. 主流浏览器都支持

## 缺点
1. Base64并不是真正意义上的加密算法，只是另一种呈现数据的方式。Base64字符串很容易在发送文本的时候被解码。弱安全性会遭来各种各样的攻击，因此`HTTPS/SSL`非常有必要
2. 在每个请求都必须携带凭证
3. 用户只能通过写入一个错误的凭证用于登出

## 第三方依赖(Python下同)
- [Flask-HTTPAuth](https://flask-httpauth.readthedocs.io/)
- [django-basicauth](https://github.com/hirokiky/django-basicauth/)
- [FastAPI:HTTP Basic Auth](https://fastapi.tiangolo.com/advanced/security/http-basic-auth/)

## 代码
用`Flask-HTTP`可以使用Flask非常轻松的实现基本认证
```python
from flask import Flask
from flask_httpauth import HTTPBasicAuth
from werkzeug.security import generate_password_hash, check_password_hash

app = Flask(__name__)
auth = HTTPBasicAuth()

users = {
    "username": generate_password_hash("password"),
}


@auth.verify_password
def verify_password(username, password):
    if username in users and check_password_hash(users.get("username"), password):
        return username


@app.route("/")
@auth.login_required
def index():
    return f"You have successfully logged in, {auth.current_user()}"


if __name__ == "__main__":
    app.run()
```

## 资料
- [IETF: The 'Basic' HTTP Authentication Scheme](https://tools.ietf.org/html/rfc7617)
- [RESTful Authentication with Flask](https://blog.miguelgrinberg.com/post/restful-authentication-with-flask)
- [DRF Basic Authentication Guide](https://www.django-rest-framework.org/api-guide/authentication/#basicauthentication)
- [FastAPI Basic Authentication Example](https://gist.github.com/nilsdebruin/8b36cd98c9949a1a87e3a582f70146f1)

# HTTP摘要验证(HTTP Digest Authentication)
HTTP摘要验证(或称HTTP摘要权限验证)是一种比HTTP基本认证更安全的方式。这两种方法的最大的区别就是传输的密码是用md5算法的Hash字符串

## 流程
1. 未验证的客户端请求一个受限资源
2. 服务端生成一个随机的值用于提示并通过返回一个值为`Digest`的`WWW-Authenticate`Header并返回401未验证状态码，整个Header为`WWW-Authenticate: Digest nonce="44f0437004157342f50f935906ad46fc"`
3. `WWW-Authenticate`会让浏览器提示输入帐号密码
4. 在输入凭证之后，密码会进行Hash处理并在每一个请求的Header中添加一条声明信息`Authorization: Digest username="username",nonce="16e30069e45a7f47b4e2606aeeb7ab62", response="89549b93e13d438cd0946c6d93321c52"`
5. 通过账户名，服务端获取对应的密码进行Hash处理并对声明中的信息进行比较是否相同
![HTTP Digest Authentication]https://testdriven.io/static/images/blog/web-authentication-methods/digest_auth.png)

## 优点
1. 比基本验证拥有更强的安全性
2. 实现简单
3. 主流浏览器基本都支持

## 缺点
1. 必须在每个请求都携带凭证
2. 用户只能重新提交个不正确的凭证来登出
3. 与基本验证相比，由于密码不能加密保存因此在服务端安全性较差
4. 很容易收到中间人袭击

## 第三方依赖
- [Flask-HTTPAuth](https://flask-httpauth.readthedocs.io/)

## 代码
用`Flask-HTTP`可以使用Flask非常轻松的实现
```python
from flask import Flask
from flask_httpauth import HTTPDigestAuth

app = Flask(__name__)
app.config["SECRET_KEY"] = "change me"
auth = HTTPDigestAuth()

users = {
    "username": "password"
}


@auth.get_password
def get_user(username):
    if username in users:
        return users.get(username)


@app.route("/")
@auth.login_required
def index():
    return f"You have successfully logged in, {auth.current_user()}"


if __name__ == "__main__":
    app.run()
```

## 资料
- [IETF: HTTP Digest Access Authentication](https://tools.ietf.org/html/rfc7616)
- [Digest Authentication from the Requests library](https://2.python-requests.org/en/latest/user/authentication/#digest-authentication)

# 会话认证(Session-based Auth)
通过会话，用户状态可以存储在服务端上。这并不需要用户每次请求都携带账户密码信息。而是在他们登录过后，在服务端会验证登录凭证，如果是合法的，就会创建一个会话并存储在会话池中，然后返回该会话的唯一标识ID(Session ID)给浏览器，浏览器会将该ID当作cookie保存在浏览器中并在每次请求时携带该cookie

## 流程
![Session-based Auth](https://testdriven.io/static/images/blog/web-authentication-methods/session_auth.png)

## 优点
1. 由于不需要携带登录凭证，在后续的判断登录会十分迅速
2. 提高了用户体验
3. 容易实现，很多框架(例如Django)提供了这种验证方法并开箱即用

## 缺点
1. Session是有状态的。服务端会跟踪每一个会话。用于存储会话信息的会话池需要给众多服务提供验证服务。因此，这对RESTful服务来说并不友好，因为REST是一种无状态协议
2. 每一个请求都会携带cookie信息，就算是不需要验证的请求也如此
3. 对CSRF攻击的防护不足([什么是CSRF，如何在Flask中抵御CSRF攻击]https://testdriven.io/blog/csrf-flask/))

## 第三方依赖
1. [Flask-Login](https://flask-login.readthedocs.io/)
2. [Flask-HTTPAuth](https://flask-httpauth.readthedocs.io/)
3. [User authentication in Django](https://docs.djangoproject.com/en/3.1/topics/auth/)
4. [FastAPI-Login](https://github.com/MushroomMaula/fastapi_login)
5. [FastAPI-Users](https://github.com/frankie567/fastapi-users)

## 代码
`Flask-login`对会话验证非常合适，该包负责登录和注销，并可以在一段时间内记住用户信息
```python
from flask import Flask, request
from flask_login import (
    LoginManager,
    UserMixin,
    current_user,
    login_required,
    login_user,
)
from werkzeug.security import generate_password_hash, check_password_hash


app = Flask(__name__)
app.config.update(
    SECRET_KEY="change_this_key",
)

login_manager = LoginManager()
login_manager.init_app(app)


users = {
    "username": generate_password_hash("password"),
}


class User(UserMixin):
    ...


@login_manager.user_loader
def user_loader(username: str):
    if username in users:
        user_model = User()
        user_model.id = username
        return user_model
    return None


@app.route("/login", methods=["POST"])
def login_page():
    data = request.get_json()
    username = data.get("username")
    password = data.get("password")

    if username in users:
        if check_password_hash(users.get(username), password):
            user_model = User()
            user_model.id = username
            login_user(user_model)
        else:
            return "Wrong credentials"
    return "logged in"


@app.route("/")
@login_required
def protected():
    return f"Current user: {current_user.id}"


if __name__ == "__main__":
    app.run()
```

## 资料
- [IETF: Cookie-based HTTP Authentication](https://tools.ietf.org/id/draft-broyer-http-cookie-auth-00.html)
- [How To Add Authentication to Your App with Flask-Login](https://www.digitalocean.com/community/tutorials/how-to-add-authentication-to-your-app-with-flask-login)
- [Session-based Auth with Flask for Single Page Apps](https://testdriven.io/blog/flask-spa-auth/)
- [CSRF Protection in Flask](https://testdriven.io/blog/csrf-flask/)
- [Django Login and Logout Tutorial](https://learndjango.com/tutorials/django-login-and-logout-tutorial)
- [Django Session-based Auth for Single Page Apps](https://testdriven.io/blog/django-spa-auth/)
- [FastAPI-Users: Cookie Auth](https://frankie567.github.io/fastapi-users/configuration/authentication/cookie/)

# 令牌验证(Token-Based Authentication)
这种方法是将cookies用令牌认证替代。用户提供登录凭证服务端返回一种标识令牌，这个令牌将会在后续的请求中携带

最常见并广泛应用的令牌是`JSON Web Token`([JWT]https://jwt.io/))。JWT由三部分组成
- Header(携带令牌的类型和使用的Hash算法类型)
- Payload(携带对该令牌信息的声明结构体)
- Signature(用于验证信息是否在传输过程中发生错误)

以上三部分都将使用base64编码，并且每部分都会进行Hash处理，不同类型中间用`.`分割。由于这些信息都被编码过，任何人都能通过解码来查看携带的信息。但是只有验证的用户才能生成合法的令牌。令牌将会通过签名来验证合法性。

> JWT是一种紧凑的、URL安全的方法，用于表示在两方互相传输的声明。在JWT中，声明信息被编码成JSON信息，用于JWS或者JWE，从而使声明可以进行数字签名或完整性保护[IETE]https://tools.ietf.org/html/rfc7519)(翻译的不是很好)

服务端并不需要保存令牌，令牌只能使用签名来验证。最近，由于RESTful API和单页应用程序的兴起，令牌验证这种方法的采用有所增加

## 流程
![Token-Based Auth](https://testdriven.io/static/images/blog/web-authentication-methods/token_auth.png)

## 优点
1. 无状态。服务端不需要存储令牌，只需要令牌中的签名用于验证。这方每个请求不需要调用数据库从而更加快速
2. 对多个服务需要验证的微服务架构来说非常合适。我们只需要配置如何处理令牌和令牌密钥。

## 缺点
1. 客户端需要考虑如何存储令牌，这可能导致XSS或者CSRF攻击
2. 令牌不可删除，只能过期。这意味着如果令牌泄漏，攻击者可以使用这个令牌直至过期。因此非常有必要对令牌设置一个较短的过期时间例如15分钟
3. 需要在令牌快过期时进行令牌的自动刷新
4. 删除令牌的唯一方式是创建一个数据库或者黑名单。这会对微服务架构来说增加额外的负担和潜在的问题。

## 第三方依赖
- [Flask-JWT-Extended](https://github.com/vimalloc/flask-jwt-extended)
- [Flask-HTTPAuth](https://flask-httpauth.readthedocs.io/)
- [Simple JWT for Django REST Framework](https://github.com/SimpleJWT/django-rest-framework-simplejwt)
- [FastAPI JWT Auth](https://github.com/IndominusByte/fastapi-jwt-auth)

## 代码
`Flask-JWT-Extended`可以提供处理JWT的许多方法
```python
from flask import Flask, request, jsonify
from flask_jwt_extended import (
    JWTManager,
    jwt_required,
    create_access_token,
    get_jwt_identity,
)
from werkzeug.security import check_password_hash, generate_password_hash

app = Flask(__name__)
app.config.update(
    JWT_SECRET_KEY="please_change_this",
)

jwt = JWTManager(app)

users = {
    "username": generate_password_hash("password"),
}


@app.route("/login", methods=["POST"])
def login_page():
    username = request.json.get("username")
    password = request.json.get("password")

    if username in users:
        if check_password_hash(users.get(username), password):
            access_token = create_access_token(identity=username)
            return jsonify(access_token=access_token), 200

    return "Wrong credentials", 400


@app.route("/")
@jwt_required
def protected():
    return jsonify(logged_in_as=get_jwt_identity()), 200


if __name__ == "__main__":
    app.run()
```

## 资料
- [Introduction to JSON Web Tokens](https://jwt.io/introduction)
- [IETF: JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519)
- [How to Use JWT Authentication with Django REST Framework](https://simpleisbetterthancomplex.com/tutorial/2018/12/19/how-to-use-jwt-authentication-with-django-rest-framework.html)
- [Securing FastAPI with JWT Token-based Authentication](https://testdriven.io/blog/fastapi-jwt-auth/)
- [JWT Authentication Best Practices](https://blog.asayer.io/jwt-authentication-best-practices)

# 一次性密码(One Time Passwords)
一次性密码通常用于确认认证消息。OTPs会随机生成代码用于验证用户是否是声明中的。也经常用于APP中的二次验证

为了使用OTPs，必须要有一个可靠可信赖的系统。这个系统将会用于验证email或电话

现代OTPs是无状态的，并且有很多方法来验证。尽管有很多类型的OTPs，但是基于时间的OTPs(TOTPs)是目前毫无争议使用最多的方法。验证代码在生成过后会在一段时间后过期。

由于你添加了一层安全检查，OTPS可以用于处理一些安全性高度需要的APP例如银行或其他金融服务

## 流程
传统的OTPS验证流程如下
- 客户端发送用户密码
- 在认证过后，服务端生成随机的代码，保存在服务端并返回给可信赖的系统
- 客户通过可信赖系统获取生成的代码，在APP中输入返回给客户端
- 服务端验证代码是否和保存的一致，如果符合则放行

基于时间的TOTPS流程如下
- 客户端发送用户密码
- 在认证后，服务端通过随机种子生成随机代码，保存在服务端并返回给可信赖的系统
- 客户通过可信赖系统获取生成的代码，在APP中输入返回给客户端
- 服务端验证代码是否和存储的代码一致并且没有过期，如果符合则放行

例如[Google Authenticator](https://en.wikipedia.org/wiki/Google_Authenticator), [Microsoft Authenticator](https://www.microsoft.com/en-us/account/authenticator), [FreeOTP](https://en.wikipedia.org/wiki/FreeOTP)的OTP代理商工作流程
- 在代理生认证2FA后，服务器生成一个随机种子，并以唯一的二维码形式发送给用户
- 用户通过2FA应用扫描二维码用于验证信赖设备
- 当OTP需要时，用户会检查设备中的代码，并且在Web中输入
- 服务端验证输入的代码，如果符合则放行

## 优点
- 添加了一层防护层
- 密码泄漏的危险性降低，或者说服务只会通过OTP验证

## 缺点
- 你需要存储OTPs生成的种子
- 如果你丢失了回复代码，OTP代理商例如Google Authenticator非常难设置
- 当可信赖设备不可使用时会产生诸多问题(没电，网络错误等等)。因此，添加备用设备非常有必要

## 第三方依赖
- [PyOTP - The Python One-Time Password Library](https://pyauth.github.io/pyotp/)
- [django-otp](https://github.com/django-otp/django-otp)

## 代码
`PyOTP`提供了基于时间和基于计数器的OTPs
```python
from time import sleep

import pyotp

if __name__ == "__main__":
    otp = pyotp.TOTP(pyotp.random_base32())
    code = otp.now()
    print(f"OTP generated: {code}")
    print(f"Verify OTP: {otp.verify(code)}")
    sleep(30)
    print(f"Verify after 30s: {otp.verify(code)}")
```
```
OTP generated: 474771
Verify OTP: True
Verify after 30s: False
```

## 资料
- [IETF: TOTP: Time-Based One-Time Password Algorithm](https://tools.ietf.org/html/rfc6238)
- [IETF: A One-Time Password System](https://tools.ietf.org/html/rfc2289)
- [Implementing 2FA: How Time-Based One-Time Password Actually Works (With Python Examples)](https://hackernoon.com/implementing-2fa-how-time-based-one-time-password-actually-works-with-python-examples-cm1m3ywt)

# OAuth和OpenID

OAuth/OAuth2和OpenID时目前比较火的认证方法。他们都要求实现社交登录(一种单点登录形式)，使用来自社交网络服务(例如Facebook，Twitter或谷歌)现有的信息介入第三方网站，而不是创建一个新的登录账户

这种类型的认证授权可以在你需要高度安全认证时使用。一些提供商有组够多的资源来处理认证。利用这些久经考验的认证系统可以让你的应用更为的安全。

这种方法通常与会话认证搭配使用

## 流程
你浏览一些需要登录的网站，你进入登录界面并发现一个`使用谷歌登录`的按钮。你点击按钮后会引导你进入谷歌登录界面。一旦登录认证完成，你会被重定向至刚刚浏览的需要登录的网站。这是一种使用OpenID认证的例子。它会让你用一个现成的账户登录而并不需要再创建一个新帐号

最著名的OpenID提供商有谷歌，Facebook，Github

在登录过后，你进入网站的下载服务，在下载大文件时直接接入到谷歌云中。网站是如何访问你的谷歌云的呢？这就是OAuth发挥作用的时候。你可以在授予访问其他网站上资源的权限，比如此时的谷歌云访问权限。

## 优点
- 提高了安全性
- 更容易且快速的登录流程，不需要额外创建一个账户
- 由于认证时无密码的，一旦出现安全漏洞，不会对第三方造成损害

## 缺点
- 你需要依赖不在你掌控的第三方APP。如果OpenID服务宕机，用户就不能再进行登录
- 人通常会护士OAuth应用的授权信息
- 没有配置OpenID提供上账户的用户将无法访问你的应用。最好的方法是本站注册和OAuth认证同时做。

## 第三方依赖
社交登录
- [Authomatic](https://authomatic.github.io/authomatic/)
- [Python Social Auth](https://python-social-auth.readthedocs.io/)
- [Flask-Dance](https://flask-dance.readthedocs.io/)
- [django-allauth](https://django-allauth.readthedocs.io/)

创建自己的OAuth或OpenID服务
- [Authlib](https://authlib.org/)
- [OAuthLib](https://oauthlib.readthedocs.io/)
- [Flask-OAuthlib](https://flask-oauthlib.readthedocs.io/)
- [Django OAuth Toolkit](https://django-oauth-toolkit.readthedocs.io/)
- [Django OIDC Provider](https://django-oidc-provider.readthedocs.io/)
- [FastAPI: Simple OAuth2 with Password and Bearer](https://fastapi.tiangolo.com/tutorial/security/simple-oauth2/)
- [FastAPI: OAuth2 with Password (and hashing), Bearer with JWT tokens](https://fastapi.tiangolo.com/tutorial/security/oauth2-jwt/#advanced-usage-with-scopes)

## 代码
你可以通过`Flask-Dance`接入Github服务
```python
from flask import Flask, url_for, redirect
from flask_dance.contrib.github import make_github_blueprint, github

app = Flask(__name__)
app.secret_key = "change me"
app.config["GITHUB_OAUTH_CLIENT_ID"] = "1aaf1bf583d5e425dc8b"
app.config["GITHUB_OAUTH_CLIENT_SECRET"] = "dee0c5bc7e0acfb71791b21ca459c008be992d7c"

github_blueprint = make_github_blueprint()
app.register_blueprint(github_blueprint, url_prefix="/login")


@app.route("/")
def index():
    if not github.authorized:
        return redirect(url_for("github.login"))
    resp = github.get("/user")
    assert resp.ok
    return f"You have successfully logged in, {resp.json()['login']}"


if __name__ == "__main__":
    app.run()
```

## 资料
- [An Illustrated Guide to OAuth and OpenID Connect](https://developer.okta.com/blog/2019/10/21/illustrated-guide-to-oauth-and-oidc)
- [Introduction to OAuth 2.0 and OpenID Connect](https://mherman.org/presentations/node-oauth-openid/#1)
- [Create a Flask Application With Google Login](https://realpython.com/flask-google-login/)
- [Django-allauth Tutorial](https://learndjango.com/tutorials/django-allauth-tutorial)
- [FastAPI — Google as an external authentication provider](https://medium.com/data-rebels/fastapi-google-as-an-external-authentication-provider-3a527672cf33)

# 结论
在这篇文章，我们认识了一些不同的Web认证方法，每一种都有各自的优点和缺点

在何时使用呢？这就需要具体问题具体分析，通常有几条规则
- 对利用服务器模板的Web应用程序，通过用户密码进行会话认证是最合适的。也可以增加OAuth和OpenID作为另一种方式
- 对RESTfulAPIs来说，无状态的令牌验证方法是最合适的
- 如果你想处理高度敏感的数据，你可能需要在认证过程中添加OTPs

最后你要知道，上述的例子只是接触到了认证的表皮。在生产环境中需要更深层的配置。