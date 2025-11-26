#ifndef DFL_H
#define DFL_H

#include <stdexcept>
#include <algorithm>
#include <cstddef>
#include <string>

using namespace std;

struct DFNode {
    string key;
    DFNode* next = nullptr;
    DFNode* prev = nullptr;
};

class DoubleList {
private:
    DFNode* head = nullptr;
    DFNode* tail = nullptr;
    size_t length = 0;

    void validateIndex(int index, bool allowEnd = false) const;
    DFNode* getNodeAt(int index) const;

public:
    ~DoubleList();
    
    void addAfter(const string& key, int index);
    void addBefore(const string& key, int index);
    void addHead(const string& key);
    void addTail(const string& key);

    void deleteAt(int index);
    void deleteHead();
    void deleteTail();
    void deleteByValue(const string& key);

    string getElement(int index) const;
    string popElement(int index);
    DFNode* findByValue(const string& key) const;

    bool isEmpty() const;
    size_t getLength() const;
};

#endif