CXX=g++
CFLAGS=-std=c++11 -g -O2 -Wall -fPIC -Wno-deprecated 

INC=-I../../include
LIB=-L../../lib -llreactor 
OBJS = $(addsuffix .o, $(basename $(wildcard *.cc)))

all:
	@$(CXX) -o lars_reactor $(CFLAGS)  lars_reactor.cpp $(INC) $(LIB)

	@-rm -r *.dSYM

	@./lars_reactor  

clean:
	-rm -f *.o lars_reactor
