#include "Queue.h"

int main() {
    Queue a = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    a.print();
    a.writeBinary("test");
    Queue b(0);
    b.readBinary("test");
    b.print();

    Queue c = {"HELP", "I CANT", "HOLD IT", "ANYMORE"};
    c.print();
    c.writeText("test2");
    Queue d(0);
    d.readText("test2");
    d.print();

    return 0;
}