# Bill
#以Hyperlerger Fabric 为底层的票据项目

#项目前提

##开发环境搭建
1. docker
2. docker-compose
3. git
4. nodejs
5. go
```
#安装docker
curl -fsSL https://get.docker.com | bash -s docker --mirror Aliyun
#安装docker-compose
curl -L https://github.com/docker/compose/releases/download/1.20.1/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
#安装gitwdewfedawfasfaef
apt-get update
apt-get install git
#安装nodejs
curl -sL https://deb.nodesource.com/setup_8.x | sudo -E bash -
sudo apt-get install -y nodejs
#安装golang 以及go环境搭建
wget https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz(1.10.3)
 tar -zxvf go1.10.3.linux-amd64.tar.gz -C /usr/local/
 sudo vim /etc/profile
 export GOPATH=$HOME/go
 export GOBIN=$HOME/go/bin
 export GOROOT=/usr/local/go
 export PATH=$GOROOT/bin:$PATH
 ##source /etc/profile(很重要！！！！)
 #docker fabric images下载
 https://github.com/hyperledger/fabric/blob/release-1.2/scripts/bootstrap.sh => 运行sh文件（fabric-1.2.1）
 生成 bin 目录，添加到环境变量
 export PATH=/root/fabric-samples/bin:$PATH
 ##修复阿里云超时bug(很重要！！！)
vim /etc/resolv.conf 注释掉 options timeout:2 attempts:3 rotate single-request-reopen 
Fabric CA 应用与配置
sudo apt install libtool libltdl-dev(很重要！！！)
go get -u github.com/hyperledger/fabric-ca/cmd/...（fabric-ca）
cd $GOPATH/src/github.com/hyperledger/fabric-ca/
# 使用make命令编译：
$ make fabric-ca-server 
$ make fabric-ca-clien
生成 bin 目录, 目录中包含 fabric-ca-client 与 fabric-ca-server 两个可执行文件
设置环境变量
export PATH=$GOPATH/src/github.com/hyperledger/fabric-ca/bin:$PATH
 go get -u --tags nopkcs11 github.com/hyperledger/fabric/core/chaincode/shim（fabric）
 go get -u github.com/hyperledger/fabric-sdk-go（fabric-sdk-go）
 cd $GOPATH/src/github.com/hyperledger/fabric-sdk-go/
 #make depend 
 以上命令会下载如下依赖包并安装至$GOBIN目录下
github.com/axw/gocov/...
github.com/AlekSi/gocov-xml
github.com/client9/misspell/cmd/misspell
github.com/golang/lint/golint
golang.org/x/tools/cmd/goimports
github.com/golang/mock/mockgen
......
安装完成后检查$GOPATH/bin目录下文件

返回用户目录（命令:cd ~）
 vim .bashrc 在文件末添加:  export PATH=$PATH:$GOPATH/bin
##source .bashrc(很重要！！！！)

# make populate
安装vendor
返回结果
Populate script last ran 07-21-2018 on revision e230c04e with Gopkg.lock revision d489eba9
Populating vendor ...
Populating dockerd vendor ...
Cloning into 'scripts/_go/src/chaincoded/vendor/github.com/hyperledger/fabric'...
remote: Counting objects: 4530, done.
remote: Compressing objects: 100% (3778/3778), done.
remote: Total 4530 (delta 543), reused 2596 (delta 376), pack-reused 0
Receiving objects: 100% (4530/4530), 16.51 MiB | 120.00 KiB/s, done.
Resolving deltas: 100% (543/543), done.
(fabric已经更新到1.3，下载fabric-sdk-go;会匹配1.3的环境，可以修改Makefile文件，以符合当前运行环境)

进入到first-net文件夹，运行命令：(./byfn generate 以及 ./byfn up) 成功实例：
```
![OK (2)](C:\Users\wuzhanfly\Desktop\bill\img\OK (2).png)
# 主菜
票据背书的应用开发实例会对票据的应用场最进行简化，实现的业务逻辑包括 票据发布、票据背书、票据签收，票据拒收、票据查询等操作
##fixtures文件（区块链底层平台）
区块链底层平台: 提供分布式共享账本的维护、状态数据库维护、智能合约的全 生命周期管理等区块链功能，实现数据的不可篡改和智能合约的业务逻辑。
单独启动网络：进入文件加，docker-compose -f docker-compose.yaml up -d：例子：
![fabric](C:\Users\wuzhanfly\Desktop\bill\img\fabric.png)
注释：fabric网路是简单的“1+1+2“模式，以first-network为原型;象生成组织结构，ca文件，使用的是crypto-config.yaml文件等等

##blockchain文件
主要是配合mainbo处理链码的初始化和实例化两步骤
#chaincode文件：智能合约(链码)
智能合约: 智能合约通过链码来实现，包括票据发布、票据背书、票据背书签 收、票据背书拒绝等链码调用功能，链码查询包括查询持票人票据、查询待签 收票据、根据链码号码查询票据信息等。
链码即商业逻辑因为使用的是levelDB,查询使用了复合键（CreateCompositeKey）
注释：couchDB可以复查询，不用这麽麻烦

##service文件（业务层）
应用程序的后端服务，给Web应用提供RESTful的接口，处理 前端的业务请求。后端服务的基本功能包括用户管理和票据管理，通过 Hyperledger Fabri提供的Go SDK和区块链网络进行通信
##web文件（ 应用层）
Web应用采用jQuery+HTML+CSS 的前端架构编写页面，提供用户交 互的界面操作，包括用户操作的功能 业务操作的功能。用户是内置的，只提供用户登录和用户退出操作。业务操作 包括发布 查询持票人持有的票据、发起票据背书、查询待签收票据、签收票据 背书、拒绝票据背书等功能

##ps:各个层之间采用不同的接口，业务层的Go SDK、智能合约和区块链底层平台之间用gRPC的接口。

