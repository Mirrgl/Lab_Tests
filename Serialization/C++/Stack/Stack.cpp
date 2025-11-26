#include "Stack.h"

SNode* Stack::createSNode(string key){
    SNode* newSNode = new SNode;
    newSNode->key = key;
    newSNode->next = nullptr;
    return newSNode;
}

Stack::Stack() : head(nullptr)
               , size(0) {}

Stack::Stack(initializer_list<string> list) : head(nullptr)
                                            , size(0) {
    for (string item : list) {
        push(item);
    }
}

void Stack::push(string data){
    if (this->size >= this->MAX_SIZE){
        throw overflow_error("Stack overflow: Maximum size reached");
    }
    SNode* newSNode = createSNode(data);
    newSNode->next = this->head;
    this->head = newSNode;
    this->size++;
}

string Stack::pop(){
    if (this->head == nullptr){
        throw underflow_error("Stack underflow: Stack is empty");
    }
    
    SNode* oldhead = this->head;
    string data = oldhead->key;
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

void Stack::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->size), sizeof(this->size));

    SNode* current = head;
    while (current != nullptr) {
        size_t keyLength = current->key.length();
        file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
        file.write(current->key.c_str(), keyLength);
        current = current->next;
    }
    
    file.close();
}

void Stack::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (!isEmpty()) {
        pop();
    }

    int fileSize;
    file.read(reinterpret_cast<char*>(&fileSize), sizeof(fileSize));

    string* tempArray = new string[fileSize];
    
    for (int i = fileSize - 1; i >= 0; --i) {
        size_t keyLength;
        if (!file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength))) {
            delete[] tempArray;
            throw runtime_error("Ошибка чтения длины строки из файла");
        }
        
        tempArray[i].resize(keyLength);
        if (!file.read(&tempArray[i][0], keyLength)) {
            delete[] tempArray;
            throw runtime_error("Ошибка чтения строки из файла");
        }
    }

    for (int i = 0; i < fileSize; ++i) {
        push(tempArray[i]);
    }

    delete[] tempArray;
    file.close();
}

void Stack::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->size << endl;

    SNode* current = head;
    while (current != nullptr) {
        file << current->key << endl;
        current = current->next;
    }
    
    file.close();
}

void Stack::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    while (!isEmpty()) {
        pop();
    }

    string sizeStr;
    getline(file, sizeStr);
    int fileSize = stoi(sizeStr);

    if (fileSize > MAX_SIZE) {
        throw runtime_error("Размер стека в файле превышает максимально допустимый");
    }

    string* tempArray = new string[fileSize];
    
    for (int i = fileSize - 1; i >= 0; --i) {
        getline(file, tempArray[i]);
    }

    for (int i = 0; i < fileSize; ++i) {
        push(tempArray[i]);
    }

    delete[] tempArray;
    file.close();
}

void Stack::print() {
    if (isEmpty()) {
        cout << "Стек пуст" << endl;
        return;
    }
    
    SNode* current = head;
    while (current != nullptr) {
        cout << current->key << " ";
        current = current->next;
    }
    cout << endl;
}