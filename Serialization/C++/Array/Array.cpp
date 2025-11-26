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

Array::Array(int size) : len(0)
                       , size(size) {
    if (size < 1) { 
        throw runtime_error("Невозможно создать массив нулевого размера"); 
    }
    this->head = new ArNode[size];
}

Array::Array(initializer_list<string> list) : len(0)
                                            , size(list.size()) {
    this->head = new ArNode[size];
    for (string item : list) {
        addElementEnd(item);
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

void Array::setElement(string key, int index) {
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

void Array::addElementAtIndex(string key, int index) {
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

void Array::addElementEnd(string key) {
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

void Array::print() {
    for (size_t i = 0; i < this->len - 1; i++) {
        cout << this->head[i].data << " ";
    }
    cout << this->head[this->len - 1].data << endl;
}

void Array::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->len), sizeof(this->len));

    for (size_t i = 0; i < this->len; i++) {
        size_t stringLength = this->head[i].data.length();
        file.write(reinterpret_cast<const char*>(&stringLength), sizeof(stringLength));
        file.write(this->head[i].data.c_str(), stringLength);
    }
    
    file.close();
}

void Array::readBinary(const string& filename){
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    file.read(reinterpret_cast<char*>(&this->size), sizeof(this->size));
    delete[] head;
    this->head = new ArNode[this->size];

    for (int i = 0; i < this->size; i++) {
        size_t stringLength;
        file.read(reinterpret_cast<char*>(&stringLength), sizeof(stringLength));
        
        string inputString;
        inputString.resize(stringLength);
        file.read(inputString.data(), stringLength);
        
        this->head[i].data = inputString;
        this->len++;
    }

    file.close();
}

void Array::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->len << endl;

    for (size_t i = 0; i < this->len; i++) {
        file << this->head[i].data << endl;
    }
    
    file.close();
}

void Array::readText(const string& filename){
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    string inputSize;
    getline(file, inputSize);
    this->size = stoi(inputSize);
    this->head = new ArNode[this->size];

    for (int i = 0; i < this->size; i++) {
        string inputString;
        getline(file, inputString);
        
        this->head[i].data = inputString;
        this->len++;
    }

    file.close();
}