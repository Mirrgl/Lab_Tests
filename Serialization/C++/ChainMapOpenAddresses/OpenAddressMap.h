#ifndef OAM_H
#define OAM_H

#include <cstddef>
#include <fstream>
#include <string>
#include <sstream>
#include <functional>
#include <iostream>
#include <stdexcept>

using namespace std;

enum CellState {
    EMPTY,
    FULL,
    DELETED
};

struct OAMapNode {
    string key;
    int data;
    CellState state;
    
    OAMapNode();
};

class OpenAddressingMap {
private:
    OAMapNode* table;
    size_t capacity;
    size_t size;

    size_t hashFunction(const string& key) const;
    size_t linearProbing(size_t index, size_t attempt, size_t capacity) const;
    size_t findIndex(const string& key, bool forInsert = false) const;
    void rehash();

public:
    OpenAddressingMap(size_t initial_capacity = 16);
    ~OpenAddressingMap();

    void add(const string& key, int data);
    int* find(const string& key);
    bool contains(const string& key);
    bool remove(const string& key);
    bool update(const string& key, int new_data);
    
    size_t getSize() const;
    size_t getCapacity() const;
    
    void printContents() const;

    void writeBinary(const string& filename);
    void readBinary(const string& filename);
    void writeText(const string& filename);
    void readText(const string& filename);
};

#endif