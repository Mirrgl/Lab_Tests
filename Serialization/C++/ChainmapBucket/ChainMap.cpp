#include "ChainMap.h"

ChainNode::ChainNode(const string& k, int d) : key(k), data(d), next(nullptr) {}

Bucket::Bucket() : head(nullptr) {}

ChainMap::ChainMap(size_t initial_capacity) : capacity(initial_capacity)
                                            , size(0)
                                            , table(new Bucket[initial_capacity]) {}

ChainMap::~ChainMap() {
    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            ChainNode* next = currentNode->next;
            delete currentNode;
            currentNode = next;
        }
    }
    delete[] this->table;
}

size_t ChainMap::hashFunction(const string& key) {
    return hash<string>{}(key) % this->capacity;
}

void ChainMap::rehash() {
    size_t newCapacity = this->capacity * 2;
    Bucket* newTable = new Bucket[newCapacity];
    
    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            ChainNode* nextNode = currentNode->next;
            size_t new_index = hash<string>{}(currentNode->key) % newCapacity;
            
            currentNode->next = newTable[new_index].head;
            newTable[new_index].head = currentNode;
            
            currentNode = nextNode;
        }
    }
    
    delete[] this->table;
    this->table = newTable;
    this->capacity = newCapacity;
}

void ChainMap::getAllKeys(ChainMap& result) const {
    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            result.add(currentNode->key, 1);
            currentNode = currentNode->next;
        }
    }
}

string ChainMap::getAllKeysAsString() const {
    string result = "";
    ChainMap tempKeys(this->capacity);
    getAllKeys(tempKeys);
    
    for (size_t i = 0; i < tempKeys.capacity; ++i) {
        ChainNode* currentNode = tempKeys.table[i].head;
        while (currentNode != nullptr) {
            result += currentNode->key;
            currentNode = currentNode->next;
        }
    }
    return result;
}

void ChainMap::add(const string& key, int data) {
    if (this->size >= this->capacity * 0.75) {
        rehash();
    }
    
    size_t index = hashFunction(key);
    ChainNode* currentNode = this->table[index].head;
    
    while (currentNode != nullptr) {
        if (currentNode->key == key) {
            currentNode->data = data;
            return;
        }
        currentNode = currentNode->next;
    }
    
    ChainNode* newNode = new ChainNode(key, data);
    newNode->next = this->table[index].head;
    this->table[index].head = newNode;
    this->size++;
}

void ChainMap::del(const string& key) {
    size_t index = hashFunction(key);
    ChainNode* currentNode = this->table[index].head;
    ChainNode* prevNode = nullptr;
    
    while (currentNode != nullptr) {
        if (currentNode->key == key) {
            if (prevNode == nullptr) {
                this->table[index].head = currentNode->next;
            } else {
                prevNode->next = currentNode->next;
            }
            delete currentNode;
            this->size--;
            return;
        }
        prevNode = currentNode;
        currentNode = currentNode->next;
    }
}

bool ChainMap::isContain(const string& key) {
    size_t index = hashFunction(key);
    ChainNode* current = this->table[index].head;
    
    while (current != nullptr) {
        if (current->key == key) {
            return true;
        }
        current = current->next;
    }
    
    return false;
}

int ChainMap::find(const string& key) {
    size_t index = hashFunction(key);
    ChainNode* current = this->table[index].head;
    
    while (current != nullptr) {
        if (current->key == key) {
            return current->data;
        }
        current = current->next;
    }
    
    throw runtime_error("В словаре нет такого ключа");
}

void ChainMap::printContents() const {
    cout << "Содержимое хеш-таблицы:" << endl;
    for (size_t i = 0; i < this->capacity; ++i) {
        cout << "[" << i << "]: ";
        ChainNode* currentNode = this->table[i].head;
        if (currentNode != nullptr) {
            while (currentNode->next != nullptr) {
                cout << currentNode->key << " -> " << currentNode->data << ", ";
                currentNode = currentNode->next;
            }
            cout << currentNode->key << " -> " << currentNode->data;
        }
        cout << endl;
    }
    cout << endl;
}

void ChainMap::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->capacity), sizeof(this->capacity));
    file.write(reinterpret_cast<const char*>(&this->size), sizeof(this->size));

    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            size_t keyLength = currentNode->key.length();
            file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
            file.write(currentNode->key.c_str(), keyLength);
            
            file.write(reinterpret_cast<const char*>(&currentNode->data), sizeof(currentNode->data));
            
            currentNode = currentNode->next;
        }
    }
    
    file.close();
}

void ChainMap::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            ChainNode* next = currentNode->next;
            delete currentNode;
            currentNode = next;
        }
    }
    delete[] this->table;

    file.read(reinterpret_cast<char*>(&this->capacity), sizeof(this->capacity));
    file.read(reinterpret_cast<char*>(&this->size), sizeof(this->size));

    this->table = new Bucket[this->capacity];

    for (size_t i = 0; i < this->size; ++i) {
        size_t keyLength;
        file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength));
        
        string key;
        key.resize(keyLength);
        file.read(&key[0], keyLength);
        
        int data;
        file.read(reinterpret_cast<char*>(&data), sizeof(data));
        
        size_t index = hashFunction(key);
        ChainNode* newNode = new ChainNode(key, data);
        newNode->next = this->table[index].head;
        this->table[index].head = newNode;
    }

    file.close();
}

void ChainMap::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->capacity << " " << this->size << endl;

    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            file << currentNode->key << " " << currentNode->data << endl;
            currentNode = currentNode->next;
        }
    }
    
    file.close();
}

void ChainMap::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    for (size_t i = 0; i < this->capacity; ++i) {
        ChainNode* currentNode = this->table[i].head;
        while (currentNode != nullptr) {
            ChainNode* next = currentNode->next;
            delete currentNode;
            currentNode = next;
        }
    }
    delete[] this->table;

    file >> this->capacity >> this->size;
    
    string dummy;
    getline(file, dummy);

    this->table = new Bucket[this->capacity];

    for (size_t i = 0; i < this->size; ++i) {
        string line;
        getline(file, line);
        
        if (line.empty()) continue;
        
        size_t spacePos = line.find_last_of(' ');
        if (spacePos == string::npos) {
            throw runtime_error("Неверный формат файла");
        }
        
        string key = line.substr(0, spacePos);
        int data = stoi(line.substr(spacePos + 1));
        
        size_t index = hashFunction(key);
        ChainNode* newNode = new ChainNode(key, data);
        newNode->next = this->table[index].head;
        this->table[index].head = newNode;
    }

    file.close();
}