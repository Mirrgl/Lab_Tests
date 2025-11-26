#ifndef CM_H
#define CM_H

#include <cstddef>
#include <fstream>
#include <string>
#include <sstream>
#include <iostream>
#include <functional>

using namespace std;

struct ChainNode {
    string key;
    int data;
    ChainNode* next;
    ChainNode(const string& k, int d);
};

struct Bucket {
    ChainNode* head;
    Bucket();
};

class ChainMap {
private:
    Bucket* table;
    size_t capacity;
    size_t size;

    size_t hashFunction(const string& key);
    void rehash();

public:
    ChainMap(size_t initial_capacity);
    ~ChainMap();

    void getAllKeys(ChainMap& result) const;
    string getAllKeysAsString() const;
    void add(const string& key, int data);
    void del(const string& key);
    bool isContain(const string& key);
    int find(const string& key);
    void printContents() const;
};

#endif