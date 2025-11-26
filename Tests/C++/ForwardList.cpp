#include "ForwardList.h"

using namespace std;

ForwardList::ForwardList() : head(nullptr)
                           , tail(nullptr)
                           , size(0) {}

ForwardList::~ForwardList() {
    FWNode* current = head;
    while (current) {
        FWNode* next = current->next;
        delete current;
        current = next;
    }
}

void ForwardList::validatePosition(int position, bool allowEnd) const {
    int maxPos = allowEnd ? size : size - 1;
    if (position < 0 || position > maxPos) {
        throw out_of_range("Индекс за пределами списка");
    }
}

FWNode* ForwardList::getNodeAt(int position) const {
    validatePosition(position, false);
    
    FWNode* current = head;
    for (int i = 0; i < position; i++) {
        current = current->next;
    }
    return current;
}

void ForwardList::pushBack(const string& key) {
    FWNode* newNode = new FWNode{key};
    
    if (!head) {
        head = tail = newNode;
    } else {
        tail->next = newNode;
        tail = newNode;
    }
    size++;
}

void ForwardList::pushFront(const string& key) {
    FWNode* newNode = new FWNode{key};
    newNode->next = head;
    
    if (!head) {
        tail = newNode;
    }
    head = newNode;
    size++;
}

void ForwardList::insertBefore(const string& key, int position) {
    if (position == 0) {
        pushFront(key);
        return;
    }
    
    validatePosition(position, false);
    FWNode* prev = getNodeAt(position - 1);
    FWNode* newNode = new FWNode{key};
    newNode->next = prev->next;
    prev->next = newNode;
    
    if (!newNode->next) {
        tail = newNode;
    }
    size++;
}

void ForwardList::insertAfter(const string& key, int position) {
    validatePosition(position, false);
    
    if (position == size - 1) {
        pushBack(key);
        return;
    }
    
    FWNode* current = getNodeAt(position);
    FWNode* newNode = new FWNode{key};
    newNode->next = current->next;
    current->next = newNode;
    size++;
}

void ForwardList::popFront() {
    if (!head) {
        throw runtime_error("Список пустой");
    }
    
    FWNode* temp = head;
    head = head->next;
    
    if (!head) {
        tail = nullptr;
    }
    
    delete temp;
    size--;
}

void ForwardList::popBack() {
    if (!head) {
        throw runtime_error("Список пустой");
    }
    
    if (head == tail) {
        delete head;
        head = tail = nullptr;
    } else {
        FWNode* current = head;
        while (current->next != tail) {
            current = current->next;
        }
        
        delete tail;
        current->next = nullptr;
        tail = current;
    }
    size--;
}

void ForwardList::removeAfter(FWNode* prevNode) {
    if (!prevNode || !prevNode->next) {
        throw invalid_argument("Нода не имеет дочерней. Нечего удалять");
    }
    
    FWNode* toDelete = prevNode->next;
    prevNode->next = toDelete->next;
    
    if (toDelete == tail) {
        tail = prevNode;
    }
    
    delete toDelete;
    size--;
}

bool ForwardList::removeByValue(const string& value) {
    if (!head) return false;
    
    if (head->key == value) {
        popFront();
        return true;
    }
    
    FWNode* current = head;
    while (current->next && current->next->key != value) {
        current = current->next;
    }
    
    if (current->next) {
        FWNode* toDelete = current->next;
        current->next = toDelete->next;
        
        if (toDelete == tail) {
            tail = current;
        }
        
        delete toDelete;
        size--;
        return true;
    }
    
    return false;
}

string ForwardList::front() const {
    if (!head) {
        throw runtime_error("Список пустой");
    }
    return head->key;
}

string ForwardList::back() const {
    if (!tail) {
        throw runtime_error("Список пустой");
    }
    return tail->key;
}

string ForwardList::getAt(size_t index) const {
    return getNodeAt(index)->key;
}

FWNode* ForwardList::findByValue(const string& value) const {
    FWNode* current = head;
    while (current) {
        if (current->key == value) return current;
        current = current->next;
    }
    return nullptr;
}

bool ForwardList::isEmpty() const {
    return !head;
}

size_t ForwardList::getSize() const {
    return size;
}