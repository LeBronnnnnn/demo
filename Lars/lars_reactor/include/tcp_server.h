#pragma once
#include <sys/socket.h> // Include the header file for socklen_t
#include <netinet/in.h> // Include the header file for sockaddr_in
#include <iostream>

class tcp_server{
public:
    //server构造函数
    tcp_server(const char *ip, uint16_t port);

    //开始接收连接
    void do_accept();

    //释放连接
    ~tcp_server();

private:
    int _sockfd;//套接字
    struct sockaddr_in _connaddr;//客户端连接地址
    socklen_t _connlen;//客户端连接地址长度
};