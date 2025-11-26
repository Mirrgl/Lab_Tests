#ifndef BRT_H
#define BRT_H

#include <string>
#include <iostream>
#include <stdexcept>
#include <fstream>
#include <sstream>

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

        void writeBinaryNode(ofstream& file, RBTNode* node);
        
        void readBinaryNode(ifstream& file, RBTNode*& node, RBTNode* parent);
        
        void writeTextNode(ofstream& file, RBTNode* node);
        
        void readTextNode(ifstream& file, RBTNode*& node, RBTNode* parent);

        int countNodes(RBTNode* node);
    
    public:
        RBTree();

        ~RBTree();

        void insert(int key);

        void del(int key);

        int get(int key);

        void writeBinary(const string& filename);

        void readBinary(const string& filename);
        
        void writeText(const string& filename);
        
        void readText(const string& filename);

        void printInOrder(RBTNode* node);

        void print();
};

#endif