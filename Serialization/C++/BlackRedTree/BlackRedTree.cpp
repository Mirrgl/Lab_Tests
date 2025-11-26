#include "BlackRedTree.h"

RBTNode* RBTree::findNodePosition(RBTNode* currentNode, int key) {
    if (key < currentNode->key) {
        return (currentNode->left == nullptr) ? currentNode : findNodePosition(currentNode->left, key);
    }
    if (key > currentNode->key) {
        return (currentNode->right == nullptr) ? currentNode : findNodePosition(currentNode->right, key);
    }
    throw runtime_error("Key already exists in this");
}

void RBTree::leftRotate(RBTNode* x) {    
    RBTNode* y = x->right;
    x->right = y->left;
    
    if (y->left != nullptr) {
        y->left->parent = x;
    }
    
    y->parent = x->parent;
    
    if (x->parent == nullptr) {
        this->root = y;
    } else if (x == x->parent->left) {
        x->parent->left = y;
    } else {
        x->parent->right = y;
    }
    
    y->left = x;
    x->parent = y;
}

void RBTree::rightRotate(RBTNode* y) {    
    RBTNode* x = y->left;
    y->left = x->right;
    
    if (x->right != nullptr) {
        x->right->parent = y;
    }
    
    x->parent = y->parent;
    
    if (y->parent == nullptr) {
        this->root = x;
    } else if (y == y->parent->right) {
        y->parent->right = x;
    } else {
        y->parent->left = x;
    }
    
    x->right = y;
    y->parent = x;
}

void RBTree::fixViolation(RBTNode* z) {
    if (z == nullptr) return;
    
    while (z != this->root && z->parent != nullptr && z->parent->color == 1) {
        // родитель является левым ребенком дедушки
        if (z->parent == z->parent->parent->left) {
            RBTNode* y = z->parent->parent->right; // дядя
            
            // дядя красный
            if (y != nullptr && y->color == 1) {
                z->parent->color = 0;
                y->color = 0;
                z->parent->parent->color = 1;
                z = z->parent->parent;
            } else {
                // дядя черный и z - правый ребенок
                if (z == z->parent->right) {
                    z = z->parent;
                    leftRotate(z);
                }
                
                // дядя черный и z - левый ребенок
                z->parent->color = 0;
                z->parent->parent->color = 1;
                rightRotate(z->parent->parent);
            }
        } else {
            // родитель является правым ребенком дедушки
            RBTNode* y = z->parent->parent->left; // дядя
            
            // дядя красный
            if (y != nullptr && y->color == 1) {
                z->parent->color = 0;
                y->color = 0;
                z->parent->parent->color = 1;
                z = z->parent->parent;
            } else {
                // дядя черный и z - левый ребенок
                if (z == z->parent->left) {
                    z = z->parent;
                    rightRotate(z);
                }
                
                // дядя черный и z - правый ребенок
                z->parent->color = 0;
                z->parent->parent->color = 1;
                leftRotate(z->parent->parent);
            }
        }
    }
    
    this->root->color = 0;
}

void RBTree::insert(int key) {    
    RBTNode* newNode = new RBTNode{key, nullptr, nullptr, nullptr, 1}; 
    
    if (this->root == nullptr) {
        this->root = newNode;
        newNode->color = 0; 
        return;
    }
    
    RBTNode* parent = findNodePosition(this->root, key);
    newNode->parent = parent;
    
    if (key < parent->key) {
        parent->left = newNode;
    } else {
        parent->right = newNode;
    }
    
    fixViolation(newNode);
}

RBTNode* RBTree::findNode(int key) {
    RBTNode* current = this->root;
    while (current != nullptr) {
        if (key == current->key) {
            return current;
        } else if (key < current->key) {
            current = current->left;
        } else {
            current = current->right;
        }
    }
    return nullptr;
}

void RBTree::transplant(RBTNode* u, RBTNode* v) {
    if (u->parent == nullptr) {
        this->root = v;
    } else if (u == u->parent->left) {
        u->parent->left = v;
    } else {
        u->parent->right = v;
    }
    
    if (v != nullptr) {
        v->parent = u->parent;
    }
}

void RBTree::fixDelete(RBTNode* x, RBTNode* xParent) {    
    while (x != this->root && (x == nullptr || x->color == 0)) {
        if (xParent == nullptr) {  
            break;
        }
        
        if (x == xParent->left) {
            RBTNode* w = xParent->right; // брат x
            
            if (w != nullptr && w->color == 1) {
                //  брат красный
                w->color = 0;
                xParent->color = 1;
                leftRotate(xParent);
                w = xParent->right;
            }
            
            if (w == nullptr) break;  
            
            if ((w->left == nullptr || w->left->color == 0) && 
                (w->right == nullptr || w->right->color == 0)) {
                // оба ребенка брата черные
                w->color = 1;
                x = xParent;
                xParent = x->parent;
            } else {
                if (w->right == nullptr || w->right->color == 0) {
                    //правый ребенок брата черный
                    if (w->left != nullptr) {
                        w->left->color = 0;
                    }
                    w->color = 1;
                    rightRotate(w);
                    w = xParent->right;
                }
                
                // правый ребенок брата красный
                if (w != nullptr) {
                    w->color = xParent->color;
                    if (w->right != nullptr) {
                        w->right->color = 0;
                    }
                }
                xParent->color = 0;
                leftRotate(xParent);
                x = this->root;
            }
        } else {
            // симметричный случай
            RBTNode* w = xParent->left;
            
            if (w != nullptr && w->color == 1) {
                w->color = 0;
                xParent->color = 1;
                rightRotate(xParent);
                w = xParent->left;
            }
            
            if (w == nullptr) break;  
            
            if ((w->right == nullptr || w->right->color == 0) && 
                (w->left == nullptr || w->left->color == 0)) {
                w->color = 1;
                x = xParent;
                xParent = x->parent;
            } else {
                if (w->left == nullptr || w->left->color == 0) {
                    if (w->right != nullptr) {
                        w->right->color = 0;
                    }
                    w->color = 1;
                    leftRotate(w);
                    w = xParent->left;
                }
                
                if (w != nullptr) {
                    w->color = xParent->color;
                    if (w->left != nullptr) {
                        w->left->color = 0;
                    }
                }
                xParent->color = 0;
                rightRotate(xParent);
                x = this->root;
            }
        }
    }
    
    if (x != nullptr) {
        x->color = 0;
    }
}

void RBTree::del(int key) {    
    if (this->root == nullptr) {
        throw runtime_error("this is empty");
    }
    
    RBTNode* z = findNode(key);
    if (z == nullptr) {
        throw runtime_error("Key not found in this");
    }
    
    RBTNode* y = z;
    RBTNode* x = nullptr;
    RBTNode* xParent = nullptr;
    bool yOriginalColor = y->color;
    
    if (z->left == nullptr) {
        // нет левого потомка
        x = z->right;
        xParent = z->parent;
        transplant(z, z->right);
    } else if (z->right == nullptr) {
        // нет правого потомка
        x = z->left;
        xParent = z->parent;
        transplant(z, z->left);
    } else {
        // есть оба потомка
        y = z->right;
        while (y != nullptr && y->left != nullptr) {
            y = y->left;
        }

        yOriginalColor = y->color;
        x = y->right;
        xParent = y;
        
        if (y->parent == z) {
            if (x != nullptr) {
                x->parent = y;
            }
            xParent = y;
        } else {
            transplant(y, y->right);
            y->right = z->right;
            if (y->right != nullptr) {
                y->right->parent = y;
            }
            xParent = y->parent;
        }
        
        transplant(z, y);
        y->left = z->left;
        if (y->left != nullptr) {
            y->left->parent = y;
        }
        y->color = z->color;
    }
    
    delete z;
    
    if (yOriginalColor == 0) {
        fixDelete(x, xParent);
    }
}

int RBTree::get(int key) {
    if (this->root == nullptr) {
        throw runtime_error("Key not found"); 
    }
    RBTNode* currentNode = findNode(key);
    return (currentNode) ? currentNode->key : throw runtime_error("Key not found");
}

void RBTree::deleteBRTreeNode(RBTNode* currentNode) {
    if (currentNode != nullptr) {
        RBTNode* left = currentNode->left;
        RBTNode* right = currentNode->right;
        if (left != nullptr) {
            deleteBRTreeNode(left);
        }
        if (right != nullptr) {
            deleteBRTreeNode(right);
        }
    }
    delete currentNode;
}

RBTree::RBTree() : root(nullptr){}

RBTree::~RBTree() {
    deleteBRTreeNode(this->root);
}

void RBTree::writeBinary(const string& filename) {
    ofstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    int nodeCount = countNodes(root);
    file.write(reinterpret_cast<const char*>(&nodeCount), sizeof(nodeCount));

    writeBinaryNode(file, root);
    
    file.close();
}

void RBTree::writeBinaryNode(ofstream& file, RBTNode* node) {
    if (node == nullptr) {
        return;
    }

    file.write(reinterpret_cast<const char*>(&node->key), sizeof(node->key));
    file.write(reinterpret_cast<const char*>(&node->color), sizeof(node->color));

    bool hasLeft = (node->left != nullptr);
    bool hasRight = (node->right != nullptr);
    
    file.write(reinterpret_cast<const char*>(&hasLeft), sizeof(hasLeft));
    if (hasLeft) {
        writeBinaryNode(file, node->left);
    }
    
    file.write(reinterpret_cast<const char*>(&hasRight), sizeof(hasRight));
    if (hasRight) {
        writeBinaryNode(file, node->right);
    }
}

void RBTree::readBinary(const string& filename) {
    ifstream file(filename, ios::binary);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    deleteBRTreeNode(root);
    root = nullptr;

    int nodeCount;
    file.read(reinterpret_cast<char*>(&nodeCount), sizeof(nodeCount));

    if (nodeCount > 0) {
        readBinaryNode(file, root, nullptr);
    }

    file.close();
}

void RBTree::readBinaryNode(ifstream& file, RBTNode*& node, RBTNode* parent) {
    node = new RBTNode();
    node->parent = parent;

    file.read(reinterpret_cast<char*>(&node->key), sizeof(node->key));
    file.read(reinterpret_cast<char*>(&node->color), sizeof(node->color));

    bool hasLeft, hasRight;
    
    file.read(reinterpret_cast<char*>(&hasLeft), sizeof(hasLeft));
    if (hasLeft) {
        readBinaryNode(file, node->left, node);
    } else {
        node->left = nullptr;
    }
    
    file.read(reinterpret_cast<char*>(&hasRight), sizeof(hasRight));
    if (hasRight) {
        readBinaryNode(file, node->right, node);
    } else {
        node->right = nullptr;
    }
}

void RBTree::writeText(const string& filename) {
    ofstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл для записи: " + filename);
    }

    file << countNodes(root) << endl;

    writeTextNode(file, root);
    
    file.close();
}

void RBTree::writeTextNode(ofstream& file, RBTNode* node) {
    if (node == nullptr) {
        return;
    }

    file << node->key << " " << node->color << " ";
    
    file << (node->left != nullptr ? 1 : 0) << " ";
    file << (node->right != nullptr ? 1 : 0) << endl;

    writeTextNode(file, node->left);
    writeTextNode(file, node->right);
}

void RBTree::readText(const string& filename) {
    ifstream file(filename);
    if (!file) {
        throw runtime_error("Не удалось открыть файл: " + filename);
    }

    deleteBRTreeNode(root);
    root = nullptr;

    int nodeCount;
    file >> nodeCount;
    file.ignore();

    if (nodeCount > 0) {
        readTextNode(file, root, nullptr);
    }

    file.close();
}

void RBTree::readTextNode(ifstream& file, RBTNode*& node, RBTNode* parent) {
    string line;
    getline(file, line);
    
    if (line == "NULL") {
        node = nullptr;
        return;
    }

    istringstream iss(line);
    node = new RBTNode();
    node->parent = parent;

    int hasLeft, hasRight;
    iss >> node->key >> node->color >> hasLeft >> hasRight;

    if (hasLeft) {
        readTextNode(file, node->left, node);
    } else {
        node->left = nullptr;
    }
    
    if (hasRight) {
        readTextNode(file, node->right, node);
    } else {
        node->right = nullptr;
    }
}

int RBTree::countNodes(RBTNode* node) {
    if (node == nullptr) {
        return 0;
    }
    return 1 + countNodes(node->left) + countNodes(node->right);
}

void RBTree::printInOrder(RBTNode* node) {
    if (node != nullptr) {
        printInOrder(node->left);
        string colorStr = node->color ? "R" : "B";
        cout << node->key << "(" << colorStr << ") ";
        printInOrder(node->right);
    }
}

void RBTree::print() {
    if (root == nullptr) {
        cout << "Дерево пустое" << endl;
        return;
    }
    cout << "In-order обход: ";
    printInOrder(root);
    cout << endl;
}