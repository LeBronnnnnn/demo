#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>

#include <unistd.h>
#include <signal.h>
#include <sys/types.h>          /* See NOTES */
#include <sys/socket.h>
#include <arpa/inet.h>
#include <errno.h>

#include "tcp_server.h"

tcp_server::tcp_server(const char *ip, uint16_t port){
    bzero(&_connaddr, sizeof(_connaddr));

    //忽略一些信号 SIGHUP，SIGPIPE
    //SIGHUP：如果终端关闭了，会给当前进程发送该信号
    //SIGPIPE：如果客户端关闭了，服务器再次write，就会产生该信号
    if (signal(SIGHUP, SIG_IGN) == SIG_ERR){ //处理SIGHUP信号，处理方式是忽略(SIG_IGN)
        fprintf(stderr, "signal ignore SIGHUP\n"); //如果处理失败，返回值为SIG_ERR，在标准错误输出中格式化的打印错误信息
    }
    if (signal(SIGPIPE, SIG_IGN) == SIG_ERR){
        fprintf(stderr, "signal ignore SIGPIPE\n");
    }

    //1.创建socket
    _sockfd = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);
    if (_sockfd == -1){
        fprintf(stderr, "tcp_server::socket()\n");
        exit(1);
    }

    //2 初始化地址
    struct sockaddr_in server_addr;
    bzero(&server_addr, sizeof(server_addr));
    server_addr.sin_family = AF_INET; 
    inet_aton(ip, &server_addr.sin_addr); //inet_aton()将点分十进制的IP地址转换为网络字节序的IP地址
    server_addr.sin_port = htons(port); //htons()将主机字节序转换为网络字节序

    //2-1可以多次监听，设置REUSE属性(服务器重启后，立刻能重新使用之前的端口)
    int op =1; //                       设置该属性    设置的值(1)     
    if (setsockopt(_sockfd, SOL_SOCKET, SO_REUSEADDR, &op, sizeof(op)) < 0){
        fprintf(stderr, "setsocketopt SO_REUSEADDR\n");
        exit(1);
    }

    //3.绑定端口
    if (bind(_sockfd, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0){
        fprintf(stderr, "bind error\n");
        exit(1);
    }

    //4.监听ip+端口
    if (listen(_sockfd, 500) < 0){
        fprintf(stderr, "listen error\n");
        exit(1);
    }
}

//开始提供创建连接服务
void tcp_server::do_accept(){
    int connfd;
    while(true){
        //accept阻塞等待客户端连接
        printf("begin accept\n");
        connfd = accept(_sockfd, (struct sockaddr *)&_connaddr, &_connlen);
        if(connfd == -1){
            //解析错误
            if (errno == EINTR) {//中断错误，不是致命的，跳过本次任务，服务器继续运行
                fprintf(stderr, "accept errno=EINTR\n");
                continue;
            }
            else if (errno == EMFILE) {//文件描述符用完了，无法再接受新的连接
                //建立链接过多，资源不够
                fprintf(stderr, "accept errno=EMFILE\n");
            }
            else if (errno == EAGAIN) {//非阻塞模式下，比如ET循环读(一次将fd里所有数据读完)，没数据之后，阻塞等待
            //因此改变fd文件描述符属性，添加非阻塞属性，当缓冲区没数据之后，返回-1，errno解析为EGAING，说明读到空，应该break，跳出循环
                fprintf(stderr, "accept errno=EAGAIN\n");
                break;
            }
            else {
                fprintf(stderr, "accept error");
                exit(1);
            }
        }else{
            //accept succ
            //TODO：添加心跳机制

            //TODO：消息队列机制

            int writed;
            char *data = "Hello Lars\n";
            do{
                writed = write(connfd, data, strlen(data)+1);
            }while(writed == -1 && errno == EINTR);
            //如果返回-1，并且errno=EINTR，慢速系统调用write被信号中断
            //处理完信号后，不会执行write，因此需要while循环，再次重启write系统调用

            if(writed>0){
                //写成功，返回写入字节数
                printf("write succ!\n");
            }

            if(writed == -1 && errno == EAGAIN ){
                //EAGAIN：写缓冲区满了，写不进去了
                writed = 0;//表示本次写入0字节的数据，不视为错误
            }
        }
    }
}

//释放连接
tcp_server::~tcp_server(){
    close(_sockfd);
}

