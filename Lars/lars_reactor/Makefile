TARGET=lib/liblreactor.a   #生成静态库
CXX=g++
CFLAGS=-g -O2 -Wall -fPIC -Wno-deprecated
#调试信息、优化编译、输出警告、构建库必选、兼容老版本库(不发出警告)

SRC=./src  #源文件

INC=-I./include  #头文件

#wildcard：展开SRC目录下所有.cpp文件
#basename: 取文件名，不包括后缀.cpp   src/main.cpp  ->  src/main
#addsuffix：给文件名添加后缀.o
#将源文件下，所有的.cpp文件，换成了.o文件
OBJS = $(addsuffix .o, $(basename $(wildcard $(SRC)/*.cpp)))


$(TARGET): $(OBJS)#最终目标静态库，依赖所有.o文件
#创建lib文件夹，如果不存在的话
	mkdir -p lib  
#创建静态库语法：ar cqs libmymath.a add.o sub.o
#$@:目标文件.a，$^:所有依赖文件.o
	ar cqs $@ $^

%.o: %.cpp
#所有.o文件，依赖对应的.cpp文件
#编译语法：g++ -c -o add.o add.cpp -I./include
#-c:只编译不链接，这样只需编译修改的部分，最后统一链接   -o:指定输出文件名
#应用编译器标志：调试信息、优化编译...
	$(CXX) $(CFLAGS) -c -o $@ $< $(INC)
# $@:目标文件 %.o，$<:第一个依赖文件 %.cpp
#只要.o文件和.cpp文件名字相同，就可以用这个.cpp生成这个.o

#伪目标，不是真正的目标，只是一个标记，不会生成文件，解决文件夹中存在clean时冲突问题
.PHONY: clean
clean:
	-rm -f src/*.o $(TARGET)
#-f:强制执行，及时文件不存在，也不会中断


