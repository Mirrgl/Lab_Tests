#include "BlackRedTree.h"

int main() {
    RBTree a;
    a.insert(5);
    a.insert(7);
    a.insert(2);
    a.insert(6);
    a.insert(8);
    a.print();
    a.writeBinary("test");
    RBTree b;
    b.readBinary("test");
    b.print();

    RBTree c;
    c.insert(5);
    c.insert(7);
    c.insert(2);
    c.insert(6);
    c.insert(8);
    c.print();
    c.writeText("test2");
    RBTree d;
    d.readText("test2");
    d.print();

    return 0;
}