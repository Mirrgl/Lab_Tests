#include "Stack.h"

SNode* Stack::createSNode(int key){
    SNode* newSNode = new SNode;
    newSNode->key = key;
    newSNode->next = nullptr;
    return newSNode;
}

Stack::Stack() : head(nullptr)
               , size(0) {}

void Stack::push(int data){
    if (this->size >= this->MAX_SIZE){
        throw overflow_error("Stack overflow: Maximum size reached");
    }
    SNode* newSNode = createSNode(data);
    newSNode->next = this->head;
    this->head = newSNode;
    this->size++;
}

int Stack::pop(){
    if (this->head == nullptr){
        throw underflow_error("Stack underflow: Stack is empty");
    }
    
    SNode* oldhead = this->head;
    int data = oldhead->key;
    this->head = this->head->next;
    delete oldhead;
    this->size--;
    return data;
}

bool Stack::isEmpty() {
    return this->head == nullptr;
}

int Stack::getSize() {
    return this->size;
}