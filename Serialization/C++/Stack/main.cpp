#include "Stack.h"

int main() {
    Stack a = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    a.print();
    a.writeBinary("test");
    Stack b;
    b.readBinary("test");
    b.print();

    Stack c = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    c.print();
    c.writeText("test2");
    Stack d;
    d.readText("test2");
    d.print();

    return 0;
}