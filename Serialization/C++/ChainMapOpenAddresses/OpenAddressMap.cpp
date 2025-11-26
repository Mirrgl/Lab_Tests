#include "OpenAddressMap.h"

using namespace std;

OAMapNode::OAMapNode() : state(CellState::EMPTY) {}

OpenAddressingMap::OpenAddressingMap(size_t initial_capacity) : capacity(initial_capacity), size(0) {
    this->table = new OAMapNode[initial_capacity];
}

OpenAddressingMap::~OpenAddressingMap() {
    delete[] this->table;
}

size_t OpenAddressingMap::hashFunction(const string& key) const {
    return hash<string>{}(key) % this->capacity;
}

size_t OpenAddressingMap::linearProbing(size_t index, size_t attempt, size_t capacity) const {
    return (index + attempt) % capacity;
}

size_t OpenAddressingMap::findIndex(const string& key, bool forInsert) const {
    size_t index = hashFunction(key);
    size_t attempt = 0;

    while (attempt < this->capacity) {
        size_t currentIndex = linearProbing(index, attempt, this->capacity);
        
        if (forInsert == true && this->table[currentIndex].state == CellState::EMPTY) {
            return currentIndex;
        }
        
        if (forInsert == false && this->table[currentIndex].state == CellState::FULL && this->table[currentIndex].key == key) {
            return currentIndex;
        }
        
        if (forInsert == true && this->table[currentIndex].state == CellState::DELETED) {
            return currentIndex;
        }
        
        attempt++;
    }

    throw runtime_error("Не найден индекс");
}

void OpenAddressingMap::rehash() {
    size_t newCapacity = this->capacity * 2;
    OAMapNode* newTable = new OAMapNode[newCapacity];
    
    for (size_t i = 0; i < this->capacity; ++i) {
        if (this->table[i].state == CellState::FULL) {
            const string& key = this->table[i].key;
            const int& data = this->table[i].data;
            
            size_t new_index = hashFunction(key);
            size_t attempt = 0;
            
            while (attempt < newCapacity) {
                size_t current_index = linearProbing(new_index, attempt, newCapacity);
                
                if (newTable[current_index].state == CellState::EMPTY) {
                    newTable[current_index].key = key;
                    newTable[current_index].data = data;
                    newTable[current_index].state = CellState::FULL;
                    break;
                }
                
                attempt++;
            }
        }
    }
    
    delete[] this->table;
    this->table = newTable;
    this->capacity = newCapacity;
}

void OpenAddressingMap::add(const string& key, int data) {
    if (this->size >= this->capacity * 0.75) {
        rehash();
    }
    
    size_t index = findIndex(key, true);
            
    if (this->table[index].state == CellState::FULL && this->table[index].key == key) {
        this->table[index].data = data;
    }
    else {
        this->table[index].key = key;
        this->table[index].data = data;
        this->table[index].state = CellState::FULL;
        this->size++;
    }
}

int* OpenAddressingMap::find(const string& key) {
    size_t index = findIndex(key);
    
    if (this->table[index].state == CellState::FULL) {
        return &(this->table[index].data);
    }
    
    throw runtime_error("Ключ отсутствует в словаре");
}

bool OpenAddressingMap::contains(const string& key) {
    try {
        return find(key) != nullptr;
    } catch (const runtime_error&) {
        return false;
    }
}

bool OpenAddressingMap::remove(const string& key) {
    try {
        size_t index = findIndex(key);
        
        if (this->table[index].state == CellState::FULL) {
            this->table[index].state = CellState::DELETED;
            this->size--;
            return true;
        }
    } catch (const runtime_error&) {
        // Ключ не найден
    }
    
    return false;
}

bool OpenAddressingMap::update(const string& key, int new_data) {
    try {
        int* existing_data = find(key);
        if (existing_data) {
            *existing_data = new_data;
            return true;
        }
    } catch (const runtime_error&) {
        // Ключ не найден
    }
    return false;
}

size_t OpenAddressingMap::getSize() const {
    return this->size;
}

size_t OpenAddressingMap::getCapacity() const {
    return this->capacity;
}

void OpenAddressingMap::printContents() const {
    cout << "Содержимое хеш-таблицы:" << endl;
    for (size_t i = 0; i < this->capacity; ++i) {
        cout << "[" << i << "]: ";
        switch (this->table[i].state) {
            case CellState::FULL:
                cout << this->table[i].key << " -> " << this->table[i].data;
                break;
            case CellState::DELETED:
                cout << "[DELETED]";
                break;
            case CellState::EMPTY:
                cout << "[EMPTY]";
                break;
        }
        cout << endl;
    }
    cout << endl;
}

void OpenAddressingMap::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file.write(reinterpret_cast<const char*>(&this->capacity), sizeof(this->capacity));
    file.write(reinterpret_cast<const char*>(&this->size), sizeof(this->size));

    for (size_t i = 0; i < this->capacity; ++i) {
        if (this->table[i].state == CellState::FULL) {
            size_t keyLength = this->table[i].key.length();
            file.write(reinterpret_cast<const char*>(&keyLength), sizeof(keyLength));
            file.write(this->table[i].key.c_str(), keyLength);
            
            file.write(reinterpret_cast<const char*>(&this->table[i].data), sizeof(this->table[i].data));
        }
    }
    
    file.close();
}

void OpenAddressingMap::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    delete[] this->table;

    size_t fileCapacity, fileSize;
    file.read(reinterpret_cast<char*>(&fileCapacity), sizeof(fileCapacity));
    file.read(reinterpret_cast<char*>(&fileSize), sizeof(fileSize));

    this->capacity = fileCapacity;
    this->size = 0;
    this->table = new OAMapNode[this->capacity];

    for (size_t i = 0; i < fileSize; ++i) {
        size_t keyLength;
        file.read(reinterpret_cast<char*>(&keyLength), sizeof(keyLength));
        
        string key;
        key.resize(keyLength);
        file.read(&key[0], keyLength);
        
        int data;
        file.read(reinterpret_cast<char*>(&data), sizeof(data));
        
        add(key, data);
    }

    file.close();
}

void OpenAddressingMap::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << this->capacity << " " << this->size << endl;

    for (size_t i = 0; i < this->capacity; ++i) {
        if (this->table[i].state == CellState::FULL) {
            file << this->table[i].key << " : " << this->table[i].data << endl;
        }
    }
    
    file.close();
}

void OpenAddressingMap::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    delete[] this->table;

    size_t fileCapacity, fileSize;
    file >> fileCapacity >> fileSize;
    
    string dummy;
    getline(file, dummy);

    this->size = 0;
    this->table = new OAMapNode[this->capacity];

    for (size_t i = 0; i < fileSize; ++i) {
        string line;
        getline(file, line);
        
        if (line.empty()) continue;
        
        size_t separatorPos = line.find(" : ");
        if (separatorPos == string::npos) {
            throw runtime_error("Неверный формат файла");
        }
        
        string key = line.substr(0, separatorPos);
        int data = stoi(line.substr(separatorPos + 3));
        
        add(key, data);
    }

    file.close();
}