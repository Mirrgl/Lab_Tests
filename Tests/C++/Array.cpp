#include "Array.h"

void Array::extendArray() {
    int newSize = this->size * 2;
    ArNode* newHead = new ArNode[newSize];
    
    for (int i = 0; i < this->len; i++) {
        newHead[i].data = this->head[i].data;
    }
    
    delete[] this->head;
    this->head = newHead;
    this->size = newSize;
}

Array::Array(int size) : len(0), size(size) {
    if (size < 1) { 
        throw runtime_error("Невозможно создать массив нулевого размера"); 
    }
    this->head = new ArNode[size];
}

Array::Array(const string& inputString) : len(0), size(inputString.size() + 1) {
    if (size < 1) { 
        throw runtime_error("Невозможно создать массив нулевого размера"); 
    }
    this->head = new ArNode[size];
    
    for (char c : inputString) {
        addElementEnd(string(1, c));
    }
}

Array::~Array() {
    delete[] this->head;
    this->head = nullptr;  
    this->len = 0;
    this->size = 0;
}

string Array::getElement(int index) {
    if (index < 0 || index >= this->len) { 
        throw range_error("Выход за пределы массива"); 
    }
    return this->head[index].data;
}

void Array::setElement(const string& key, int index) {
    if (index < 0 || index >= this->len) { 
        throw range_error("Выход за пределы массива"); 
    }
    this->head[index].data = key;
}

void Array::deleteElement(int index) {
    if (index < 0 || index >= this->len) { 
        throw range_error("Выход за пределы массива"); 
    }
    for (int i = index; i < this->len - 1; i++) {
        this->head[i].data = this->head[i + 1].data;  
    }
    this->len -= 1;
}

void Array::addElementAtIndex(const string& key, int index) {
    if (index < 0 || index > this->len) { 
        throw range_error("Выход за пределы массива"); 
    }
    if (this->len >= this->size) { 
        extendArray(); 
    }
    for (int i = this->len; i > index; i--) {
        this->head[i].data = this->head[i - 1].data;  
    }
    this->head[index].data = key;
    this->len += 1;
}

void Array::addElementEnd(const string& key) {
    if (this->len >= this->size) { 
        extendArray(); 
    }
    this->head[this->len].data = key;
    this->len += 1;
}

int Array::getLength() const {
    return this->len;
}

ArNode* Array::getHead() {
    return this->head;
}

int Array::getSize() const {
    return this->size;
}

int Array::isInArray(const string& key) const {
    for (int i = 0; i < this->len; i++) {
        if (this->head[i].data == key) {
            return i;
        }
    }
    return -1;
}