#ifndef ST_H
#define ST_H
#include <string>
#include <stdexcept>

using namespace std;

struct SNode {
    int key;
    SNode* next;
};

class Stack{
    private:
        SNode* head;
        int size;
        static const int MAX_SIZE = 2;
        
        SNode* createSNode(int key);

    public:
        Stack();

        void push(int data);

        int pop();

        bool isEmpty();

        int getSize();
};

#endif