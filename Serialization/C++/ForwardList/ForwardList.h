#ifndef FL_H
#define FL_H

#include <stdexcept>
#include <algorithm>
#include <cstddef>
#include <string>
#include <fstream>
#include <iostream>

using namespace std;

struct FWNode {
    string key;
    FWNode* next = nullptr;
};

class ForwardList {
private:
    FWNode* head = nullptr;
    FWNode* tail = nullptr;
    size_t size = 0;

    void validatePosition(int position, bool allowEnd = false) const;
    FWNode* getNodeAt(int position) const;

public:
    ForwardList();
    ForwardList(initializer_list<string> list);
    ~ForwardList();
    
    ForwardList(const ForwardList&) = delete;
    ForwardList& operator=(const ForwardList&) = delete;

    void pushBack(const string& key);
    void pushFront(const string& key);
    void insertBefore(const string& key, int position);
    void insertAfter(const string& key, int position);

    void popFront();
    void popBack();
    void removeAfter(FWNode* prevNode);
    bool removeByValue(const string& value);

    string front() const;
    string back() const;
    string getAt(size_t index) const;

    FWNode* findByValue(const string& value) const;

    bool isEmpty() const;
    size_t getSize() const;

    void writeBinary(const string& filename);
    void readBinary(const string& filename);
    void writeText(const string& filename);
    void readText(const string& filename);
    void print() const;
};

#endif