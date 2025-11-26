#ifndef ARRAY_H
#define ARRAY_H

#include <string>
#include <stdexcept>

using namespace std;

struct ArNode {
    string data;
};

class Array {
private:
    ArNode* head;
    int len;
    int size;

    void extendArray();

public:
    Array(int size);
    Array(const string& inputString);
    ~Array();

    string getElement(int index);
    void setElement(const string& key, int index);
    void deleteElement(int index);
    void addElementAtIndex(const string& key, int index);
    void addElementEnd(const string& key);
    int getLength() const;
    ArNode* getHead();
    int getSize() const;
    int isInArray(const string& key) const;
};

#endif