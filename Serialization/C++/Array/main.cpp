#include "Array.h"

int main() {
    Array a = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    a.print();
    a.writeBinary("test");
    Array b(1);
    b.readBinary("test");
    b.print();

    Array c = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    c.print();
    c.writeText("test2");
    Array d(1);
    d.readText("test2");
    d.print();

    return 0;
}