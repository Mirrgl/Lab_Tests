#ifndef BRT_H
#define BRT_H

#include <string>
#include <iostream>
#include <stdexcept>

using namespace std;

struct RBTNode {
    int key;
    RBTNode* right;
    RBTNode* left;
    RBTNode* parent;
    bool color;
};

class RBTree {
    private:
        RBTNode* root;

        void leftRotate(RBTNode* x);

        void rightRotate(RBTNode* y);

        void fixViolation(RBTNode* z);

        RBTNode* findNodePosition(RBTNode* currentBRTNode, int key);

        void transplant(RBTNode* u, RBTNode* v);

        void fixDelete(RBTNode* x, RBTNode* xParent);

        RBTNode* findNode(int key);

        void deleteBRTreeNode(RBTNode* currentNode);
    
    public:
        RBTree();

        ~RBTree();

        void insert(int key);

        void del(int key);

        int get(int key);
};

#endif