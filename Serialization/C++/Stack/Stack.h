#ifndef ST_H
#define ST_H
#include <string>
#include <stdexcept>
#include <iostream>
#include <fstream>

using namespace std;

struct SNode {
    string key;
    SNode* next;
};

class Stack{
    private:
        SNode* head;
        int size;
        static const int MAX_SIZE = 10;
        
        SNode* createSNode(string key);

    public:
        Stack();
        Stack(initializer_list<string> list);

        void push(string data);

        string pop();

        bool isEmpty();

        int getSize();

        void writeBinary(const string& filename);
        void readBinary(const string& filename);
        void writeText(const string& filename);
        void readText(const string& filename);
        void print();
};

#endif