# Go语言编写的适用于Linux和Windows的多线程agent软件

## 项目简介
使用Go语言开发的多线程agent软件<br>
**主要功能：**<br>
1. 使用http方式调用<br>
2. 支持加密传输方式<br>
3. 支持调用底层命令与脚本<br>
4. 支持同步与异步调用方式<br>
5. 支持心跳功能<br>
6. 支持windows和linux多操作系统（需分别编译）<br>

## 代码架构：
totoroAgent.go  //入口方法，初始化程序 <br>
/totoroAgent    //主要处理程序 <br>
/mahonia        //多语言支持

## 编译

### LINUX编译

安装好golang后执行build.sh <br>
编译成功后，执行/bin目录下的totoroAgent运行 <br>


### Windows编译

配置好golang和goPath后 <br>
执行
```
go install totoroAgent totoroAgent.go
```
编译完成后，执行totoroAgent.exe <br>


## 运行
1.直接执行totoroAgent，使用默认配置（默认端口10099）<br>
2.自定义配置(使用json配置启动参数) <br>
```
./totoroAgent -c config.json
```
config中可配置启动端口，进程地址，日志地址，加密key等等

## 使用
启动应用后，访问http://{agent地址}:10099/version，返回版本号即成功 <br>

### 执行命令
POST 访问http://{agent地址}:10099/exec, body内容为 <br>
```
{
  "actionType" : "exec",   
  "cmd" : "ls /export"
}
```
返回内容 <br>
```
Command exit code: 0
code
go
logs
var
```
第一行返回命令执行结果，后面返回命令执行的内容 <br>

### 异步命令
POST 访问http://{agent地址}:10099/tasks, body内容为 <br>

```
[{              //task数组，可以同时传多个task任务
"id":0,
"taskId":"1111",
"actionType":"exec",
"cmdType":"query",
"cmd":"ls /export",
"status":0,
"resultCode":0,
"resultInfo":"123",
"url":"http://127.0.0.1"   //完成后回调函数
}]
```

### 加密传输
支持des加密方式，java端加密的代码如下
```
public static void main(String[] args) {
    String content = "{\"actionType\" : \"exec\",\"cmd\" : \"ls /export\"}";
    // 加密的key
    String KEY = "totoro&&agent%#.*&$agent";
    byte[] crypted = null;
    try {
        byte[] keyBytes = KEY.getBytes();
        DESedeKeySpec desKeySpec = new DESedeKeySpec(keyBytes);
        IvParameterSpec ivSpec = new IvParameterSpec(KEY.substring(0, 8).getBytes());
        SecretKeyFactory factory = SecretKeyFactory.getInstance("DESede");
        SecretKey secKey = factory.generateSecret(desKeySpec);
        Cipher cipher = Cipher.getInstance("DESede/CBC/PKCS5Padding");
        cipher.init(Cipher.ENCRYPT_MODE, secKey, ivSpec);
        crypted =  cipher.doFinal(content.getBytes());
    }  catch (Exception e) {
        e.printStackTrace();
    }
    String encodeContent = new String(Base64.encodeBase64(crypted));
    System.out.println(encodeContent);
}
```
加密后访问连接不变，但需要设置Header，在Header里增加
```
Secure-Type = TRUE  //TRUE都是大写
```

上述例子中加密字符串结果为
```
mt4O2nVGYe5kkIEIg9Ttoygw8VnDyBYFGlG7bY7aVgxCbsu4rK+FYDJeYaCfeos5
```

调用成功返回结果
```
Command exit code: 0
code
go
logs
var

```
