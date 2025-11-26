#ifndef QUEUE_H
#define QUEUE_H

#include <string>
#include <cstddef>
#include <stdexcept>
#include <iostream>
#include <fstream>

using namespace std;

struct QNode {
    string data;
    QNode* next;
    QNode* prev;
};

class Queue {
    private:
        QNode* head;
        QNode* tail;
        size_t size;
        static const size_t MAX_SIZE = 1000;

    public: 
        Queue(size_t initSize);
        Queue(initializer_list<string> list);
        ~Queue();
    
        void enqueue(const string& value);
        string dequeue();
        void del(const string& key);
        size_t getSize();
        QNode* getHead();

        void writeBinary(const string& filename);
        void readBinary(const string& filename);
        void writeText(const string& filename);
        void readText(const string& filename);
        void print() const;
};

#endif