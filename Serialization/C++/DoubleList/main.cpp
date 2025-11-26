#include "DoubleList.h"

int main() {
    DoubleList a = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    a.print();
    a.writeBinary("test");
    DoubleList b;
    b.readBinary("test");
    b.print();

    DoubleList c = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    c.print();
    c.writeText("test2");
    DoubleList d;
    d.readText("test2");
    d.print();

    return 0;
}