#include "DoubleList.h"

using namespace std;

DoubleList::~DoubleList() {
    DFNode* current = head;
    while (current) {
        DFNode* next = current->next;
        delete current;
        current = next;
    }
}

void DoubleList::validateIndex(int index, bool allowEnd) const {
    int maxIndex = allowEnd ? length : length - 1;
    if (index < 0 || index > maxIndex) {
        throw out_of_range("Индекс больше возможного");
    }
}

DFNode* DoubleList::getNodeAt(int index) const {
    validateIndex(index, false);
    
    if (index <= static_cast<int>(length / 2)) {
        DFNode* current = head;
        for (int i = 0; i < index; i++) {
            current = current->next;
        }
        return current;
    } else {
        DFNode* current = tail;
        for (int i = length - 1; i > index; i--) {
            current = current->prev;
        }
        return current;
    }
}

void DoubleList::addAfter(const string& key, int index) {
    validateIndex(index, false);
    
    DFNode* newNode = new DFNode{key};
    DFNode* current = getNodeAt(index);
    
    newNode->next = current->next;
    newNode->prev = current;
    
    if (current->next) {
        current->next->prev = newNode;
    } else {
        tail = newNode;
    }
    current->next = newNode;
    
    length++;
}

void DoubleList::addBefore(const string& key, int index) {
    if (index == 0) {
        addHead(key);
    } else {
        addAfter(key, index - 1);
    }
}

void DoubleList::addHead(const string& key) {
    DFNode* newNode = new DFNode{key};
    newNode->next = head;
    
    if (head) {
        head->prev = newNode;
    } else {
        tail = newNode;
    }
    head = newNode;
    length++;
}

void DoubleList::addTail(const string& key) {
    DFNode* newNode = new DFNode{key};
    newNode->prev = tail;
    
    if (tail) {
        tail->next = newNode;
    } else {
        head = newNode;
    }
    tail = newNode;
    length++;
}

void DoubleList::deleteAt(int index) {
    validateIndex(index, false);
    
    if (index == 0) {
        deleteHead();
    } else if (index == length - 1) {
        deleteTail();
    } else {
        DFNode* toDelete = getNodeAt(index);
        toDelete->prev->next = toDelete->next;
        toDelete->next->prev = toDelete->prev;
        delete toDelete;
        length--;
    }
}

void DoubleList::deleteHead() {
    if (!head) return;
    
    DFNode* toDelete = head;
    head = toDelete->next;
    
    if (head) {
        head->prev = nullptr;
    } else {
        tail = nullptr;
    }
    
    delete toDelete;
    length--;
}

void DoubleList::deleteTail() {
    if (!tail) return;
    
    DFNode* toDelete = tail;
    tail = toDelete->prev;
    
    if (tail) {
        tail->next = nullptr;
    } else {
        head = nullptr;
    }
    
    delete toDelete;
    length--;
}

void DoubleList::deleteByValue(const string& key) {
    DFNode* current = head;
    int index = 0;
    
    while (current) {
        if (current->key == key) {
            deleteAt(index);
            return;
        }
        current = current->next;
        index++;
    }
    throw runtime_error("Ключ не найден");
}

string DoubleList::getElement(int index) const {
    return getNodeAt(index)->key;
}

string DoubleList::popElement(int index) {
    string value = getElement(index);
    deleteAt(index);
    return value;
}

DFNode* DoubleList::findByValue(const string& key) const {
    DFNode* current = head;
    while (current) {
        if (current->key == key) return current;
        current = current->next;
    }
    return nullptr;
}

bool DoubleList::isEmpty() const {
    return length == 0;
}

size_t DoubleList::getLength() const {
    return length;
}