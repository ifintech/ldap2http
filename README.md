# 简介

将ldap协议的请求转发给后端的HTTP协议的用户授权中心认证

# 依赖

* [企业的统一认证授权中心](https://github.com/ifintech/auth)

# 使用

## 命令  

```bash
ldap-auth -host=0.0.0.0 -port=10389 -auth_url=https://auth.com -auth_token=demo_token &
```

## 参数解释

| 参数名 | 意义 | 默认值 | 备注 | 
| - | - | - | - |
| host | 监听host地址 | 127.0.0.1 | ldap服务监听的ip，例：127.0.0.1 |
| port  | 监听端口 | 10389 | ldap服务监听的端口号，建议使用10389 |
| auth_url | 统一认证授权中心服务路径 | https://127.0.0.1/ | 员工授权服务地址的倒数第二级路径，例：假设服务路径为 https://auth.com/api/auth ，则auth_url为 https://auth.com/api （去掉/auth）|
| auth_token | 用户信息查询服务密钥 | token | 配置在统一认证授权中心conf/security/api.php文件中的ldap.password项 |


## gitlab使用

gitlab.yml配置示范
```bash
gitlab_rails['ldap_enabled'] = true
gitlab_rails['ldap_servers'] = YAML.load <<-EOS # remember to close this block with 'EOS' below
main: # 'main' is the GitLab 'provider ID' of this LDAP server
  label: 'LDAP'
  host: '0.0.0.0'
  port: 10389`
  uid: 'sAMAccountName'
  method: 'plain' # "tls" or "ssl" or "plain"
  base: 'OU=People,DC=auth,DC=democompany,DC=com'
  bind_dn: 'CN=sAMAccountName@democompany.com,OU=People,DC=auth,DC=democompany,DC=com'
  password: 'password'
  user_filter: ''
  active_directory: false
EOS
```

参考文档
* [Gitlab LDAP集成](https://docs.gitlab.com/ee/administration/auth/ldap.html)