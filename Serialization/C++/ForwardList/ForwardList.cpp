#include "ForwardList.h"

using namespace std;

ForwardList::ForwardList() : head(nullptr)
                           , tail(nullptr)
                           , size(0) {}

ForwardList::ForwardList(initializer_list<string> list) : head(nullptr)
                                                        , tail(nullptr)
                                                        , size(0) {
    for (string item : list) {
        pushBack(item);
    }
}

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

void ForwardList::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->size), sizeof(this->size));

    FWNode* current = head;
    while (current != nullptr) {
        size_t keyLength = current->key.length();
        file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
        file.write(current->key.c_str(), keyLength);
        current = current->next;
    }
    
    file.close();
}

void ForwardList::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        popFront();
    }

    size_t fileSize;
    file.read(reinterpret_cast<char*>(&fileSize), sizeof(fileSize));

    for (size_t i = 0; i < fileSize; ++i) {
        size_t keyLength;
        file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength));

        string key;
        key.resize(keyLength);
        file.read(key.data(), keyLength);
        pushBack(key);
    }

    file.close();
}

void ForwardList::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->size << endl;

    FWNode* current = head;
    while (current != nullptr) {
        file << current->key << endl;
        current = current->next;
    }
    
    file.close();
}

void ForwardList::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        popFront();
    }

    string sizeStr;
    getline(file, sizeStr);
    size_t fileSize = stoul(sizeStr);

    for (size_t i = 0; i < fileSize; ++i) {
        string keyStr;
        getline(file, keyStr);
        if (!keyStr.empty()) {
            pushBack(keyStr);
        }
    }

    file.close();
}

void ForwardList::print() const {
    if (isEmpty()) {
        cout << "Список пуст" << endl;
        return;
    }
    
    FWNode* current = head;
    while (current != nullptr) {
        cout << current->key << " ";
        current = current->next;
    }
    cout << endl;
}