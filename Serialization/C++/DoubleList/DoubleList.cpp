#include "DoubleList.h"

using namespace std;

DoubleList::DoubleList(initializer_list<string> list) {
    for (string item : list) {
        addTail(item);
    }
}

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

void DoubleList::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->length), sizeof(this->length));

    DFNode* current = head;
    while (current != nullptr) {
        size_t keyLength = current->key.length();
        file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
        file.write(current->key.c_str(), keyLength);
        current = current->next;
    }
    
    file.close();
}

void DoubleList::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        deleteHead();
    }

    size_t newLength;
    file.read(reinterpret_cast<char*>(&newLength), sizeof(newLength));

    for (size_t i = 0; i < newLength; ++i) {
        size_t keyLength;
        file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength));

        string key;
        key.resize(keyLength);
        file.read(key.data(), keyLength);
        addTail(key);
    }

    file.close();
}

void DoubleList::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->length << endl;

    DFNode* current = head;
    while (current != nullptr) {
        file << current->key << endl;
        current = current->next;
    }
    
    file.close();
}

void DoubleList::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        deleteHead();
    }

    string lengthStr;
    getline(file, lengthStr);
    size_t newLength = stoul(lengthStr);

    for (size_t i = 0; i < newLength; ++i) {
        string key;
        getline(file, key);
        if (!key.empty()) {
            addTail(key);
        }
    }

    file.close();
}

void DoubleList::print() const {
    if (isEmpty()) {
        cout << "Список пуст" << endl;
        return;
    }
    
    DFNode* current = head;
    while (current != nullptr) {
        cout << current->key << " ";
        current = current->next;
    }
    cout << endl;
}