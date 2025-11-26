#include "Queue.h"

Queue::Queue(size_t initSize) : head(nullptr)
                              , tail(nullptr)
                              , size(initSize) {}

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