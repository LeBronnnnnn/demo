CXX=g++
CFLAGS=-g -O2 -Wall -fPIC -Wno-deprecated 

INC=-I ../../include  #头文件路径
LIB=-L ../../lib -l lreactor  #静态库文件路径，及库

#将所有.cc文件，换成.o文件
OBJS = $(addsuffix .o, $(basename $(wildcard *.cc)))

#all:默认make时，执行的目标
#@：不输出命令本身
all:
#编译命令：g++ -o hello_lars hello_lars.cpp -I (头文件位置) -L (静态库路径) -l (静态库去头尾liblreactor.a -> lreactor)
	@$(CXX) -o hello_lars $(CFLAGS)  hello_lars.cpp $(INC) $(LIB)
	
#强制删除所有dSYM过程文件，及时没有也不报错
	@-rm -r *.dSYM

#执行生成的可执行文件，配合code-runner插件使用，修改其配置文件，运行make命令，完成编译和运行
	@./hello_lars  

clean:
	-rm -f *.o hello_lars



#makefile  在code-runner的配置json文件中（设置搜索Code-runner: Executor Map），搜索cpp，修改如下
#
    # //同一文件夹下的多文件联合编译
    # //"cpp": "cd $dir && g++ *.cpp -o $fileNameWithoutExt && $dir$fileNameWithoutExt",
    # //通过makefile联合编译不同文件夹下的文件
    # "cpp":"cd $dir && make",
	