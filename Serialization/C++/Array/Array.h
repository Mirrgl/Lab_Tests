#ifndef ARRAY_H
#define ARRAY_H

#include <stdexcept>
#include <fstream>
#include <sstream>
#include <string>
#include <iostream>

using namespace std;

struct ArNode {
    string data = "";
};

class Array {
private:
    ArNode* head;
    int len;
    int size;

    void extendArray();

public:
    Array(int size);
    Array(initializer_list<string> list);
    ~Array();

    string getElement(int index);
    void setElement(string key, int index);
    void deleteElement(int index);
    void addElementAtIndex(string key, int index);
    void addElementEnd(string key);
    int getLength() const;
    ArNode* getHead();
    int getSize() const;
    int isInArray(const string& key) const;

    void print();

    void writeBinary(const string& filename);
    void readBinary(const string& filename);

    void writeText(const string& filename);
    void readText(const string& filename);
};

#endif