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