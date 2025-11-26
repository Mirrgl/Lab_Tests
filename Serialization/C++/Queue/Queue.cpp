#include "Queue.h"

Queue::Queue(size_t initSize) : head(nullptr)
                              , tail(nullptr)
                              , size(initSize) {}

Queue::Queue(initializer_list<string> list) : head(nullptr)
                                            , tail(nullptr)
                                            , size(0) {
    for (string item : list) {
        enqueue(item);
    }
}

Queue::~Queue() {
    QNode* currentNode = head;
    QNode* prevNode;
    while (currentNode != nullptr) {
        prevNode = currentNode;
        currentNode = currentNode->next;
        delete prevNode;
    }
}

void Queue::enqueue(const string& value) {
    if (this->size >= MAX_SIZE) {
        throw overflow_error("Переполнение очереди");
    }

    QNode* newNode = new QNode{value, nullptr, this->tail};
    if (size == 0) {
        this->head = newNode;
        this->tail = newNode;
    }
    else {
        this->tail->next = newNode;
        this->tail = newNode;
    }
    this->size++;
}

string Queue::dequeue() {
    if (this->size == 0) {
        throw underflow_error("Очередь пустая");
    }
    
    QNode* currentNode;
    if (this->size == 1) {
        string data = this->head->data;
        currentNode = this->head;
        this->head = nullptr;
        this->tail = nullptr;
        this->size--;
        delete currentNode;
        return data;
    }

    string data = this->head->data;
    currentNode = this->head;
    this->head = this->head->next;
    this->head->prev = nullptr;
    this->size--;
    delete currentNode;
    return data;
}

void Queue::del(const string& key) {
    QNode* currentNode = this->head;
    while (currentNode != nullptr && currentNode->data != key) {
        currentNode = currentNode->next;
    }
    
    if (currentNode == nullptr) {
        return;
    }
    
    if (currentNode->prev == nullptr) {
        this->head = currentNode->next;
        if (this->head != nullptr) {
            this->head->prev = nullptr;
        }
    } else {
        currentNode->prev->next = currentNode->next;
    }
    
    if (currentNode->next == nullptr) {
        this->tail = currentNode->prev;
        if (this->tail != nullptr) {
            this->tail->next = nullptr;
        }
    } else {
        currentNode->next->prev = currentNode->prev;
    }
    
    this->size--;
    delete currentNode;
}

size_t Queue::getSize() {
    return this->size;
}

QNode* Queue::getHead() {
    return this->head;
}

void Queue::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->size), sizeof(this->size));

    QNode* current = head;
    while (current != nullptr) {
        size_t keyLength = current->data.length();
        file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
        file.write(current->data.c_str(), keyLength);
        current = current->next;
    }
    
    file.close();
}

void Queue::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        dequeue();
    }

    size_t fileSize;
    file.read(reinterpret_cast<char*>(&fileSize), sizeof(fileSize));

    for (size_t i = 0; i < fileSize; ++i) {
        size_t keyLength;
        file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength));

        string key;
        key.resize(keyLength);
        file.read(key.data(), keyLength);
        enqueue(key);
    }

    file.close();
}

void Queue::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->size << endl;

    QNode* current = head;
    while (current != nullptr) {
        file << current->data << endl;
        current = current->next;
    }
    
    file.close();
}

void Queue::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (head != nullptr) {
        dequeue();
    }

    string sizeStr;
    getline(file, sizeStr);
    size_t fileSize = stoul(sizeStr);

    if (fileSize > MAX_SIZE) {
        throw runtime_error("Размер очереди в файле превышает максимально допустимый");
    }

    for (size_t i = 0; i < fileSize; ++i) {
        string value;
        getline(file, value);
        enqueue(value);
    }

    file.close();
}

void Queue::print() const {
    if (size == 0) {
        cout << "Очередь пуста" << endl;
        return;
    }
    
    QNode* current = head;
    while (current != nullptr) {
        cout << current->data << " ";
        current = current->next;
    }
    cout << endl;
}